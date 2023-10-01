# ipfinder

Go utility package for finding your local and public IP addresses.

Run `go get hermannm.dev/ipfinder` to add it to your project!

## Usage

```go
import (
	"context"
	"fmt"

	"hermannm.dev/ipfinder"
)

func main() {
	publicIP, err := ipfinder.FindPublicIP(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Public IP: %s\n", publicIP.String())

	localIPs, err := ipfinder.FindLocalIPs()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Local IPs:")
	for _, ip := range localIPs {
		fmt.Printf("%s (interface %s)\n", ip.Address.String(), ip.NetworkInterface.Name)
	}
}
```
