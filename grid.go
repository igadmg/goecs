package ecs

import "github.com/igadmg/raylib-go/raymath/vector2"

type CellGrid[T any] struct {
	size  vector2.Int
	Cells []T
}

func MakeCellGrid[T any](size vector2.Int) CellGrid[T] {
	return CellGrid[T]{
		size:  size,
		Cells: make([]T, size.Product()),
	}
}

func (g CellGrid[T]) Size() vector2.Int {
	return g.size.SubXY(2, 2)
}

func (g CellGrid[T]) IsValid(xy vector2.Int) bool {
	return vector2.GreaterEq(xy, vector2.Zero[int]()) && vector2.Less(xy, g.Size())
}

func (g CellGrid[T]) XY(i int) vector2.Int {
	return vector2.New(
		i%g.size.X(),
		i/g.size.X(),
	)
}

func (g CellGrid[T]) Tile(xy vector2.Int) T {
	return g.Cells[g.size.X()*(xy.Y()+1)+(xy.X()+1)]
}

func (g CellGrid[T]) SetTile(xy vector2.Int, v T) {
	g.Cells[g.size.X()*(xy.Y()+1)+(xy.X()+1)] = v
}

func (g CellGrid[T]) TileClamped(xy vector2.Int) T {
	xy = xy.AddXY(1, 1).Clamp0V(g.size.AddXY(-1, -1))
	return g.Cells[g.size.X()*(xy.Y())+(xy.X())]
}

func (g *CellGrid[T]) CellGridRegion(xy vector2.Int) CellGridRegion[T] {
	return CellGridRegion[T]{
		XY:   xy,
		grid: g,
	}
}

type CellGridRegion[T any] struct {
	XY   vector2.Int
	grid *CellGrid[T]
}

func MakeCellGridRegion[T any](grid *CellGrid[T], xy vector2.Int) CellGridRegion[T] {
	return CellGridRegion[T]{
		XY:   xy,
		grid: grid,
	}
}

func (r CellGridRegion[T]) Center() T {
	if !r.grid.IsValid(r.XY) {
		return r.grid.TileClamped(r.XY)
	}
	return r.grid.Tile(r.XY)
}

func (r CellGridRegion[T]) Tile(xy vector2.Int) T {
	gxy := r.XY.Add(xy)
	if !r.grid.IsValid(r.XY) {
		return r.grid.TileClamped(gxy)
	}
	return r.grid.Tile(gxy)
}

func (r CellGridRegion[T]) TileXY(x, y int) T {
	return r.Tile(vector2.New(x, y))
}

func (r CellGridRegion[T]) SetTile(xy vector2.Int, v T) {
	if !r.grid.IsValid(r.XY) {
		return
	}
	gxy := r.XY.Add(xy)
	r.grid.SetTile(gxy, v)
}
