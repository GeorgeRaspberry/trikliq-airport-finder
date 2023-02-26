package encoders

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"bookbox-backend/pkg/logger/filters"

	logSecurity "bookbox-backend/pkg/logger/security"
	libTransform "bookbox-backend/pkg/transform"

	"github.com/fatih/color"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

func init() {
	zap.RegisterEncoder("secureConsole", func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		enc := NewSecureConsoleEncoder(cfg)
		return enc, nil
	})
}

// NewSecureConsoleEncoder exports secureConsole encoder
func NewSecureConsoleEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return newSecureConsoleEncoder(cfg)
}

func newSecureConsoleEncoder(cfg zapcore.EncoderConfig) *secureConsoleEncoder {
	return &secureConsoleEncoder{
		EncoderConfig: &cfg,
		buf:           bufferpool.Get(),
	}
}

//lint:ignore GLOBAL this is okay
var (
	colorRegular = color.New(color.FgBlue).SprintfFunc()
	colorWarning = color.New(color.Bold, color.FgRed).SprintfFunc()
	colorGood    = color.New(color.Bold, color.FgGreen).SprintfFunc()

	_secureConsolePool = sync.Pool{
		New: func() interface{} {
			return &secureConsoleEncoder{}
		},
	}

	ErrUnsupportedValueType = errors.New("unsupported value type")
)

func getEncoder() *secureConsoleEncoder {
	return _secureConsolePool.Get().(*secureConsoleEncoder)
}

func putEncoder(enc *secureConsoleEncoder) {
	enc.EncoderConfig = nil
	enc.buf = nil
	enc.namespaces = nil
	_secureConsolePool.Put(enc)
}

type secureConsoleEncoder struct {
	*zapcore.EncoderConfig
	buf        *buffer.Buffer
	namespaces []string
}

func (enc *secureConsoleEncoder) AddArray(key string, arr zapcore.ArrayMarshaler) error {
	key, valueNew := logSecurity.ObfuscateField(key, arr)
	valueObfuscated, _ := valueNew.(string)

	enc.addKey(key)
	enc.AppendString(valueObfuscated)

	return nil
}

func (enc *secureConsoleEncoder) AddObject(key string, obj zapcore.ObjectMarshaler) error {
	enc.addKey(key)
	return enc.AppendObject(obj)
}

func (enc *secureConsoleEncoder) AddBinary(key string, value []byte) {
	enc.AddString(key, base64.StdEncoding.EncodeToString(value))
}

func (enc *secureConsoleEncoder) AddByteString(key string, value []byte) {
	enc.addKey(key)
	enc.AppendByteString(value)
}

func (enc *secureConsoleEncoder) AddBool(key string, value bool) {
	enc.addKey(key)
	enc.AppendBool(value)
}

func (enc *secureConsoleEncoder) AddComplex128(key string, value complex128) {
	enc.addKey(key)
	enc.AppendComplex128(value)
}

func (enc *secureConsoleEncoder) AddDuration(key string, value time.Duration) {
	enc.addKey(key)
	enc.AppendDuration(value)
}

func (enc *secureConsoleEncoder) AddFloat64(key string, value float64) {
	enc.addKey(key)
	enc.AppendFloat64(value)
}

func (enc *secureConsoleEncoder) AddInt64(key string, value int64) {
	key, valueNew := logSecurity.ObfuscateField(key, value)
	valueObfuscated, _ := valueNew.(string)

	enc.addKey(key)
	enc.AppendString(valueObfuscated)
}

func (enc *secureConsoleEncoder) AddReflected(key string, value interface{}) error {
	rvalue := reflect.ValueOf(value)

	switch rvalue.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Struct:
		val := libTransform.AnyToString(value)

		key, valueNew := logSecurity.ObfuscateField(key, val)
		value, _ = valueNew.(string)

		enc.addKey(key)
		enc.AppendString(val)
		return nil

	case reflect.Array, reflect.Slice, reflect.Ptr:
		if rvalue.IsNil() {
			enc.AddByteString(key, nil)
			return nil
		}

		val := libTransform.AnyToString(value)

		key, valueNew := logSecurity.ObfuscateField(key, val)
		value, _ = valueNew.(string)

		enc.addKey(key)
		enc.AppendString(val)
		return nil
	}

	enc.AddString(key, fmt.Sprint(value))
	return nil
}

