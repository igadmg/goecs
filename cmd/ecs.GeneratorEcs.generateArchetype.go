package cmd

import (
	"fmt"
	"io"
)

func (g *GeneratorEcs) generateArchetype(wr io.Writer, id int, e *Type) {
	g.genAs(wr, e)

	wr.Write([]byte(`

func _`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`_constraints() {
	var _ ecs.Id = `))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`{}.Id
}

type storage_`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(` struct {
	ecs.BaseStorage

`))

	for c := range EnumFieldsSeq(e.StructComponentsSeq()) {

		wr.Write([]byte(`	s_`))
		wr.Write([]byte(fmt.Sprintf("%v", c.Name)))
		wr.Write([]byte(` []`))
		wr.Write([]byte(fmt.Sprintf("%v", c.GetTypeName())))
		wr.Write([]byte(`
`))

	}

	wr.Write([]byte(`}

var s_`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(` = storage_`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`{
	BaseStorage: ecs.MakeBaseStorage(`))
	wr.Write([]byte(fmt.Sprintf("%v", id)))
	wr.Write([]byte(`),
}

func Match`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`(id ecs.Id) (ecs.Ref[`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`], bool) {
	if id.GetType() == s_`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`.TypeId() {
		ref := ecs.Ref[`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`]{Id: id}
		_ = ref.Get()

		return ref, true
	}
`))

	for s := range EnumTypes(e.Subclasses) {

		wr.Write([]byte(`	if id.GetType() == s_`))
		wr.Write([]byte(fmt.Sprintf("%v", s.Name)))
		wr.Write([]byte(`.TypeId() {
		ref := ecs.Ref[`))
		wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
		wr.Write([]byte(`]{Id: id}
		_ = ref.Get()

		return ref, true
	}
`))

	}

	wr.Write([]byte(`
	return ecs.Ref[`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`]{}, false
}

func (e `))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`) Ref() ecs.Ref[`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`] {
	return ecs.Ref[`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`] {
		Id: e.Id,
		Age: s_`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`.Age(),
		Ptr: e,
	}
}

func (e *`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`) Allocate() ecs.Ref[`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`] {
	s := &s_`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`
	age, id := s.BaseStorage.AllocateId()
	index := (int)(id.GetId() - 1)
	_ = index

`))

	for c := range EnumFieldsSeq(e.StructComponentsSeq()) {

		wr.Write([]byte(`	s.s_`))
		wr.Write([]byte(fmt.Sprintf("%v", c.Name)))
		wr.Write([]byte(` = slicesex.Reserve(s.s_`))
		wr.Write([]byte(fmt.Sprintf("%v", c.Name)))
		wr.Write([]byte(`, index+1)
`))

	}

	wr.Write([]byte(`
	ref := ecs.Ref[`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`]{
		Age: age - 1,
		Id:  id,
	}
	_ = ref.Get()

	if e != nil {
`))

	for c := range EnumFieldsSeq(e.StructComponentsSeq()) {
		if ct, ok := CastType(c.Type); ok {
			if ct.IsTransient() {
				continue
			}

			wr.Write([]byte(`		if e.`))
			wr.Write([]byte(fmt.Sprintf("%v", c.Name)))
			wr.Write([]byte(` != nil {
			*ref.Ptr.`))
			wr.Write([]byte(fmt.Sprintf("%v", c.Name)))
			wr.Write([]byte(` = *e.`))
			wr.Write([]byte(fmt.Sprintf("%v", c.Name)))
			wr.Write([]byte(`
		}
`))

		}
	}

	wr.Write([]byte(`		*e = ref.Ptr
	}

	return ref
}

func (e *`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`) Free() {
	Free`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`(e.Id)
}

`))

	g.fnLoad(wr, e)
	if !e.IsTransient() {
		g.fnStore(wr, e)
		g.fnRestore(wr, e)
	}

	wr.Write([]byte(`
func Allocate`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`() (ref ecs.Ref[`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`], entity `))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`) {
	var e *`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(` = nil
	ref = e.Allocate()
	return ref, ref.Ptr
}

func Free`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`(id ecs.Id) {
	s := &s_`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`
	index := (int)(id.GetId() - 1)
	_ = index

`))

	for c := range EnumFieldsSeq(e.StructComponentsSeq()) {

		wr.Write([]byte(`	s.s_`))
		wr.Write([]byte(fmt.Sprintf("%v", c.Name)))
		wr.Write([]byte(`[index] = `))
		wr.Write([]byte(fmt.Sprintf("%v", c.GetTypeName())))
		wr.Write([]byte(`{}
`))

	}

	wr.Write([]byte(`
	s.Free(id)
}

func Update`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`Id(id ecs.Id) {
	tid := id.GetType()
	if s := s_`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`; s.TypeId() == tid {
		index := (int)(id.GetId() - 1)

		s_`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`.Ids[index] = id
	}
}
`))

	if _, ok := g.queries[e.Name+"Query"]; !ok {

		wr.Write([]byte(`
// Auto-generated query for `))
		wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
		wr.Write([]byte(` entity
type `))
		wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
		wr.Write([]byte(`Query struct {
	_ ecs.MetaTag ` + "`" + `ecs:"query: {`))
		wr.Write([]byte(fmt.Sprintf("%v", e.QueryTags)))
		wr.Write([]byte(`}"` + "`" + `

	Id ecs.Id
`))

		for c := range EnumFieldsSeq(e.QueryComponentsSeq()) {

			wr.Write([]byte(`	`))
			wr.Write([]byte(fmt.Sprintf("%v", c.Name)))
			wr.Write([]byte(` *`))
			wr.Write([]byte(fmt.Sprintf("%v", c.GetTypeName())))
			wr.Write([]byte(`
`))

		}

		wr.Write([]byte(`}
`))

	}
}

