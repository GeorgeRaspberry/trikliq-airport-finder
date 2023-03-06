package pdf

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"sync"
	"syscall"
	"unsafe"
)

// ListDir defines list dir structure
type ListDir struct {
	Path string
	os.FileInfo
}

type readDirOpts struct {
	// The maximum number of entries to return
	count int
	// Follow directory symlink
	followDirSymlink bool
}

// The buffer must be at least a block long.
// refer https://github.com/golang/go/issues/24015
const blockSize = 8 << 10 // 8192

//lint:ignore GLOBAL this is okay
var (
	direntPool = sync.Pool{
		New: func() interface{} {
			buf := make([]byte, blockSize*128)
			return &buf
		},
	}

	direntNamePool = sync.Pool{
		New: func() interface{} {
			buf := make([]byte, blockSize)
			return &buf
		},
	}
)

// List will retrieve recursive list of all files within path using getdents
func List(path string, includeDir, includeFileInfo bool) (results []ListDir, err error) {
	opts := readDirOpts{
		count:            -1,
		followDirSymlink: true,
	}

	files, err := readDirWithOpts(path, opts)
	if err != nil {
		return
	}

	for _, file := range files {

		obj := ListDir{}
		obj.Path = file

		if includeFileInfo {
			fileInfo, err := os.Stat(file)
			if err != nil {
				continue
			}
			obj.FileInfo = fileInfo
			mode := fileInfo.Mode()
			if includeDir {
				results = append(results, obj)
			} else {
				if mode.IsRegular() {
					results = append(results, obj)
				}
			}
		} else {
			results = append(results, obj)
		}

	}

	return
}

// Return count entries at the directory dirPath and all entries
// if count is set to -1
func readDirWithOpts(dirPath string, opts readDirOpts) (entries []string, err error) {
	entries = make([]string, 0, 4096)

	f, err := os.Open(dirPath)
	if err != nil {
		if os.IsPermission(err) {
			return nil, nil
		}
		return nil, err

	}
	defer f.Close()

	bufp := direntPool.Get().(*[]byte)
	defer direntPool.Put(bufp)
	buf := *bufp

	nameTmp := direntNamePool.Get().(*[]byte)
	defer direntNamePool.Put(nameTmp)
	tmp := *nameTmp

	// starting read position in buf
	boff := 0

	// end valid data in buf
	nbuf := 0

	count := opts.count

	for count != 0 {
		if boff >= nbuf {
			boff = 0
			nbuf, err = syscall.ReadDirent(int(f.Fd()), buf)
			if err != nil {
				return nil, err
			}

			if nbuf <= 0 {
				break
			}
		}
		consumed, name, typ, err := parseDirEnt(buf[boff:nbuf])
		if err != nil {
			return nil, err
		}
		boff += consumed

		// if no name, skip it
		if len(name) == 0 {
			continue
		}

		// if . skip it
		if bytes.Equal(name, []byte{'.'}) {
			continue
		}

		// if .. skip it
		if bytes.Equal(name, []byte{'.', '.'}) {
			continue
		}

		// Fallback for filesystems (like old XFS) that don't
		// support Dirent.Type and have DT_UNKNOWN (0) there
		// instead.
		if typ == unexpectedFileMode || typ&os.ModeSymlink == os.ModeSymlink {
			fi, err := os.Stat(path.Join(dirPath, string(name)))
			if err != nil {
				// It got deleted in the meantime, not found
				// or returns too many symlinks ignore this
				// file/directory.

				if errors.Is(err, os.ErrNotExist) {
					continue
				}

				var pathErr *os.PathError
				if errors.As(err, &pathErr) {
					if pathErr.Err == syscall.ENOENT {
						continue
					}
				}

				if errors.Is(err, syscall.ELOOP) {
					continue
				}

				return nil, err
			}

			// Ignore symlinked directories.
			if !opts.followDirSymlink && typ&os.ModeSymlink == os.ModeSymlink && fi.IsDir() {
				continue
			}

			typ = fi.Mode() & os.ModeType
		}

		var nameStr string
		if typ.IsRegular() {
			nameStr = path.Join(dirPath, string(name))

		} else if typ.IsDir() {
			// Use temp buffer to append a slash to avoid string concat.
			tmp = tmp[:len(name)+1]
			copy(tmp, name)
			tmp[len(tmp)-1] = '/' // SlashSeparator

			nameStr = path.Join(dirPath, string(tmp))

			ent, _ := readDirWithOpts(nameStr, opts)
			entries = append(entries, ent...)
		}

		count--
		entries = append(entries, nameStr)
	}

	return
}

// unexpectedFileMode is a sentinel (and bogus) os.FileMode
// value used to represent a syscall.DT_UNKNOWN Dirent.Type.
const unexpectedFileMode os.FileMode = os.ModeNamedPipe | os.ModeSocket | os.ModeDevice

// parseDirEnt parses dir entry from buffer
func parseDirEnt(buf []byte) (consumed int, name []byte, typ os.FileMode, err error) {
	// golang.org/issue/15653
	dirent := (*syscall.Dirent)(unsafe.Pointer(&buf[0]))
	if v := unsafe.Offsetof(dirent.Reclen) + unsafe.Sizeof(dirent.Reclen); uintptr(len(buf)) < v {
		return consumed, nil, typ, fmt.Errorf("buf size of %d smaller than dirent header size %d", len(buf), v)
	}
	if len(buf) < int(dirent.Reclen) {
		return consumed, nil, typ, fmt.Errorf("buf size %d < record length %d", len(buf), dirent.Reclen)
	}
	consumed = int(dirent.Reclen)
	if dirent.Ino == 0 { // File absent in directory.
		return
	}
	switch dirent.Type {
	case syscall.DT_REG:
		typ = 0
	case syscall.DT_DIR:
		typ = os.ModeDir
	case syscall.DT_LNK:
		typ = os.ModeSymlink
	default:
		// Skip all other file types. Revisit if/when this code needs
		// to handle such files, MinIO is only interested in
		// files and directories.
		typ = unexpectedFileMode
	}

	nameBuf := (*[unsafe.Sizeof(dirent.Name)]byte)(unsafe.Pointer(&dirent.Name[0]))
	nameLen, err := direntNamlen(dirent)
	if err != nil {
		return consumed, nil, typ, err
	}

	return consumed, nameBuf[:nameLen], typ, nil
}

// direntNamlen gets name len from sys call
func direntNamlen(dirent *syscall.Dirent) (uint64, error) {
	const fixedHdr = uint16(unsafe.Offsetof(syscall.Dirent{}.Name))
	nameBuf := (*[unsafe.Sizeof(dirent.Name)]byte)(unsafe.Pointer(&dirent.Name[0]))
	const nameBufLen = uint16(len(nameBuf))
	limit := dirent.Reclen - fixedHdr
	if limit > nameBufLen {
		limit = nameBufLen
	}
	// Avoid bugs in long file names
	// https://github.com/golang/tools/commit/5f9a5413737ba4b4f692214aebee582b47c8be74
	nameLen := bytes.IndexByte(nameBuf[:limit], 0)
	if nameLen < 0 {
		return 0, fmt.Errorf("failed to find terminating 0 byte in dirent")
	}
	return uint64(nameLen), nil
}
