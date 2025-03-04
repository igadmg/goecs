<?go
package cmd

import (
	
	"fmt"
	"io"
)

func (g *GeneratorEcs) generateQuery(wr io.Writer, q *Type, es []*Type) {
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

func Do<?= q.Name ?>() iter.Seq[<?= q.Name ?>] {
	return func(yield func(<?= q.Name ?>) bool) {
<?
	for  _, e := range es {
		i := 0
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
			i++
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

type <?= q.Name ?>Result struct {
	Age    uint64
	Result []<?= q.Name ?>
}

func (r *<?= q.Name ?>Result) Query() bool {
	if r.Age != Age<?= q.Name ?>() {
		r.Age = Age<?= q.Name ?>()
		r.Result = slices.Collect(Do<?= q.Name ?>())

		return true
	}

	return false
}

<?
}
?>