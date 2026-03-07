package template

import "github.com/vikingpingvin/chainrun/internal/types"

// TextRenderer implements Renderer using Go's text/template package.
// Template data exposes: .Steps.*, .Env.*, .Trigger.*
type TextRenderer struct{}

// NewRenderer returns a new TextRenderer.
func NewRenderer() *TextRenderer {
	return &TextRenderer{}
}

// Render executes tmpl as a Go text/template with ctx as data.
func (r *TextRenderer) Render(tmpl string, ctx *types.RunContext) (string, error) {
	panic("not implemented")
}

// RenderMap renders each value in m as a template and returns the resulting map.
func (r *TextRenderer) RenderMap(m map[string]string, ctx *types.RunContext) (map[string]string, error) {
	panic("not implemented")
}
