package exchange

import (
	"html/template"

	"github.com/shmel1k/exchangego/exchange/session/context"
	"github.com/shmel1k/exchangego/base/contextlog"
)

func PathToName(name string) TemplateType {
	return TemplateType{
		name: name,
		path: "exchangego/template/controller/" + name + ".html",
	}
}

type TemplateType struct {
	name string
	path string
}

var baseTmpl = TemplateType{
	name: "index.html",
	path: "exchangego/template/index.html",
}

var (
	GameTmpl TemplateType = PathToName("game")
	AuthTmpl TemplateType = PathToName("auth")
	RegTmpl  TemplateType = PathToName("reg")
)

type TemplateData struct {
	IsAuth bool
	UserName string

	ModuleName string
}

func ReturnTemplate(ctx *context.ExContext, tmpl TemplateType) {
	t := template.New(baseTmpl.name)

	t, err := t.ParseFiles(baseTmpl.path, tmpl.path)
	if err != nil {
		contextlog.Println(ctx, err)
		return
	}

	data := TemplateData{
		IsAuth: false,
		UserName: "test",

		ModuleName: tmpl.name,
	}

	t.ExecuteTemplate(ctx.HTTPResponseWriter(), "index.html", data)
}