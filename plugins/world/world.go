package world

import (
	"github.com/mjolnir-mud/engine"
	"github.com/mjolnir-mud/engine/plugins/ecs"
	"github.com/mjolnir-mud/engine/plugins/world/internal/controller_registry"
	"github.com/mjolnir-mud/engine/plugins/world/internal/session"
	"github.com/mjolnir-mud/engine/plugins/world/pkg/controller"
	session2 "github.com/mjolnir-mud/engine/plugins/world/pkg/systems/session"
	"github.com/spf13/cobra"
)

type world struct {
	stop chan bool
}

func (w world) Name() string {
	return "world"
}

var command = &cobra.Command{
	Use:   "world",
	Short: "Mjolnir MUD Engine world service",
	Long:  "Mjolnir MUD Engine world service powers the game world",
}

func (w world) Registered() error {
	controller_registry.Start()

	engine.RegisterCLICommand(command)

	ecs.RegisterSystem(session2.System)

	session.RegisterSessionStartedHandler(func(id string) error {
		err := ecs.AddEntityWithID(id, "session", map[string]interface{}{})

		if err != nil {
			return err
		}

		return session2.Start(id)
	})

	return nil
}

// GetController returns the controller of the given name. If the controller is not found then an error is returned.
func GetController(name string) (controller.Controller, error) {
	return controller_registry.Get(name)
}

var Plugin = world{
	stop: make(chan bool),
}
