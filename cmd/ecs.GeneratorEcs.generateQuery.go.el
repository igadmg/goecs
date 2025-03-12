<?go
package cmd

import (
	
	"fmt"
	"io"
)

func (g *GeneratorEcs) generateQuery(wr io.Writer, q *Type, es []*Type) {
	g.genAs(wr, q)

?>

func _<?= q.Name ?>_constraints() {
	var _ ecs.Id = <?= q.Name ?>{}.Id
}

func Age<?= q.Name ?>() (age uint64) {
	age = 0
<?
	for _, e := range es {
?>
	age += s_<?= e.Name ?>.Age
<?
	}
?>
	return
}

func Get<?= q.Name ?>(id ecs.Id) (<?= q.Name ?>, bool) {
	t := id.GetType()
	index := (int)(id.GetId() - 1)
	_ = index

<?
	for  _, e := range es {
?>
	if s := &s_<?= e.Name ?>; s.TypeId == t {
		return <?= q.Name ?>{
			Id:      id,
<?
		for iq := range EnumFieldsSeq(q.StructComponentsSeq()) {
?>
			<?= iq.Name ?>: &s.s_<?= iq.Name ?>[index],
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

func Do<?= q.Name ?>() iter.Seq[<?= q.Name ?>] {
	return func(yield func(<?= q.Name ?>) bool) {
<?
	for  _, e := range es {
?>
{
	s := &s_<?= e.Name ?>
	for id := range s.EntityIds() {
		index := (int)(id.GetId() - 1)
		_ = index
		if !yield(<?= q.Name ?>{
			Id:       id,
<?
		for iq := range EnumFieldsSeq(q.StructComponentsSeq()) {
?>
			<?= iq.Name ?>: &s.s_<?= iq.Name ?>[index],
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