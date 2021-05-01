package httpclient

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type header struct {
	key   string
	value string
}

var generalHeaders = [...]header{
	{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"},
	{"Accept-Encoding", "gzip"},
	{"Accept-Language", "zh-CN,zh;q=0.9"},
	{"Connection", "keep-alive"},
	{"User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36"},
}

func postWithContext(ctx context.Context, url string, data url.Values) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		url,
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, err
	}
	setGeneralHeader(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, err
}

func getWithContext(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		url,
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}
	setGeneralHeader(req)
	return req, err
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

func trimSuffixSpace(data []byte) []byte {
	start := 0
	for start < len(data) && asciiSpace[data[start]] == 1 {
		start++
	}
	return data[start:]
}

// scanLine scan a line
func scanLine(reader *bufio.Reader) (string, error) {

	data, isPrefix, err := reader.ReadLine() // data is not a copy, use it carefully
	res := string(trimSuffixSpace(data))     // copy the data to string(remove the leading space)
	for isPrefix {
		_, isPrefix, err = reader.ReadLine()
	}

	return res, err
}

func setCookies(req *http.Request, cookies []*http.Cookie) {
	for i := range cookies {
		req.AddCookie(cookies[i])
	}
}

func setGeneralHeader(req *http.Request) {
	for i := range generalHeaders {
		req.Header.Set(generalHeaders[i].key, generalHeaders[i].value)
	}
}

type closeFunc func() error

type resReader struct {
	io.Reader
	close []closeFunc
}

func (r *resReader) Close() error {
	for i := range r.close {
		r.close[i]()
	}
	return nil
}

// responseReader provide the response reader,
// if occur an error, it will close the response.Body
func responseReader(res *http.Response) (io.ReadCloser, error) {
	var err error
	defer func() {
		if err != nil {
			res.Body.Close()
		}
	}()
	r := &resReader{}
	encoding := res.Header.Get("content-encoding")
	switch encoding {
	case "gzip":
		var reader *gzip.Reader
		reader, err := gzip.NewReader(res.Body)
		if err != nil {
			return nil, err
		}
		r.Reader = reader
		r.close = []closeFunc{reader.Close, res.Body.Close}
	case "":
		r.Reader = res.Body
		r.close = []closeFunc{res.Body.Close}
	default:
		return nil, fmt.Errorf("reader: unsupported encoding: %s", encoding)
	}
	return r, nil
}
