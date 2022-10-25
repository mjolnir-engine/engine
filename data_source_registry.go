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
	"github.com/mjolnir-engine/engine/errors"
	"github.com/mjolnir-engine/engine/uid"
	"github.com/rs/zerolog"
)

type dataSourceRegistry struct {
	dataSources map[string]DataSource
	engine      *Engine
	logger      zerolog.Logger
}

func newDataSourceRegistry(engine *Engine) *dataSourceRegistry {
	return &dataSourceRegistry{
		dataSources: make(map[string]DataSource),
		engine:      engine,
		logger:      engine.logger.With().Str("component", "data_sources_registry").Logger(),
	}
}

func (r *dataSourceRegistry) register(dataSource DataSource) {
	r.logger.Debug().Str("name", dataSource.Name()).Msg("registering data source")
	r.dataSources[dataSource.Name()] = dataSource
}

func (r *dataSourceRegistry) start() {
	for _, dataSource := range r.dataSources {
		err := dataSource.Start()

		if err != nil {
			r.logger.Fatal().Err(err).Str("name", dataSource.Name()).Msg("error starting data source")
			panic(err)
		}
	}
}

func (r *dataSourceRegistry) stop() {
	for _, dataSource := range r.dataSources {
		err := dataSource.Stop()

		if err != nil {
			r.logger.Fatal().Err(err).Str("name", dataSource.Name()).Msg("error stopping data source")
			panic(err)
		}
	}
}

// RegisterDataSource registers a data source with the engine. If the data source is already registered, it is replaced.
func (e *Engine) RegisterDataSource(dataSource DataSource) {
	e.dataSourceRegistry.register(dataSource)
}

// AllInDataSource finds all records in a data source. If the data source does not exist, an error is returned.
func (e *Engine) AllInDataSource(dataSourceName string, entities interface{}) error {
	dataSource, ok := e.dataSourceRegistry.dataSources[dataSourceName]

	if !ok {
		return errors.DataSourceNotFoundError{
			Name: dataSourceName,
		}
	}

	return dataSource.All(entities)
}

// DeleteFromDataSource deletes a record in a data source. If the data source does not exist, an error is returned. If the
// record does not exist, an error is returned.
func (e *Engine) DeleteFromDataSource(dataSourceName string, entity interface{}) error {
	dataSource, ok := e.dataSourceRegistry.dataSources[dataSourceName]

	if !ok {
		return errors.DataSourceNotFoundError{
			Name: dataSourceName,
		}
	}

	return dataSource.Delete(entity)
}

// FindInDataSource finds a record in a data source. If the data source does not exist, an error is returned. If the
// record does not exist, an error is returned.
func (e *Engine) FindInDataSource(dataSourceName string, search interface{}, entity interface{}) error {
	dataSource, ok := e.dataSourceRegistry.dataSources[dataSourceName]

	if !ok {
		return errors.DataSourceNotFoundError{
			Name: dataSourceName,
		}
	}

	return dataSource.Find(search, entity)
}

// FindOneInDataSource finds a single record in a data source. If the data source does not exist, an error is returned.
// If the record does not exist, an error is returned.
func (e *Engine) FindOneInDataSource(dataSourceName string, search interface{}, entity interface{}) error {
	dataSource, ok := e.dataSourceRegistry.dataSources[dataSourceName]

	if !ok {
		return errors.DataSourceNotFoundError{
			Name: dataSourceName,
		}
	}

	return dataSource.FindOne(search, entity)
}

// CountInDataSource counts the number of records in a data source. If the data source does not exist, an error is
// returned.
func (e *Engine) CountInDataSource(dataSourceName string, search interface{}) (int64, error) {
	dataSource, ok := e.dataSourceRegistry.dataSources[dataSourceName]

	if !ok {
		return 0, errors.DataSourceNotFoundError{
			Name: dataSourceName,
		}
	}

	return dataSource.Count(search)
}

// SaveInDataSource saves a record in a data source. If the data source does not exist, an error is returned. If the
// record does not exist, an error is returned.
func (e *Engine) SaveInDataSource(dataSourceName string, entity interface{}) (uid.UID, error) {
	dataSource, ok := e.dataSourceRegistry.dataSources[dataSourceName]

	if !ok {
		return uid.New(), errors.DataSourceNotFoundError{
			Name: dataSourceName,
		}
	}

	return dataSource.Save(entity)
}
