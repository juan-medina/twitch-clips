/*
 * MIT License
 *
 * Copyright (c) 2023 Juan Medina
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package extract

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/rs/zerolog/log"
)

func rateLimitCallback(lastResponse *helix.Response) error {
	if lastResponse.GetRateLimitRemaining() > 0 {
		return nil
	}

	reset64 := int64(lastResponse.GetRateLimitReset())
	currentTime := time.Now().Unix()

	if currentTime < reset64 {
		if timeDiff := time.Duration(reset64 - currentTime); timeDiff > 0 {
			log.Warn().Int64("seconds", int64(timeDiff.Seconds())).Msg("waiting on rate limit")
			time.Sleep(timeDiff * time.Second)
			log.Warn().Msg("done waiting on rate limit")
		}
	}

	return nil
}

func maskField(s string) string {
	masked := ""
	for range s {
		masked += "*"
	}
	return masked
}

func getClient(clientId string, secret string) (*helix.Client, error) {
	log.Info().Str("client_id", clientId).Str("secret", maskField(secret)).Msg("getting a twitch client")
	if client, err := helix.NewClient(&helix.Options{
		ClientID:      clientId,
		ClientSecret:  secret,
		RateLimitFunc: rateLimitCallback}); err == nil {
		log.Info().Msg("requesting app access token")
		if resp, err := client.RequestAppAccessToken([]string{"user:read:email"}); err == nil {
			client.SetAppAccessToken(resp.Data.AccessToken)
			client.SetRefreshToken(resp.Data.RefreshToken)
			return client, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func getBroadcasterId(client *helix.Client, channel string) (string, error) {
	log.Info().Str("channel", channel).Msg("getting broadcaster id")
	if resp, err := client.GetUsers(&helix.UsersParams{
		Logins: []string{channel},
	}); err == nil {
		if resp.StatusCode != 200 {
			err = errors.New("fail to get broadcaster id")
			log.Err(err).Int("status_code", resp.ErrorStatus).Str("error", resp.Error).Str("description", resp.ErrorMessage).Msg("request failed")
			return "", err
		}
		if len(resp.Data.Users) == 0 {
			return "", fmt.Errorf("channel %s not found", channel)
		}
		return resp.Data.Users[0].ID, nil
	} else {
		return "", err
	}
}

type clipInfo struct {
	url   string
	title string
	game  string
	date  time.Time
}

func getClipsAfter(client *helix.Client, broadcasterId string, after string) ([]clipInfo, string, error) {
	result := []clipInfo{}
	newAfter := after
	log.Info().Str("broadcaster_id", broadcasterId).Str("after", after).Msg("getting clips after")

	if resp, err := client.GetClips(&helix.ClipsParams{
		BroadcasterID: broadcasterId,
		First:         100,
		After:         after,
	}); err == nil {
		if resp.StatusCode == 200 {
			for _, clip := range resp.Data.Clips {
				if date, err := time.Parse(time.RFC3339, clip.CreatedAt); err == nil {
					result = append(result, clipInfo{
						url:   clip.URL,
						title: clip.Title,
						game:  clip.GameID,
						date:  date,
					})
				} else {
					log.Err(err).Str("date", clip.CreatedAt).Msg("fail to parse date")
				}
			}
			newAfter = resp.Data.Pagination.Cursor
		} else {
			err = errors.New("fail to get clips")
			log.Err(err).Int("status_code", resp.ErrorStatus).Str("error", resp.Error).Str("description", resp.ErrorMessage).Msg("request failed")
			return result, "", err
		}
	} else {
		return nil, "", err
	}

	return result, newAfter, nil
}

func getClips(client *helix.Client, broadcasterId string) ([]clipInfo, error) {
	result := []clipInfo{}
	log.Info().Str("broadcaster_id", broadcasterId).Msg("getting clips")

	continueGet := true
	after := ""

	for continueGet {
		if clips, newAfter, err := getClipsAfter(client, broadcasterId, after); err == nil {
			result = append(result, clips...)
			if newAfter == "" {
				continueGet = false
			}
			after = newAfter
		} else {
			return nil, err
		}
	}

	// create a map of game ids to game names
	gameNames := make(map[string]string)

	for _, clip := range result {
		gameNames[clip.game] = clip.game
	}

	// get the keys from map
	gameIds := make([]string, 0, len(gameNames))
	for k := range gameNames {
		gameIds = append(gameIds, k)
	}

	if resp, err := client.GetGames(&helix.GamesParams{
		IDs: gameIds,
	}); err == nil {
		if resp.StatusCode == 200 {
			for _, game := range resp.Data.Games {
				gameNames[game.ID] = game.Name
			}
		} else {
			err = errors.New("fail to get games")
			log.Err(err).Int("status_code", resp.ErrorStatus).Str("error", resp.Error).Str("description", resp.ErrorMessage).Msg("request failed")
			return result, err
		}
	}

	for i, clip := range result {
		result[i].game = gameNames[clip.game]
	}

	// sort by date, desc
	sort.Slice(result, func(i, j int) bool {
		return result[i].date.After(result[j].date)
	})

	return result, nil
}
func writeClipInfoToCSV(filename string, data []clipInfo) error {
	log.Info().Str("filename", filename).Msg("writing clip info to csv")
	// Create a new CSV file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)

	// Write header row
	header := []string{"Title", "Game", "Date", "URL"}
	if err = writer.Write(header); err != nil {
		return err
	}

	// Write data rows
	for _, clip := range data {
		row := []string{clip.title, clip.game, clip.date.Format("2006-01-02T15:04:05"), clip.url}
		if err = writer.Write(row); err != nil {
			return err
		}
	}

	// Flush the writer to write any buffered data to the file
	writer.Flush()

	return writer.Error()
}

func Execute(clientId string, secret string, channel string, filename string) error {
	log.Info().Str("client_id", clientId).Str("channel", channel).Msg("extract clips")
	if client, err := getClient(clientId, secret); err == nil {
		if broadcasterId, err := getBroadcasterId(client, channel); err == nil {
			if clips, err := getClips(client, broadcasterId); err == nil {
				log.Info().Int("clips", len(clips)).Msg("got clips")
				if err := writeClipInfoToCSV(filename, clips); err != nil {
					return fmt.Errorf("fail to write clips to csv: %w", err)
				}
			} else {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}

	return nil
}
