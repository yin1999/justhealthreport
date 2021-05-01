package httpclient

import (
	"net/http"
	"net/url"
)

type customCookieJar interface {
	http.CookieJar
	GetCookieByDomain(domain string) []*http.Cookie
	SameSite(flag bool)
	SetNewFilter(name []string)
	GetCookieByName(name []string) []*http.Cookie
}

type cookieJar struct {
	cookies   []*http.Cookie
	cookieMap map[string]bool
	sameSite  bool
}

// newCookieJar sameSite default to true
func newCookieJar(name []string) customCookieJar {
	m := make(map[string]bool, len(name))
	for i := range name {
		m[name[i]] = true
	}
	return &cookieJar{
		cookieMap: m,
		sameSite:  true,
	}
}

func (j *cookieJar) SameSite(flag bool) {
	j.sameSite = flag
}

func (j *cookieJar) SetNewFilter(name []string) {
	j.cookieMap = make(map[string]bool, len(name))
	for i := range name {
		j.cookieMap[name[i]] = true
	}
}

func (j *cookieJar) GetCookieByName(name []string) (res []*http.Cookie) {
	m := make(map[string]bool, len(name))
	for i := range name {
		m[name[i]] = true
	}
	for i := range j.cookies {
		if m[j.cookies[i].Name] {
			res = append(res, j.cookies[i])
		}
	}
	return
}

func (j *cookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	for i := range cookies {
		if ok := j.cookieMap[cookies[i].Name]; ok {
			cookies[i].Domain = u.Hostname()
			j.cookies = append(j.cookies, cookies[i])
		}
	}
}

func (j *cookieJar) GetCookieByDomain(domain string) (res []*http.Cookie) {
	for i := range j.cookies {
		if j.cookies[i].Domain == domain {
			res = append(res, j.cookies[i])
		}
	}
	return
}

func (j *cookieJar) Cookies(u *url.URL) (res []*http.Cookie) {
	if j.sameSite {
		for i := range j.cookies {
			if j.cookies[i].Domain == u.Hostname() {
				res = append(res, j.cookies[i])
			}
		}
	} else {
		res = j.cookies
	}
	return
}
