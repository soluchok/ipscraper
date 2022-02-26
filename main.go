package ipscraper

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/sync/semaphore"
)

const (
	timeout       = 5 * time.Second
	maxGoroutines = 32
)

type Provider interface {
	Get() ([]string, error)
}

func New() *Scraper {
	return &Scraper{
		providers: []Provider{NewFreeProxyList(), NewOpenProxyList(), NewGeonodeList()},
		sem:       semaphore.NewWeighted(maxGoroutines),
		wg:        &sync.WaitGroup{},
		cache:     &sync.Map{},
	}
}

type Scraper struct {
	providers []Provider
	sem       *semaphore.Weighted
	wg        *sync.WaitGroup
	cache     *sync.Map
}

func (s *Scraper) Get() ([]string, error) {
	var allProxies []string
	for i := range s.providers {
		proxies, err := s.providers[i].Get()
		if err != nil {
			return nil, fmt.Errorf("provider get: %w", err)
		}

		allProxies = append(allProxies, proxies...)
	}

	var res = make(chan string, len(allProxies))
	for i := range allProxies {
		s.wg.Add(1)
		// TODO: Check error
		s.sem.Acquire(context.Background(), 1)

		go func(ip string) {
			defer s.wg.Done()
			defer s.sem.Release(1)

			val, err := s.checkIPWithCache(ip)
			if err != nil {
				log.Printf("check IP with cache: %v\n", err)
				return
			}

			if val {
				res <- ip
			}
		}(allProxies[i])
	}

	s.wg.Wait()

	close(res)

	var result []string
	for ip := range res {
		result = append(result, ip)
	}

	return result, nil
}

func (s *Scraper) checkIPWithCache(ip string) (bool, error) {
	val, ok := s.cache.Load(ip)
	if ok {
		return val.(bool), nil
	}

	valid, err := checkIP(ip)
	if err != nil {
		return valid, fmt.Errorf("check IP: %w", err)
	}

	s.cache.Store(ip, valid)

	return valid, nil
}

func checkIP(ip string) (bool, error) {
	const myIp = "https://api.myip.com"

	_ip, err := url.Parse(ip)
	if err != nil {
		return false, err
	}

	resp, err := (&http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(_ip),
		},
	}).Get(myIp)
	if err != nil {
		return false, fmt.Errorf("http get: %w", err)
	}

	jsonResp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("read all: %w", err)
	}

	if !strings.Contains(string(jsonResp), strings.Split(_ip.Host, ":")[0]) {
		return false, nil
	}

	return true, nil
}
