package proxmox

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/lwch/logging"
)

// Client proxmox client
type Client struct {
	cli       *http.Client
	url       string
	authToken string
	debug     bool
}

// New create client
func New(url, user, token string, timeout time.Duration) *Client {
	tr := http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var httpCli *http.Client
	if timeout > 0 {
		httpCli = &http.Client{
			Timeout:   timeout,
			Transport: &tr,
		}
	} else {
		httpCli = &http.Client{
			Transport: &tr,
		}
	}
	return &Client{
		cli:       httpCli,
		url:       url,
		authToken: fmt.Sprintf("PVEAPIToken=%s=%s", user, token),
	}
}

// SetDebug set debug flag
func (cli *Client) SetDebug(v bool) {
	cli.debug = v
}

// get send get request
func (cli *Client) get(uri string, args url.Values, value any) error {
	url := fmt.Sprintf("%s/api2/json/"+uri, cli.url)
	enc := args.Encode()
	if len(enc) > 0 {
		url += "?" + enc
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", cli.authToken)
	rep, err := cli.cli.Do(req)
	if err != nil {
		return err
	}
	defer rep.Body.Close()
	if cli.debug {
		data, _ := httputil.DumpResponse(rep, true)
		logging.Info(string(data))
	}
	return json.NewDecoder(rep.Body).Decode(value)
}
