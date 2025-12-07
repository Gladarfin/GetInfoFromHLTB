package utils

import (
	"GetInfoFromHLTB/models"
	"fmt"
)

func SecondsToHours(seconds int) float64 {
	return float64(seconds) / 3600.0
}

func FormatHours(hours float64) string {
	return fmt.Sprintf("%.1f h", hours)
}

// FormatGameInfo into string
func FormatGameInfo(game models.GameData) string {
	var result string

	result += fmt.Sprintf("Game: %s; ", game.GameName)

	if game.CompMain > 0 {
		result += fmt.Sprintf("Main story: %s; ", FormatHours(SecondsToHours(game.CompMain)))
	}

	if game.ReviewScore > 0 {
		result += fmt.Sprintf("Rating: %d/100.", game.ReviewScore)
	}

	return result
}

func FormatGamesList(games []models.GameData) string {
	var result string

	for i, game := range games {
		result += fmt.Sprintf("\nResult %d:\n", i+1)
		result += FormatGameInfo(game)
	}

	return result
}
