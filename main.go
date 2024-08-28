package main

import (
	"fmt"
	"log"

	"github.com/vsomera/focusmode/hosts"
)

func main() {
	store := hosts.NewHostsStore()
	defer store.Close()

	newDomains := []string{"www.facebook.com", "www.instagram.com"}
	err := store.AddDomainsToHost(newDomains)
	if err != nil {
		log.Fatal(err)
	}

	domains, err := store.GetDomainsFromHost()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(domains)
}
