package main

import (
	"log"
	"net/http"
	"time"

	"github.com/vikingpingvin/chainrun/action"
	"github.com/vikingpingvin/chainrun/action/file"
	httpaction "github.com/vikingpingvin/chainrun/action/http"
	"github.com/vikingpingvin/chainrun/action/llm"
	"github.com/vikingpingvin/chainrun/action/notify"
	"github.com/vikingpingvin/chainrun/action/shell"
	"github.com/vikingpingvin/chainrun/cmd"
	"github.com/vikingpingvin/chainrun/config"
	"github.com/vikingpingvin/chainrun/engine"
	"github.com/vikingpingvin/chainrun/internal/types"
	"github.com/vikingpingvin/chainrun/secrets"
	"github.com/vikingpingvin/chainrun/template"
	"github.com/vikingpingvin/chainrun/trigger"
	trigcron "github.com/vikingpingvin/chainrun/trigger/cron"
	"github.com/vikingpingvin/chainrun/trigger/manual"
	"github.com/vikingpingvin/chainrun/trigger/webhook"
)

func main() {
	secretsRes := secrets.NewEnvResolver()
	renderer := template.NewRenderer()
	cfgLoader := config.NewYAMLLoader()
	httpClient := &http.Client{Timeout: 30 * time.Second}

	actReg := action.NewRegistry()
	actReg.Register("shell", func() action.Executor { return shell.New() })
	actReg.Register("http", func() action.Executor { return httpaction.New(httpClient, renderer) })
	actReg.Register("llm", func() action.Executor { return llm.New(httpClient, secretsRes, renderer) })
	actReg.Register("notify", func() action.Executor { return notify.New(httpClient, renderer) })
	actReg.Register("file", func() action.Executor { return file.New(renderer) })

	trgReg := trigger.NewRegistry()
	trgReg.Register("cron", func(def types.TriggerDef) (trigger.Source, error) { return trigcron.New(def) })
	trgReg.Register("manual", func(def types.TriggerDef) (trigger.Source, error) { return manual.New(def) })
	trgReg.Register("webhook", func(def types.TriggerDef) (trigger.Source, error) { return webhook.New(def) })

	factory := func(cfgPath string) (engine.Engine, error) {
		cfg, err := cfgLoader.Load(cfgPath)
		if err != nil {
			return nil, err
		}
		return engine.New(engine.Deps{
			Workflows:       cfg.Workflows,
			ActionRegistry:  actReg,
			TriggerRegistry: trgReg,
			Renderer:        renderer,
			Secrets:         secretsRes,
			Loader:          cfgLoader,
			Logger:          newLogger(),
		}), nil
	}

	if err := cmd.Execute(factory, cfgLoader); err != nil {
		log.Fatal(err)
	}
}

// newLogger returns a stdlib-backed Logger implementation.
func newLogger() types.Logger {
	return &stdLogger{l: log.Default()}
}

type stdLogger struct {
	l *log.Logger
}

func (s *stdLogger) Info(msg string, fields ...interface{})  { s.l.Printf("[INFO]  "+msg, fields...) }
func (s *stdLogger) Warn(msg string, fields ...interface{})  { s.l.Printf("[WARN]  "+msg, fields...) }
func (s *stdLogger) Error(msg string, fields ...interface{}) { s.l.Printf("[ERROR] "+msg, fields...) }
func (s *stdLogger) Debug(msg string, fields ...interface{}) { s.l.Printf("[DEBUG] "+msg, fields...) }
