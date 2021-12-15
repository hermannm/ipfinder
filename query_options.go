package ipfinder

import (
	"errors"
	"net/url"
)

// Options to customize the public IP querying with your own list of API URLs,
// or your own timeout limit.
type QueryOptions struct {
	APIs    []string
	Timeout int
}

// The default list of APIs to query to find your public IP.
// These have almost guaranteed uptime, and essentially no usage limit.
var defaultAPIs = []string{
	"https://api.ipify.org/",
	"https://ip.seeip.org/",
	"https://myexternalip.com/raw", // Be aware: this API will block you for >30 calls/minute, do not abuse
}

// Default number of milliseconds to wait for API responses before timing out.
const defaultTimeout int = 2000

// Checks for the validity of provided query options.
// Returns error if a provided option was invalid.
// If not all query options are provided, injects the default options for those.
func (opts *QueryOptions) validate() error {
	if opts.Timeout == 0 {
		opts.Timeout = defaultTimeout
	}

	if len(opts.APIs) == 0 {
		opts.APIs = defaultAPIs

		// Assumes that default APIs work correctly, and so does not validate further.
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
