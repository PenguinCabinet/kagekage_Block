package main

import (
	"fmt"
	_ "image/png"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

var block_img *ebiten.Image
var background_img *ebiten.Image

func init() {
	var err error
	block_img, _, err = ebitenutil.NewImageFromFile("images/block.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	background_img, _, err = ebitenutil.NewImageFromFile("images/frame.png", ebiten.FilterDefault)
}

func Draw_image(screen, img *ebiten.Image, x, y, s1, s2 float64) {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Scale(s1, s2)

	op.GeoM.Translate(x, y)

	screen.DrawImage(img, op)
}

func Block_Draw_image(screen *ebiten.Image, x, y float64) {
	Draw_image(screen, block_img, x, y, 2.8, 1.8)
}

const (
	Game_S_Make = iota
	Game_S_Down
	Game_S_Down_End
)

type Game struct {
	Down_X        int     //y,x
	Down_Y        int     //y,x
	Down_Data     [][]int //y,x
	Data          [][]int //y,x
	time_data     float64
	old_time_data float64
	Game_S        int
}

func (g *Game) Init() error {
	g.Data = make([][]int, 30)
	for y := 0; y < len(g.Data); y++ {
		g.Data[y] = make([]int, 10)
		/*
			for x := 0; x < len(g.Data[y]); x++ {
				g.Data[y][x] = 1
			}
		*/
	}
	g.old_time_data = float64(time.Now().UnixNano() / 1000000)
	g.time_data = g.old_time_data
	return nil
}

func (g *Game) Make_Block(screen *ebiten.Image) error {
	pattern_data := [][][]int{
		{
			{0, 1, 0},
			{0, 1, 0},
			{1, 1, 1},
		},
	}
	g.Down_Y = 0
	g.Down_X = 4
	g.Down_Data = make([][]int, 3)
	for y := 0; y < len(pattern_data[0]); y++ {
		g.Down_Data[y] = make([]int, 3)
		for x := 0; x < len(pattern_data[0][y]); x++ {
			g.Down_Data[y][x] = pattern_data[0][y][x]
		}
	}
	return nil
}

func (g *Game) Can_Move_Block(dx, dy int) bool {
	for y := 0; y < len(g.Down_Data); y++ {
		for x := 0; x < len(g.Down_Data[y]); x++ {
			if g.Down_Data[y][x] == 1 {
				if 0 <= g.Down_Y+dy+y && g.Down_Y+dy+y < len(g.Data) && 0 <= g.Down_X+dx+x && g.Down_X+dx+x < len(g.Data[0]) {
					if g.Data[g.Down_Y+dy+y][g.Down_X+dx+x] == 1 {
						return false
					}
				} else {
					return false
				}
			}
		}
	}
	return true
}

func (g *Game) Set_Move_Block() {
	for y := 0; y < len(g.Down_Data); y++ {
		for x := 0; x < len(g.Down_Data[y]); x++ {
			g.Data[g.Down_Y+y][g.Down_X+x] = g.Down_Data[y][x]
		}
	}
}

func (g *Game) Update(screen *ebiten.Image) error {
	temp := float64(time.Now().UnixNano() / 1000000)
	g.time_data = temp - g.old_time_data

	if inpututil.IsKeyJustPressed((ebiten.KeyRight)) {
		if g.Can_Move_Block(1, 0) {
			g.Down_X += 1
		}
	}
	if inpututil.IsKeyJustPressed((ebiten.KeyLeft)) {
		if g.Can_Move_Block(-1, 0) {
			g.Down_X -= 1
		}
	}

	switch g.Game_S {
	case Game_S_Down:
		if g.time_data >= 100 {
			if g.Can_Move_Block(0, 1) {
				g.Down_Y += 1
				g.old_time_data = temp
			} else {
				g.Set_Move_Block()
				g.Game_S = Game_S_Make
				g.old_time_data = temp
			}
		}
	case Game_S_Make:
		if g.time_data >= 1000 {
			fmt.Println("TRUE2")
			g.Make_Block(screen)
			g.old_time_data = temp
			g.Game_S = Game_S_Down
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	Draw_image(screen, background_img, 0, 0, 1, 1)
	offset_x := 155
	offset_y := 110
	one_offset_x := 29
	one_offset_y := 20
	for y := 0; y < len(g.Data); y++ {
		for x := 0; x < len(g.Data[y]); x++ {
			if g.Data[y][x] == 1 {
				//fmt.Println("TRUE")
				Block_Draw_image(screen, float64(one_offset_x*x+offset_x), float64(one_offset_y*y+offset_y))
			}
		}
	}

	for y := 0; y < len(g.Down_Data); y++ {
		for x := 0; x < len(g.Down_Data[y]); x++ {
			if g.Down_Data[y][x] == 1 {
				Block_Draw_image(screen, float64(one_offset_x*(x+g.Down_X)+offset_x), float64(one_offset_y*(y+g.Down_Y)+offset_y))
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 600, 800
}

func main() {
	ebiten.SetWindowSize(600, 800)
	ebiten.SetWindowTitle("KageKage_tetris")
	g := &Game{Data: [][]int{}}
	g.Init()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
