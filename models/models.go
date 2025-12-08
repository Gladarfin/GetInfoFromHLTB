package models

// GameData contains info about the game
type GameData struct {
	GameID       int    `json:"game_id"`
	GameName     string `json:"game_name"`
	ReleaseWorld int64  `json:"release_world"`
	GameType     string `json:"game_type"` // "game", "dlc", "mod", "hack"
	CompMain     int    `json:"comp_main"` // in seconds
	CompPlus     int    `json:"comp_plus"`
	Comp100      int    `json:"comp_100"`
	CompAll      int    `json:"comp_all"`
	GameImage    string `json:"game_image"`
	ReviewScore  int    `json:"review_score"`
}

// SearchResponse contains a response from the API
type SearchResponse struct {
	Color       string     `json:"color"`
	Title       string     `json:"title"`
	Category    string     `json:"category"`
	Count       int        `json:"count"`
	PageCurrent int        `json:"pageCurrent"`
	PageTotal   int        `json:"pageTotal"`
	PageSize    int        `json:"pageSize"`
	Data        []GameData `json:"data"`
}

// SearchOptions contains search parameters
type SearchOptions struct {
	FilterDLC  bool
	FilterMods bool
	MaxResults int
}

func DefaultOptions() SearchOptions {
	return SearchOptions{
		FilterDLC:  true,
		FilterMods: true,
		MaxResults: 2,
	}
}
