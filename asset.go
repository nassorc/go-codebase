package gandalf

// rl "github.com/gen2brain/raylib-go/raylib"
import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Animation struct {
	Src         image.Rectangle // sprite
	FrmSize     Vec2
	FrmOffset   Vec2 // X and Y offset normalized to 1 unit
	Texture     *ebiten.Image
	TotalFrames int
	CurFrame    int
	Speed       float32
	Rotation    float32
	Scale       float32
}

func (anim *Animation) Draw(screen *ebiten.Image, op *ebiten.DrawImageOptions) {
	screen.DrawImage(anim.Texture.SubImage(anim.Src).(*ebiten.Image), op)
}

func (anim *Animation) Update(n int) {
	// assumes that the animation is laid out from left to right
	var curframe = (n / int(anim.Speed)) % anim.TotalFrames

	// wraps frame to next row each time we go beyond the texture's width
	var nrow = (curframe * int(anim.Src.Dx())) / int(anim.Texture.Bounds().Dx())
	anim.CurFrame = curframe

	// move sprite to next frame
	row := (nrow + int(anim.FrmOffset.Y)) * anim.Src.Dy()
	// wraps col back to zero, if col goes beyond the texture's width
	col := (curframe + int(anim.FrmOffset.X)) * anim.Src.Dx() % int(anim.Texture.Bounds().Dx())

	anim.Src = image.Rect(col, row, col+int(anim.FrmSize.X), row+int(anim.FrmSize.Y))
	// anim.Src = image.Rect(32*((n/5)%8), int(anim.FrmOffset.Y)*int(anim.FrmSize.Y), 32*((n/5)%8)+32, 64)
}

func NewAssetManager() *AssetManager {
	return &AssetManager{
		textures:   make(map[string]*ebiten.Image),
		animations: make(map[string]*Animation),
	}
}

type AssetManager2[Texture any, Font any] struct {
	textures map[string]Texture
}

type AssetManager struct {
	textures   map[string]*ebiten.Image
	animations map[string]*Animation
	tickCount  int

	// animSize  int
	// sparse set
	// animations   []*Animation
	// entityToAnim []int // EntityId -> animation index
}

func (mgr *AssetManager) loadTexture(name string, img image.Image) error {
	_, ok := mgr.textures[name]

	if ok {
		return fmt.Errorf("texture %s already exists", name)
	}

	mgr.textures[name] = ebiten.NewImageFromImage(img)

	return nil
}

func (mgr *AssetManager) getTexture(name string) (*ebiten.Image, bool) {
	Texture, ok := mgr.textures[name]
	return Texture, ok
}

func (mgr *AssetManager) getAnimation(name string) (*Animation, bool) {
	anim, ok := mgr.animations[name]
	return anim, ok
}

func (mgr *AssetManager) loadAnimation(
	animName string,
	textName string,
	totalFrames int,
	src image.Rectangle,
	fmrSize Vec2,
	frmOffset Vec2,
	Scale float32,
	Rotation float32,
	Speed float32,
) bool {
	var texture, ok = mgr.getTexture(textName)
	if !ok {
		return false
	}
	var anim = Animation{
		Texture:     texture,
		Src:         src,
		FrmSize:     fmrSize,
		FrmOffset:   frmOffset,
		TotalFrames: totalFrames,
		Scale:       Scale,
		Rotation:    Rotation,
		Speed:       Speed,
	}

	mgr.animations[animName] = &anim

	return true
}

func (mgr *AssetManager) Update() {
	for k := range mgr.animations {
		mgr.animations[k].Update(mgr.tickCount)
	}

	mgr.tickCount += 1
}

// func (mgr *AssetManager) getAnimationId(entity EntityId) int {
// 	return mgr.entityToAnim[entity]
// }

// func (mgr *AssetManager) hasAnimation(entity EntityId) bool {
// 	return mgr.animations[mgr.entityToAnim[entity]] != nil && mgr.entityToAnim[entity] < mgr.animSize
// }
