package main

import (
	"fmt"

	"github.com/vsomera/focusmode/hosts"
)

func main() {
	domains, err := hosts.GetDomainsFromHost()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(domains)
	fmt.Println(len(domains))
}
