package httpclient

import (
	"bufio"
	"compress/gzip"
	"context"
	"errors"
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
	{"Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
	{"Accept-Encoding", "gzip"},
	{"Accept-Language", "zh-CN,zh;q=0.9"},
	{"Connection", "keep-alive"},
	{"User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36"},
}

func postFormWithContext(ctx context.Context, url string, data url.Values) (*http.Request, error) {
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

var isAsciiSpace = [256]bool{'\t': true, '\n': true, '\v': true, '\f': true, '\r': true, ' ': true}

func trimSuffixSpace(data []byte) []byte {
	start := 0
	for start < len(data) && isAsciiSpace[data[start]] {
		start++
	}
	return data[start:]
}

// scanLine scan a line
func scanLine(reader *bufio.Reader) (string, error) {
	data, isPrefix, err := reader.ReadLine() // data is not a copy, use it carefully
	res := string(trimSuffixSpace(data))     // copy the data to string(remove the leading space)
	for isPrefix {                           // discard the remaining runes in the line
		_, isPrefix, err = reader.ReadLine()
	}
	return res, err
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
	var (
		err error
		r   io.ReadCloser

		encoding = res.Header.Get("Content-Encoding")
	)

	switch encoding {
	case "gzip":
		var reader *gzip.Reader
		reader, err = gzip.NewReader(res.Body)
		if err == nil {
			r = &resReader{
				Reader: reader,
				close:  []closeFunc{reader.Close, res.Body.Close},
			}
		}
	case "", "identity":
		r = res.Body
	default:
		err = errors.New("reader: unsupported encoding: " + encoding)
	}
	if err != nil {
		res.Body.Close()
	}
	return r, err
}
