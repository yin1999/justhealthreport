package httpclient

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
)

var generalHeaders = http.Header{
	"Accept":          []string{"text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"},
	"Accept-Language": []string{"zh-CN,zh;q=0.9"},
	"Connection":      []string{"keep-alive"},
	"User-Agent":      []string{"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.0.0 Safari/537.36"},
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
	req.Header = generalHeaders.Clone()
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
	req.Header = generalHeaders.Clone()
	return req, err
}

// scanLine scan a line
func scanLine(reader *bufio.Reader) (string, error) {
	data, isPrefix, err := reader.ReadLine() // data is not a copy, use it carefully
	res := string(textproto.TrimBytes(data)) // copy the data to string(remove the leading space)
	for isPrefix {                           // discard the remaining runes in the line
		_, isPrefix, err = reader.ReadLine()
	}
	return res, err
}

func drainBody(body io.ReadCloser) {
	io.Copy(io.Discard, body)
	body.Close()
}
