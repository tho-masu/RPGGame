package systems

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"math/rand"
	"os"
	"time"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
)

// Spritesheet タイルの画像
var Spritesheet *common.Spritesheet

// camEntity カメラシステムのエンティティ
var camEntity *common.CameraSystem

// 敵の画像の大きさ
var tileRadius float32 = 7

// ObstaclePoints 障害物のある座標(タイルベース)
var ObstaclePoints map[int][]int

// cellLength セル一辺のピクセル数（必ず16の倍数にすること）
var cellLength = 32

// Tile タイル一つ一つを表す構造体
type Tile struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// SceneSystem シーンシステム
type SceneSystem struct {
	world   *ecs.World
	texture *common.Texture
}

// タイルシステムのエンティティのインスタンス
var tileEntities []*Tile

// Remove 削除する
func (*SceneSystem) Remove(ecs.BasicEntity) {}

// Update アップデートする
func (ss *SceneSystem) Update(dt float32) {
}

// New 作成時に呼び出される
func (ss *SceneSystem) New(w *ecs.World) {
	rand.Seed(time.Now().UnixNano())

	ss.world = w
	// 素材シートの読み込み
	loadTxt := "pics/overworld_tileset_grass.png"
	Spritesheet = common.NewSpritesheetWithBorderFromFile(loadTxt, 16, 16, 0, 0)
	Tiles := make([]*Tile, 0)
	ObstaclePoints = map[int][]int{}
	file, err := os.Open("./stages/main.csv")
	if err != nil {
		fmt.Println(err)
	}
	reader := csv.NewReader(file)
	reader.Comma = ','
	reader.LazyQuotes = true
	i := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		j := 0
		for _, r := range record {
			tileNum, err := strconv.Atoi(r)
			if err != nil {
				fmt.Println("CSVファイルがおかしい")
			}
			if tileNum == 95 {
				// 障害物として、タイルベースで座標を記録（曖昧化のために、前後の複数点を記録）
				ObstaclePoints[i] = append(ObstaclePoints[i], j)
			} // Tileエンティティの作成
			tile := &Tile{BasicEntity: ecs.NewBasic()}
			// 描画位置の指定
			tile.SpaceComponent.Position = engo.Point{
				X: float32(i * cellLength),
				Y: float32(j * cellLength),
			}
			// 見た目の設定
			tile.RenderComponent = common.RenderComponent{
				Drawable: Spritesheet.Cell(tileNum),
				Scale:    engo.Point{X: float32(cellLength / 16), Y: float32(cellLength / 16)},
			}

			tile.RenderComponent.SetZIndex(0)
			Tiles = append(Tiles, tile)

			j++

		}
		i++
	}

	for _, system := range ss.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for _, v := range Tiles {
				tileEntities = append(tileEntities, v)
				sys.Add(&v.BasicEntity, &v.RenderComponent, &v.SpaceComponent)
			}
		}
	}
	// カメラエンティティの取得
	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.CameraSystem:
			camEntity = sys
		}
	}
}
