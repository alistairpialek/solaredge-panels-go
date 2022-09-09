package main

import (
	"fmt"
	"os"

	"github.com/alistairpialek/solaredge-panels-go/solaredge"
)

func main() {
	username := os.Getenv("SOLAREDGE_USERNAME")
	password := os.Getenv("SOLAREDGE_PASSWORD")

	// You may optionally include your own http client
	client := solaredge.NewClient(nil, username, password)
	panels, err := client.Sites.PanelsEnergy()
	if err != nil {
		panic(err)
	}

	fmt.Println(panels)
}
