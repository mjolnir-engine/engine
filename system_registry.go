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
	engineEvents "github.com/mjolnir-mud/engine/events"
	"github.com/rs/zerolog"
)

type systemRegistry struct {
	logger  zerolog.Logger
	systems map[string]System
	engine  *Engine
}

func newSystemRegistry(e *Engine) *systemRegistry {
	registry := &systemRegistry{
		logger:  e.Logger.With().Str("component", "system_registry").Logger(),
		systems: make(map[string]System),
		engine:  e,
	}

	registry.logger.Info().Msg("initializing system registry")

	return registry
}

// Register a system with the engine. If a system with the same name is already registered, it will be overwritten.
func (r *systemRegistry) Register(s System) {
	r.logger.Info().Str("system", s.Name()).Msg("registering system")
	r.systems[s.Name()] = s
}

// Start starts the system registry.
func (r *systemRegistry) Start() {
	r.logger.Info().Msg("starting")
	r.engine.PSubscribe(engineEvents.ComponentAddedEvent{}, r.onComponentAdded)
	r.engine.PSubscribe(engineEvents.ComponentRemovedEvent{}, r.onComponentRemoved)
	r.engine.PSubscribe(engineEvents.ComponentUpdatedEvent{}, r.onComponentUpdated)
}

// Stop stops the system registry.
func (r *systemRegistry) Stop() {
	r.logger.Info().Msg("stopping system registry")
	r.engine.PUnsubscribe(engineEvents.ComponentAddedEvent{})
	r.engine.PUnsubscribe(engineEvents.ComponentRemovedEvent{})
	r.engine.PUnsubscribe(engineEvents.ComponentUpdatedEvent{})
}

func (r *systemRegistry) onComponentAdded(e EventMessage) {
	event := &engineEvents.ComponentAddedEvent{}
	if err := e.Unmarshal(event); err != nil {
		r.logger.Error().Err(err).Msg("unable to unmarshal event")
		return
	}

	for _, s := range r.systems {
		cName := s.Component()
		exists, err := r.engine.HasComponent(event.EntityId, cName)

		if err != nil {
			r.logger.Error().Err(err).Msg("unable to check if component exists")
		}

		matches := s.Match(event.Name, event.Value)

		if exists && matches {
			err := s.ComponentAdded(event.EntityId, event.Name, event.Value)

			if err != nil {
				r.logger.
					Error().
					Str("system", s.Name()).
					Str("entity", event.EntityId.String()).
					Str("entityComponentName", event.Name).
					Err(err).
					Msg("unable to call ComponentAdded")
			}
		}

		if matches {
			err := s.MatchingComponentAdded(event.EntityId, event.Value)

			if err != nil {
				r.logger.
					Error().
					Str("system", s.Name()).
					Str("entity", event.EntityId.String()).
					Str("entityComponentName", event.Name).
					Err(err).
					Msg("unable to call MatchingComponentAdded")
			}
		}

	}
}

func (r *systemRegistry) onComponentRemoved(e EventMessage) {
	event := &engineEvents.ComponentRemovedEvent{}
	if err := e.Unmarshal(event); err != nil {
		r.logger.Error().Err(err).Msg("unable to unmarshal event")
		return
	}

	for _, s := range r.systems {
		cName := s.Component()
		exists, err := r.engine.HasComponent(event.EntityId, cName)

		if err != nil {
			r.logger.Error().Err(err).Msg("unable to check if component exists")
		}

		matches := s.Match(event.Name, event.Value)

		if exists && matches {
			err := s.ComponentRemoved(event.EntityId, event.Name, event.Value)

			if err != nil {
				r.logger.
					Error().
					Str("system", s.Name()).
					Str("entity", event.EntityId.String()).
					Str("entityComponentName", event.Name).
					Err(err).
					Msg("unable to call ComponentRemoved")
			}
		}

		if matches {
			err := s.MatchingComponentRemoved(event.EntityId, event.Value)

			if err != nil {
				r.logger.
					Error().
					Str("system", s.Name()).
					Str("entity", event.EntityId.String()).
					Str("entityComponentName", event.Name).
					Err(err).
					Msg("unable to call MatchingComponentRemoved")
			}
		}

	}
}

func (r *systemRegistry) onComponentUpdated(e EventMessage) {
	event := &engineEvents.ComponentUpdatedEvent{}
	if err := e.Unmarshal(event); err != nil {
		r.logger.Error().Err(err).Msg("unable to unmarshal event")
		return
	}

	for _, s := range r.systems {
		cName := s.Component()
		exists, err := r.engine.HasComponent(event.EntityId, cName)

		if err != nil {
			r.logger.Error().Err(err).Msg("unable to check if component exists")
		}

		matches := s.Match(event.Name, event.Value)

		if exists && matches {
			err := s.ComponentUpdated(event.EntityId, event.Name, event.Value, event.PreviousValue)

			if err != nil {
				r.logger.
					Error().
					Str("system", s.Name()).
					Str("entity", event.EntityId.String()).
					Str("entityComponentName", event.Name).
					Err(err).
					Msg("unable to call ComponentUpdated")
			}
		}

		if matches {
			err := s.MatchingComponentUpdated(event.EntityId, event.Value, event.PreviousValue)

			if err != nil {
				r.logger.
					Error().
					Str("system", s.Name()).
					Str("entity", event.EntityId.String()).
					Str("entityComponentName", event.Name).
					Err(err).
					Msg("unable to call MatchingComponentUpdated")
			}
		}

	}
}
