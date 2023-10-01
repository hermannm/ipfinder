package ipfinder_test

import (
	"context"
	"testing"
	"time"

	"hermannm.dev/ipfinder"
)

func TestFindPublicIP(t *testing.T) {
	publicIP, err := ipfinder.FindPublicIP(context.Background())
	if err != nil {
		t.Fatalf("received error: %v", err)
	}
	if publicIP == nil {
		t.Fatal("public IP was nil")
	}
}

func TestFindPublicIPWithTimeout(t *testing.T) {
	ctx, cleanupCtx := context.WithTimeout(context.Background(), time.Microsecond)

	_, err := ipfinder.FindPublicIP(ctx)
	cleanupCtx()
	if err == nil {
		t.Fatal("expected FindPublicIP to error with timeout")
	}
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
	}
}
