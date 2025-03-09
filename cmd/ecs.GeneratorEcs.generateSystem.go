package cmd

import (
	"fmt"
	"io"
)

func (g *GeneratorEcs) generateSystem(wr io.Writer, t *Type) {

	wr.Write([]byte(`

var _ ecs.System = (*`))
	wr.Write([]byte(fmt.Sprintf("%v", t.Name)))
	wr.Write([]byte(`)(nil)
`))

}
