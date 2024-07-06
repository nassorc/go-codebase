package scenes

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/nassorc/gandalf"
	"github.com/nassorc/gandalf/cmd/components"
)

func NewMenuScene(color rl.Color, selected rl.Color) *MenuScene {

	return &MenuScene{
		color:         color,
		selectedColor: selected,
		mainmenu:      &Menu{},
	}
}

type MenuScene struct {
	color         rl.Color
	selectedColor rl.Color
	world         *gandalf.World
	mainmenu      *Menu
}

func (scene *MenuScene) Setup(world *gandalf.World) {
	scene.world = world

	var onPlay = func() {
		scene.world.ChangeScene(&Level1{})
	}
	var onExit = func() {
		fmt.Println("exiting")
	}
	var menuItems = []string{"Play", "Exit"}
	var menuActions = []func(){onPlay, onExit}

	for idx := 0; idx < len(menuItems); idx++ {
		scene.mainmenu.AddItem(menuItems[idx], menuActions[idx])
	}

	world.RegisterComponents(&components.Transform{})
	world.RegisterSystem(scene.MenuLogic)
	world.RegisterSystem(scene.RenderMenuItem, &components.Transform{})
}

func (scene *MenuScene) MenuLogic(entities []gandalf.EntityHandle) {
	// update selected scene
	if rl.IsKeyPressed(rl.KeyS) {
		// scene.SelectedIdx = (scene.SelectedIdx + 1) % scene.Size
		scene.mainmenu.Next()
	}
	if rl.IsKeyPressed(rl.KeyW) {
		scene.mainmenu.Prev()
	}

	// selecting scene
	if rl.IsKeyPressed(rl.KeyEnter) {
		scene.mainmenu.Select()
	}
}

func (scene *MenuScene) RenderMenuItem(entities []gandalf.EntityHandle) {
	for idx, item := range scene.mainmenu.Items() {
		var color rl.Color

		if scene.mainmenu.Current() == idx {
			color = scene.selectedColor
		} else {
			color = scene.color
		}

		rl.DrawText(item, 24, int32(24+38*idx), 32, color)
	}
}

func NewMenu(items []string, actions []func()) *Menu {
	return &Menu{
		0,
		items,
		actions,
	}
}

type Menu struct {
	selected int
	items    []string
	actions  []func()
}

func (scene *Menu) AddItem(item string, action func()) {
	scene.items = append(scene.items, item)
	scene.actions = append(scene.actions, action)
}

func (scene *Menu) Next() {
	scene.selected = (scene.selected + 1) % len(scene.items)
}

func (scene *Menu) Prev() {
	scene.selected = scene.selected - 1
	if scene.selected < 0 {
		scene.selected = len(scene.items) - 1
	}
}

func (scene *Menu) Select() bool {
	scene.actions[scene.selected]()
	return true
}

func (scene *Menu) Selectn(n int) bool {
	scene.actions[n]()
	return true
}

func (scene *Menu) Current() int {
	return scene.selected
}

func (scene *Menu) Items() []string {
	return scene.items[:]
}
