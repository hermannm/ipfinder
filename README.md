# ipfinder

Go utility package for finding your local and public IP addresses.

Run `go get hermannm.dev/ipfinder` to add it to your project!

**Docs:** [pkg.go.dev/hermannm.dev/ipfinder](https://pkg.go.dev/hermannm.dev/ipfinder)

**Contents:**

- [Usage](#usage)
- [Maintainer's guide](#maintainers-guide)

## Usage

Example:

<!-- @formatter:off -->
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
<!-- @formatter:on -->

## Maintainer's guide

### Publishing a new release

- Run tests and linter ([`golangci-lint`](https://golangci-lint.run/)):
  ```
  go test ./... && golangci-lint run
  ```
- Add an entry to `CHANGELOG.md` (with the current date)
    - Remember to update the link section, and bump the version for the `[Unreleased]` link
- Create commit and tag for the release (update `TAG` variable in below command):
  ```
  TAG=vX.Y.Z && git commit -m "Release ${TAG}" && git tag -a "${TAG}" -m "Release ${TAG}" && git log --oneline -2
  ```
- Push the commit and tag:
  ```
  git push && git push --tags
  ```
    - Our release workflow will then create a GitHub release with the pushed tag's changelog entry
