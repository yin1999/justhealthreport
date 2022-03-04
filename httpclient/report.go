package httpclient

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var ErrCannotParseData = errors.New("httpclient: cannot parse data")

func getFormData(ctx context.Context, jar customCookieJar) (form *HealthForm, err error) {
	var req *http.Request
	req, err = getWithContext(ctx, "http://ehall.just.edu.cn/default/work/jkd/jkxxtb/jkxxcj.jsp")
	if err != nil {
		return
	}
	client := http.Client{
		Jar: jar,
	}

	var res *http.Response
	if res, err = client.Do(req); err != nil {
		return
	}

	defer drainBody(res.Body)

	reader := bufio.NewReader(res.Body)

	var line string
	for !strings.HasPrefix(line, "$(\"") && err == nil {
		line, err = scanLine(reader)
	}

	var e *element
	form = &HealthForm{
		Form: Entity{
			Time:                 time.Now().In(timeZone).Format("2006-01-02 15:04"),
			MorningTemperature:   fmt.Sprintf("36.%d", randomNumberGen.Int31n(10)),
			LastNightTemperature: fmt.Sprintf("36.%d", randomNumberGen.Int31n(10)),
			Ext:                  "{}", // default
		},
	}
	var filler *structFiller
	filler, err = newFiller(&(form.Form), "fill")
	if err != nil {
		form = nil
		return
	}
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
		form = nil
		err = ErrCannotParseData
	} else {
		err = nil
	}
	return
}

func postForm(ctx context.Context, jar http.CookieJar, form *HealthForm) error {
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
	req.Header = generalHeaders.Clone()
	req.Header.Set("Content-Type", "text/json")

	client := &http.Client{Jar: jar}

	var res *http.Response
	if res, err = client.Do(req); err != nil {
		return err
	}
	defer drainBody(res.Body)

	r := &response{}

	if err = json.NewDecoder(res.Body).Decode(r); err != nil || !r.Res {
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
