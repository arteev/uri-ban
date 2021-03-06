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

var DefaultOptions = []Option{
	WithOption(Password, ModeStarred(6)),
	WithOption(Username, ModeValue("user")),
}

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
		return ""
	}
	if m, exists := opts[cur]; exists {
		return m(cur, s)
	}
	return s
}

func modes(opts ...Option) map[Part]Mode {
	mo := make(map[Part]Mode)
	for _, opt := range opts {
		p, m := opt()
		mo[p] = m
	}
	return mo
}

//Replace returns a string in which the replacement part of the URL in the selected mode
func Replace(s string, opts ...Option) string {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		mo := modes(opts...)
		return replaceByOpt(All, s, mo)
	}
	ur := ReplaceURL(u, opts...)
	res, _ := url.PathUnescape(ur.String())
	return res
}

//ReplaceURL returns a url.URL in which the replacement part of the URL in the selected mode
func ReplaceURL(u *url.URL, opts ...Option) url.URL {
	if len(opts) == 0 {
		opts = append([]Option{}, DefaultOptions...)
	}
	currentModes := modes(opts...)
	if u.User != nil {
		user := u.User
		username := replaceByOpt(Username, u.User.Username(), currentModes)
		if p, ok := u.User.Password(); ok {
			pwd := replaceByOpt(Password, p, currentModes)
			if pwd == "" {
				user = url.User(username)
			} else {
				user = url.UserPassword(username, pwd)
			}
		}
		u.User = user
	}
	u.Scheme = replaceByOpt(Scheme, u.Scheme, currentModes)
	u.Path = replaceByOpt(Path, u.Path, currentModes)
	u.RawQuery = replaceByOpt(Query, u.RawQuery, currentModes)
	u.Fragment = replaceByOpt(Fragment, u.Fragment, currentModes)
	u.Host = replaceByOpt(Host, u.Host, currentModes)
	return *u
}