func (g *GeneratorEcs) genFieldEcsCall(wr io.Writer, f *Field, call string) {
	if f.IsArray {
		if f.isEcsRef {

			wr.Write([]byte(`	for i := range e.`))
			wr.Write([]byte(fmt.Sprintf("%v", f.Name)))
			wr.Write([]byte(` {
		`))
			wr.Write([]byte(fmt.Sprintf("%v", call)))
			wr.Write([]byte(`(&e.`))
			wr.Write([]byte(fmt.Sprintf("%v", f.Name)))
			wr.Write([]byte(`[i])
}
`))

		} else if f.Type.CanCall(call) {

			wr.Write([]byte(`	for i := range e.`))
			wr.Write([]byte(fmt.Sprintf("%v", f.Name)))
			wr.Write([]byte(` {
		e.`))
			wr.Write([]byte(fmt.Sprintf("%v", f.Name)))
			wr.Write([]byte(`[i].`))
			wr.Write([]byte(fmt.Sprintf("%v", call)))
			wr.Write([]byte(`()
}
`))

		}
	} else {
		if f.isEcsRef {

			wr.Write([]byte(`	`))
			wr.Write([]byte(fmt.Sprintf("%v", call)))
			wr.Write([]byte(`(&e.`))
			wr.Write([]byte(fmt.Sprintf("%v", f.Name)))
			wr.Write([]byte(`)
`))

		} else if f.Type.CanCall(call) {

			wr.Write([]byte(`	e.`))
			wr.Write([]byte(fmt.Sprintf("%v", f.Name)))
			wr.Write([]byte(`.`))
			wr.Write([]byte(fmt.Sprintf("%v", call)))
			wr.Write([]byte(`()
`))

		}
	}
}

