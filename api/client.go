package api

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	fmt "github.com/jhunt/go-ansi"
)

const DefaultAPIVersion = "2.14"

type Client struct {
	URL        string
	Username   string
	Password   string
	SkipVerify bool
	Timeout    int

	Debug bool
	Trace bool

	APIVersion string

	ua *http.Client
}

func (c *Client) init() {
	if c.APIVersion == "" {
		c.APIVersion = DefaultAPIVersion
	}

	if c.ua == nil {
		roots, _ := x509.SystemCertPool()
		c.ua = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				TLSClientConfig: &tls.Config{
					RootCAs:            roots,
					InsecureSkipVerify: c.SkipVerify,
				},
			},
			Timeout: time.Duration(c.Timeout) * time.Second,
		}
	}
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	c.init()
	req.Header.Set("X-Broker-API-Version", c.APIVersion)
	req.SetBasicAuth(c.Username, c.Password)

	if c.Trace {
		b, err := httputil.DumpRequest(req, true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "@Y{failed to trace outgoing %s %s request:} @R{%s}\n", req.Method, req.URL, err)
		} else {
			fmt.Fprintf(os.Stderr, "@M{%s}\n\n", string(b))
		}

		res, err := c.ua.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "@Y{request failed:} @R{%s}\n", err)
		}
		if res != nil {
			b, err1 := httputil.DumpResponse(res, true)
			if err1 != nil {
				fmt.Fprintf(os.Stderr, "@Y{failed to trace response:} @R{%s}\n", err)
			} else {
				switch res.StatusCode / 100 {
				case 1:
					fmt.Fprintf(os.Stderr, "@W{%s}\n\n", string(b))
				case 2:
					fmt.Fprintf(os.Stderr, "@G{%s}\n\n", string(b))
				case 3:
					fmt.Fprintf(os.Stderr, "@C{%s}\n\n", string(b))
				default:
					fmt.Fprintf(os.Stderr, "@R{%s}\n\n", string(b))
				}
			}
		}

		return res, err
	}

	return c.ua.Do(req)
}

func (c *Client) parse(res *http.Response, out interface{}) error {
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, out)
}

func (c *Client) url(rel string) string {
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(c.URL, "/"), strings.TrimPrefix(rel, "/"))
}

func (c *Client) get(path string) (res *http.Response, err error) {
	req, err := http.NewRequest("GET", c.url(path), nil)
	if err != nil {
		return nil, err
	}

	return c.do(req)
}

func (c *Client) put(path string, in interface{}) (res *http.Response, err error) {
	b, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", c.url(path), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	return c.do(req)
}

func (c *Client) del(path string) (res *http.Response, err error) {
	req, err := http.NewRequest("DELETE", c.url(path), nil)
	if err != nil {
		return nil, err
	}

	return c.do(req)
}
