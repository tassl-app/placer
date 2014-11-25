package main

import (
	"encoding/json"
	"flag"
	"github.com/tassl/placer"
	"log"
)

var keyFlag = flag.String("key", "", "Set google places API key")

func main() {
	flag.Parse()
	if len(*keyFlag) == 0 {
		log.Fatalln("Key required")
	}

	log.Println("Testing nearby places")
	nse := placer.NewNearbySearchEndpoint(*keyFlag, -33.8670522, 151.1957362, 500)
	nse.Request.Types = []string{"food"}
	nse.Request.Name = "cruise"
	err := nse.Fetch()
	if err != nil {
		log.Fatalln(err.Error())
	}
	b, err := json.MarshalIndent(nse.Response, "", "  ")
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(string(b))
	log.Println("Testing autocomplete")
	ace := placer.NewAutocompleteEndpoint(*keyFlag, "403 Broad Street Philadelphia PA", 39.945267, -75.165011, placer.OneMileRadius)
	err = ace.Fetch()
	if err != nil {
		log.Fatalln(err.Error())
	}
	b, err = json.MarshalIndent(ace.Response, "", "  ")
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Println(string(b))
}
