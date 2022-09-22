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

package data_sources

import (
	"github.com/mjolnir-mud/engine/plugins/data_sources/internal/plugin"
	"github.com/mjolnir-mud/engine/plugins/data_sources/internal/registry"
	"github.com/mjolnir-mud/engine/plugins/data_sources/pkg/data_source"
)

// All loads all entities from a data source. It will call `ecs.Create` passing the map returned by the data source
// for each entity, and return a map of entities keyed by their ids.
func All(source string) (map[string]map[string]interface{}, error) {
	return registry.All(source)
}

// Count returns the number of entities in a data source using the provided map as a filter. If the data source does not
// exist, an error will be returned.
func Count(source string, filter map[string]interface{}) (int64, error) {
	return registry.Count(source, filter)
}

// CreateEntity creates a new entity in a data source. If the data source does not exist, an error will be thrown. If the
// entityType does not exist, an error will be thrown. It returns the id map representing the entities components as
// well as the entity metadata.
func CreateEntity(dataSource string, entityType string, data map[string]interface{}) (string, map[string]interface{}, error) {
	return registry.CreateEntity(dataSource, entityType, data)
}

// CreateEntityWithId creates a new entity in a data source with the provided id. If the data source does not exist, an
// error will be thrown. If the entityType does not exist, an error will be thrown. It returns the id map representing
// the entities components as well as the entity metadata.
func CreateEntityWithId(dataSource string, entityType string, id string, data map[string]interface{}) (map[string]interface{}, error) {
	return registry.CreateEntityWithId(dataSource, entityType, id, data)
}

// Delete deletes an entity from a data source. If the data source does not exist, an error will be thrown.
func Delete(source string, entityId string) error {
	return registry.Delete(source, entityId)
}

// Find returns all entities in a data source that match the provided filter. If the data source does not exist, an
// error will be thrown.
func Find(source string, filter map[string]interface{}) (map[string]map[string]interface{}, error) {
	return registry.Find(source, filter)
}

// FindOne returns the first entity in a data source that matches the provided filter. If the data source does not
// exist, an error will be thrown.
func FindOne(source string, filter map[string]interface{}) (string, map[string]interface{}, error) {
	return registry.FindOne(source, filter)
}

// FindAndDelete deletes all entities in a data source that match the provided filter. If the data source does not exist,
// an error will be thrown.
func FindAndDelete(source string, filter map[string]interface{}) error {
	return registry.FindAndDelete(source, filter)
}

// Register registers a data source with the registry. If a data source with the same name is already registered,
//i it will be overwritten.
func Register(source data_source.DataSource) {
	registry.Register(source)
}

// Save saves data to a data source for a given entity. If the entity does not have a valid metadata field an error will
// be thrown. If the data source does not exist, an error will be thrown. If the metadata field does not have a type
// set, an error will be thrown. If the entity exists in the data source, it will be overwritten.
func Save(source string, entityId string, entity map[string]interface{}) error {
	return registry.SaveWithId(source, entityId, entity)
}

// SaveWithId saves data to a data source for a given entity. If the data source does not exist, an error will be thrown.
func SaveWithId(source string, entityId string, entity map[string]interface{}) error {
	return registry.SaveWithId(source, entityId, entity)
}

var Plugin = plugin.Plugin
