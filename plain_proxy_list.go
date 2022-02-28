package ipscraper

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type PlainProxyList struct {
	url string
}

// Works for an URL that returns a proxy list in a text format IP:PORT

func NewPlainProxyList(url string) *PlainProxyList {
	return &PlainProxyList{url: url}
}

func (list *PlainProxyList) Get() ([]string, error) {
	resp, err := http.Get(list.url)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	var ips []string
	scanner := bufio.NewScanner(strings.NewReader(string(body)))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		ips = append(ips, "http://"+scanner.Text())
	}

	return ips, nil
}
