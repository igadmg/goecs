package ecs

import "slices"

type WorldI interface {
	AddSystem(system System)
	RemoveSystem(system System)
	Update(dt float32)
}

type World struct {
	Systems []System
}

var _ WorldI = (*World)(nil)

func (w *World) AddSystem(system System) {
	w.Systems = append(w.Systems, system)
}

func (w *World) RemoveSystem(system System) {
	si := slices.Index(w.Systems, system)
	if si < 0 {
		return
	}

	w.Systems = slices.Delete(w.Systems, si, si+1)
}

func (w *World) Update(dt float32) {
	for _, s := range w.Systems {
		s.Update(dt)
	}
}
