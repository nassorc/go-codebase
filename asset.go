package gandalf

import (
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Animation struct {
	texture     rl.Texture2D
	row         int
	col         int
	offsetRow   int
	offsetCol   int
	totalFrames int
	width       int
	height      int
	tileWidth   int
	tileHeight  int
	curFrame    int
	speed       int
	rotation    int
}

func (a *Animation) update(n int) {
	// assumes that the animation is laid out from left to right
	var curframe = (n / a.speed) % a.totalFrames
	// if current frame exceeds texture's width, advance to next row
	var nrow = (curframe * a.tileWidth) / int(a.texture.Width)
	a.curFrame = curframe
	a.row = (nrow + a.offsetRow) * a.tileHeight
	a.col = ((curframe + a.offsetCol) * a.tileWidth) % int(a.texture.Width)
}

type AssetManager struct {
	textures map[string]rl.Texture2D
}

func (mgr *AssetManager) getTexture(name string) (rl.Texture2D, bool) {
	texture, ok := mgr.textures[name]
	return texture, ok
}

func (mgr *AssetManager) loadTexture(name string, path string) error {
	_, ok := mgr.textures[name]

	if ok {
		return fmt.Errorf("texture %s already exists", name)
	}

	if _, err := os.Stat(path); err != nil {
		return err
	}

	mgr.textures[name] = rl.LoadTexture(path)

	return nil
}
