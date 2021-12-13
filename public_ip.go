package ipfinder

import (
	"errors"
	"io"
	"net"
	"net/http"
	"time"
)

func FindPublicIP() (string, error) {
	return CustomFindPublicIP(QueryOptions{})
}

func CustomFindPublicIP(opts QueryOptions) (string, error) {
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

func queryAPI(url string, results chan<- string, errs chan<- error) {
	resp, err := http.Get(url)

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

func startTimeout(milliseconds int, timeout chan<- bool) {
	time.Sleep(time.Millisecond * time.Duration(milliseconds))
	timeout <- true
}
