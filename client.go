package gl

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/spacemonkeygo/errors"
	"github.com/spacemonkeygo/errors/errhttp"
)

const (
	privateToken = "PRIVATE-TOKEN"
	paramSudo    = "SUDO"

	APIv3 = "/api/v3"
)

var (
	unknownError    = errors.NewClass("unknown")
	networkError    = errors.NewClass("network")
	invalidURLError = errors.NewClass("Invalid URL")
	gitlabError     = errors.NewClass("gitlab")
	jsonFormatError = errors.NewClass("jsonformat")

	jsonUnmarshal = errors.GenSym()
)

// Page is used for list queries in gitlab
type Page struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

type Pagination struct {
	FirstPage *Page
	LastPage  *Page
	NextPage  *Page
	PrevPage  *Page
}

// Representation of a gitlab server.
type Client struct {
	hostURL *url.URL
	token   string
	sudo    *string
	params  url.Values
	apiPath string
	client  *http.Client
	log     *log.Logger
}

// Opens a connection to a gitlab server with the v3 api path.
func OpenV3(hosturl string) (*Client, error) {
	return New(hosturl, APIv3, true)
}

// Opens a connection to the given gitlab server. SSL certificates
// are verified.
func Open(hosturl, apipath string) (*Client, error) {
	return New(hosturl, apipath, true)
}

// Create a new Gitlab Client with the given url and api-path. If
// certcheck is false, the SSL certificate will not be verified.
func New(hosturl, apiPath string, certcheck bool) (*Client, error) {
	config := &tls.Config{InsecureSkipVerify: !certcheck}
	tr := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: config,
	}
	client := &http.Client{Transport: tr}

	u, e := url.Parse(hosturl)
	if e != nil {
		return nil, invalidURLError.Wrap(e)
	}
	return &Client{
		hostURL: u,
		params:  make(map[string][]string),
		apiPath: apiPath,
		client:  client,
	}, nil
}

func (c *Client) Host() string {
	return strings.Split(c.hostURL.Host, ":")[0]
}

// Returns a client to gitlab with a copy of all values in the original
// client.
func (c *Client) Child() *Client {
	return &Client{
		hostURL: c.hostURL,
		params:  copyMap(c.params),
		apiPath: c.apiPath,
		client:  c.client,
	}
}

// Sets the privatetoken for the given client.
func (c *Client) Token(t string) {
	c.token = t
}

// Sets a sudo user to be used by the client.
func (c *Client) Sudo(uid string) {
	id := uid
	c.sudo = &id
}

func (c *Client) SetLogger(l *log.Logger) {
	c.log = l
}

func (g *Client) httpexecute(method, u string, params url.Values, paramInbody bool, body []byte, pg *Page) ([]byte, *Pagination, error) {

	var req *http.Request
	var err error

	newurl := *g.hostURL

	parms := make(url.Values)
	if !paramInbody && params != nil && len(params) > 0 {
		for k, v := range params {
			parms[k] = v
		}
	}
	if pg == nil {
		pg = &Page{Page: 1, PerPage: 100}
	}
	parms.Set("page", strconv.Itoa(pg.Page))
	parms.Set("per_page", strconv.Itoa(pg.PerPage))
	newurl.RawQuery = parms.Encode()
	newurl.Opaque = "//" + g.hostURL.Host + g.apiPath + u

	// if no body is given but the params should be in the body
	// overwrite the body value
	if paramInbody && len(params) > 0 && body == nil {
		body = []byte(params.Encode())
		newurl.RawQuery = ""
	}
	if body != nil {
		reader := bytes.NewReader(body)
		req, err = http.NewRequest(method, newurl.String(), reader)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req, err = http.NewRequest(method, newurl.String(), nil)
	}
	req.URL.Opaque = newurl.Opaque
	req.URL.Path = ""
	if err != nil {
		return nil, nil, unknownError.Wrap(err)
	}
	// don't use Add-method, it canonicalizes header names
	req.Header[privateToken] = []string{g.token}
	if g.sudo != nil {
		req.Header.Add(paramSudo, *g.sudo)
	}
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, nil, networkError.Wrap(err)
	}
	defer resp.Body.Close()
	lnk := resp.Header.Get("Link")
	p := parseLinkHeaders(lnk)
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, networkError.Wrap(err)
	}

	if resp.StatusCode >= 400 {
		msg := fmt.Sprintf("%s %s (%d): %s", method, req.URL.String(), resp.StatusCode, strings.TrimSpace(string(contents)))
		if g.log != nil {
			g.log.Printf("%s", msg)
		}
		return nil, nil, gitlabError.NewWith(msg)
	}
	return contents, p, nil
}

func (g *Client) execute(method, u string, params url.Values, paramInbody bool, body []byte, pg *Page, target interface{}) (*Pagination, error) {
	buf, pag, err := g.httpexecute(method, u, params, paramInbody, body, pg)
	if g.log != nil {
		g.log.Printf("%s %s [%+v], buf = %s, err=%s\n", method, u, params, string(buf), err)
	}
	if err != nil {
		return nil, err
	}
	if target != nil {
		err = json.Unmarshal(buf, target)
		if err != nil {
			return nil, jsonFormatError.New("cannont unmarshal json: %s", string(buf))
		}
	}
	return pag, nil
}

func (g *Client) get(u string, params url.Values, pg *Page, target interface{}) (*Pagination, error) {
	return g.execute("GET", u, params, false, nil, pg, target)
}
func (g *Client) put(u string, params url.Values, target interface{}) error {
	_, err := g.execute("PUT", u, params, true, nil, nil, target)
	return err
}
func (g *Client) delete(u string, params url.Values, target interface{}) error {
	_, err := g.execute("DELETE", u, params, false, nil, nil, target)
	return err
}
func (g *Client) post(u string, params url.Values, target interface{}) error {
	_, err := g.execute("POST", u, params, true, nil, nil, target)
	return err
}

func GetStatusCode(err error, default_code int) int {
	return errhttp.GetStatusCode(err, default_code)
}

func GetErrorBody(err error) string {
	return errhttp.GetErrorBody(err)
}
