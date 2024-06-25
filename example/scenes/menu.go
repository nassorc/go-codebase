package scenes

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nassorc/gandalf"
	"github.com/nassorc/gandalf/cmd/components"
)

var MenuItems = []string{
	"Play",
	"Exit",
}

var MenuActions = []func(*gandalf.World){
	func(world *gandalf.World) {
		world.ChangeScene(&Game{})
	},
	func(world *gandalf.World) {
	},
}

func NewMenuScene(items []string, actions []func(*gandalf.World), Color rl.Color, selectedColor rl.Color) *MenuScene {
	return &MenuScene{
		items:         items,
		actions:       actions,
		SelectedIdx:   0,
		Size:          len(items),
		Color:         Color,
		SelectedColor: selectedColor,
	}
}

type MenuScene struct {
	world         *gandalf.World
	actions       []func(*gandalf.World)
	items         []string
	SelectedIdx   int
	Size          int
	Color         rl.Color
	SelectedColor rl.Color
}

func (menu *MenuScene) Setup(world *gandalf.World) {
	menu.world = world

	world.RegisterComponents(&components.Transform{})

	world.RegisterSystem(menu.MenuLogic)
	world.RegisterSystem(menu.RenderMenuItem, &components.Transform{})
}

func (menu *MenuScene) MenuLogic(world *gandalf.World, entities []*gandalf.EntityHandle) {
	// update selected menu
	if rl.IsKeyPressed(rl.KeyS) {
		menu.SelectedIdx = (menu.SelectedIdx + 1) % menu.Size
	}
	if rl.IsKeyPressed(rl.KeyW) {
		menu.SelectedIdx = menu.SelectedIdx - 1
		if menu.SelectedIdx < 0 {
			menu.SelectedIdx = menu.Size - 1
		}
	}

	// selecting menu
	if rl.IsKeyPressed(rl.KeyEnter) {
		menu.actions[menu.SelectedIdx](world)
	}
}

func (menu *MenuScene) RenderMenuItem(world *gandalf.World, entities []*gandalf.EntityHandle) {
	for idx, item := range menu.items {
		var color rl.Color

		if menu.SelectedIdx == idx {
			color = menu.SelectedColor
		} else {
			color = menu.Color
		}

		rl.DrawText(item, int32((500 / 2)), int32((500/2)+32*idx), 32, color)
	}

}
