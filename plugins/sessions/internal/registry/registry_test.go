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

package registry

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/pkg/event"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	ecsTesting "github.com/mjolnir-mud/engine/plugins/ecs/pkg/testing"
	"github.com/mjolnir-mud/engine/plugins/sessions/entities/session"
	events2 "github.com/mjolnir-mud/engine/plugins/sessions/events"
	engineTesting "github.com/mjolnir-mud/engine/testing"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() {
	engineTesting.RegisterSetupCallback("sessions", func() {
		engine.RegisterPlugin(ecs.Plugin)
		ecs.RegisterEntityType(session.Type)
		ecsTesting.Setup()
	})
	engineTesting.Setup("world")

	Start()
}

func teardown() {
	Stop()
	ecsTesting.Teardown()
	engineTesting.Teardown()
}

func TestPlayerConnected(t *testing.T) {
	setup()
	defer teardown()

	ch := make(chan bool)

	RegisterSessionStartedHandler(func(id string) error {
		go func() { ch <- true }()
		return nil
	})

	err := engine.Publish(events2.PlayerConnectedEvent{})

	assert.NoError(t, err)

	<-ch

	assert.Len(t, sessionHandlers, 1)
}

func TestRegisterSessionStartedHandler(t *testing.T) {
	setup()
	defer teardown()

	RegisterSessionStartedHandler(func(id string) error {
		return nil
	})

	assert.Len(t, sessionStartedHandlers, 1)
}

func TestRegisterSessionStoppedHandler(t *testing.T) {
	setup()
	defer teardown()

	RegisterSessionStoppedHandler(func(id string) error {
		return nil
	})

	assert.Len(t, sessionStoppedHandlers, 1)
}

func TestRegisterLineHandler(t *testing.T) {
	setup()
	defer teardown()

	RegisterReceiveLineHandler(func(id string, line string) error {
		return nil
	})

	assert.Len(t, receiveLineHandlers, 1)
}

func TestSendLine(t *testing.T) {
	setup()
	defer teardown()

	ch := make(chan bool)

	err := ecs.AddEntityWithID("session", "testSession", map[string]interface{}{})

	assert.Nil(t, err)

	add(NewSessionHandler("testSession"))

	sub := engine.Subscribe(events2.PlayerOutputEvent{Id: "testSession"}, func(event event.EventPayload) {
		go func() { ch <- true }()
	})

	defer sub.Stop()

	err = SendLine("testSession", "testing")

	assert.NoError(t, err)

	<-ch
}

func TestRegisterSendLineHandler(t *testing.T) {
	setup()
	defer teardown()

	RegisterSendLineHandler(func(id string, line string) error {
		return nil
	})

	assert.Len(t, sendLineHandlers, 1)
}
