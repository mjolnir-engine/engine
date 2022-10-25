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

package testing

import (
	"github.com/google/uuid"
	"github.com/mjolnir-engine/engine"
)

// Setup sets up the engine for testing. It accepts a callback function where the engine instance is passed as a
// parameter. The callback function should be to execute any code that should be run before the engine is started.
func Setup(cb func(e *engine.Engine)) *engine.Engine {
	puid, _ := uuid.NewRandom()
	prefix := puid.String()

	e := engine.New(&engine.Configuration{
		Redis: &engine.RedisConfiguration{
			Host: "localhost",
			Port: 6379,
			DB:   1,
		},
		InstanceId:        prefix,
		DefaultController: "test",
		Environment:       "test",
	})

	e.RegisterService("test")
	cb(e)

	e.Start("test")

	return e
}
