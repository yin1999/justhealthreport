package httpclient

import (
	"context"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

var (
	// timeZone is used for set DataTime in HealthForm,
	// default: CTS(China Standard Time)
	timeZone = time.FixedZone("CST", 8*3600)

	randomNumberGen = rand.New(rand.NewSource(time.Now().Unix()))
)

// LoginConfirm 验证账号密码
func LoginConfirm(ctx context.Context, account [2]string, timeout time.Duration) error {
	var cc context.CancelFunc
	ctx, cc = context.WithTimeout(ctx, timeout)
	_, err := login(ctx, account)
	cc()
	return parseURLError(err)
}

// Punch 打卡
func Punch(ctx context.Context, account [2]string, timeout time.Duration) (err error) {
	var cc context.CancelFunc
	ctx, cc = context.WithTimeout(ctx, timeout)
	defer cc()

	defer func() {
		err = parseURLError(err)
	}()

	var jar customCookieJar
	jar, err = login(ctx, account)
	if err != nil {
		return
	}

	var form *HealthForm
	var cookie []*http.Cookie
	form, cookie, err = getFormData(ctx, jar)
	if err != nil {
		return
	}

	err = postForm(ctx, cookie, form)
	return
}

// SetTimeZone 设置时区
func SetTimeZone(tz *time.Location) {
	if tz != nil {
		timeZone = tz
	}
}

// parseURLError 解析URL错误
func parseURLError(err error) error {
	if err == nil {
		return err
	}
	if v, ok := err.(*url.Error); ok {
		return v.Err
	}
	return err
}
