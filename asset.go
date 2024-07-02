package gandalf

import (
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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
