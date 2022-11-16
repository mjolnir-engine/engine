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
	engineErrors "github.com/mjolnir-engine/engine/errors"
	engineEvents "github.com/mjolnir-engine/engine/events"
	"github.com/mjolnir-engine/engine/uid"
	"github.com/rs/zerolog"
)

type sessionHandler struct {
	id                                uid.UID
	logger                            zerolog.Logger
	engine                            *Engine
	registry                          *sessionRegistry
	handleSessionSendDataSubscription uid.UID
}

// Start starts the session handler.
func (h *sessionHandler) Start() {
	h.logger.Info().Msg("starting session handler")
	err := h.engine.Publish(engineEvents.SessionStartedEvent{
		Id: h.id,
	})

	if err != nil {
		h.registry.remove(h.id)
	}

	h.handleSessionSendDataSubscription = h.engine.Subscribe(engineEvents.SessionSendDataEvent{
		Id: h.id,
	}, h.handleSessionSendData)

}

// Stop stops the session handler.
func (h *sessionHandler) Stop() {
	h.logger.Info().Msg("stopping session registry")
	err := h.engine.Publish(engineEvents.SessionStoppedEvent{
		Id: h.id,
	})

	if err != nil {
		h.logger.Err(err).Msg("failed to publish session stopped event")
	}
}

func (h *sessionHandler) getController() Controller {
	var controllerName string
	err := h.engine.GetComponent(h.id, "Controller", &controllerName)

	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get controller name")
		h.Stop()

		return nil
	}

	controller, err := h.engine.GetController(controllerName)

	if err != nil {
		h.logger.Error().Err(err).Msg("failed to get controller")
		h.Stop()

		return nil
	}

	return controller
}

func (h *sessionHandler) handleSessionSendData(event EventMessage) {
	sessionSendDataEvent := &engineEvents.SessionSendDataEvent{}

	if err := event.Unmarshal(sessionSendDataEvent); err != nil {
		h.logger.Error().Err(err).Msg("failed to unmarshal session send data event")
		return
	}

	h.logger.Debug().Str("sessionId", string(sessionSendDataEvent.Id)).Msg("session sent data")
	controller := h.getController()

	if controller == nil {
		h.logger.Warn().Msg("failed to get controller")
		return
	}

	err := controller.HandleInput(
		&ControllerContext{SessionId: sessionSendDataEvent.Id, Engine: h.engine},
		string(sessionSendDataEvent.Data),
	)

	if err != nil {
		h.logger.Error().Err(err).Msg("failed to handle input")
	}
}

type sessionRegistry struct {
	logger                   zerolog.Logger
	sessions                 map[uid.UID]*sessionHandler
	engine                   *Engine
	sessionStartSubscription uid.UID
	sessionStopSubscription  uid.UID
}

func newSessionRegistry(engine *Engine) *sessionRegistry {
	return &sessionRegistry{
		engine:   engine,
		sessions: make(map[uid.UID]*sessionHandler),
		logger:   engine.logger.With().Str("component", "session_registry").Logger(),
	}
}

func (r *sessionRegistry) start() {
	r.logger.Info().Msg("starting session registry")

	r.sessionStartSubscription = r.engine.Subscribe(engineEvents.SessionStartEvent{}, r.handleSessionStartEvent)
	r.sessionStopSubscription = r.engine.Subscribe(engineEvents.SessionStopEvent{}, r.handleSessionStopEvent)
}

func (r *sessionRegistry) stop() {
	r.logger.Info().Msg("stopping session registry")

	r.engine.Unsubscribe(r.sessionStartSubscription)
	r.engine.Unsubscribe(r.sessionStopSubscription)
}

func (r *sessionRegistry) add(session *Session) {
	handler := &sessionHandler{
		id:     session.Id,
		logger: r.logger.With().Str("sessionId", string(session.Id)).Logger(),
		engine: r.engine,
	}

	r.sessions[session.Id] = handler
	handler.Start()
}

func (r *sessionRegistry) remove(id uid.UID) {
	handler := r.sessions[id]

	if handler == nil {
		return
	}

	handler.Stop()

	delete(r.sessions, id)

	err := r.engine.RemoveEntity(id)

	if err != nil {
		r.logger.Error().Err(err).Msg("failed to remove session from entity registry")
	}
}

func (r *sessionRegistry) handleSessionStartEvent(event EventMessage) {
	sessionStartEvent := &engineEvents.SessionStartedEvent{}

	if err := event.Unmarshal(sessionStartEvent); err != nil {
		r.logger.Error().Err(err).Msg("failed to unmarshal session start event")
		return
	}

	r.logger.Debug().Str("sessionId", string(sessionStartEvent.Id)).Msg("session started")
	session := &Session{
		Id:         sessionStartEvent.Id,
		Controller: r.engine.Config.DefaultController,
	}

	err := r.engine.AddEntity(session)

	if err != nil {
		r.logger.Error().Err(err).Msg("failed to add session to entity registry")
		return
	}

	r.add(session)
}

func (r *sessionRegistry) handleSessionStopEvent(event EventMessage) {
	sessionStopEvent := &engineEvents.SessionStopEvent{}

	if err := event.Unmarshal(sessionStopEvent); err != nil {
		r.logger.Error().Err(err).Msg("failed to unmarshal session stop event")
		return
	}

	r.logger.Debug().Str("sessionId", string(sessionStopEvent.Id)).Msg("session stopped")
	r.remove(sessionStopEvent.Id)
}

// SendToSession sends data to a session. The data can be anything that can be marshalled to JSON. This will dispatch
// a SessionSendDataEvent. If the session does not exist, this will return an error.
func (e *Engine) SendToSession(sessionId uid.UID, data []byte) error {
	exists, err := e.HasEntity(sessionId)

	if err != nil {
		return err
	}

	if !exists {
		return engineErrors.SessionNotFoundError{
			Id: sessionId,
		}
	}

	return e.Publish(engineEvents.SessionSendDataEvent{
		Id:   sessionId,
		Data: data,
	})
}

// GetSessionController gets the controller for a session. If the session does not exist, this will return an error.
func (e *Engine) GetSessionController(id uid.UID) (Controller, error) {
	var name string
	err := e.GetComponent(id, "Controller", &name)

	if err != nil {
		return nil, err
	}

	controller, err := e.GetController(name)

	if err != nil {
		return nil, err
	}

	return controller, nil

}
