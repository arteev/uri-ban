package uriban

import (
	"net/url"
	"strings"
)

//Part of uri
type Part int

//Mode - Replacement mode
type Mode func(Part, string) string

//Types of parts uri
const (
	Username Part = iota
	Password
	Scheme
	Host
	Path
	Query
	Fragment
	All
)

//ModeHidden returns empty-line replacement mode
func ModeHidden() Mode {
	return func(Part, string) string {
		return ""
	}
}

//ModeNothing returns mode without any replacement
func ModeNothing() Mode {
	return func(p Part, s string) string {
		return s
	}
}

//ModeStarred returns the substitution mode with the symbol *
func ModeStarred(count int) Mode {
	return func(Part, string) string {
		return strings.Repeat("*", count)
	}
}

//ModeValue returns the substitution mode with a value
func ModeValue(value string) Mode {
	return func(Part, string) string {
		return value
	}
}

//ModeFunc returns the mode that defines the function "f"
func ModeFunc(f func(string) string) Mode {
	return func(p Part, s string) string {
		return f(s)
	}
}

//Option when replacing in URI
type Option func() (Part, Mode)

//WithOption returns the option to replace it in the URI of its part and the selected replacement mode
func WithOption(p Part, mode Mode) Option {
	return func() (Part, Mode) {
		return p, mode
	}
}

func replaceByOpt(cur Part, s string, opts map[Part]Mode) string {
	if s == "" {
		return s
	}
	if m, exists := opts[cur]; exists {
		return m(cur, s)
	}
	return s
}

//Replace returns a string in which the replacement part of the URL in the selected mode
func Replace(s string, opts ...Option) string {
	mo := make(map[Part]Mode)
	for _, opt := range opts {
		p, m := opt()
		mo[p] = m
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
