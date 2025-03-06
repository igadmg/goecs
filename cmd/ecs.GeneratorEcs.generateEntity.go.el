<?go
package cmd

import (
	"fmt"
	"io"
)

func (g *GeneratorEcs) generateEntity(wr io.Writer, id int, e *Type) {
	g.genAs(wr, e)

?>

type storage_<?= e.Name ?> struct {
	ecs.BaseStorage

<?
	for c := range EnumFieldsSeq(e.StructComponentsSeq()) {
?>
	s_<?= c.Name ?> []<?= c.Type.GetName() ?>
<?
	}
?>
}

var s_<?= e.Name ?> = storage_<?= e.Name ?>{
	BaseStorage: ecs.BaseStorage{TypeId: <?= id  ?>},
}

func Match<?= e.Name ?>(id ecs.Id) (ecs.Ref[<?= e.Name ?>], bool) {
	if id.GetType() == s_<?= e.Name ?>.TypeId {
		ref := ecs.Ref[<?= e.Name ?>]{Id: id}
		_ = ref.Get()

		return ref, true
	}

	return ecs.Ref[<?= e.Name ?>]{}, false
}

func (e <?= e.Name ?>) Ref() ecs.Ref[<?= e.Name ?>] {
	return ecs.Ref[<?= e.Name ?>] { Id: e.Id }
}

func (e *<?= e.Name ?>) Allocate() ecs.Ref[<?= e.Name ?>] {
	s := &s_<?= e.Name ?>
	age, id := s.BaseStorage.AllocateId()
	index := (int)(id.GetId() - 1)
	_ = index

<?
	for c := range EnumFieldsSeq(e.StructComponentsSeq()) {
?>
	s.s_<?= c.Name ?> = slicesex.Reserve(s.s_<?= c.Name ?>, index+1)
<?
	}
?>

	ref := ecs.Ref[<?= e.Name ?>]{
		Age: age - 1,
		Id:  id,
	}
	_ = ref.Get()

	if e != nil {
<?

	for c := range EnumFieldsSeq(e.StructComponentsSeq()) {
		if ct, ok := CastType(c.Type); ok {
			if ct.IsTransient() {
				continue
			}

?>
		if e.<?= c.Name ?> != nil {
			*ref.Ptr.<?= c.Name ?> = *e.<?= c.Name ?>
		}
<?
		}
 	}
?>
		*e = ref.Ptr
	}

	return ref
}

func (e *<?= e.Name ?>) Free() {
	Free<?= e.Name ?>(e.Id)
}

func (e <?= e.Name ?>) Load(age uint64, id ecs.Id) (uint64, <?= e.Name ?>) {
	index := (int)(id.GetId() - 1)
	tid := id.GetType()
	_ = index

<?
 	for _, s := range e.Subclasses {
		switch sc := s.(type) {
		case *Type:
?>
	if s := &s_<?= sc.Name ?>; s.TypeId == tid {
		if age != s.Age {
			e.Id = id
<?
			for field := range EnumFields(e.Fields) {
				if field.Tag.HasField(Tag_Virtual) || field.Tag.HasField(Tag_Abstract) {
?>
			e.<?= field.Name ?> = &s.s_<?= field.Name ?>[index].<?= field.Type.GetName() ?>
<?
				} else {
?>
			e.<?= field.Name ?> = &s.s_<?= field.Name ?>[index]
<?
				}
			}
			for c := range e.ComponentOverridesSeq() {
?>
			e.<?= c.Base.Name ?>.<?= c.Field.Name ?> = &e.<?= c.Field.Name ?>.<?= c.Field.Type.GetName() ?>
<?
			}
?>
			age = s.Age
		}

		return age, e
	}
<?
		}
	}
?>
	if s := s_<?= e.Name ?>; s.TypeId == tid {
		if age != s.Age {
			e.Id = id
<?
	for c := range EnumFieldsSeq(e.StructComponentsSeq()) {
?>
			e.<?= c.Name ?> = &s.s_<?= c.Name ?>[index]
<?
 	}
	for c := range e.ComponentOverridesSeq() {
?>
			e.<?= c.Base.Name ?>.<?= c.Field.Name ?> = &e.<?= c.Field.Name ?>.<?= c.Field.Type.GetName() ?>
<?
	}
?>
			age = s.Age
		}

		return age, e
	}

	panic("Wrong type requested.")
}
<?
	if !e.IsTransient() {
		g.fnStore(wr, e)
		g.fnRestore(wr, e)
	}
?>

func Allocate<?= e.Name ?>() (ref ecs.Ref[<?= e.Name ?>], entity <?= e.Name ?>) {
	var e *<?= e.Name ?> = nil
	ref = e.Allocate()
	return ref, ref.Ptr
}

func Free<?= e.Name ?>(id ecs.Id) {
	s := &s_<?= e.Name ?>
	index := (int)(id.GetId() - 1)
	_ = index

<?
	for c := range EnumFieldsSeq(e.StructComponentsSeq()) {
?>
	s.s_<?= c.Name ?>[index] = <?= c.Type.GetName() ?>{}
<?
 	}
?>

	s.Free(id)
}

func Update<?= e.Name ?>Id(id ecs.Id) {
	tid := id.GetType()
	if s := s_<?= e.Name ?>; s.TypeId == tid {
		index := (int)(id.GetId() - 1)

		s_<?= e.Name ?>.Ids[index] = id
	}
}
<?
	if _, ok := g.queries[e.Name+"Query"]; !ok {
?>

// Auto-generated query for <?= e.Name ?> entity
type <?= e.Name ?>Query struct {
	_ MetaTag `ecs:"ecsq"`

	Id ecs.Id
<?
	for c := range EnumFieldsSeq(e.StructComponentsSeq()) {

?>
	<?= c.Name ?> *<?= c.Type.GetName() ?>
<?
	}
?>
}
<?
	}


	qt := NewType()
	qt.Name= e.Name + "Query"
	qt.Fields = e.Fields
	g.generateQuery(wr, qt, []*Type{e})
}

