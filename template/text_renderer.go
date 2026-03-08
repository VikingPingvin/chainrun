package template

import (
	"bytes"
	"text/template"
	"time"

	"github.com/vikingpingvin/chainrun/internal/types"
)

// TextRenderer implements Renderer using Go's text/template package.
// Template data exposes: .Steps.*, .Env.*, .Trigger.*
type TextRenderer struct{}

type renderData struct {
	Env     map[string]string
	Steps   map[string]types.StepResult
	Trigger types.TriggerEvent
	Time    time.Time
}

// NewRenderer returns a new TextRenderer.
func NewRenderer() *TextRenderer {
	return &TextRenderer{}
}

// Render executes tmpl as a Go text/template with ctx as data.
func (r *TextRenderer) Render(tmpl string, ctx *types.RunContext) (string, error) {
	renderData := renderData{
		Env:     ctx.Env,
		Steps:   ctx.Steps,
		Trigger: ctx.TriggerEvent,
		Time:    time.Now(),
	}

	tmplObj, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmplObj.Execute(&buf, renderData); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// RenderMap renders each value in m as a template and returns the resulting map.
func (r *TextRenderer) RenderMap(m map[string]string, ctx *types.RunContext) (map[string]string, error) {
	res := make(map[string]string, len(m))
	for k, v := range m {
		rendered, err := r.Render(v, ctx)
		if err != nil {
			return nil, err
		}
		res[k] = rendered
	}

	return res, nil
}
