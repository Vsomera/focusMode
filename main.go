package main

import (
	"fmt"

	"github.com/vsomera/focusmode/hosts"
)

func main() {
	domains := []string{"www.instagram.com", "www.facebook.com"}
	err := hosts.AddDomainsToHost(domains)
	if err != nil {
		fmt.Println(err)
	}
	newDomains, err := hosts.GetDomainsFromHost()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(newDomains)
}
