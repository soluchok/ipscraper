package ipscraper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Does not seem to provide good results at the moment, but the site reports an ongoing issue.

type GeonodeList struct{}

func NewGeonodeList() *GeonodeList {
	return &GeonodeList{}
}

type Response struct {
	Data []struct {
		Ip        string   `json:"ip"`
		Port      string   `json:"port"`
		Protocols []string `json:"protocols"`
	} `json:"data"`
}

func (f *GeonodeList) Get() ([]string, error) {
	resp, err := http.Get("https://proxylist.geonode.com/api/proxy-list?limit=100&page=1&sort_by=lastChecked&sort_type=desc") // 100 could be enough? The result is normally smaller.
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	var ips []string
	for _, rec := range result.Data {
		if len(rec.Protocols) > 0 && strings.HasPrefix(rec.Protocols[0], "http") {
			var url = rec.Protocols[0] + "://" + rec.Ip + ":" + rec.Port
			ips = append(ips, url)
		}
	}

	return ips, nil
}
