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

package engine

import (
	"fmt"
	"github.com/rueian/rueidis"
)

type RedisConfiguration struct {
	Host string
	Port int
	DB   int
}

func newRedisClient(engine *Engine) (rueidis.Client, error) {
	logger := engine.logger.With().
		Str("component", "redis").
		Str("host", engine.config.Redis.Host).
		Int("port", engine.config.Redis.Port).
		Int("db", engine.config.Redis.DB).
		Logger()

	logger.Info().Msg("starting redis client")
	c, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{
			fmt.Sprintf("%s:%d", engine.config.Redis.Host, engine.config.Redis.Port),
		},
		SelectDB: engine.config.Redis.DB,
	})

	if err != nil {
		return nil, err
	}

	return c, nil
}
