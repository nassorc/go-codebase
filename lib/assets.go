package animation

import (
	"bytes"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

func NewFrameInfo(Cell, Duration int) FrameInfo {
	return FrameInfo{Cell, Duration}
}

type FrameInfo struct{ Cell, Duration int }

type Sheet []FrameInfo

type Animation struct {
	Sheet       Sheet
	TextureName string
	TotalFrames int
	TileWidth   int
	TileHeight  int
}

type AssetMgr struct {
	Textures   map[string]*ebiten.Image
	Animations map[string]Animation
	TickCount  float32
}

func (mgr *AssetMgr) LoadTexture(name string, buf []byte) bool {
	image, _, err := image.Decode(bytes.NewReader(buf))

	if err != nil {
		return false
	}

	texture := ebiten.NewImageFromImage(image)

	if mgr.Textures == nil {
		mgr.Textures = make(map[string]*ebiten.Image)
	}

	mgr.Textures[name] = texture

	return true
}

func (mgr *AssetMgr) LoadAnimation(name string, anim Animation) bool {
	if mgr.Animations == nil {
		mgr.Animations = make(map[string]Animation)
	}

	if _, ok := mgr.Textures[anim.TextureName]; !ok {
		return false
	}

	mgr.Animations[name] = anim
	return true
}
