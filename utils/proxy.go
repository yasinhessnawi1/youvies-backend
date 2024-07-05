package utils

import (
	"errors"
	"net/http"
	"net/url"
	"sync"
)

var (
	proxies = []string{
		"http://203.145.179.119:8080",
		"http://117.54.114.101:80",
		"http://103.106.193.236:3128",
		// Add more proxies as needed
	}
	currentProxyIndex int
	mu                sync.Mutex
	requestCount      int
	maxRequests       = 10
)
var count = 99
var client *http.Client

func getClientWithProxy() (*http.Client, error) {
	mu.Lock()
	defer mu.Unlock()

	if len(proxies) == 0 {
		return nil, errors.New("no proxies available")
	}

	// Check if we need to rotate the proxy
	if requestCount >= maxRequests {
		currentProxyIndex = (currentProxyIndex + 1) % len(proxies)
		requestCount = 0
	}

	proxyURL, err := url.Parse(proxies[currentProxyIndex])
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	client := &http.Client{Transport: transport}

	requestCount++
	return client, nil
}

func fetchWithRotatingProxy(url string) (*http.Response, error) {
	count++
	var err error
	if count >= 100 {
		client, err = getClientWithProxy()
		if err != nil {
			return nil, err
		}
		count = 0
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
