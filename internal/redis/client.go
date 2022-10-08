/*
 * Copyright (c) 2022 eightfivefour llc. All rights reserved.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
 * documentation files (the "Software"), to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
 * permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
 * Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
 * WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 * OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package redis

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rueian/rueidis"

	engineLogger "github.com/mjolnir-mud/engine/logger"
)

var logger zerolog.Logger

// New creates a new redis client. It accepts a `redis.Configuration` struct which defines the host, port, and database
func New(config *Configuration) (rueidis.Client, error) {
	logger = engineLogger.Instance.With().
		Str("component", "redis").
		Str("host", config.Host).
		Int("port", config.Port).
		Int("db", config.DB).
		Logger()

	logger.Info().Msg("starting redis client")
	c, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{
			fmt.Sprintf("%s:%d", config.Host, config.Port),
		},
		SelectDB: config.DB,
	})

	if err != nil {
		return nil, err
	}

	return c, nil
}