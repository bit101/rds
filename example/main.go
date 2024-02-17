// Package main renders an image, gif or video
package main

import (
	"math"

	"github.com/bit101/bitlib/blmath"
	cairo "github.com/bit101/blcairo"
	"github.com/bit101/blcairo/render"
	"github.com/bit101/blcairo/target"
)

var (
	grid *rds.grid
	size = 400
)

func main() {
	sizef := float64(size)
	renderTarget := target.Image

	switch renderTarget {
	case target.Image:
		render.Image(sizef, sizef, "out/out.png", scene1, 0.0)
		render.ViewImage("out/out.png")
		break

	case target.Video:
		program := render.NewProgram(sizef, sizef, 30)
		// program.AddSceneWithFrames(scene1, 200)
		program.AddSceneWithFrames(scene3, 500)
		// program.AddSceneWithFrames(fadeOut, 60)
		program.RenderVideo("out/frames", "out/out.mp4")
		render.PlayVideo("out/out.mp4")
		break
	}
}

func init() {
	grid = rds.NewCellGrid(size, size)
	// grid.Feed = 0.0367
	// grid.Kill = 0.0649
	grid.Feed = 0.0545
	grid.Kill = 0.062
	// grid.Feed = 0.024
	// grid.Kill = 0.060
	// for i := 0; i < 100; i++ {
	// 	x := random.IntRange(0, size)
	// 	y := random.IntRange(0, size)
	// 	grid.SetCellB(x, y, 1)
	// }
	// for x := 50; x < 150; x++ {
	// 	for y := 120; y < 140; y++ {
	// 		grid.SetCellB(x, y, 1)
	// 	}
	// }
	// for x := 70; x < 130; x++ {
	// 	for y := 20; y < 40; y++ {
	// 		grid.SetCellB(x, y, 1)
	// 	}
	// }
	for x := 80; x < 120; x++ {
		for y := 80; y < 120; y++ {
			grid.SetCellB(x+100, y+100, 1)
		}
	}
}

//revive:disable-next-line:unused-parameter
func scene1(context *cairo.Context, width, height, percent float64) {
	grid.Update(50, true)
	context.Surface.SetData(grid.GetImageDataA())
}

//revive:disable-next-line:unused-parameter
func scene3(context *cairo.Context, width, height, percent float64) {
	// grid.Feed = 0.0545
	// grid.Kill = 0.062
	grid.UpdateAdvanced(
		50,
		func(x, y int) float64 {
			// return grid.Feed
			t := math.Sin(float64(x) * 0.2)
			return blmath.Map(t, -1, 1, 0.0543, 0.0547)
			// t := float64(x) / width
			// return blmath.Lerp(t, 0.03, 0.07)
		},
		func(x, y int) float64 {
			// return grid.Kill
			// t := math.Sin(float64(y) * 0.2)
			// return blmath.Map(t, -1, 1, 0.061, 0.063)
			t := float64(y) / height
			if percent > 0.5 {
				t = 1 - t
			}
			return blmath.Lerp(t, 0.060, 0.065)
		},
		func(x, y int) float64 {
			return grid.DiffuseA
			// t := math.Sin(float64(x) * 0.2)
			// return blmath.Map(t, -1, 1, 0.75, 1.1)

		},
		func(x, y int) float64 {
			return grid.DiffuseB
			// t := float64(x-y) / (height - width)
			// return blmath.Lerp(t, 0.125, 0.6)
		},
		true,
	)

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
