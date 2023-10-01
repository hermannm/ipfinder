package ipfinder_test

import (
	"context"
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
