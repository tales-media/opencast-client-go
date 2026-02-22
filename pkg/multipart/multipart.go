/*
Copyright 2025 shio solutions GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package multipart

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/textproto"
	"os"
	"path"
	"slices"
	"strings"
	"sync"
)

type Multipart struct {
	parts    []Part
	boundary string
}

func New() *Multipart {
	return &Multipart{
		parts:    make([]Part, 0),
		boundary: RandomBoundary(),
	}
}

func (mp *Multipart) Boundary() string {
	return mp.boundary
}

func (mp *Multipart) AddPart(p Part) {
	mp.parts = append(mp.parts, p)
}

func (mp *Multipart) AddParts(p ...Part) {
	mp.parts = append(mp.parts, p...)
}

func (mp *Multipart) Len() int64 {
	n := int64(
		len(mp.parts)*(2+len(mp.boundary)+2) + // `--` boundary `\r\n`
			(2 + len(mp.boundary) + 4), // `--` boundary `--\r\n`
	)
	for _, part := range mp.parts {
		partLen := part.Len()
		if partLen < 0 {
			return -1
		}
		n = n + partLen
	}
	return n
}

type reader struct {
	mp *Multipart
	pr *io.PipeReader
	pw *io.PipeWriter
}

var _ io.ReadCloser = &reader{}

func Reader(mp *Multipart) *reader {
	pr, pw := io.Pipe()
	r := &reader{
		mp: mp,
		pr: pr,
		pw: pw,
	}
	go func() { _ = r.write() }()
	return r
}

func (r *reader) write() error {
	for _, part := range r.mp.parts {
		// boundary
		if _, err := fmt.Fprintf(r.pw, "--%s\r\n", r.mp.boundary); err != nil {
			return r.handleErr(err)
		}

		// header
		header := part.GetHeader()
		for _, k := range slices.Sorted((maps.Keys(header))) {
			for _, v := range header[k] {
				if _, err := fmt.Fprintf(r.pw, "%s: %s\r\n", k, v); err != nil {
					return r.handleErr(err)
				}
			}
		}
		if _, err := fmt.Fprintf(r.pw, "\r\n"); err != nil {
			return r.handleErr(err)
		}

		// body
		body, err := part.GetBody()
		if err != nil {
			return r.handleErr(err)
		}
		defer func() { _ = body.Close() }()
		if _, err = io.Copy(r.pw, body); err != nil {
			return r.handleErr(err)
		}
		if _, err = fmt.Fprintf(r.pw, "\r\n"); err != nil {
			return r.handleErr(err)
		}
	}

	// finishing boundary
	if _, err := fmt.Fprintf(r.pw, "--%s--\r\n", r.mp.boundary); err != nil {
		return r.handleErr(err)
	}

	// EOF
	return r.pw.Close()
}

func (r *reader) handleErr(err error) error {
	r.pr.CloseWithError(err)
	return err
}

func (r *reader) Read(p []byte) (int, error) {
	return r.pr.Read(p)
}

func (r *reader) Close() error {
	return r.pr.Close()
}

type Part interface {
	GetHeader() textproto.MIMEHeader
	GetBody() (io.ReadCloser, error)
	Len() int64
}

type BytePart struct {
	Header textproto.MIMEHeader
	Body   []byte
}

var _ Part = &BytePart{}

func FormField(fieldname string, value []byte) *BytePart {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"`, EscapeQuotes(fieldname)))
	return &BytePart{
		Header: h,
		Body:   value,
	}
}

func (p *BytePart) GetHeader() textproto.MIMEHeader {
	return p.Header
}

func (p *BytePart) GetBody() (io.ReadCloser, error) {
	r := bytes.NewReader(p.Body)
	return io.NopCloser(r), nil
}

func (p *BytePart) Len() int64 {
	return HeaderLen(p.Header) + // header
		2 + // `\r\n`
		int64(len(p.Body)) + // body
		2 // `\r\n`
}

type StringPart struct {
	Header textproto.MIMEHeader
	Body   string
}

var _ Part = &StringPart{}

func FormFieldString(fieldname, value string) *StringPart {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"`, EscapeQuotes(fieldname)))
	return &StringPart{
		Header: h,
		Body:   value,
	}
}

func (p *StringPart) GetHeader() textproto.MIMEHeader {
	return p.Header
}

func (p *StringPart) GetBody() (io.ReadCloser, error) {
	r := strings.NewReader(p.Body)
	return io.NopCloser(r), nil
}

func (p *StringPart) Len() int64 {
	return HeaderLen(p.Header) + // header
		2 + // `\r\n`
		int64(len(p.Body)) + // body
		2 // `\r\n`
}

type FilePart struct {
	Header textproto.MIMEHeader
	File   string
}

var _ Part = &FilePart{}

func File(fieldname, file string) *FilePart {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, EscapeQuotes(fieldname), EscapeQuotes(path.Base(file))))
	h.Set("Content-Type", "application/octet-stream")
	return &FilePart{
		Header: h,
		File:   file,
	}
}

func (p *FilePart) GetHeader() textproto.MIMEHeader {
	return p.Header
}

func (p *FilePart) GetBody() (io.ReadCloser, error) {
	return os.Open(p.File)
}

func (p *FilePart) Len() int64 {
	fi, err := os.Stat(p.File)
	if err != nil {
		return -1
	}
	return HeaderLen(p.Header) + // header
		2 + // `\r\n`
		fi.Size() + // body
		2 // `\r\n`
}

type StreamPart struct {
	Header textproto.MIMEHeader
	Reader io.ReadCloser
	once   sync.Once
}

var _ Part = &StreamPart{}

func Stream(fieldname, filename string, r io.ReadCloser) *StreamPart {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, EscapeQuotes(fieldname), EscapeQuotes(filename)))
	h.Set("Content-Type", "application/octet-stream")
	return &StreamPart{
		Header: h,
		Reader: r,
		once:   sync.Once{},
	}
}

func (p *StreamPart) GetHeader() textproto.MIMEHeader {
	return p.Header
}

func (p *StreamPart) GetBody() (io.ReadCloser, error) {
	var r io.ReadCloser
	p.once.Do(func() { r = p.Reader })
	if r == nil {
		return nil, errors.New("multipart: body of stream part can only be accessed once")
	}
	return r, nil
}

func (p *StreamPart) Len() int64 {
	return -1
}

func HeaderLen(h textproto.MIMEHeader) int64 {
	n := 0
	for key, list := range h {
		for _, val := range list {
			n = n +
				len(key) + // `key`
				2 + // `: `
				len(val) + // value
				2 // `\r\n`
		}
	}
	return int64(n)
}

// Code below adopted from mime/multipart

var quoteEscaper = strings.NewReplacer(`\`, `\\`, `"`, `\"`)

func EscapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func RandomBoundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}
