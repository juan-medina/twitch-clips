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

package twitch

import (
	"errors"
	"math"

	"github.com/nicklaw5/helix/v2"
	"github.com/rs/zerolog/log"
)

func getMapOfGameIdsToGameNames(client *helix.Client, clips []ClipInfo) (map[string]string, error) {
	uniqueGameIds := getUniqueGameIds(clips)
	result := make(map[string]string)

	log.Info().Strs("game_ids", uniqueGameIds).Msg("getting games names")

	// we can get 100 games at max, so we need to get them in batches
	batchSize := 100
	batches := int(math.Ceil(float64(len(uniqueGameIds)) / float64(batchSize)))
	for i := 0; i < batches; i++ {
		log.Info().Int("batch", i+1).Int("batches", batches).Msg("getting games names")
		end := int(math.Min(float64(i+1)*float64(batchSize), float64(len(uniqueGameIds))))
		gameIds := uniqueGameIds[i*batchSize : end]
		if resp, err := client.GetGames(&helix.GamesParams{
			IDs: gameIds,
		}); err == nil {
			if resp.StatusCode == 200 {
				for _, game := range resp.Data.Games {
					result[game.ID] = game.Name
				}
			} else {
				err = errors.New("fail to get games")
				log.Err(err).Int("status_code", resp.ErrorStatus).Str("error", resp.Error).Str("description", resp.ErrorMessage).Msg("request failed")
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	//output to log all game names
	for gameId, gameName := range result {
		log.Info().Str("game_id", gameId).Str("game_name", gameName).Msg("game name")
	}
	return result, nil
}

func getUniqueGameIds(clips []ClipInfo) []string {
	gameIds := []string{}
	for _, clip := range clips {
		if !contains(gameIds, clip.Game) {
			gameIds = append(gameIds, clip.Game)
		}
	}
	return gameIds
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
