package main

import (
	"image"
	"image/color"
	"image/gif"
	"log"
	"os"
)

// Open Questions:
// How big is a tile? 1000ft?
// Are there Z-levels?

type TileType int8

const (
	Unknown TileType = iota
	Dirt
	Rock
	Forest
	Gold //4
	Iron
	FreshWater
	SaltWater
	Road
	Farm
	Town
	City
	Industry
	_
	_
	_ // 15
)

type Tile struct {
	Types    []TileType
	Percents []int
}

type Board struct {
	T      []TileType
	Height int
	Width  int
}

func (b *Board) Set(x, y int, t TileType) {
	b.T[x+y*b.Width] = t
}

func NewBoard(height, width int) *Board {
	b := &Board{
		T:      make([]TileType, height*width),
		Height: height,
		Width:  width,
	}
	for i := range b.T {
		b.T[i] = Forest
	}
	return b
}

const NumPoints = 400

func main() {
	tmpfile, err := os.Create("output.gif") //ioutil.TempFile("", "worldsim")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Starting!, Output at %s", tmpfile.Name())

	frameCount := 10
	b := NewBoard(150, 150)
	points := genPoints(NumPoints, b.Height, b.Width)
	for _, p := range points {
		b.Set(int(p.X()), int(p.Y()), Gold)
	}
	log.Print(points)
	buildGraph(points)
	gif.EncodeAll(tmpfile, ToGif(b, UpdateBoard, frameCount))
	log.Print("Done!")
}

func UpdateBoard(b *Board, frame int) {
	//log.Print(b.T)
	b.T[frame*(b.Width+1)] = Dirt
}

const FrameTime = 30 // hundredths of a second

func ToGif(b *Board, updateState func(*Board, int), n int) *gif.GIF {
	g := &gif.GIF{
		Image: []*image.Paletted{ToImage(b)},
		Delay: []int{FrameTime},
	}
	for frame := 0; frame < n; frame++ {
		updateState(b, frame)
		g.Image = append(g.Image, ToImage(b))
		g.Delay = append(g.Delay, FrameTime)
	}

	g.Image = append(g.Image, ToImage(b))
	g.Delay = append(g.Delay, 300)
	return g
}

func toUint8(i []TileType) []uint8 {
	o := make([]uint8, len(i))
	for idx := range i {
		o[idx] = uint8(i[idx])
	}
	return o
}

func ToImage(b *Board) *image.Paletted {

	return &image.Paletted{
		Pix:     toUint8(b.T),
		Rect:    image.Rect(0, 0, b.Width, b.Height),
		Stride:  b.Width,
		Palette: Palette,
	}
}

var Palette = color.Palette{
	color.RGBA{0x00, 0xff, 0x00, 0xff}, // Unknown
	color.RGBA{0xa0, 0x52, 0x2d, 0xff}, // Dirt
	color.RGBA{0x69, 0x69, 0x69, 0xff}, // Rock
	color.RGBA{0x00, 0x64, 0x00, 0xff}, // Forest
	color.RGBA{0xff, 0xd7, 0x00, 0xff}, // Gold
	color.RGBA{0x2f, 0x4f, 0x4f, 0xff}, // Iron
	color.RGBA{0x41, 0x69, 0xe1, 0xff}, // FreshWater
	color.RGBA{0x87, 0xce, 0xeb, 0xff}, // SaltWater
	color.RGBA{0x11, 0x11, 0x11, 0xff}, // Road
	color.RGBA{0x80, 0x80, 0x00, 0xff}, // Farm
	color.RGBA{0xdd, 0xa0, 0xdd, 0xff}, // Town
	color.RGBA{0xba, 0x55, 0xd3, 0xff}, // City
	color.RGBA{0xff, 0xff, 0x00, 0xff}, // Industry
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
	color.RGBA{0x00, 0x00, 0x00, 0xff},
}
