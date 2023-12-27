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

	"github.com/juan-medina/twitch-clips/internal/cmd/extract"
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

	flag.Parse()

	if *clientId == "" || *channel == "" || *secret == "" {
		flag.Usage()
		os.Exit(1)
		return
	}

	if err := extract.Execute(*clientId, *secret, *channel, *output); err != nil {
		log.Error().Err(err).Msg("failed to extract clips")
		os.Exit(1)
		return
	}
}
