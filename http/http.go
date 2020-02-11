package http

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, time.Second*30)
}

func typeSwitcher(t interface{}) string {
	switch v := t.(type) {
	case int:
		return strconv.Itoa(v)
	case string:
		return v
	case int64:
		return strconv.Itoa(int(v))
	case []string:
		return "typeArray"
	case map[string]interface{}:
		return "typeMap"
	default:
		return ""
	}
}

func ParamsToStr(params map[string]interface{}) string {
	requestUrl := ""
	for k, v := range params {
		if strings.Contains(k, "_") {
			strings.Replace(k, ".", "_", -1)
		}
		v := typeSwitcher(v)
		requestUrl = requestUrl + k + "=" + url.QueryEscape(v) + "&"
	}
	requestUrl = strings.TrimRight(requestUrl, "&")
	return requestUrl
}

func httpGet(url string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr, Timeout: time.Duration(3) * time.Second}
	fmt.Println(url)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("http wrong err")
		return err.Error()
	}

	return string(body)
}

func httpPost(requestUrl string, params map[string]interface{}) string {
	b, err := json.Marshal(params)
	if err != nil {
		fmt.Errorf("json.Marshal failed[%v]", err)
		return err.Error()
	}

	req, err := http.NewRequest("POST", requestUrl, strings.NewReader(string(b)))
	req.Header.Set("Content-Type", "application/json")

	transport := http.Transport{
		Dial:              dialTimeout,
		DisableKeepAlives: true,
	}

	client := &http.Client{Transport: &transport, Timeout: time.Duration(30) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("http wrong err")
		return err.Error()
	}

	return string(body)
}

func SendRequest(requestUrl string, params map[string]interface{}, method string) string {
	response := ""
	if method == "GET" {
		paramsStr := "?" + ParamsToStr(params)
		requestUrl = requestUrl + paramsStr
		response = httpGet(requestUrl)
	} else if method == "POST" {
		response = httpPost(requestUrl, params)
	} else {
		fmt.Println("unsuppported http method")
		return "unsuppported http method"
	}

	return response
}