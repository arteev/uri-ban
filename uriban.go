package uriban

import (
	"net/url"
	"strings"
)

type UrlPart int

type BanMode func(UrlPart, string) string

//Url part types
const (
	Username UrlPart = iota
	Password
	Scheme
	Host
	Path
	Query
	Fragment
	All
)

func ModeHidden() BanMode {
	return func(UrlPart, string) string {
		return ""
	}
}

func ModeNothing() BanMode {
	return func(p UrlPart, s string) string {
		return s
	}
}

func ModeStarred(count int) BanMode {
	return func(UrlPart, string) string {
		return strings.Repeat("*", count)
	}
}

func ModeValue(value string) BanMode {
	return func(UrlPart, string) string {
		return value
	}
}

func ModeFunc(f func(string) string) BanMode {
	return func(p UrlPart, s string) string {
		return f(s)
	}
}

type Option func(*map[UrlPart]BanMode)

func WithOption(p UrlPart, mode BanMode) Option {
	return func(m *map[UrlPart]BanMode) {
		(*m)[p] = mode
	}
}

func replaceByOpt(cur UrlPart, s string, opts map[UrlPart]BanMode) string {
	if s == "" {
		return s
	}
	if m, exists := opts[cur]; exists {
		return m(cur, s)
	}
	return s
}

func Ban(s string, opts ...Option) string {
	//init opts
	mo := make(map[UrlPart]BanMode)
	for _, opt := range opts {
		opt(&mo)
	}

	u, err := url.ParseRequestURI(s)
	if err != nil {
		return replaceByOpt(All, s, mo)
	}

	if u.User != nil {
		user := u.User
		username := replaceByOpt(Username, u.User.Username(), mo)

		if p, ok := u.User.Password(); ok {
			pwd := replaceByOpt(Password, p, mo)

			if pwd == "" {
				user = url.User(username)
			} else {
				user = url.UserPassword(username, pwd)
			}
		}

		u.User = user
	}

	u.Scheme = replaceByOpt(Scheme, u.Scheme, mo)
	u.Path = replaceByOpt(Path, u.Path, mo)
	u.RawQuery = replaceByOpt(Query, u.RawQuery, mo)
	u.Fragment = replaceByOpt(Fragment, u.Fragment, mo)

	res, err := url.PathUnescape(u.String())
	if err != nil {
		return replaceByOpt(All, s, mo)
	}

	return res
}
