package ipprinter

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/hermannm/ipfinder"
)

func PrintIPs() {
	CustomPrintIPs(ipfinder.QueryOptions{})
}

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
