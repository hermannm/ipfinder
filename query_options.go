package ipfinder

import (
	"errors"
	"net/url"
)

// Options to customize the public IP querying with your own list of API URLs,
// or your own timeout limit.
type QueryOptions struct {
	// URLs to API endpoints that return your public IP as a simple string.
	// See DefaultAPIs for examples.
	APIs []string

	// Time in milliseconds before the API calls should time out.
	Timeout int
}

// The default list of APIs to query to find your public IP.
// These have almost guaranteed uptime, and essentially no usage limit.
var DefaultAPIs = []string{
	"https://api.ipify.org/",       // No usage limit.
	"https://ip.seeip.org/",        // No usage limit.
	"https://myexternalip.com/raw", // This API will block you for >30 requests/minute, do not abuse.
}

// Default time in milliseconds to wait for API responses before timing out.
const DefaultTimeout int = 2000

// Checks for the validity of provided query options.
// Returns error if a provided option was invalid.
// If not all query options are provided, injects the default options for those.
func (opts *QueryOptions) validate() error {
	if opts.Timeout == 0 {
		opts.Timeout = DefaultTimeout
	}

	if len(opts.APIs) == 0 {
		opts.APIs = DefaultAPIs

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
