package client

import (
	"GetInfoFromHLTB/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// Search searches for a game with the specified options
func (c *Client) Search(gameName string, options models.SearchOptions) (*models.SearchResponse, error) {
	ctx := context.Background()

	if err := c.getAuthToken(ctx); err != nil {
		return nil, fmt.Errorf("failed to get auth token: %v", err)
	}

	if err := c.getSearchURL(ctx); err != nil {
		log.Printf("Warning: could not get search URL: %v", err)
		c.searchURL = hltbBaseURL + "/api/search"
	}

	rawResponse, err := c.doSearch(ctx, gameName)
	if err != nil {
		return nil, err
	}

	var response models.SearchResponse
	if err := json.Unmarshal(rawResponse, &response); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	if options.FilterDLC || options.FilterMods {
		c.filterResults(&response, options.FilterDLC, options.FilterMods)
	}

	if options.MaxResults > 0 && len(response.Data) > options.MaxResults {
		response.Data = response.Data[:options.MaxResults]
		response.Count = len(response.Data)
	}

	return &response, nil
}

func (c *Client) filterResults(response *models.SearchResponse, filterDLC, filterMods bool) {
	if response == nil || len(response.Data) == 0 {
		return
	}

	var filteredData []models.GameData

	for _, game := range response.Data {
		if filterDLC && (game.GameType == "dlc" || strings.Contains(strings.ToLower(game.GameName), "dlc")) {
			continue
		}

		if filterMods && (game.GameType == "mod" || game.GameType == "hack" ||
			strings.Contains(strings.ToLower(game.GameName), "mod") ||
			strings.Contains(strings.ToLower(game.GameName), "hack")) {
			continue
		}

		filteredData = append(filteredData, game)
	}

	response.Data = filteredData
	response.Count = len(filteredData)
}

func (c *Client) getAuthToken(ctx context.Context) error {
	if c.authToken != "" {
		return nil
	}

	tokenURL := hltbBaseURL + "/api/search/init"
	timestamp := time.Now().UnixMilli()
	fullURL := fmt.Sprintf("%s?t=%d", tokenURL, timestamp)

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return err
	}

	c.setCommonHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var tokenResponse struct {
		Token string `json:"token"`
	}

	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return fmt.Errorf("failed to parse token response: %v", err)
	}

	if tokenResponse.Token == "" {
		return fmt.Errorf("auth token not found in response")
	}

	c.authToken = tokenResponse.Token
	return nil
}

// getSearchURL from javascript file
func (c *Client) getSearchURL(ctx context.Context) error {
	if c.searchURL != "" {
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, "GET", hltbBaseURL, nil)
	if err != nil {
		return err
	}

	c.setCommonHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`src=["']([^"']*_app-[^"']*\.js)["']`)
	matches := re.FindStringSubmatch(string(body))

	if len(matches) < 2 {
		re = regexp.MustCompile(`src=["']([^"']*main[^"']*\.js)["']`)
		matches = re.FindStringSubmatch(string(body))
	}

	if len(matches) < 2 {
		return fmt.Errorf("script not found in HTML")
	}

	scriptPath := matches[1]
	if !strings.HasPrefix(scriptPath, "http") {
		scriptPath = hltbBaseURL + scriptPath
	}

	req, err = http.NewRequestWithContext(ctx, "GET", scriptPath, nil)
	if err != nil {
		return err
	}

	c.setCommonHeaders(req)

	resp, err = c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scriptBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	re = regexp.MustCompile(`fetch\s*\(\s*["']/api/([a-zA-Z0-9_/]+)[^"']*["']\s*,\s*\{[^}]*method:\s*["']POST["'][^}]*\}`)
	matches = re.FindStringSubmatch(string(scriptBody))

	if len(matches) > 1 {
		pathSuffix := matches[1]
		if idx := strings.Index(pathSuffix, "/"); idx != -1 {
			pathSuffix = pathSuffix[:idx]
		}

		if pathSuffix != "find" {
			c.searchURL = hltbBaseURL + "/api/" + pathSuffix
			return nil
		}
	}

	return fmt.Errorf("search URL not found in script")
}

func (c *Client) doSearch(ctx context.Context, gameName string) ([]byte, error) {
	searchBody := map[string]interface{}{
		"searchType":  "games",
		"searchTerms": strings.Fields(gameName),
		"searchPage":  1,
		"size":        20,
		"searchOptions": map[string]interface{}{
			"games": map[string]interface{}{
				"userId":        0,
				"platform":      "",
				"sortCategory":  "popular",
				"rangeCategory": "main",
				"rangeTime": map[string]interface{}{
					"min": 0,
					"max": 0,
				},
				"gameplay": map[string]interface{}{
					"perspective": "",
					"flow":        "",
					"genre":       "",
					"difficulty":  "",
				},
				"rangeYear": map[string]interface{}{
					"min": "",
					"max": "",
				},
				"modifier": "",
			},
			"users": map[string]interface{}{
				"sortCategory": "postcount",
			},
			"lists": map[string]interface{}{
				"sortCategory": "follows",
			},
			"filter":     "",
			"sort":       0,
			"randomizer": 0,
		},
		"useCache": true,
	}

	jsonData, err := json.Marshal(searchBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search body: %v", err)
	}

	searchURL := c.searchURL
	if searchURL == "" {
		searchURL = hltbBaseURL + "/api/search"
	}

	req, err := http.NewRequestWithContext(ctx, "POST", searchURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	c.setSearchHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	return body, nil
}

func (c *Client) setCommonHeaders(req *http.Request) {
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Referer", hltbBaseURL+"/")
	req.Header.Set("DNT", "1")
}

func (c *Client) setSearchHeaders(req *http.Request) {
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Referer", hltbBaseURL+"/")
	req.Header.Set("Origin", hltbBaseURL)
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")

	if c.authToken != "" {
		req.Header.Set("x-auth-token", c.authToken)
	}
}
