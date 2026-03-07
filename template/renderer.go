package template

import "github.com/vikingpingvin/chainrun/internal/types"

// Renderer renders Go text/template strings against a RunContext.
type Renderer interface {
	Render(tmpl string, ctx *types.RunContext) (string, error)
	RenderMap(m map[string]string, ctx *types.RunContext) (map[string]string, error)
}
