package main

import (
	"context"
	"log"

	"github.com/sfomuseum/go-geojson-show"
)

func main() {

	ctx := context.Background()
	err := show.Run(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
