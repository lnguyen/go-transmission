package client

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type ApiClient struct {
	url      string
	username string
	password string
	token    string
	client   http.Client
}

func NewClient(url string,
	username string, password string) ApiClient {
	ac := ApiClient{url: url + "/transmission/rpc", username: username, password: password}

	return ac
}

func (ac *ApiClient) CreateClient(apiToken string) {
	ac.client = http.Client{}
}

func (ac *ApiClient) Post(body string) ([]byte, error) {
	authRequest, err := ac.authRequest("POST", body)
	if err != nil {
		return make([]byte, 0), err
	}
	res, err := ac.client.Do(authRequest)
	if err != nil {
		return make([]byte, 0), err
	}
	if res.StatusCode == 409 {
		ac.getToken()
		authRequest, err := ac.authRequest("POST", body)
		if err != nil {
			return make([]byte, 0), err
		}
		res, err = ac.client.Do(authRequest)
		if err != nil {
			return make([]byte, 0), err
		}
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return make([]byte, 0), err
	}
	return resBody, nil
}

func (ac *ApiClient) getToken() error {
	req, err := http.NewRequest("POST", ac.url, strings.NewReader(""))
	if err != nil {
		return err
	}

	req.SetBasicAuth(ac.username, ac.password)
	res, err := ac.client.Do(req)
	if err != nil {
		return err
	}
	ac.token = res.Header.Get("X-Transmission-Session-Id")
	return nil
}

func (ac *ApiClient) authRequest(method string, body string) (*http.Request, error) {
	if ac.token == "" {
		err := ac.getToken()
		if err != nil {
			return &http.Request{}, err
		}
	}
	req, err := http.NewRequest(method, ac.url, strings.NewReader(body))
	if err != nil {
		return &http.Request{}, err
	}
	req.Header.Add("X-Transmission-Session-Id", ac.token)

	req.SetBasicAuth(ac.username, ac.password)
	return req, nil
}
