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

// Plugin is the interface that must be implemented by a Mjolnir plugin.
type Plugin interface {
	// Name returns the name of the plugin. Plugin names must be unique. If two plugins are registered with the same
	// name, an error will be returned.
	Name() string

	// BeforeStart is called before the engine starts. This is the appropriate time to create any connections to any
	// external services that the plugin needs to use, or other setup items. This is effectively the plugin's
	// constructor.
	BeforeStart() error

	// AfterStart is called after the engine starts. This is the appropriate time register any resources that the plugin
	// may provide to the engine. For example, a plugin may register a data source, or a controller.
	AfterStart() error

	// BeforeStop is called before the engine stops. This is the appropriate time to close any connections to any
	// external services that the plugin needs to use, or other cleanup items. This is effectively the plugin's
	// destructor.
	BeforeStop() error
}
