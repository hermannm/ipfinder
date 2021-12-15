package ipprinter

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/hermannm/ipfinder"
)

// Calls CustomPrintIPs with default query options.
func PrintIPs() {
	CustomPrintIPs(ipfinder.QueryOptions{})
}

// Uses the ipfinder package to find your public IP (using provided query options)
// and local IPs, then prints them to the console in a structured format.
// Prints errors if it fails to find IPs.
func CustomPrintIPs(opts ipfinder.QueryOptions) {
	publicIP, err := ipfinder.CustomFindPublicIP(opts)

	fmt.Println("PUBLIC IP")
	if err == nil {
		fmt.Println(publicIP)
	} else {
		fmt.Println(err)
	}
	fmt.Println()

	localIPs, err := ipfinder.FindLocalIPs()

	fmt.Println("LOCAL IPs")
	if err == nil {
		// Prints local IPs with corresponding interface names in two columns
		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)

		for network, ips := range localIPs {
			for _, ip := range ips {
				fmt.Fprintf(w, "%v\t(%v)\n", ip, network)
			}
		}

		w.Flush()
	} else {
		fmt.Println(err)
	}
}
