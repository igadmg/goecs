<?go
package cmd

import (
	"fmt"
	"io"

	"github.com/igadmg/gogen/core"
)

func (g *GeneratorEcs) genFunctions(wr io.Writer, typ core.TypeI) {
	for f := range core.EnumFuncsSeq(typ.FuncsSeq()) {
		if f.Tag.IsEmpty() {
			continue
		}
		if ecsf, ok := f.Tag.GetObject("ecs"); ok {
			if ecsf.HasField(Tag_Fn_RefCall) {
?>

func (o *<?= typ.GetName() ?>) <?= f.Name ?>_ref(id ecs.Id) func(<?= f.DeclArguments() ?>) {
	return func(<?= f.DeclArguments() ?>) {
		_, o := ecs.GetT[<?= typ.GetName() ?>](id)
		o.<?= f.Name ?>(<?= f.CallArguments() ?>)
	}
}
<?
			}
		}
	}
}
?>