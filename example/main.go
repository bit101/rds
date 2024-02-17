// Package main renders an image, gif or video
package main

import (
	cairo "github.com/bit101/blcairo"
	"github.com/bit101/blcairo/render"
	"github.com/bit101/blcairo/target"
	"github.com/bit101/rds"
)

var (
	grid *rds.Grid
	size = 200
)

func main() {
	sizef := float64(size)
	renderTarget := target.Video

	switch renderTarget {
	case target.Image:
		render.Image(sizef, sizef, "out/out.png", scene1, 0.0)
		render.ViewImage("out/out.png")
		break

	case target.Video:
		program := render.NewProgram(sizef, sizef, 30)
		program.AddSceneWithFrames(scene1, 200)
		program.AddSceneWithFrames(fadeOut, 60)
		program.RenderVideo("out/frames", "out/out.mp4")
		render.PlayVideo("out/out.mp4")
		break
	}
}

func init() {
	grid = rds.NewCellGrid(size, size)
	grid.Feed = 0.0545
	grid.Kill = 0.062
	for x := 80; x < 120; x++ {
		for y := 80; y < 120; y++ {
			grid.SetCellB(x, y, 1)
		}
	}
}

//revive:disable-next-line:unused-parameter
func scene1(context *cairo.Context, width, height, percent float64) {
	grid.Update(50, false)
	context.Surface.SetData(grid.GetImageDataA())
}

func fadeOut(context *cairo.Context, width, height, percent float64) {
	data := grid.GetImageDataA()
	for i := 0; i < len(data); i += 4 {
		b := float64(data[i])
		b += (255 - b) * percent
		data[i] = byte(b)
		data[i+1] = byte(b)
		data[i+2] = byte(b)

	}
	context.Surface.SetData(data)
}
