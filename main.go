package main

import (
	"fmt"
	"sort"

	"github.com/DTreshy/TeamWork_recrutation/csvimporter"
)

func main() {
	hostnameMap := csvimporter.Import("customers.csv")
	hostnames := make([]string, 0, len(hostnameMap))

	for hostname := range hostnameMap {
		hostnames = append(hostnames, hostname)
	}

	sort.Strings(hostnames)

	for _, k := range hostnames {
		fmt.Println(k, hostnameMap[k])
	}
}
