// Package rds is for reaction diffusion simulations.
// While this works, it is really quite slow. This would be better done using shaders.
package rds

import (
	"fmt"

	"github.com/bit101/bitlib/blmath"
)

// PositionalValueFunc is the signature of a function used to get the feed, kill, or diffusion at a given position.
type PositionalValueFunc func(int, int) float64

// Cell holds one cell of a grid.
type Cell struct {
	a float64
	b float64
}

// NewCell creates a new cell.
func NewCell(a, b float64) *Cell {
	return &Cell{a, b}
}

// Grid is a grid of cells.
type Grid struct {
	grid     []*Cell
	buffer   []*Cell
	width    int
	height   int
	DiffuseA float64
	DiffuseB float64
	Feed     float64
	Kill     float64
}

// NewCellGrid creates a new cell grid.
func NewCellGrid(width, height int) *Grid {
	grid := &Grid{}
	grid.grid = make([]*Cell, width*height)
	grid.buffer = make([]*Cell, width*height)
	grid.width = width
	grid.height = height
	grid.DiffuseA = 1.0
	grid.DiffuseB = 0.5
	grid.Feed = 0.0545
	grid.Kill = 0.062
	for i := 0; i < len(grid.grid); i++ {
		grid.grid[i] = NewCell(1, 0)
		grid.buffer[i] = NewCell(1, 0)
	}
	return grid
}

// SetCellA sets the a of a given cell.
func (g *Grid) SetCellA(x, y int, a float64) {
	index := y*g.width + x
	g.grid[index].a = a
}

// SetCellB sets the b of a given cell.
func (g *Grid) SetCellB(x, y int, b float64) {
	index := y*g.width + x
	g.grid[index].b = b
}

// setBufferA sets the a of a given cell.
func (g *Grid) setBufferA(x, y int, a float64) {
	index := y*g.width + x
	g.buffer[index].a = a
}

// setBufferB sets the b of a given cell.
func (g *Grid) setBufferB(x, y int, b float64) {
	index := y*g.width + x
	g.buffer[index].b = b
}

// getCellA gets the a of a given cell.
func (g *Grid) getCellA(x, y int) float64 {
	index := y*g.width + x
	return g.grid[index].a
}

// getCellB gets the b of a given cell.
func (g *Grid) getCellB(x, y int) float64 {
	index := y*g.width + x
	return g.grid[index].b
}

// Update updates the grid the specified number of iterations.
// A count value of around 50 produces a decent change between frames in an animation.
// If making a still image, setting this value high is more efficient than calling Update multiple times.
// If feedback is true, the current iteration will be output to stdout.
func (g *Grid) Update(count int, feedback bool) {
	f := g.Feed
	kf := g.Kill + f
	da := g.DiffuseA
	db := g.DiffuseB
	for i := 0; i < count; i++ {
		if feedback {
			fmt.Printf("\rIteration: %d/%d", i+1, count)
		}
		for x := 1; x < g.width-1; x++ {
			for y := 1; y < g.height-1; y++ {
				index := y*g.width + x
				a := g.grid[index].a
				b := g.grid[index].b
				abb := a * b * b
				lapA := g.laplace(x, y, g.getCellA)
				lapB := g.laplace(x, y, g.getCellB)
				a1 := a + da*lapA - abb + f*(1-a)
				b1 := b + db*lapB + abb - kf*b
				g.buffer[index].a = blmath.Clamp(a1, 0, 1)
				g.buffer[index].b = blmath.Clamp(b1, 0, 1)
			}
		}
		g.buffer, g.grid = g.grid, g.buffer
	}
	if feedback {
		fmt.Println()
	}
}

// UpdateAdvanced updates the grid the specified number of iterations.
// The advanced version lets you supply PositionalValueFuncs to be used to get
// the values for feed, kill, diffuse A and B.
// A count value of around 50 produces a decent change between frames in an animation.
// If making a still image, setting this value high is more efficient than calling Update multiple times.
// If feedback is true, the current iteration will be output to stdout.
func (g *Grid) UpdateAdvanced(count int, getFeed, getKill, getDiffuseA, getDiffuseB PositionalValueFunc, feedback bool) {
	for i := 0; i < count; i++ {
		if feedback {
			fmt.Printf("\rIteration: %d/%d", i+1, count)
		}
		for x := 1; x < g.width-1; x++ {
			for y := 1; y < g.height-1; y++ {
				f := getFeed(x, y)
				kf := getKill(x, y) + f
				da := getDiffuseA(x, y)
				db := getDiffuseB(x, y)
				index := y*g.width + x
				a := g.grid[index].a
				b := g.grid[index].b
				abb := a * b * b
				lapA := g.laplace(x, y, g.getCellA)
				lapB := g.laplace(x, y, g.getCellB)
				a1 := a + da*lapA - abb + f*(1-a)
				b1 := b + db*lapB + abb - kf*b
				g.buffer[index].a = blmath.Clamp(a1, 0, 1)
				g.buffer[index].b = blmath.Clamp(b1, 0, 1)
			}
		}
		g.buffer, g.grid = g.grid, g.buffer
	}
	if feedback {
		fmt.Println()
	}
}

func (g *Grid) laplace(x, y int, getter func(int, int) float64) float64 {
	total := 0.0
	counter := 0
	mults := []float64{
		0.05, 0.2, 0.05,
		0.2, -1.0, 0.2,
		0.05, 0.2, 0.05,
	}
	for xx := x - 1; xx <= x+1; xx++ {
		for yy := y - 1; yy <= y+1; yy++ {
			total += mults[counter] * getter(xx, yy)
			counter++
		}
	}
	return total
}

// GetImageDataA returns a byte array of image data
func (g *Grid) GetImageDataA() []byte {
	data := make([]byte, g.width*g.height*4)
	i := 0
	for _, cell := range g.grid {
		g := byte(cell.a * 255)
		data[i] = g
		data[i+1] = g
		data[i+2] = g
		data[i+3] = 255
		i += 4
	}
	return data
}

// GetImageDataB returns a byte array of image data
func (g *Grid) GetImageDataB() []byte {
	data := make([]byte, g.width*g.height*4)
	i := 0
	for _, cell := range g.grid {
		g := byte(cell.b * 255)
		data[i] = g
		data[i+1] = g
		data[i+2] = g
		data[i+3] = 255
		i += 4
	}
	return data
}
