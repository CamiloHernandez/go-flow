package flow

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"

	"github.com/json-iterator/go"
	"github.com/pkg/errors"
)

// requestError is an error response.
type requestError struct {
	// Message is the error detail.
	Message string `json:"message"`

	// Code is the error's code.
	Code    int    `json:"code"`
}

// buildPOST parses an URL and prepares the data body. It automatically adds the verification hash and the API Key.
func (c Client) buildPOST(endpoint string, data map[string]interface{}) (*url.URL, string) {
	data["apiKey"] = c.APIKey

	rqURL, _ := url.Parse(c.URL)
	rqURL.Path = path.Join(rqURL.Path, endpoint)

	var form = url.Values{}
	form.Set("s", c.signData(data))
	for key, value := range data {
		form.Set(key, fmt.Sprintf("%v", value))
	}

	return rqURL, form.Encode()
}

// buildPOST parses an URL and adds the data as query parameters. It automatically adds the verification hash and the API Key.
func (c Client) buildGET(endpoint string, data map[string]interface{}) *url.URL {
	data["apiKey"] = c.APIKey

	rqURL, _ := url.Parse(c.URL)
	rqURL.Path = path.Join(rqURL.Path, endpoint)

	query := rqURL.Query()
	query.Set("s", c.signData(data))

	for key, value := range data {
		query.Set(key, fmt.Sprintf("%v", value))
	}

	rqURL.RawQuery = query.Encode()

	return rqURL
}

// do executes a request.
func (c Client) do(method string, rqURL *url.URL, body string) (data []byte, err error) {
	var rq *http.Request
	if body != "" {
		rq, err = http.NewRequest(method, rqURL.String(), strings.NewReader(body))
	} else {
		rq, err = http.NewRequest(method, rqURL.String(), nil)
	}
	if err != nil {
		return nil, err
	}

	rq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := http.Client{}
	resp, err := httpClient.Do(rq)
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var rqError requestError
		err = jsoniter.Unmarshal(data, &rqError)
		if err != nil {
			return nil, fmt.Errorf("server error: http error %s (%d)", resp.Status, resp.StatusCode)
		}

		return nil, errors.New(rqError.Message)
	}

	return data, nil
}

// get is short hand for do with "GET" as the method.
func (c Client) get(rqURL *url.URL) (data []byte, err error) {
	return c.do("GET", rqURL, "")
}

// get is short hand for do with "POST" as the method.
func (c Client) post(rqURL *url.URL, body string) (data []byte, err error) {
	return c.do("POST", rqURL, body)
}

// signData creates a verification hashed using the client's secret key.
func (c Client) signData(data map[string]interface{}) string {
	// https://www.flow.cl/docs/api.html#section/Introduccion/Como-firmar-con-su-SecretKey
	ordered := sortMapKeys(data)

	var hashValues string
	for _, key := range ordered {
		hashValues += fmt.Sprintf("%v%v", key, data[key])
	}

	return hmac256Sign(hashValues, c.SecretKey)
}

// sortMapKeys creates an alphabetically sorted slice of keys.
func sortMapKeys(m map[string]interface{}) []string {
	mk := make([]string, len(m))
	i := 0
	for k := range m {
		mk[i] = k
		i++
	}
	sort.Strings(mk)

	return mk
}

// hmac256Sign signs a string with a key using HMAC with SHA256.
func hmac256Sign(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return fmt.Sprintf("%x", h.Sum(nil))
}
