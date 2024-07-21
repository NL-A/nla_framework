package templates

import (
	"fmt"
	"text/template"

	"github.com/NL-A/nla_framework/types"
	"github.com/NL-A/nla_framework/utils"
)

func docIsRecursionProccess(p types.ProjectType, d *types.DocType) {
	sourcePath := fmt.Sprintf("%s/webClient/quasar_%v/doc/comp/recursiveChildList.vue", getCurrentDir(), p.GetQuasarVersion())
	// проверяем возможность того, что путь к шаблону был переопределен внутри документа
	if d.TemplatePathOverride != nil {
		if tmpl, ok := d.TemplatePathOverride["recursiveChildList.vue"]; ok {
			if len(tmpl.Source) > 0 {
				sourcePath = tmpl.Source
			}
		}
	}
	t, err := template.New("recursiveChildList.vue").Funcs(funcMap).Delims("[[", "]]").ParseFiles(sourcePath)
	utils.CheckErr(err, "recursiveChildList.vue")
	docRouteName := d.Name
	if len(d.Vue.Path) > 0 {
		docRouteName = d.Vue.Path
	}
	distPath := fmt.Sprintf("%s/webClient/src/app/components/%s/comp", p.DistPath, docRouteName)
	d.Templates["webClient_comp_recursiveChildList.vue"] = &types.DocTemplate{Tmpl: t, DistPath: distPath, DistFilename: "recursiveChildList.vue"}
}
