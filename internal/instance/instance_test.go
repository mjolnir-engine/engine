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

package instance

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {
	setup()
	defer teardown()

	assert.Equal(t, "testing", GetEnv())
	assert.NotNil(t, log)
}

func TestGetGameName(t *testing.T) {
	setup()
	defer teardown()

	assert.Equal(t, "testing", GetGameName())
}

func TestGetEnv(t *testing.T) {
	setup()
	defer teardown()

	assert.Equal(t, "testing", GetEnv())
}

func TestStart(t *testing.T) {
	setup()
	defer teardown()

	beforeStartCalled := make(chan bool)
	afterStartCalled := make(chan bool)
	beforeStartCallbacksForEnvCalled := make(chan bool)
	afterStartCallbacksForEnvCalled := make(chan bool)

	RegisterBeforeStartCallback(func() {
		go func() {
			beforeStartCalled <- true
		}()
	})

	RegisterBeforeStartCallbackForEnv("testing", func() {
		go func() {
			beforeStartCallbacksForEnvCalled <- true
		}()
	})

	RegisterAfterStartCallback(func() {
		go func() {
			afterStartCalled <- true
		}()
	})

	RegisterAfterStartCallbackForEnv("testing", func() {
		go func() {
			afterStartCallbacksForEnvCalled <- true
		}()
	})

	Start()
	defer Stop()

	<-beforeStartCalled
	<-afterStartCalled
	<-beforeStartCallbacksForEnvCalled
	<-afterStartCallbacksForEnvCalled

	assert.True(t, true)
}

func TestStop(t *testing.T) {
	setup()
	defer teardown()

	beforeStopCalled := make(chan bool)
	afterStopCalled := make(chan bool)
	beforeStopCallbacksForEnvCalled := make(chan bool)

	RegisterBeforeStopCallback(func() {
		go func() {
			beforeStopCalled <- true
		}()
	})

	RegisterAfterStopCallback(func() {
		go func() {
			afterStopCalled <- true
		}()
	})

	RegisterBeforeStopCallbackForEnv("testing", func() {
		go func() {
			beforeStopCallbacksForEnvCalled <- true
		}()
	})

	Start()

	Stop()
	<-beforeStopCalled
	<-afterStopCalled
	<-beforeStopCallbacksForEnvCalled
	assert.True(t, true)
}

func TestStartService(t *testing.T) {
	setup()
	defer teardown()

	beforeStartCalled := make(chan bool)
	beforeStartCalledForEnv := make(chan bool)
	afterStartCalled := make(chan bool)
	afterStartCalledForEnv := make(chan bool)

	RegisterBeforeServiceStartCallback("testing", func() {
		go func() {
			beforeStartCalled <- true
		}()
	})

	RegisterBeforeServiceStartCallbackForEnv("testing", "testing", func() {
		go func() {
			beforeStartCalledForEnv <- true
		}()
	})

	RegisterAfterServiceStartCallback("testing", func() {
		go func() {
			afterStartCalled <- true
		}()
	})

	RegisterAfterServiceStartCallbackForEnv("testing", "testing", func() {
		go func() {
			afterStartCalledForEnv <- true
		}()
	})

	Start()

	StartService("testing")
	defer StopService("testing")
	defer Stop()

	<-beforeStartCalled
	<-beforeStartCalledForEnv
	<-afterStartCalled
	<-afterStartCalledForEnv

	assert.True(t, true)
}

func TestStopService(t *testing.T) {
	setup()
	defer teardown()

	beforeStopCalled := make(chan bool)
	beforeStopCalledForEnv := make(chan bool)
	afterStopCalled := make(chan bool)
	afterStopCalledForEnv := make(chan bool)

	RegisterBeforeServiceStopCallback("testing", func() {
		go func() {
			beforeStopCalled <- true
		}()
	})

	RegisterBeforeServiceStopCallbackForEnv("testing", "testing", func() {
		go func() {
			beforeStopCalledForEnv <- true
		}()
	})

	RegisterAfterServiceStopCallback("testing", func() {
		go func() {
			afterStopCalled <- true
		}()
	})

	RegisterAfterServiceStopCallbackForEnv("testing", "testing", func() {
		go func() {
			afterStopCalledForEnv <- true
		}()
	})

	Start()

	StartService("testing")
	StopService("testing")
	defer Stop()

	<-beforeStopCalled
	<-beforeStopCalledForEnv
	<-afterStopCalled
	<-afterStopCalledForEnv

	assert.True(t, true)
}