func (g *GeneratorEcs) genFieldEcsCall(wr io.Writer, f *Field, call string) {
	if f.IsArray {
		if f.isEcsRef {
?>
	for i := range e.<?= f.Name ?> {
		<?= call ?>(&e.<?= f.Name ?>[i])
}
<?
		} else if f.Type.CanCall(call) {
?>
	for i := range e.<?= f.Name ?> {
		e.<?= f.Name ?>[i].<?= call ?>()
}
<?
		}
	} else {
		if f.isEcsRef {
?>
	<?= call ?>(&e.<?= f.Name ?>)
<?
		} else if f.Type.CanCall(call) {
?>
	e.<?= f.Name ?>.<?= call ?>()
<?
		}
	}
}

func (g *GeneratorEcs) fnStore(wr io.Writer, typ *Type) {
	if !typ.NeedStore() {
		return
	}
?>

func (e *<?= typ.Name ?>) Store() {
<?
	for field := range EnumFieldsSeq(typ.StoreComponentsSeq()) {
		if field.IsArray {
			} else {
				if field.isEcsRef {
				} else {
?>
	c_<?= field.Name ?> := *e.<?= field.Name ?>
	e.<?= field.Name ?> = &c_<?= field.Name ?>
<?
				}
			}

			g.genFieldEcsCall(wr, field, "Store")
	}
?>
	Update<?= typ.Name ?>Id(e.Id.Store())
}
<?
}

func (g *GeneratorEcs) fnRestore(wr io.Writer, typ *Type) {
	if !typ.NeedRestore() {
		return
	}
?>

func (e *<?= typ.Name ?>) Restore() {
<?
	for field := range EnumFieldsSeq(typ.StoreComponentsSeq()) {
		if field.IsArray {
		} else {
		}

		g.genFieldEcsCall(wr, field, "Restore")
	}

	if typ.CanCall("Construct") {
?>
		e.Construct()
<?
	}

	/*
	// Not sure about that, but that fixes problem.
	for b, p := range typ.baseParameters.AllFromFront() {
		if bt, ok := CastField(p.Base); !ok || !bt.IsTransient() {
			continue
		}

		Need to decouple Create and Construct here and call only parameterless

?>
	e.<?= b ?>.Create(<?= strings.Join(p.Params, ",") ?>)
<?
	}

		Need to call base construct from consturct
	*/
?>
	Update<?= typ.Name ?>Id(e.Id.Restore())
}
<?
}
?>