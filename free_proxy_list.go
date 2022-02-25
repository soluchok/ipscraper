package ipscraper

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type FreeProxyList struct{}

func NewFreeProxyList() *FreeProxyList {
	return &FreeProxyList{}
}

func (f *FreeProxyList) Get() ([]string, error) {
	resp, err := http.Get("https://free-proxy-list.net/")
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	tokenizer := html.NewTokenizer(resp.Body)
	for {
		if tokenizer.Next() == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				return nil, nil
			}

			return nil, tokenizer.Err()
		}

		tag, _ := tokenizer.TagName()

		if string(tag) == "textarea" {
			list := strings.Split(string(tokenizer.Buffered()), "\n")
			if len(list) < 5 {
				return nil, errors.New("empty list")
			}

			const http_ = "http://"

			for i := range list {
				list[i] = http_ + list[i]
			}

			return list[3 : len(list)-1], nil
		}
	}
}
