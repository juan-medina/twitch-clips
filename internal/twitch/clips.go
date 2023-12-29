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
	"sort"
	"time"

	"github.com/juan-medina/twitch-clips/internal/times"
	"github.com/nicklaw5/helix/v2"
	"github.com/rs/zerolog/log"
)

func GetClips(client *helix.Client, broadcasterId string, from time.Time, to time.Time) ([]ClipInfo, error) {
	result := []ClipInfo{}
	log.Info().Str("broadcaster_id", broadcasterId).Msg("getting clips")

	continueGet := true
	after := ""

	for continueGet {
		if clips, newAfter, err := getClipsAfter(client, broadcasterId, after, from, to); err == nil {
			result = append(result, clips...)
			if newAfter == "" {
				continueGet = false
			}
			after = newAfter
		} else {
			return nil, err
		}
	}

	// get the game names
	if gameNames, err := getMapOfGameIdsToGameNames(client, result); err == nil {
		for i, clip := range result {
			result[i].Game = gameNames[clip.Game]
		}
	} else {
		return nil, err
	}

	// sort by date, desc
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.After(result[j].Date)
	})

	return result, nil
}

func getClipsAfter(client *helix.Client, broadcasterId string, after string, from time.Time, to time.Time) ([]ClipInfo, string, error) {
	result := []ClipInfo{}
	newAfter := after
	l := log.Info().Str("broadcaster_id", broadcasterId).Str("after", after)

	if from != times.NilDateTime {
		l = l.Time("from", from)
	}

	if to != times.NilDateTime {
		l = l.Time("to", to)
	}

	l.Msg("getting clips after")

	params := helix.ClipsParams{
		BroadcasterID: broadcasterId,
		First:         100,
		After:         after,
	}

	if from != times.NilDateTime {
		params.StartedAt = helix.Time{Time: from}
	}

	if to != times.NilDateTime {
		params.EndedAt = helix.Time{Time: to}
	}

	if resp, err := client.GetClips(&params); err == nil {
		if resp.StatusCode == 200 {
			for _, clip := range resp.Data.Clips {
				if date, err := time.Parse(time.RFC3339, clip.CreatedAt); err == nil {
					result = append(result, ClipInfo{
						URL:   clip.URL,
						Title: clip.Title,
						Game:  clip.GameID,
						Date:  date,
						Views: clip.ViewCount,
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
