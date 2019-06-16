package systems

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// Shade シーン切り替え時に画面を覆うタイル
type Shade struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// IntermissionSystem intermisson
type IntermissonSystem struct {
	world        *ecs.World
	playerEntity *Player
	texture      *common.Texture
}

// New 新規作成
func (is *IntermissonSystem) New(w *ecs.World) {
	is.world = w
	camX := camEntity.X()
	camY := camEntity.Y()
	// 画面を黒く覆う
	shadePic, _ := common.LoadedSprite("pics/black_bk.png")
	for i := 0; i < 20; i++ {
		for j := 0; j < 20; j++ {
			shade := &Shade{BasicEntity: ecs.NewBasic()}
			// 描画位置の指定
			shade.SpaceComponent.Position = engo.Point{
				X: float32(i*16*tileMultiply + int(camX)),
				Y: float32(j*16*tileMultiply + int(camY)),
			}
			// 見た目の設定
			shade.RenderComponent = common.RenderComponent{
				Drawable: shadePic,
				Scale:    engo.Point{X: float32(tileMultiply), Y: float32(tileMultiply)},
			}
			shade.RenderComponent.SetZIndex(3)
		}
	}
}

// Remove 削除する
func (is *IntermissonSystem) Remove(entity ecs.BasicEntity) {}

// Update アップデートする
func (is *IntermissonSystem) Update(dt float32) {}