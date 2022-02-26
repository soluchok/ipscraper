package ipscraper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type OpenProxyList struct{}

func NewOpenProxyList() *OpenProxyList {
	return &OpenProxyList{}
}

func (f *OpenProxyList) Get() ([]string, error) {
	re := regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}:(\d{2,5})`)
	resp, err := http.Get("https://openproxy.space/list/http")
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	data := string(body)

	var result []string

	submatchall := re.FindAllString(data, -1)
	for _, element := range submatchall {
		result = append(result, "http://"+element)
	}
	return result[0 : len(result)-2], nil
}
