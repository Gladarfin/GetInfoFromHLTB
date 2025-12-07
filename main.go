package main

import (
	"fmt"
	"log"

	client "github.com/Gladarfin/GetInfoFromHLTB/client"
	models "github.com/Gladarfin/GetInfoFromHLTB/models"
	utils "github.com/Gladarfin/GetInfoFromHLTB/utils"
)

func main() {

	hltbClient := client.New()

	singleGameResponse, err := hltbClient.Search("Space Rangers", models.SearchOptions{
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
