package gl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type HandlerStub func(*testing.T, url.Values) (interface{}, error, int)

func StubHandler(t *testing.T, method string, asbody bool, h HandlerStub) (*httptest.Server, *Client) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if method != r.Method {
			t.Errorf("'%s' expected, but '%s' arrived", method, r.Method)
			return
		}
		b, _ := ioutil.ReadAll(r.Body)
		var v url.Values
		if asbody {
			v, _ = url.ParseQuery(string(b))
		} else {
			v = r.URL.Query()
		}
		res, err, code := h(t, v)
		if err != nil {
			http.Error(w, err.Error(), code)
		}
		if res != nil {
			json.NewEncoder(w).Encode(res)
		}
	}))
	gitlab, _ := Open(ts.URL, "")
	return ts, gitlab
}

func has(actual interface{}, expected ...interface{}) string {
	mp := actual.(url.Values)
	for _, e := range expected {
		se := e.(string)
		_, ok := mp[se]
		if !ok {
			return fmt.Sprintf("'%s' is not in the set of parameters", se)
		}
	}
	return ""
}

func hasnot(actual interface{}, expected ...interface{}) string {
	mp := actual.(url.Values)
	for _, e := range expected {
		se := e.(string)
		_, ok := mp[se]
		if ok {
			return fmt.Sprintf("'%s' is in the set of parameters but should not", se)
		}
	}
	return ""
}
