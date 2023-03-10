package model

import "net/textproto"

type Request struct {
	Data map[string]any `json:"data"`
}

type Response struct {
	Status bool     `json:"status"`
	Errors []string `json:"errors"`
	Data   any      `json:"data,omitempty"`
}

type MultipartFile struct {
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Header   textproto.MIMEHeader
	MimeType string `json:"mimeType"`
	Content  []byte `json:"content"`
}
