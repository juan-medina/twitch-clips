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

package csv

import (
	"encoding/csv"
	"os"

	"github.com/juan-medina/twitch-clips/internal/twitch"
	"github.com/rs/zerolog/log"
)

func WriteClipInfoToCSV(filename string, data []twitch.ClipInfo) error {
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
		row := []string{clip.Title, clip.Game, clip.Date.Format("2006-01-02T15:04:05"), clip.URL}
		if err = writer.Write(row); err != nil {
			return err
		}
	}

	// Flush the writer to write any buffered data to the file
	writer.Flush()

	return writer.Error()
}
