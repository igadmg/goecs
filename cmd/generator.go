package ecs

import (
	"bytes"
	"go/ast"

	"github.com/igadmg/gogen/core"
)

type GeneratorEcs struct {
	core.GeneratorBaseT //[Type, Field, core.Func]

	pwd string

	components map[string]*Type
	entities   map[string]*Type
	queries    map[string]*Type
	features   map[string]*Type

	EntitesByQueries map[*Type][]*Type
}

var _ core.Generator = (*GeneratorEcs)(nil)

func NewGeneratorEcs(pwd string) *GeneratorEcs {
	core.TagNames = []string{"ecs"}

	g := &GeneratorEcs{
		GeneratorBaseT:   core.MakeGeneratorB("0.gen_ecs.go"),
		pwd:              pwd,
		components:       map[string]*Type{},
		entities:         map[string]*Type{},
		queries:          map[string]*Type{},
		features:         map[string]*Type{},
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

		et.Etype = Tag(et.Tag).GetEcsTag()

		switch Tag(et.Tag).GetEcsTag() {
		case EcsEntity:
			g.entities[et.Name] = et
		case EcsFeature:
			g.features[et.Name] = et
		case EcsComponent:
			g.components[et.Name] = et
		case EcsQuery:
			g.queries[et.Name] = et
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

func (g *GeneratorEcs) Generate() bytes.Buffer {
	/*
		EntitesByQueries: func() map[string][]string {
			r := map[string][]string{}

			for _, e := range p_entites {
				ec := map[string]bool{}
				for _, c := range g.gen_allComponents(e) {
					ec[c.Type.GetName()] = true
				}
				for _, q := range p_queries {
					for _, c := range g.gen_allComponents(q) {
						if _, ok := ec[c.Type.GetName()]; !ok {
							goto skip_query
						}
					}

					r[q.GetName()] = append(r[q.GetName()], e.GetName())
				skip_query:
				}
			}

			return r
		}
	*/

	source := bytes.Buffer{}
	pkg := "game"
	g.generate(&source, pkg)

	return source
}

/*
func (g *GeneratorEcs) gen_getEntity(name string) (e core.TypeI, ok bool) {
	e, ok = g.p_entites_by_name[name]
	return
}

func (g *GeneratorEcs) gen_getQuery(name string) (e core.TypeI, ok bool) {
	e, ok = g.p_queries_by_name[name]
	return
}

func (g *GeneratorEcs) gen_allComponents(e core.TypeI) (r []*Field) {
	t, ok := CastType(e)
	if !ok {
		return
	}

	for base := range EnumFields(t.Bases) {
		r = append(r, g.gen_allComponents(base.Type)...)
	}

	lcm := map[string]*Field{}
	lc := g.gen_Components(e)
	for _, c := range lc {
		lcm[c.Name] = c
	}
	r = slices.Collect(xiter.Filter(slices.Values(r),
		func(i *Field) bool {
			_, ok := lcm[i.Name]
			return !(ok && (i.Tag.HasField(Tag_Virtual) || i.Tag.HasField(Tag_Abstract)))
		}))
	r = append(r, lc...)

	return r
}

func (g *GeneratorEcs) gen_Components(e core.TypeI) (r []*Field) {
	t, ok := CastType(e)
	if !ok {
		return
	}

	for _, c := range t.Fields {
		switch field := c.(type) {
		case *Field:
			if field.Tag.HasField(Tag_Abstract) {
				continue
			}

			r = append(r, field)
		}
	}

	return r
}

func (g *GeneratorEcs) gen_allOverrides(e *Type) []cgo_gen_h_p_o {
	var r []cgo_gen_h_p_o
	for base := range EnumFields(e.Bases) {
		be, ok := CastType(base.Type)
		if !ok {
			continue
		}

		for field := range EnumFields(be.Fields) {
			if field.Tag.HasField(Tag_Virtual) || field.Tag.HasField(Tag_Abstract) {
				r = append(r, cgo_gen_h_p_o{
					Base: be.GetName(),
					Name: field.Name,
					Type: field.Type,
				})
			}
		}
	}

	return r
}
*/
