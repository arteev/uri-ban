package uriban

import "testing"

func TestBan(t *testing.T) {

	ft := func(s string) string { return "123" }
	cases := []struct {
		url   string
		want  string
		modes []Option
	}{
		{"", "", nil},
		{"url", "url", nil},

		//BEGIN TEST MODES
		{"", "", []Option{WithOption(All, ModeNothing())}},
		{"url", "url", []Option{WithOption(All, ModeNothing())}},
		{"url", "", []Option{WithOption(All, ModeHidden())}},
		{"url", "*****", []Option{WithOption(All, ModeStarred(5))}},
		{"url", "HIDDEN", []Option{WithOption(All, ModeValue("HIDDEN"))}},
		{"url", "123", []Option{WithOption(All, ModeFunc(ft))}},
		//END TEST MODES

		{"scheme://userinfo:pwd@host/path?query#fragment", "scheme://userinfo:pwd@host/path?query#fragment", nil},

		//PASSWORD
		{"scheme://userinfo:pwd@host/path?query#fragment", "scheme://userinfo@host/path?query#fragment", []Option{WithOption(Password, ModeHidden())}},
		{"scheme://userinfo:pwd@host/path?query#fragment", "scheme://userinfo:pwd@host/path?query#fragment", []Option{WithOption(Password, ModeNothing())}},
		{"scheme://userinfo:pwd@host/path?query#fragment", "scheme://userinfo:******@host/path?query#fragment", []Option{WithOption(Password, ModeStarred(6))}},
		{"scheme://userinfo:pwd@host/path?query#fragment", "scheme://userinfo:password@host/path?query#fragment", []Option{WithOption(Password, ModeValue("password"))}},
		{"scheme://userinfo:pwd@host/path?query#fragment", "scheme://userinfo:123@host/path?query#fragment", []Option{WithOption(Password, ModeFunc(ft))}},
		{"scheme://userinfo@host/path?query#fragment", "scheme://userinfo@host/path?query#fragment", []Option{WithOption(Password, ModeStarred(6))}},
		//END PASSWORD

		//USER
		{"scheme://userinfo:pwd@host/path?query#fragment", "scheme://usr1:pwd@host/path?query#fragment", []Option{WithOption(Username, ModeValue("usr1"))}},

		//SCHEME
		{"scheme://userinfo:pwd@host/path?query#fragment", "//userinfo:pwd@host/path?query#fragment", []Option{WithOption(Scheme, ModeHidden())}},

		//PATH
		{"scheme://userinfo:pwd@host/path?query#fragment", "scheme://userinfo:pwd@host/newpath?query#fragment", []Option{WithOption(Path, ModeValue("newpath"))}},
		//QUERY
		{"scheme://userinfo:pwd@host/path?query#fragment", "scheme://userinfo:pwd@host/path?newquery#fragment", []Option{WithOption(Query, ModeValue("newquery#fragment"))}},

		//TODO: FRAGMENT. Bug?
		//{"scheme://userinfo:pwd@host/path?query#fragment", "scheme://userinfo:pwd@host/path?query#newfragment", []Option{WithOption(Fragment, ModeValue("newfragment"))}},

		//COMPLEX

		{
			"scheme://userinfo:pwd@host/path?query#fragment",
			"scheme://usr1:****@host/newpath?query#fragment",
			[]Option{
				WithOption(Username, ModeValue("usr1")),
				WithOption(Password, ModeStarred(4)),
				WithOption(Path, ModeValue("newpath")),
			},
		},
		//TODO: with err escape symbols
	}

	for _, c := range cases {
		got := Ban(c.url, c.modes...)
		if got != c.want {
			t.Errorf("Expected %q, got %q", c.want, got)
		}
	}

	//  scheme:opaque?query#fragment
	//  scheme://userinfo@host/path?query#fragment
}
