package main

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/world"
)

func main() {
	engine.RegisterPlugin(world.Plugin)

	engine.Start("example")
	engine.ExecuteCLI()
}
