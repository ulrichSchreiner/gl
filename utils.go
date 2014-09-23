package gl

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

func copyMap(m map[string][]string) map[string][]string {
	res := make(map[string][]string)
	for k, v := range m {
		res[k] = v
	}
	return res
}

func pageFromURL(u string) (*Page, error) {
	ur, e := url.Parse(u)
	if e != nil {
		return nil, e
	}
	vals := ur.Query()
	var p Page
	fmt.Sscanf(vals.Get("page"), "%d", &p.Page)
	fmt.Sscanf(vals.Get("per_page"), "%d", &p.PerPage)
	return &p, nil
}

func parseLinkHeaders(lnk string) *Pagination {
	var p Pagination
	lnks := strings.Split(lnk, ",")
	for _, l := range lnks {
		if strings.Contains(l, "rel=\"first\"") {
			u := strings.Split(l, ";")[0]
			u = u[1 : len(u)-1]
			p.FirstPage, _ = pageFromURL(u)
		}
		if strings.Contains(l, "rel=\"next\"") {
			u := strings.Split(l, ";")[0]
			u = u[1 : len(u)-1]
			p.NextPage, _ = pageFromURL(u)
		}
		if strings.Contains(l, "rel=\"prev\"") {
			u := strings.Split(l, ";")[0]
			u = u[1 : len(u)-1]
			p.PrevPage, _ = pageFromURL(u)
		}
		if strings.Contains(l, "rel=\"last\"") {
			u := strings.Split(l, ";")[0]
			u = u[1 : len(u)-1]
			p.LastPage, _ = pageFromURL(u)
		}
	}
	return &p
}

func expandUrl(u string, params map[string]interface{}) string {
	if params != nil {
		for key, val := range params {
			sval := fmt.Sprintf("%v", val)
			u = strings.Replace(u, key, sval, -1)
		}
	}

	return u
}

//type projectFetcher func(*Page) (Projects, *Pagination, error)
type fetchFunc func(pg *Page) (interface{}, *Pagination, error)

func fetchAll(ff fetchFunc, result interface{}) error {
	var pg *Page
	ptr := reflect.ValueOf(result)
	targ := reflect.Indirect(ptr)
	for {
		vals, pag, err := ff(pg)
		if err != nil {
			return err
		}
		targ = reflect.AppendSlice(targ, reflect.ValueOf(vals))
		if pag.NextPage == nil {
			break
		}
		pg = pag.NextPage
	}
	ptr.Elem().Set(targ)
	return nil
}
