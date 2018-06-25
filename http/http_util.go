package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/v-zhidu/orb/logging"
)

//Get returns response body that send a GET request to the url
func Get(url string, params map[string]string) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			logging.Error("HTTP GET request failed", logging.Fields{
				"url": url,
			}, err.(error))
		}
	}()

	if len(params) > 0 {
		url = url + "?"
		for key, value := range params {
			url = url + key + "=" + value + "&"
		}
		url = strings.TrimRight(url, "&")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	return doRequest(req)
}

//PostJSON returns response body that send a http POST request to the url
func PostJSON(url string, body []byte, headers map[string]string) ([]byte, error) {
	defer func() {
		if err := recover(); err != nil {
			logging.Error("HTTP POST request failed", logging.Fields{
				"url": url,
			}, err.(error))
		}
	}()

	b := bytes.NewBuffer([]byte(body))
	req, err := http.NewRequest(http.MethodPost, url, b)
	if err != nil {
		return nil, err
	}

	//Add request headers
	req.Header.Add("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	return doRequest(req)
}

//PostMap wapper of HTTPPost method
func PostMap(url string, body map[string]interface{}, headers map[string]string) ([]byte, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return PostJSON(url, data, headers)
}

func doRequest(req *http.Request) ([]byte, error) {
	reqURL := req.Host + req.URL.RequestURI()
	logging.Debug("Exeute HTTP request", logging.Fields{
		"url":  reqURL,
		"body": req.Body,
	})
	start := time.Now()
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		logging.Error("excute HTTP request failed", logging.Fields{
			"url":  reqURL,
			"body": req.Body,
		}, err)
		return nil, err
	}
	defer res.Body.Close()

	logging.Debug("HTTP returns response", logging.Fields{
		"url":      reqURL,
		"body":     res.Body,
		"length":   res.ContentLength,
		"duration": time.Since(start),
	})

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logging.WithError("read response body err", err)
		return nil, err
	}

	return data, nil
}
