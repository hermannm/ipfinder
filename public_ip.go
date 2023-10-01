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

// PublicIPAPIs are the URLs which FindPublicIP calls to find your public IP.
// These have almost guaranteed uptime, and no usage limit.
var PublicIPAPIs = []string{
	"https://api.ipify.org/",
	"https://ip.seeip.org/",
}

// FindPublicIP queries the URLs listed in PublicIPAPIs for your public IP. It errors if all API
// calls failed, or if the given context canceled before a result was received.
func FindPublicIP(ctx context.Context) (net.IP, error) {
	ipChan := make(chan net.IP)
	errChan := make(chan error)

	ctx, cancelCtx := context.WithCancel(ctx)

	for _, url := range PublicIPAPIs {
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

	errs := make([]error, 0, len(PublicIPAPIs))
	for {
		select {
		case ip := <-ipChan:
			cancelCtx()
			return ip, nil
		case err := <-errChan:
			errs = append(errs, err)
			if len(errs) == len(PublicIPAPIs) {
				cancelCtx()
				return nil, wrap.Errors("all public IP API calls failed", errs...)
			}
		case <-ctx.Done():
			cancelCtx()
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
