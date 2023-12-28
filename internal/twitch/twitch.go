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
	"fmt"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/rs/zerolog/log"
)

func GetClient(clientId string, secret string) (*helix.Client, error) {
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

func GetBroadcasterId(client *helix.Client, channel string) (string, error) {
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
