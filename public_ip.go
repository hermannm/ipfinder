// Package ipfinder provides functions for finding your public and local IP addresses.
package ipfinder

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"

	"hermannm.dev/wrap"
)

// DefaultPublicIPAPIs are the default URLs used to find your public IP.
// These have almost guaranteed uptime, and no usage limit.
var DefaultPublicIPAPIs = []string{
	"https://api.ipify.org/",
	"https://ip.seeip.org/",
}

// FindPublicIP queries the given API URLs for your public IP, and returns the first to respond.
// If no URLs are given, it uses the ones in DefaultPublicIPAPIs.
// It expects the given APIs to return an IP only, in plain-text.
//
// It errors if all API calls fail, or if the given context cancels before a result is received.
func FindPublicIP(ctx context.Context, apiURLs ...string) (net.IP, error) {
	if len(apiURLs) == 0 {
		apiURLs = DefaultPublicIPAPIs
	}

	ipChan := make(chan net.IP)
	errChan := make(chan error)

	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	for _, url := range apiURLs {
		url := url // Avoids mutating loop variable

		go func() {
			ip, err := queryPublicIPAPI(ctx, url)
			if err == nil {
				select {
				case ipChan <- ip:
				case <-ctx.Done():
				}
			} else {
				select {
				case errChan <- fmt.Errorf("%s: %w", url, err):
				case <-ctx.Done():
				}
			}
		}()
	}

	errs := make([]error, 0, len(apiURLs))
	for {
		select {
		case ip := <-ipChan:
			return ip, nil
		case err := <-errChan:
			errs = append(errs, err)
			if len(errs) == len(apiURLs) {
				return nil, wrap.Errors("all public IP API calls failed", errs...)
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func queryPublicIPAPI(ctx context.Context, apiURL string) (net.IP, error) {
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		if err == nil {
			return nil, fmt.Errorf(
				"api responded with status code %d, message '%s'",
				res.StatusCode,
				string(body),
			)
		} else {
			return nil, fmt.Errorf("api responded with status code %d", res.StatusCode)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	bodyString := string(body)
	ip := net.ParseIP(bodyString)
	if ip == nil {
		return nil, fmt.Errorf("failed to parse api response '%s' as IP", bodyString)
	}

	return ip, nil
}
