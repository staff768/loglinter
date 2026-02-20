package plugin

import (
	"loglinter/pkg/analyzer"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("loglinter", New)
}

func New(conf any) (register.LinterPlugin, error) {
	return &loglintPlugin{}, nil
}

type loglintPlugin struct{}

func (p *loglintPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}

func (p *loglintPlugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
