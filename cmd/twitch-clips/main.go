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

package main

import (
	"flag"
	"os"
	"time"

	"github.com/juan-medina/twitch-clips/internal/cmd/extract"
	"github.com/juan-medina/twitch-clips/internal/times"
	"github.com/juan-medina/twitch-clips/internal/twitch"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Info().Msg("Twitch Clips by Juan Medina")

	clientId := flag.String("client_id", os.Getenv("TWITCH_CLIENT_ID"), "Twitch Client Id, or set TWITCH_CLIENT_ID environment variable")
	secret := flag.String("secret", os.Getenv("TWITCH_SECRET"), "Twitch Secret, or set TWITCH_SECRET environment variable")
	channel := flag.String("channel", os.Getenv("TWITCH_CHANNEL"), "Twitch Channel, or set TWITCH_CHANNEL environment variable")
	output := flag.String("output", "clips.csv", "Output file")

	from := flag.String("date_from", "", "Date Time from in format: "+times.DateFormat)
	to := flag.String("date_to", "", "Date Time to in format: "+times.DateFormat)

	sort := flag.String("sort_by", "DATE", "Sort by: DATE or VIEWS")

	flag.Parse()

	if *clientId == "" || *channel == "" || *secret == "" {
		flag.Usage()
		os.Exit(1)
		return
	}

	timeFrom := times.NilDateTime
	timeTo := times.NilDateTime

	if *from != "" {
		var err error
		if timeFrom, err = time.Parse(times.DateFormat, *from); err != nil {
			log.Error().Err(err).Msg("invalid from format")
			os.Exit(1)
			return
		}
	}

	if *to != "" {
		var err error
		if timeTo, err = time.Parse(times.DateFormat, *to); err != nil {
			log.Error().Err(err).Msg("invalid to format")
			os.Exit(1)
			return
		}
	}

	sortBy := twitch.SortByDate

	if *sort != "" {
		*sort = "SORT_" + *sort
		if *sort != string(twitch.SortByDate) && *sort != string(twitch.SortByViews) {
			flag.Usage()
			os.Exit(1)
			return
		} else {
			sortBy = twitch.SortBy(*sort)
		}
	}

	if err := extract.Execute(*clientId, *secret, *channel, *output, timeFrom, timeTo, sortBy); err != nil {
		log.Error().Err(err).Msg("failed to extract clips")
		os.Exit(1)
		return
	}
}
