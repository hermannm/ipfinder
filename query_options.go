package ipfinder

import (
	"errors"
	"net/url"
)

type QueryOptions struct {
	APIs    []string
	Timeout int
}

var defaultAPIs = []string{
	"https://api.ipify.org/",
	"https://ip.seeip.org/",
	"https://myexternalip.com/raw",
}

const defaultTimeout int = 2000

func (opts *QueryOptions) validate() error {
	if opts.Timeout == 0 {
		opts.Timeout = defaultTimeout
	}
	if len(opts.APIs) == 0 {
		opts.APIs = defaultAPIs
		return nil
	}

	validAPIs := make([]string, 0)
	for _, api := range opts.APIs {
		if _, err := url.ParseRequestURI(api); err == nil {
			validAPIs = append(validAPIs, api)
		}
	}

	if len(validAPIs) == 0 {
		return errors.New("no valid API URLs provided")
	}

	opts.APIs = validAPIs
	return nil
}