func (enc *secureConsoleEncoder) OpenNamespace(key string) {
	enc.namespaces = append(enc.namespaces, key)
}

func (enc *secureConsoleEncoder) AddString(key, value string) {
	key, valueNew := logSecurity.ObfuscateField(key, value)
	value, _ = valueNew.(string)

	enc.addKey(key)
	enc.AppendString(value)
}

func (enc *secureConsoleEncoder) AddTime(key string, value time.Time) {
	enc.addKey(key)
	enc.AppendTime(value)
}

func (enc *secureConsoleEncoder) AddUint64(key string, value uint64) {
	enc.addKey(key)
	enc.AppendUint64(value)
}

func (enc *secureConsoleEncoder) AppendArray(arr zapcore.ArrayMarshaler) error {
	marshaler := literalEncoder{
		EncoderConfig: enc.EncoderConfig,
		buf:           bufferpool.Get(),
	}

	err := arr.MarshalLogArray(&marshaler)
	if err == nil {
		enc.AppendString(marshaler.buf.String())
	} else {
		enc.AppendByteString(nil)
	}

	marshaler.buf.Free()
	return err

}

func (enc *secureConsoleEncoder) AppendObject(obj zapcore.ObjectMarshaler) error {
	marshaler := enc.clone()
	marshaler.namespaces = nil

	err := obj.MarshalLogObject(marshaler)
	if err == nil {
		enc.AppendString(" " + marshaler.buf.String())
	} else {
		enc.AppendByteString(nil)
	}

	marshaler.buf.Free()
	putEncoder(marshaler)
	return err
}

func (enc *secureConsoleEncoder) AppendBool(value bool) {
	if value {
		enc.AppendString("true")
	} else {
		enc.AppendString("false")
	}
}

func (enc *secureConsoleEncoder) AppendByteString(value []byte) {
	needsQuotes := bytes.IndexFunc(value, needsQuotedValueRune) != -1
	if needsQuotes {
		enc.buf.AppendByte('"')
	}

	enc.safeAddByteString(value)
	if needsQuotes {
		enc.buf.AppendByte('"')
	}
}

func (enc *secureConsoleEncoder) AppendComplex128(value complex128) {
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(value)), float64(imag(value))
	enc.buf.AppendFloat(r, 64)
	enc.buf.AppendByte('+')
	enc.buf.AppendFloat(i, 64)
	enc.buf.AppendByte('i')
}

func (enc *secureConsoleEncoder) AppendDuration(value time.Duration) {
	cur := enc.buf.Len()
	enc.EncodeDuration(value, enc)
	if cur == enc.buf.Len() {
		enc.AppendInt64(int64(value))
	}
}

func (enc *secureConsoleEncoder) AppendInt64(value int64) {
	enc.buf.AppendInt(value)
}

func (enc *secureConsoleEncoder) AppendReflected(value interface{}) error {
	rvalue := reflect.ValueOf(value)
	switch rvalue.Kind() {
	case reflect.Array, reflect.Chan, reflect.Func, reflect.Map, reflect.Slice, reflect.Struct:
		return ErrUnsupportedValueType
	case reflect.Ptr:
		if rvalue.IsNil() {
			enc.AppendByteString(nil)
			return nil
		}
		return enc.AppendReflected(rvalue.Elem().Interface())
	}

	enc.AppendString(fmt.Sprint(value))
	return nil
}

func (enc *secureConsoleEncoder) AppendString(val string) {
	enc.buf.AppendString(val)
}

func (enc *secureConsoleEncoder) AppendTime(value time.Time) {
	cur := enc.buf.Len()
	enc.EncodeTime(value, enc)
	if cur == enc.buf.Len() {
		enc.AppendInt64(value.UnixNano())
	}
}

