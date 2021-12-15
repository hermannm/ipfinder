package ipfinder

import (
	"errors"
	"io"
	"net"
	"net/http"
	"time"
)

// Calls CustomFindPublicIP with default query options.
func FindPublicIP() (string, error) {
	return CustomFindPublicIP(QueryOptions{})
}

// Calls a list of public IP APIs (given in query options) concurrently.
// Returns the IP given by the first API to respond.
// Returns error if no API responds within timeout limit (given in query options).
// Also returns error if all API calls fail before timeout.
func CustomFindPublicIP(opts QueryOptions) (string, error) {
	// Checks validity of provided options, and injects defaults for lacking options.
	err := opts.validate()
	if err != nil {
		return "", err
	}

	results := make(chan string, 1)
	errs := make(chan error, len(opts.APIs))
	timeout := make(chan bool, 1)

	go startTimeout(opts.Timeout, timeout)

	for _, url := range opts.APIs {
		go queryAPI(url, results, errs)
	}

	// Checks for result, timeout or all API queries failing.
	for {
		select {
		case ip := <-results:
			return ip, nil
		case <-timeout:
			return "", errors.New("public IP API queries timed out")
		default:
			if len(errs) == cap(errs) {
				return "", errors.New("all public IP API queries failed")
			}
		}
	}
}

// Calls the given API URL, parses the result as an IP string, and sends it on the results channel.
// If API call or IP parsing fail, instead sends error on the error channel.
func queryAPI(api string, results chan<- string, errs chan<- error) {
	resp, err := http.Get(api)

	if err != nil {
		errs <- err
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errs <- errors.New("API response not OK")
		return
	}

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		errs <- err
		return
	}

	ipString := string(ip)
	if net.ParseIP(ipString) == nil {
		errs <- errors.New("invalid IP returned from API")
		return
	}

	results <- ipString
}

// Sleeps for the provided duration, then sends true to the timeout channel.
func startTimeout(milliseconds int, timeout chan<- bool) {
	time.Sleep(time.Millisecond * time.Duration(milliseconds))
	timeout <- true
}
