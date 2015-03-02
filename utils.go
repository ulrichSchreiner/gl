package gl

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
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

func addString(mp url.Values, key string, val *string) {
	if val != nil {
		mp.Set(key, *val)
	}
}
func addInt(mp url.Values, key string, val *int) {
	if val != nil {
		mp.Set(key, strconv.Itoa(*val))
	}
}
func addBool(mp url.Values, key string, val bool) {
	mp.Set(key, fmt.Sprintf("%v", val))
}

// Some crypto helpers, copied from drone

const (
	RSA_BITS     = 2048 // Default number of bits in an RSA key
	RSA_BITS_MIN = 768  // Minimum number of bits in an RSA key
)

// helper function to generate an RSA Private Key.
func GeneratePrivateKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, RSA_BITS)
}

// helper function that marshalls an RSA Public Key to an SSH
// .authorized_keys format
func MarshalPublicKey(pubkey *rsa.PublicKey) string {
	pk, err := ssh.NewPublicKey(pubkey)
	if err != nil {
		return ""
	}

	return string(ssh.MarshalAuthorizedKey(pk))
}

// helper function that marshalls an RSA Private Key to
// a PEM encoded file.
func MarshalPrivateKey(privkey *rsa.PrivateKey) string {
	privateKeyMarshaled := x509.MarshalPKCS1PrivateKey(privkey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Headers: nil, Bytes: privateKeyMarshaled})
	return string(privateKeyPEM)
}