func (enc *secureConsoleEncoder) AppendUint64(value uint64) {
	enc.buf.AppendUint(value)
}

func (enc *secureConsoleEncoder) AddComplex64(k string, v complex64) {
	enc.AddComplex128(k, complex128(v))
}

func (enc *secureConsoleEncoder) AddFloat32(k string, v float32) {
	enc.AddFloat64(k, float64(v))
}

func (enc *secureConsoleEncoder) AddInt(k string, v int) {
	enc.AddInt64(k, int64(v))
}

func (enc *secureConsoleEncoder) AddInt32(k string, v int32) {
	enc.AddInt64(k, int64(v))
}

func (enc *secureConsoleEncoder) AddInt16(k string, v int16) {
	enc.AddInt64(k, int64(v))
}

func (enc *secureConsoleEncoder) AddInt8(k string, v int8) {
	enc.AddInt64(k, int64(v))
}

func (enc *secureConsoleEncoder) AddUint(k string, v uint) {
	enc.AddUint64(k, uint64(v))
}

func (enc *secureConsoleEncoder) AddUint32(k string, v uint32) {
	enc.AddUint64(k, uint64(v))
}

func (enc *secureConsoleEncoder) AddUint16(k string, v uint16) {
	enc.AddUint64(k, uint64(v))
}

func (enc *secureConsoleEncoder) AddUint8(k string, v uint8) {
	enc.AddUint64(k, uint64(v))
}

func (enc *secureConsoleEncoder) AddUintptr(k string, v uintptr) {
	enc.AddUint64(k, uint64(v))
}

func (enc *secureConsoleEncoder) AppendComplex64(v complex64) {
	enc.AppendComplex128(complex128(v))
}

func (enc *secureConsoleEncoder) AppendFloat64(v float64) {
	enc.appendFloat(v, 64)
}

func (enc *secureConsoleEncoder) AppendFloat32(v float32) {
	enc.appendFloat(float64(v), 32)
}

func (enc *secureConsoleEncoder) AppendInt(v int) {
	enc.AppendInt64(int64(v))
}
func (enc *secureConsoleEncoder) AppendInt32(v int32) {
	enc.AppendInt64(int64(v))
}
func (enc *secureConsoleEncoder) AppendInt16(v int16) {
	enc.AppendInt64(int64(v))
}
func (enc *secureConsoleEncoder) AppendInt8(v int8) {
	enc.AppendInt64(int64(v))
}
func (enc *secureConsoleEncoder) AppendUint(v uint) {
	enc.AppendUint64(uint64(v))
}
func (enc *secureConsoleEncoder) AppendUint32(v uint32) {
	enc.AppendUint64(uint64(v))
}
func (enc *secureConsoleEncoder) AppendUint16(v uint16) {
	enc.AppendUint64(uint64(v))
}
func (enc *secureConsoleEncoder) AppendUint8(v uint8) {
	enc.AppendUint64(uint64(v))
}
func (enc *secureConsoleEncoder) AppendUintptr(v uintptr) {
	enc.AppendUint64(uint64(v))
}

func (enc *secureConsoleEncoder) Clone() zapcore.Encoder {
	clone := enc.clone()
	clone.buf.Write(enc.buf.Bytes())
	return clone
}

func (enc *secureConsoleEncoder) clone() *secureConsoleEncoder {
	clone := getEncoder()
	clone.EncoderConfig = enc.EncoderConfig
	clone.buf = bufferpool.Get()
	clone.namespaces = enc.namespaces
	return clone
}

