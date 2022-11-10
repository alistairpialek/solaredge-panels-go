package main

import (
	"fmt"
	"os"

	"github.com/alistairpialek/solaredge-panels-go/solaredge"
)

func main() {
	username := os.Getenv("SOLAREDGE_USERNAME")
	password := os.Getenv("SOLAREDGE_PASSWORD")
	siteID := os.Getenv("SOLAREDGE_SITE_ID")

	// You may optionally include your own http client
	client := solaredge.NewClient(nil, username, password)
	panels, err := client.Site.PanelsEnergy(siteID)
	if err != nil {
		panic(err)
	}

	for _, p := range panels {
		fmt.Printf("%s, %d, %s, %f\n", p.DisplayName, p.ID, p.SerialNumber, p.Energy)
	}
}
