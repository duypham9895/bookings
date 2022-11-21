package config

import (
	"html/template"

	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
	IsUsedCache   bool
	TemplateCache map[string]*template.Template
	IsEnvProd     bool
	Session       *scs.SessionManager
}