func (enc *secureConsoleEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	final := enc.clone()

	if final.TimeKey != "" {
		final.AddTime(final.TimeKey, ent.Time)
	}

	if final.LevelKey != "" {
		final.buf.AppendString("\t")
		cur := final.buf.Len()
		final.EncodeLevel(ent.Level, final)
		if cur == final.buf.Len() {
			final.AppendString(ent.Level.String())
		}
		final.buf.AppendString("\t")
	}

	if ent.LoggerName != "" && final.NameKey != "" {
		nameEncoder := final.EncodeName
		if nameEncoder != nil {
			nameEncoder(ent.LoggerName, final)
			final.buf.AppendString("\t")
		}
	}

	if ent.Caller.Defined && final.CallerKey != "" {
		final.EncodeCaller(ent.Caller, final)
		final.buf.AppendString("\t")
	}

	if final.MessageKey != "" {
		message := ent.Message
		message = filters.MinWidth(message, " ", 32)

		message = colorGood("%s", message)

		final.AppendString(message)
		final.buf.AppendString("\t")
	}

	if enc.buf.Len() > 0 {
		final.buf.AppendByte(' ')
		final.buf.Write(enc.buf.Bytes())
	}

	addFields(final, fields)
	if ent.Stack != "" && final.StacktraceKey != "" {
		final.AddString(final.StacktraceKey, ent.Stack)
	}

	if final.LineEnding != "" {
		final.buf.AppendString(final.LineEnding)
	} else {
		final.buf.AppendString(zapcore.DefaultLineEnding)
	}

	ret := final.buf
	putEncoder(final)

	return ret,
		nil
}

func (enc *secureConsoleEncoder) addKey(key string) {
	if enc.buf.Len() > 0 {
		enc.buf.AppendByte(' ')
	}

	if strings.HasSuffix(key, "*") {
		key = key[:len(key)-1]
		key = colorWarning("%s", key)
	} else {
		key = colorRegular("%s", key)
	}

	for _, ns := range enc.namespaces {
		enc.safeAddString(ns)
		enc.buf.AppendByte('.')
	}

	enc.AppendString(key)
	enc.buf.AppendByte('=')
}

func (enc *secureConsoleEncoder) appendFloat(val float64, bitSize int) {
	switch {
	case math.IsNaN(val):
		enc.buf.AppendString(`NaN`)
	case math.IsInf(val, 1):
		enc.buf.AppendString(`+Inf`)
	case math.IsInf(val, -1):
		enc.buf.AppendString(`-Inf`)
	default:
		enc.buf.AppendFloat(val, bitSize)
	}
}

// safeAddString JSON-escapes a string and appends it to the internal buffer.
// Unlike the standard library's encoder, it doesn't attempt to protect the
// user from browser vulnerabilities or JSONP-related problems.
func (enc *secureConsoleEncoder) safeAddString(s string) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRuneInString(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.AppendString(s[i : i+size])
		i += size
	}
}

// safeAddByteString is no-alloc equivalent of safeAddString(string(s)) for s []byte.
func (enc *secureConsoleEncoder) safeAddByteString(s []byte) {
	for i := 0; i < len(s); {
		if enc.tryAddRuneSelf(s[i]) {
			i++
			continue
		}
		r, size := utf8.DecodeRune(s[i:])
		if enc.tryAddRuneError(r, size) {
			i++
			continue
		}
		enc.buf.Write(s[i : i+size])
		i += size
	}
}

// tryAddRuneSelf appends b if it is valid UTF-8 character represented in a single byte.
func (enc *secureConsoleEncoder) tryAddRuneSelf(b byte) bool {
	if b >= utf8.RuneSelf {
		return false
	}
	if 0x20 <= b && b != '\\' && b != '"' {
		enc.buf.AppendByte(b)
		return true
	}
	switch b {
	case '\\', '"':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte(b)
	case '\n':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('n')
	case '\r':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('r')
	case '\t':
		enc.buf.AppendByte('\\')
		enc.buf.AppendByte('t')
	default:
		// Encode bytes < 0x20, except for the escape sequences above.
		enc.buf.AppendString(`\u00`)
		enc.buf.AppendByte(_hex[b>>4])
		enc.buf.AppendByte(_hex[b&0xF])
	}
	return true
}

func (enc *secureConsoleEncoder) tryAddRuneError(r rune, size int) bool {
	if r == utf8.RuneError && size == 1 {
		enc.buf.AppendString(`\ufffd`)
		return true
	}
	return false
}

