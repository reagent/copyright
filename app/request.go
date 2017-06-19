package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	// "regexp"
	// "database/sql"
)

type Request struct {
	*http.Request

	body *[]byte
}

func NewRequest(r *http.Request) *Request {
	return &Request{Request: r, body: nil}
}

func (r *Request) cachedBody() []byte {
	body := []byte{}

	if r.body == nil {
		tmp, err := ioutil.ReadAll(r.Body)

		if err == nil {
			body = tmp
		}

		r.body = &body
	}

	return *r.body
}

func (r *Request) UnmarshalTo(data interface{}) (err error) {
	err = json.Unmarshal(r.cachedBody(), &data)
	return
}

// func (r *Request) MatchURL(pat string) (matches []string, ok bool) {
// 	re, _ := regexp.Compile(pat)
// 	m := re.FindStringSubmatch(r.URL.Path)

// 	if len(m) == 0 {
// 		ok = false
// 		matches = []string{}
// 	} else {
// 		ok = true
// 		matches = m[1:]
// 	}

// 	return
// }