func (g *GeneratorEcs) fnLoad(wr io.Writer, e *Type) {
	if e.HasFunction("Load") {
		return
	}

	wr.Write([]byte(`func (e `))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`) Load(age uint64, id ecs.Id) (uint64, `))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`) {
	index := (int)(id.GetId() - 1)
	tid := id.GetType()
	_ = index

`))

	for _, s := range e.Subclasses {
		switch sc := s.(type) {
		case *Type:

			wr.Write([]byte(`	if s := &s_`))
			wr.Write([]byte(fmt.Sprintf("%v", sc.Name)))
			wr.Write([]byte(`; s.TypeId() == tid {
		if age != s.Age() {
			e.Id = id
`))

			for field := range EnumFields(e.Fields) {
				if field.Tag.HasField(Tag_Virtual) || field.Tag.HasField(Tag_Abstract) {

					wr.Write([]byte(`			e.`))
					wr.Write([]byte(fmt.Sprintf("%v", field.Name)))
					wr.Write([]byte(` = &s.s_`))
					wr.Write([]byte(fmt.Sprintf("%v", field.Name)))
					wr.Write([]byte(`[index].`))
					wr.Write([]byte(fmt.Sprintf("%v", field.GetTypeName())))
					wr.Write([]byte(`
`))

				} else {

					wr.Write([]byte(`			e.`))
					wr.Write([]byte(fmt.Sprintf("%v", field.Name)))
					wr.Write([]byte(` = &s.s_`))
					wr.Write([]byte(fmt.Sprintf("%v", field.Name)))
					wr.Write([]byte(`[index]
`))

				}
			}
			for c := range e.ComponentOverridesSeq() {

				wr.Write([]byte(`			e.`))
				wr.Write([]byte(fmt.Sprintf("%v", c.Base.Name)))
				wr.Write([]byte(`.`))
				wr.Write([]byte(fmt.Sprintf("%v", c.Field.Name)))
				wr.Write([]byte(` = &e.`))
				wr.Write([]byte(fmt.Sprintf("%v", c.Field.Name)))
				wr.Write([]byte(`.`))
				wr.Write([]byte(fmt.Sprintf("%v", c.Field.GetTypeName())))
				wr.Write([]byte(`
`))

			}

			wr.Write([]byte(`			age = s.Age()
		}

		return age, e
	}
`))

		}
	}

	wr.Write([]byte(`	if s := s_`))
	wr.Write([]byte(fmt.Sprintf("%v", e.Name)))
	wr.Write([]byte(`; s.TypeId() == tid {
		if age != s.Age() {
			e.Id = id
`))

	for c := range EnumFieldsSeq(e.StructComponentsSeq()) {

		wr.Write([]byte(`			e.`))
		wr.Write([]byte(fmt.Sprintf("%v", c.Name)))
		wr.Write([]byte(` = &s.s_`))
		wr.Write([]byte(fmt.Sprintf("%v", c.Name)))
		wr.Write([]byte(`[index]
`))

	}
	for c := range e.ComponentOverridesSeq() {

		wr.Write([]byte(`			e.`))
		wr.Write([]byte(fmt.Sprintf("%v", c.Base.Name)))
		wr.Write([]byte(`.`))
		wr.Write([]byte(fmt.Sprintf("%v", c.Field.Name)))
		wr.Write([]byte(` = &e.`))
		wr.Write([]byte(fmt.Sprintf("%v", c.Field.Name)))
		wr.Write([]byte(`.`))
		wr.Write([]byte(fmt.Sprintf("%v", c.Field.GetTypeName())))
		wr.Write([]byte(`
`))

	}

	wr.Write([]byte(`			age = s.Age()
		}

		return age, e
	}

	panic("Wrong type requested.")
}
`))

}

func (g *GeneratorEcs) fnStore(wr io.Writer, typ *Type) {
	if !typ.NeedStore() {
		return
	}

	wr.Write([]byte(`
func (e *`))
	wr.Write([]byte(fmt.Sprintf("%v", typ.Name)))
	wr.Write([]byte(`) Store() {
`))

	for field := range EnumFieldsSeq(typ.StoreComponentsSeq()) {
		if field.IsArray {
		} else {
			if field.isEcsRef {
			} else {

				wr.Write([]byte(`	c_`))
				wr.Write([]byte(fmt.Sprintf("%v", field.Name)))
				wr.Write([]byte(` := *e.`))
				wr.Write([]byte(fmt.Sprintf("%v", field.Name)))
				wr.Write([]byte(`
	e.`))
				wr.Write([]byte(fmt.Sprintf("%v", field.Name)))
				wr.Write([]byte(` = &c_`))
				wr.Write([]byte(fmt.Sprintf("%v", field.Name)))
				wr.Write([]byte(`
`))

			}
		}

		g.genFieldEcsCall(wr, field, "Store")
	}

	wr.Write([]byte(`	Update`))
	wr.Write([]byte(fmt.Sprintf("%v", typ.Name)))
	wr.Write([]byte(`Id(e.Id.Store())
}
`))

}

func (g *GeneratorEcs) fnRestore(wr io.Writer, typ *Type) {
	if !typ.NeedRestore() {
		return
	}

	wr.Write([]byte(`
func (e *`))
	wr.Write([]byte(fmt.Sprintf("%v", typ.Name)))
	wr.Write([]byte(`) Restore() {
`))

	for field := range EnumFieldsSeq(typ.StoreComponentsSeq()) {
		if field.IsArray {
		} else {
		}

		g.genFieldEcsCall(wr, field, "Restore")
	}

	if typ.CanCall("Construct") {

		wr.Write([]byte(`		e.Construct()
`))

	}

	wr.Write([]byte(`	Update`))
	wr.Write([]byte(fmt.Sprintf("%v", typ.Name)))
	wr.Write([]byte(`Id(e.Id.Restore())
}
`))

}
