<?go
package cmd

import (
	"fmt"
	"io"
)

func (g *GeneratorEcs) generateQuery(wr io.Writer, q *Type, es []*Type) {
?>
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

func Execute<?= q.Name ?>(yield func(q <?= q.Name ?>) bool) {
<?
	for  _, e := range es {
		i := 0
?>
{
	s := &s_<?= e.Name ?>
	for id := range s.EntityIds() {
		_, e := ecs.GetT[<?= e.Name ?>](id)
		if !yield(<?= q.Name ?>{
			Id:       id,
<?
		for iq := range EnumFieldsSeq(q.StructComponentsSeq()) {
			i++
?>
			<?= iq.Name ?>: e.<?= iq.Name ?>,
<?
		}
?>
		}) {
			return
		}
<?
		if i == 0 {
?>
		_ = e
<?
		}
?>
	}
}
<?
	}
?>
}

type <?= q.Name ?>Result struct {
	Age    uint64
	Result []<?= q.Name ?>
}

func (r *<?= q.Name ?>Result) Query() bool {
	if r.Age != Age<?= q.Name ?>() {
		r.Age = Age<?= q.Name ?>()
		r.Result = slices.Collect(Execute<?= q.Name ?>)

		return true
	}

	return false
}

<?
}
?>