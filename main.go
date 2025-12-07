package main

import (
	"GetInfoFromHLTB/client"
	"GetInfoFromHLTB/models"
	"GetInfoFromHLTB/utils"
	"fmt"
	"log"
)

func main() {

	hltbClient := client.New()

	singleGameResponse, err := hltbClient.Search("Portal", models.SearchOptions{
		FilterDLC:  true,
		FilterMods: true,
		MaxResults: 1,
	})

	if err != nil {
		log.Printf("Error: %v", err)
	} else if len(singleGameResponse.Data) > 0 {
		game := singleGameResponse.Data[0]
		fmt.Println("\nInfo about the game:")
		fmt.Println(utils.FormatGameInfo(game))
	}
}
