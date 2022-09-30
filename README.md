# Golang SolarEdge Panels API

Go client for the _unofficial_ SolarEdge gateway API. Strangely, the [SolarEdge monitoring API][1] does not expose solar
panel data, which if you have Power Optimizers installed means you are not able to see via an API how much energy each of
your solar panels are producing. This data _is_ available via the website which is largely how I discovered the gateway
API endpoint.

The SolarEdge gateway API also exposes data for:

* Current power flow direction (import vs export)
* Live weather
* Weather forecast

**As this API is an unofficial API, please be considerate with your usage.**

## Install

```
go get github.com/alistairpialek/solaredge-panels-go
```

## Usage

```
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

	for _, p := range *panels {
		fmt.Printf("%s, %d, %s, %f\n", p.DisplayName, p.ID, p.SerialNumber, p.Energy)
	}
}
```

[1]: https://www.solaredge.com/sites/default/files/se_monitoring_api.pdf
