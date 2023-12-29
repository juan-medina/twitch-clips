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
	"fmt"
	"time"

	"github.com/juan-medina/twitch-clips/internal/csv"
	"github.com/juan-medina/twitch-clips/internal/twitch"
	"github.com/rs/zerolog/log"
)

func Execute(clientId string, secret string, channel string, filename string, from time.Time, to time.Time, sortBy twitch.SortBy) error {
	log.Info().Str("client_id", clientId).Str("channel", channel).Msg("extract clips")
	if client, err := twitch.GetClient(clientId, secret); err == nil {
		if broadcasterId, err := twitch.GetBroadcasterId(client, channel); err == nil {
			if clips, err := twitch.GetClips(client, broadcasterId, from, to, sortBy); err == nil {
				log.Info().Int("clips", len(clips)).Msg("got clips")
				if err := csv.WriteClipInfoToCSV(filename, clips); err != nil {
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
