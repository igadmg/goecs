package cmd

import (
	"bytes"
	"go/ast"

	"github.com/igadmg/gogen/core"
)

type GeneratorEcs struct {
	core.GeneratorBaseT

	components map[string]*Type
	entities   map[string]*Type
	queries    map[string]*Type
	features   map[string]*Type
	systems    map[string]*Type

	EntitesByQueries map[*Type][]*Type
}

var _ core.Generator = (*GeneratorEcs)(nil)

func NewGeneratorEcs() core.Generator {
	g := &GeneratorEcs{
		GeneratorBaseT:   core.MakeGeneratorB("ecs", "0.gen_ecs.go", "ecs"),
		components:       map[string]*Type{},
		entities:         map[string]*Type{},
		queries:          map[string]*Type{},
		features:         map[string]*Type{},
		systems:          map[string]*Type{},
		EntitesByQueries: map[*Type][]*Type{},
	}
	g.G = g
	return g
}

func (g *GeneratorEcs) NewType(t core.TypeI, spec *ast.TypeSpec) (core.TypeI, error) {
	if t == nil {
		t = NewType()
		defer func() {
			g.Types[t.GetName()] = t
		}()
	}

	switch et := t.(type) {
	case *Type:
		var err error
		_, err = g.GeneratorBaseT.NewType(&et.Type, spec)
		if err != nil {
			return nil, err
		}

		et.EType = Tag(et.Tag).GetEcsTag()

		switch Tag(et.Tag).GetEcsTag() {
		case EcsEntity:
			g.entities[et.Name] = et
		case EcsFeature:
			g.features[et.Name] = et
		case EcsComponent:
			g.components[et.Name] = et
		case EcsQuery:
			g.queries[et.Name] = et
		case EcsSystem:
			g.systems[et.Name] = et
		}
	}

	return t, nil
}

func (g *GeneratorEcs) NewField(f core.FieldI, spec *ast.Field) (core.FieldI, error) {
	if f == nil {
		f = &Field{}
		defer func() {
			g.Fields = append(g.Fields, f)
		}()
	}

	switch ef := f.(type) {
	case *Field:
		var err error

		_, err = g.GeneratorBaseT.NewField(&ef.Field, spec)
		if err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (g *GeneratorEcs) NewFunc(f core.FuncI, spec *ast.FuncDecl) (core.FuncI, error) {
	if f == nil {
		f = &core.Func{}
		defer func() {
			if id := f.GetFullTypeName(); id != "" {
				g.Funcs[id] = append(g.Funcs[id], f)
			}
		}()
	}

	switch ef := f.(type) {
	case *core.Func:
		var err error

		_, err = g.GeneratorBaseT.NewFunc(ef, spec)
		if err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (g *GeneratorEcs) GetEcsType(name string) (t EcsTypeI, ok bool) {
	if t, ok := g.GetType(name); ok {
		et, ok := t.(EcsTypeI)
		return et, ok
	}

	return nil, ok
}

func (g *GeneratorEcs) Prepare() {
	g.GeneratorBaseT.Prepare()

	for _, q := range g.queries {
		g.EntitesByQueries[q] = []*Type{}
	}
	for _, e := range g.entities {
		ec := map[*Type]struct{}{}
		for c := range EnumFieldsSeq(e.StructComponentsSeq()) {
			ct, ok := CastType(c.Type)
			if !ok {
				continue
			}

			ec[ct] = struct{}{}
		}
		for _, q := range g.queries {
			for c := range EnumFieldsSeq(q.StructComponentsSeq()) {
				ct, ok := CastType(c.Type)
				if !ok {
					continue
				}

				if _, ok := ec[ct]; !ok {
					goto skip_query
				}
			}

			g.EntitesByQueries[q] = append(g.EntitesByQueries[q], e)
		skip_query:
		}
	}
}

func (g *GeneratorEcs) Generate(pkg string) bytes.Buffer {
	source := bytes.Buffer{}
	g.generate(&source, pkg)
	return source
}
