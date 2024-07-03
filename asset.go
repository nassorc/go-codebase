package gandalf

import (
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Animation struct {
	Texture     rl.Texture2D
	Src         rl.Rectangle // sprite
	FrmOffset   rl.Vector2   // X and Y offset normalized to 1 unit
	TotalFrames int
	CurFrame    int
	Speed       float32
	Rotation    float32
	Scale       float32
}

func (a *Animation) update(n int) {
	// assumes that the animation is laid out from left to right
	var curframe = (n / int(a.Speed)) % a.TotalFrames
	// if current frame exceeds Texture's width, advance to next row
	var nrow = (curframe * int(a.Src.Width)) / int(a.Texture.Width)
	a.CurFrame = curframe
	// move sprite to next frame
	a.Src.Y = float32((nrow + int(a.FrmOffset.Y)) * int(a.Src.Height))                             // calculate next ROW
	a.Src.X = float32(((curframe + int(a.FrmOffset.X)) * int(a.Src.Width)) % int(a.Texture.Width)) // calculate next COL
}

func NewAssetManager() *AssetManager {
	return &AssetManager{
		textures:   make(map[string]rl.Texture2D),
		animations: make(map[string]*Animation),
	}
}

type AssetManager struct {
	textures   map[string]rl.Texture2D
	animations map[string]*Animation
	tickCount  int

	// animSize  int
	// sparse set
	// animations   []*Animation
	// entityToAnim []int // EntityId -> animation index

}

func (mgr *AssetManager) getTexture(name string) (rl.Texture2D, bool) {
	Texture, ok := mgr.textures[name]
	return Texture, ok
}

func (mgr *AssetManager) getAnimation(name string) (*Animation, bool) {
	anim, ok := mgr.animations[name]
	return anim, ok
}

func (mgr *AssetManager) loadTexture(name string, path string) error {
	_, ok := mgr.textures[name]

	if ok {
		return fmt.Errorf("Texture %s already exists", name)
	}

	if _, err := os.Stat(path); err != nil {
		return err
	}

	mgr.textures[name] = rl.LoadTexture(path)

	return nil
}

func (mgr *AssetManager) loadAnimation(animName string, textName string, totalFrames int, Src rl.Rectangle, frmOffset rl.Vector2, Scale float32, Rotation float32, Speed float32) bool {
	var texture, ok = mgr.getTexture(textName)
	if !ok {
		return false
	}
	var anim = Animation{
		Texture:     texture,
		Src:         Src,
		FrmOffset:   frmOffset,
		TotalFrames: totalFrames,
		Scale:       Scale,
		Rotation:    Rotation,
		Speed:       Speed,
	}

	mgr.animations[animName] = &anim

	return true
}

func (mgr *AssetManager) update() {
	for k := range mgr.animations {
		mgr.animations[k].update(mgr.tickCount)
	}

	mgr.tickCount += 1
}

// func (mgr *AssetManager) getAnimationId(entity EntityId) int {
// 	return mgr.entityToAnim[entity]
// }

// func (mgr *AssetManager) hasAnimation(entity EntityId) bool {
// 	return mgr.animations[mgr.entityToAnim[entity]] != nil && mgr.entityToAnim[entity] < mgr.animSize
// }
