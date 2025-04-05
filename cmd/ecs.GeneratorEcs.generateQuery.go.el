<?go
package cmd

import (
	
	"fmt"
	"io"
	"strings"
)

func (g *GeneratorEcs) generateQuery(wr io.Writer, q *Type, es []*Type) {
	g.genAs(wr, q)

?>

func _<?= q.Name ?>_constraints() {
	var _ ecs.Id = <?= q.Name ?>{}.Id
}

type <?= q.Name ?>Type struct {
}

func (<?= q.Name ?>Type) Age() (age uint64) {
	age = 0
<?
	for _, e := range es {
		if e.GetPackage() == q.GetPackage() {
?>
	age += S_<?= e.Name ?>.Age()
<?
		} else if strings.HasPrefix(e.GetPackage().Pkg.PkgPath, q.GetPackage().Pkg.PkgPath) {
?>
	age += <?= e.GetPackage().Name ?>.S_<?= e.Name ?>.Age()
<?
		}
	}
?>
	return
}

func (<?= q.Name ?>Type) Get(id ecs.Id) (<?= q.Name ?>, bool) {
	t := id.GetType()
	index := (int)(id.GetId() - 1)
	_ = index
	_ = t

<?
	for  _, e := range es {
		if e.GetPackage() == q.GetPackage() {
?>
	if s := &S_<?= e.Name ?>; s.TypeId() == t {
<?
		} else if strings.HasPrefix(e.GetPackage().Pkg.PkgPath, q.GetPackage().Pkg.PkgPath) {
?>
	if s := &<?= e.GetPackage().Name ?>.S_<?= e.Name ?>; s.TypeId() == t {
<?
		} else {
			continue
		}
?>
		return <?= q.Name ?>{
			Id:      id,
<?
		for iq := range EnumFieldsSeq(q.StructComponentsSeq()) {
?>
			<?= iq.Name ?>: &s.S_<?= iq.Name ?>[index],
<?
		}
?>
		}, true
	}
<?
	}
?>

	return <?= q.Name ?>{}, false
}

func (<?= q.Name ?>Type) Do() iter.Seq[<?= q.Name ?>] {
	return func(yield func(<?= q.Name ?>) bool) {
<?
	for  _, e := range es {
	if e.GetPackage() == q.GetPackage() {
?>
	{
		s := &S_<?= e.Name ?>
<?
	} else if strings.HasPrefix(e.GetPackage().Pkg.PkgPath, q.GetPackage().Pkg.PkgPath) {
?>
	{
		s := &<?= e.GetPackage().Name ?>.S_<?= e.Name ?>
<?
	} else {
		continue
	}
?>
	for id := range s.EntityIds() {
		index := (int)(id.GetId() - 1)
		_ = index
		if !yield(<?= q.Name ?>{
			Id:       id,
<?
		for iq := range EnumFieldsSeq(q.StructComponentsSeq()) {
?>
			<?= iq.Name ?>: &s.S_<?= iq.Name ?>[index],
<?
		}
?>
		}) {
			return
		}
	}
}
<?
	}
?>
	}
}

func Age<?= q.Name ?>() (age uint64) {
	return <?= q.Name ?>Type{}.Age()
}

func Get<?= q.Name ?>(id ecs.Id) (<?= q.Name ?>, bool) {
	return <?= q.Name ?>Type{}.Get(id)
}

func Do<?= q.Name ?>() iter.Seq[<?= q.Name ?>] {
	return <?= q.Name ?>Type{}.Do()
}
<?
	if qt, ok := q.Tag.GetObject(Tag_Query); ok && qt.HasField(Tag_Cached) {
?>

type <?= q.Name ?>Cache struct {
	Age    uint64
	Cache []<?= q.Name ?>
}

func (r *<?= q.Name ?>Cache) Query() bool {
	if r.Age != Age<?= q.Name ?>() {
		r.Age = Age<?= q.Name ?>()
		r.Cache = slices.Collect(Do<?= q.Name ?>())

		return true
	}

	return false
}
<?
	}
}
?>