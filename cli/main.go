package main

import (
	"fmt"
	"net/url"
	"os"
	"sort"

	"github.com/ghchinoy/ce-go/ce"
)

func main() {
	fmt.Println("test CE client")

	baseURL, err := url.Parse(os.Getenv("CE_BASE"))
	if err != nil {
		fmt.Println("Couldn't parse CE_BASE environment variable URL")
		os.Exit(1)
	}

	client := ce.Client{
		BaseURL:      baseURL,
		Organization: os.Getenv("CE_ORG"),
		User:         os.Getenv("CE_USER"),
	}

	elements, err := client.ListElements()
	if err != nil {
		fmt.Println("Couldn't call Elements")
		fmt.Println(err)
		os.Exit(1)
	}

	sort.Sort(ce.ByName(elements))

	for _, e := range elements {
		fmt.Printf("%v\n", e.Name)
	}
}
