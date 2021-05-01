package httpclient

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var ErrCannotParseData = errors.New("httpclient: cannot parse data")

func getFormData(ctx context.Context, jar customCookieJar) (form *HealthForm, cookie []*http.Cookie, err error) {
	var req *http.Request
	req, err = getWithContext(ctx, "http://ehall.just.edu.cn/default/work/jkd/jkxxtb/jkxxcj.jsp")
	if err != nil {
		return
	}
	jar.SameSite(true)
	client := http.Client{
		Jar: jar,
	}

	var res *http.Response
	if res, err = client.Do(req); err != nil {
		return
	}

	var body io.ReadCloser
	if body, err = responseReader(res); err != nil {
		return
	}
	defer body.Close()

	reader := bufio.NewReader(body)

	var line string
	for !strings.HasPrefix(line, "$(\"") && err == nil {
		line, err = scanLine(reader)
	}

	var e *element
	form = &HealthForm{
		Form: Entity{
			Time:                 time.Now().In(timeZone).Format("2006-01-02 15:04"),
			MorningTemperature:   fmt.Sprintf("%.1f", randomNumberGen.Float32()+36),
			LastNightTemperature: fmt.Sprintf("%.1f", randomNumberGen.Float32()+36),
			Ext:                  "{}", // default
		},
	}

	filler := newFiller(&(form.Form), "fill")
	for ; ; line, _ = scanLine(reader) {
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "$(\"") {
			break
		}
		e, err = parseData(line)
		if err != nil {
			break
		}
		filler.fill(e.key, e.value)
	}
	if !strings.HasPrefix(line, "$('div[name=sqrid]')") {
		err = ErrCannotParseData
		return
	}
	err = nil
	cookie = jar.GetCookieByDomain("ehall.just.edu.cn")
	return
}

func postForm(ctx context.Context, cookie []*http.Cookie, form *HealthForm) error {
	data, err := json.Marshal(form)
	if err != nil {
		return err
	}

	var req *http.Request
	req, err = http.NewRequestWithContext(ctx,
		http.MethodPost,
		"http://ehall.just.edu.cn/default/work/jkd/jkxxtb/com.sudytech.work.suda.jkxxtb.jktbSave.save.biz.ext",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return err
	}
	setGeneralHeader(req)
	setCookies(req, cookie)
	req.Header.Set("Content-Type", "text/json")

	var res *http.Response
	if res, err = http.DefaultClient.Do(req); err != nil {
		return err
	}
	var reader io.ReadCloser
	reader, err = responseReader(res)
	if err != nil {
		return err
	}
	defer reader.Close()

	var body []byte
	if body, err = io.ReadAll(reader); err != nil {
		return err
	}
	r := &response{}
	json.Unmarshal(body, r)
	if !r.Res {
		return errors.New("httpclient: post form failed")
	}
	return nil
}

type response struct {
	Res bool `json:"res"`
}

type element struct {
	key   string
	value string
}

func parseData(line string) (e *element, err error) {
	index := strings.Index(line, "\")")
	valueIndexLeft := strings.IndexByte(line, '\'')
	if index <= 13 || valueIndexLeft <= 13 {
		err = ErrCannotParseData
		return
	}
	e = &element{
		key:   line[12 : index-1],
		value: line[valueIndexLeft+1 : len(line)-3],
	}
	return
}
