package main

import (
	_ "image/png"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hugolgst/rich-go/client"
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
	Game_S_Del_check
	Game_S_Del
)

type Game struct {
	Down_X            int     //y,x
	Down_Y            int     //y,x
	Down_Data         [][]int //y,x
	Down_speed_time   float64
	Data              [][]int //y,x
	time_data         float64
	old_time_data     float64
	key_time_data     float64
	key_old_time_data float64
	key_count         int
	Game_S            int
	Can_Del_Line_data []bool
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
	g.key_old_time_data = float64(time.Now().UnixNano() / 1000000)
	g.key_time_data = g.key_old_time_data
	g.Down_speed_time = 500
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
			if g.Down_Data[y][x] == 1 {
				g.Data[g.Down_Y+y][g.Down_X+x] = g.Down_Data[y][x]
			}
		}
	}
}

func (g *Game) Can_Del_Line() ([]bool, bool) {
	Is_all_ok := false
	A := make([]bool, len(g.Data))
	for y := 0; y < len(g.Data); y++ {
		Is_ok := true
		for x := 0; x < len(g.Data[y]); x++ {
			if g.Data[y][x] == 0 {
				Is_ok = false
			}
		}
		Is_all_ok = Is_all_ok || Is_ok
		A[y] = Is_ok
	}
	return A, Is_all_ok
}

func (g *Game) set_line_Data(l, v int) {
	for x := 0; x < len(g.Data[l]); x++ {
		g.Data[l][x] = v
	}
}

func (g *Game) Del_line(lines []bool) {
	s_x := len(g.Data[0])
	New_Data := [][]int{}
	for y := 0; y < len(g.Data); y++ {
		if lines[y] == false {
			New_Data = append(New_Data, g.Data[y])
		}
	}

	for y := 0; y < len(g.Data)-len(New_Data); y++ {
		New_Data = append([][]int{make([]int, s_x)}, New_Data...)
	}
	g.Data = New_Data
}

func (g *Game) Can_Rotate() bool {
	for y := 0; y < len(g.Down_Data); y++ {
		for x := 0; x < len(g.Down_Data[y]); x++ {
			ty := y
			tx := len(g.Down_Data[y]) - 1 - x
			if g.Down_Data[tx][ty] == 1 {
				if 0 <= g.Down_Y+y && g.Down_Y+y < len(g.Data) && 0 <= g.Down_X+x && g.Down_X+x < len(g.Data[0]) {
					if g.Data[g.Down_Y+y][g.Down_X+x] == 1 {
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

func (g *Game) Rotate() {
	temp := make([][]int, len(g.Down_Data))
	for y := 0; y < len(g.Down_Data); y++ {
		temp[y] = make([]int, len(g.Down_Data[y]))
	}
	for y := 0; y < len(g.Down_Data); y++ {
		for x := 0; x < len(g.Down_Data[y]); x++ {
			temp[y][x] = g.Down_Data[y][x]
		}
	}
	for y := 0; y < len(g.Down_Data); y++ {
		for x := 0; x < len(g.Down_Data[y]); x++ {
			ty := y
			tx := len(g.Down_Data[y]) - 1 - x
			g.Down_Data[y][x] = temp[tx][ty]
		}
	}
}

func (g *Game) Update(screen *ebiten.Image) error {
	temp := float64(time.Now().UnixNano() / 1000000)
	g.time_data = temp - g.old_time_data
	g.key_time_data = temp - g.key_old_time_data

	//if g.key_time_data >= 100 && g.Game_S == Game_S_Down {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		if g.Can_Move_Block(1, 0) {
			if g.key_count == 0 || g.key_count >= 10 {
				g.Down_X += 1
				g.key_old_time_data = temp
				g.key_count = 0
			}
			g.key_count += 1
		}
	} else if ebiten.IsKeyPressed((ebiten.KeyLeft)) {
		if g.Can_Move_Block(-1, 0) {
			if g.key_count == 0 || g.key_count >= 10 {
				g.Down_X -= 1
				g.key_old_time_data = temp
				g.key_count = 0
			}
			g.key_count += 1
		}
	} else {
		g.key_count = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		if g.Can_Rotate() {
			g.Rotate()
			g.key_old_time_data = temp
		} else {
		}
	}
	if ebiten.IsKeyPressed((ebiten.KeyDown)) {
		g.Down_speed_time = 50
	} else {
		g.Down_speed_time = 500
	}
	//}

	switch g.Game_S {
	case Game_S_Down:
		if g.time_data >= g.Down_speed_time {
			if g.Can_Move_Block(0, 1) {
				g.Down_Y += 1
				g.old_time_data = temp
			} else {
				g.Set_Move_Block()
				g.Game_S = Game_S_Del_check
				g.old_time_data = temp
			}
		}
	case Game_S_Del_check:
		ok := false
		g.Can_Del_Line_data, ok = g.Can_Del_Line()
		if ok {
			g.Game_S = Game_S_Del
		} else {
			g.Game_S = Game_S_Make
		}
		g.old_time_data = temp
	case Game_S_Del:
		if g.time_data <= 200 {
			for y := 0; y < len(g.Data); y++ {
				if g.Can_Del_Line_data[y] {
					g.set_line_Data(y, 0)
				}
			}
		} else if g.time_data <= 400 {
			for y := 0; y < len(g.Data); y++ {
				if g.Can_Del_Line_data[y] {
					g.set_line_Data(y, 1)
				}
			}
		} else if g.time_data <= 600 {
			for y := 0; y < len(g.Data); y++ {
				if g.Can_Del_Line_data[y] {
					g.set_line_Data(y, 0)
				}
			}
		} else if g.time_data <= 800 {
			for y := 0; y < len(g.Data); y++ {
				if g.Can_Del_Line_data[y] {
					g.set_line_Data(y, 1)
				}
			}
		} else if g.time_data >= 1000 {
			g.Del_line(g.Can_Del_Line_data)
			g.Game_S = Game_S_Make
			g.old_time_data = temp
		}
	case Game_S_Make:
		if g.time_data >= 100 {
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

	if g.Game_S == Game_S_Down {
		for y := 0; y < len(g.Down_Data); y++ {
			for x := 0; x < len(g.Down_Data[y]); x++ {
				if g.Down_Data[y][x] == 1 {
					Block_Draw_image(screen, float64(one_offset_x*(x+g.Down_X)+offset_x), float64(one_offset_y*(y+g.Down_Y)+offset_y))
				}
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 600, 800
}

func main() {
	err := client.Login("885929853793210368")
	if err != nil {
		panic(err)
	}

	Game_Start_time := time.Now()

	err = client.SetActivity(client.Activity{
		State:      "KageKage_Tetris!!!",
		Details:    "I'm playing on KageKage_Tetris!",
		LargeImage: "icon",
		//LargeText:  "This is the large image :D",
		SmallImage: "icon",
		//SmallText:  "And this is the small image",
		Party: &client.Party{
			ID:         "-1",
			Players:    15,
			MaxPlayers: 24,
		},
		Timestamps: &client.Timestamps{
			Start: &Game_Start_time,
		},
	})

	if err != nil {
		panic(err)
	}

	ebiten.SetWindowSize(600, 800)
	ebiten.SetWindowTitle("KageKage_tetris")
	g := &Game{Data: [][]int{}}
	g.Init()
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
