<?go
package cmd

import (
	"fmt"
	"io"
)

func (g *GeneratorEcs) genAs(wr io.Writer, t *Type) {
	for f := range EnumFields(t.Fields) {
		if fet, ok := f.GetType().(EcsTypeI); ok {
			for af := range fet.AsComponentsSeq() {
				if af.IsEcsRef() {
?>

func (e <?= t.Name ?>) <?= af.GetA() ?>() <?= af.GetType().GetName() ?> {
	return e.<?= f.GetName() ?>.<?= af.GetName() ?>.Get()
}

func (e <?= t.Name ?>) <?= af.GetA() ?>Ref() ecs.Ref[<?= af.GetType().GetName() ?>] {
	return e.<?= f.GetName() ?>.<?= af.GetName() ?>
}

func (e <?= t.Name ?>) Set<?= af.GetA() ?>(v <?= af.GetTypeName() ?>) {
	e.<?= f.GetName() ?>.<?= af.GetName() ?> = v
}
<?
				} else {
?>

func (e <?= t.Name ?>) <?= af.GetA() ?>() <?= af.GetTypeName() ?> {
	return e.<?= f.GetName() ?>.<?= af.GetName() ?>
}

func (e <?= t.Name ?>) Set<?= af.GetA() ?>(v <?= af.GetTypeName() ?>) {
	e.<?= f.GetName() ?>.<?= af.GetName() ?> = v
}
<?					
				}
			}
		}
	}
}
?>