func needsQuotedValueRune(r rune) bool {
	return r <= ' ' || r == '=' || r == '"' || r == utf8.RuneError
}

func addFields(enc zapcore.ObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}

type literalEncoder struct {
	*zapcore.EncoderConfig
	buf *buffer.Buffer
}

func (enc *literalEncoder) AppendBool(value bool) {
	enc.addSeparator()
	if value {
		enc.AppendString("true")
	} else {
		enc.AppendString("false")
	}
}

func (enc *literalEncoder) AppendByteString(value []byte) {
	enc.addSeparator()
	enc.buf.AppendString(string(value))
}

func (enc *literalEncoder) AppendComplex128(value complex128) {
	enc.addSeparator()
	// Cast to a platform-independent, fixed-size type.
	r, i := float64(real(value)), float64(imag(value))
	enc.buf.AppendFloat(r, 64)
	enc.buf.AppendByte('+')
	enc.buf.AppendFloat(i, 64)
	enc.buf.AppendByte('i')
}

func (enc *literalEncoder) AppendComplex64(value complex64) {
	enc.AppendComplex128(complex128(value))
}

func (enc *literalEncoder) AppendFloat64(value float64) {
	enc.addSeparator()
	enc.buf.AppendFloat(value, 64)
}

func (enc *literalEncoder) AppendFloat32(value float32) {
	enc.addSeparator()
	enc.buf.AppendFloat(float64(value), 32)
}

func (enc *literalEncoder) AppendInt64(value int64) {
	enc.addSeparator()
	enc.buf.AppendInt(value)
}

func (enc *literalEncoder) AppendInt(v int) {
	enc.AppendInt64(int64(v))
}
func (enc *literalEncoder) AppendInt32(v int32) {
	enc.AppendInt64(int64(v))
}
func (enc *literalEncoder) AppendInt16(v int16) {
	enc.AppendInt64(int64(v))
}
func (enc *literalEncoder) AppendInt8(v int8) {
	enc.AppendInt64(int64(v))
}

func (enc *literalEncoder) AppendString(value string) {
	enc.addSeparator()
	enc.buf.AppendString(value)
}

func (enc *literalEncoder) AppendUint64(value uint64) {
	enc.addSeparator()
	enc.buf.AppendUint(value)
}

func (enc *literalEncoder) AppendUint(v uint) {
	enc.AppendUint64(uint64(v))
}
func (enc *literalEncoder) AppendUint32(v uint32) {
	enc.AppendUint64(uint64(v))
}
func (enc *literalEncoder) AppendUint16(v uint16) {
	enc.AppendUint64(uint64(v))
}
func (enc *literalEncoder) AppendUint8(v uint8) {
	enc.AppendUint64(uint64(v))
}
func (enc *literalEncoder) AppendUintptr(v uintptr) {
	enc.AppendUint64(uint64(v))
}

func (enc *literalEncoder) AppendDuration(value time.Duration) {
	cur := enc.buf.Len()
	enc.EncodeDuration(value, enc)
	if cur == enc.buf.Len() {
		enc.AppendInt64(int64(value))
	}
}

func (enc *literalEncoder) AppendTime(value time.Time) {
	cur := enc.buf.Len()
	enc.EncodeTime(value, enc)
	if cur == enc.buf.Len() {
		enc.AppendInt64(value.UnixNano())
	}
}

func (enc *literalEncoder) AppendArray(arr zapcore.ArrayMarshaler) error {
	return arr.MarshalLogArray(enc)
}

func (enc *literalEncoder) AppendObject(zapcore.ObjectMarshaler) error {
	return ErrUnsupportedValueType
}

func (enc *literalEncoder) AppendReflected(value interface{}) error {
	return ErrUnsupportedValueType
}

func (enc *literalEncoder) addSeparator() {
	if enc.buf.Len() > 0 {
		enc.buf.AppendByte(',')
	}
}
