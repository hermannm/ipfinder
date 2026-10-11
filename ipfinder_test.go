package ipfinder_test

import (
	"context"
	"regexp"
	"strings"
	"testing"
	"time"

	"hermannm.dev/ipfinder"
)

func TestFindPublicIP(t *testing.T) {
	ctx, cleanupCtx := context.WithTimeout(context.Background(), 10*time.Second)

	publicIP, err := ipfinder.FindPublicIP(ctx)
	cleanupCtx()

	if err != nil {
		t.Fatalf("received error: %v", err)
	}
	if publicIP == nil {
		t.Fatal("public IP was nil")
	}

	t.Logf("Found public IP %s", publicIP.String())
}

func TestFindPublicIPWithTimeout(t *testing.T) {
	ctx, cleanupCtx := context.WithTimeout(context.Background(), time.Microsecond)

	_, err := ipfinder.FindPublicIP(ctx)
	cleanupCtx()
	if err == nil {
		t.Fatal("expected FindPublicIP to error with timeout")
	}
}

func TestFindPublicIPWithCustomURL(t *testing.T) {
	ctx, cleanupCtx := context.WithTimeout(context.Background(), 10*time.Second)

	// Note: Rate limited to 30 requests/minute, so don't abuse
	const customURL = "https://myexternalip.com/raw"

	publicIP, err := ipfinder.FindPublicIP(ctx, customURL)
	cleanupCtx()

	if err != nil {
		t.Fatalf("received error: %v", err)
	}

	t.Logf("Found public IP %s from custom URL %s", publicIP.String(), customURL)
}

func TestFindPublicIPError(t *testing.T) {
	ctx, cleanupCtx := context.WithTimeout(context.Background(), 10*time.Second)

	invalidURLs := []string{"invalid-url", "invalid-url2"}

	publicIP, err := ipfinder.FindPublicIP(ctx, invalidURLs...)
	cleanupCtx()

	if err == nil {
		t.Fatalf("expected error, got nil error and public IP '%s'", publicIP)
	}

	unwrapper, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("expected error to implement Unwrap() []error interface")
	}

	causes := unwrapper.Unwrap()
	if len(causes) != len(invalidURLs) {
		t.Fatalf("expected unwrapped errors '%v' to match length of invalid URLs", causes)
	}

	// The URLs may be attempted in any order, since we run them concurrently - so check for both
	expectedErrorRegex1 := regexp.MustCompile(
		`all public IP API calls failed: 1\) invalid-url: .+, 2\) invalid-url2: .+`,
	)
	expectedErrorRegex2 := regexp.MustCompile(
		`all public IP API calls failed: 1\) invalid-url2: .+, 2\) invalid-url: .+`,
	)

	errMsg := err.Error()
	if !expectedErrorRegex1.MatchString(errMsg) && !expectedErrorRegex2.MatchString(errMsg) {
		t.Fatalf("unexpected error message '%s'", errMsg)
	}

CauseLoop:
	for _, cause := range causes {
		// Check that the cause error is prefixed by one of the URLs
		for _, url := range invalidURLs {
			if strings.HasPrefix(cause.Error(), url+": ") {
				continue CauseLoop
			}
		}

		t.Fatalf("expected cause error to be prefixed by URL, got '%s'", cause.Error())
	}

	t.Logf("Got expected public IP error:\n%s", errMsg)
}

func TestFindLocalIPs(t *testing.T) {
	localIPs, err := ipfinder.FindLocalIPs()
	if err != nil {
		t.Fatalf("received error: %v", err)
	}

	for _, ip := range localIPs {
		if ip.Address == nil {
			t.Fatal("local IP was nil")
		}

		t.Logf("Found local IP %s (interface: %s)", ip.Address.String(), ip.NetworkInterface.Name)
	}
}
