package gl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
)

type stub func(url.Values) (interface{}, error, int)

type testrq struct {
	method     string
	path       string
	values     url.Values
	parmasbody bool
	h          stub
}

func (rq *testrq) get(k string) string {
	return rq.values.Get(k)
}
func (rq *testrq) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadAll(r.Body)
	var v url.Values
	if rq.parmasbody {
		v, _ = url.ParseQuery(string(b))
	} else {
		v = r.URL.Query()
	}
	rq.method = r.Method
	rq.path = r.URL.Path
	rq.values = v
	res, err, code := rq.h(v)
	if err != nil {
		http.Error(w, err.Error(), code)
	}
	if res != nil {
		json.NewEncoder(w).Encode(res)
	}
}

func th(bodyparm bool, s stub) *testrq {
	return &testrq{method: "", path: "", parmasbody: bodyparm, h: s}
}

func StubHandler(rq *testrq) (*httptest.Server, *Client) {
	ts := httptest.NewServer(rq)
	gitlab, _ := Open(ts.URL, "")
	return ts, gitlab
}

/*
func xStubHandler(asbody bool, h HandlerStub) (*httptest.Server, *Client) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := ioutil.ReadAll(r.Body)
		var v url.Values
		if asbody {
			v, _ = url.ParseQuery(string(b))
		} else {
			v = r.URL.Query()
		}
		res, err, code := h(r.Method, r.URL.Path, v)
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
*/
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
