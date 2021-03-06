package httpclient

import (
	"bufio"
	"context"
	"encoding/xml"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"
)

type loginForm struct {
	Username string `url:"username"`
	Password string `url:"password"`
	Excution string `fill:"execution" url:"execution"`
	Event    string `fill:"_eventId" url:"_eventId"`
	Type     string `fill:"loginType" url:"loginType"`
}

func login(ctx context.Context, account *Account) (j customCookieJar, err error) {
	const loginURL = "http://ids2.just.edu.cn/cas/login"
	jar := newCookieJar()
	client := http.Client{
		Jar: jar,
	}
	var req *http.Request
	req, err = getWithContext(ctx, loginURL)
	if err != nil {
		return
	}

	var res *http.Response
	if res, err = client.Do(req); err != nil {
		return
	}
	defer drainBody(res.Body)

	reader := bufio.NewReaderSize(res.Body, 6000)

	var line string
	for !strings.HasPrefix(line, "<input type=\"hidden") && err == nil {
		line, err = scanLine(reader)
	}
	var element *elementInput
	form := &loginForm{}
	var filler *structFiller
	filler, err = newFiller(form, "fill")
	if err != nil {
		return
	}
	for ; strings.HasPrefix(line, "<input type=\"hidden") || line == ""; line, err = scanLine(reader) {
		if line == "" {
			continue
		}
		if element, err = elementParse(line); err != nil {
			return
		}
		filler.fill(element.Key, element.Value)
	}
	form.Username = account.Username
	form.Password = account.Password
	var value url.Values
	if value, err = query.Values(form); err != nil {
		return
	}

	if req, err = postFormWithContext(ctx, loginURL, value); err != nil {
		return
	}

	if res, err = client.Do(req); err != nil {
		return
	}
	drainBody(res.Body)
	if tmp := jar.GetCookieByName("CASTGC"); len(tmp) == 0 {
		err = CookieNotFoundErr{"CASTGC"}
		return
	}
	j = jar
	return
}

type elementInput struct {
	Key   string `xml:"name,attr"`
	Value string `xml:"value,attr"`
	ID    string `xml:"id,attr"`
}

func elementParse(v string) (*elementInput, error) {
	if len(v) < 2 {
		return nil, &xml.SyntaxError{Msg: "error format", Line: 1}
	}
	out := &elementInput{}
	data := []byte(v)
	if data[len(data)-2] != '/' {
		data = append(data[:len(data)-1], '/', '>')
	}
	err := xml.Unmarshal(data, out)
	if err != nil {
		return nil, err
	}
	if out.Key == "" {
		out.Key = out.ID
	}
	return out, err
}

type structFiller struct {
	m map[string]int
	v reflect.Value
}

// newFiller default tag: fill.
// The item must be a pointer
func newFiller(item interface{}, tag string) (*structFiller, error) {
	v := reflect.ValueOf(item).Elem()
	if !v.CanAddr() {
		return nil, errors.New("reflect: item must be a pointer")
	}
	if tag == "" {
		tag = "fill"
	}
	findTagName := func(t reflect.StructTag) (string, error) {
		if tn, ok := t.Lookup(tag); ok && len(tn) > 0 {
			return strings.Split(tn, ",")[0], nil
		}
		return "", errors.New("skip")
	}
	s := &structFiller{
		m: make(map[string]int),
		v: v,
	}
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		name, err := findTagName(typeField.Tag)
		if err != nil {
			continue
		}
		s.m[name] = i
	}
	return s, nil
}

func (s *structFiller) fill(key string, value interface{}) error {
	fieldNum, ok := s.m[key]
	if !ok {
		return errors.New("reflect: field <" + key + "> not exists")
	}
	s.v.Field(fieldNum).Set(reflect.ValueOf(value))
	return nil
}
