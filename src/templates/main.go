package templates

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/pepelazz/projectGenerator/src/types"
	"github.com/pepelazz/projectGenerator/src/utils"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

var project *types.ProjectType

func SetProject(p *types.ProjectType)  {
	project = p
}


var funcMap = template.FuncMap{
	"ToUpper":        strings.ToUpper,
	"ToLower":        strings.ToLower,
	"UpperCaseFirst": utils.UpperCaseFirst,
	"ToLowerCamel": strcase.ToLowerCamel,
	"PrintVueFldTemplate": PrintVueFldTemplate,
}

func ParseTemplates(p types.ProjectType) map[string]*template.Template {
	// парсинг общих шаблонов
	res := map[string]*template.Template{}
	readFiles := func(prefix, delimLeft, delimRight string, path ...string) {
		tmpls, err := template.New("").Funcs(funcMap).Delims(delimLeft, delimRight).ParseFiles(path...)
		utils.CheckErr(err, "ParseFiles")
		for _, t := range tmpls.Templates() {
			res[prefix + t.Name()] = t
		}
	}
	// project
	path := "../../projectGenerator/src/templates/project/"
	readFiles("project_", "{{", "}}", path+"config.toml", path+"docker-compose.yml", path+"docker-compose.dev.yml", path+"restoreDump.sh", path+"deploy.ps1")

	// webClient
	path = "../../projectGenerator/src/templates/webClient/doc/"
	readFiles("webClient_", "[[", "]]", path + "index.vue", path + "item.vue")
	// sql
	path = "../../projectGenerator/src/templates/sql/"
	readFiles("sql_", "{{", "}}", path + "main.toml")
	path = "../../projectGenerator/src/templates/sql/function/"
	readFiles("sql_function_", "{{", "}}", path + "get_by_id.sql", path + "list.sql", path + "update.sql", path + "trigger_before.sql", path + "trigger_after.sql")

	// парсинг шаблонов для конкретного документа
	for i, d := range p.Docs {
		for tName, dt := range d.Templates {
			t, err := template.New(tName).Funcs(funcMap).Delims("[[", "]]").ParseFiles(dt.Source)
			utils.CheckErr(err, fmt.Sprintf("doc: %s tmpl: %s parse template error: %s", d.Name, tName, err))
			// сохраняем template в поле структуры
			dt.Tmpl = t
		}
		// дописываем стандартные шаблоны
		baseTmplNames := []string{}
		if d.IsBaseTemapltes.Vue {
			baseTmplNames = append(baseTmplNames, "webClient_item.vue", "webClient_index.vue")
		}
		if d.IsBaseTemapltes.Sql {
			baseTmplNames = append(baseTmplNames, "sql_main.toml", "sql_function_get_by_id.sql", "sql_function_list.sql", "sql_function_update.sql", "sql_function_trigger_before.sql", "sql_function_trigger_after.sql")
		}
		for _, tName := range baseTmplNames{
			// если шаблона с таким именем нет, то добавляем стандартный
			if _, ok := d.Templates[tName]; !ok {
				if tName == "sql_function_trigger_before.sql" && !d.Sql.IsBeforeTrigger {
					continue
				}
				if tName == "sql_function_trigger_after.sql" && !d.Sql.IsAfterTrigger {
					continue
				}
				distPath, distFilename := utils.ParseDocTemplateFilename(d.Name, tName, p.DistPath, i)
				d.Templates[tName]= &types.DocTemplate{Tmpl: res[tName], DistPath: distPath, DistFilename: distFilename}
			}
		}
	}

	return res
}

func ExecuteToFile(t *template.Template, d interface{}, path, filename string) error {
	if t == nil {
		log.Fatalf("template is nil for path '%s/%s'\n", path, filename)
	}
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	var tpl bytes.Buffer
	err = t.Execute(&tpl, d)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path +"/" + filename, []byte(tpl.String()), 0644)
}

// печать vue темплейтов для
func PrintVueFldTemplate(fld types.FldType) string {
	name := fld.Vue.Name
	if len(name) == 0 {
		name = fld.Name
	}
	nameRu := fld.Vue.NameRu
	if len(nameRu) == 0 {
		nameRu = fld.NameRu
	}
	fldType := fld.Vue.Type
	if len(fldType) == 0 {
		fldType = fld.Type
		// в случае ref поля
		if fld.Type == types.FldTypeInt && len(fld.Sql.Ref) > 0 {
			fldType = "ref"
		}
	}
	// если указана функция для композиции, то меняем тип на vueComposition
	if fld.Vue.Composition != nil {
		fldType = types.FldTypeVueComposition
	}
 	switch fldType {
	case types.FldTypeString, types.FldTypeText:
		return fmt.Sprintf(`<q-input outlined type='text' v-model="item.%s" label="%s" autogrow/>`, name, nameRu)
	case types.FldTypeInt, types.FldTypeDouble:
		return fmt.Sprintf(`<q-input outlined type='number' v-model="item.%s" label="%s"/>`, name, nameRu)
	// дата
	case types.FldTypeDate:
		return fmt.Sprintf(`<comp-fld-date label="%s" :date-string="$utils.formatPgDate(item.%s)" @update="v=> item.%s = v"/>`, nameRu, name, name)
	// дата с временем
	case types.FldTypeDatetime:
		return fmt.Sprintf(`<comp-fld-date-time label="%s" :date-string="$utils.formatPgDateTime(item.%s)" @update="v=> item.%s = v"/>`, nameRu, name, name)
	// вариант ссылки на другую таблицу
	case "ref":
		// если map Ext не инициализирован, то создаем его, чтобы не было ошибки при json.Marshal
		if fld.Vue.Ext == nil  {
			fld.Vue.Ext = map[string]string{}
		}
		// если специально не определено поле для ajaxSelectTitle, то формируем стандартное [ref_table_name]_title
		ajaxSelectTitle := fld.Sql.Ref + "_title"
		if v, ok := fld.Vue.Ext["ajaxSelectTitle"]; ok {
			ajaxSelectTitle = v
		}
		extJsonStr, err := json.Marshal(fld.Vue.Ext)
		utils.CheckErr(err, fmt.Sprintf("json.Marshal(fld.Vue.Ext) fld %s", fld.Name))

		// заполняем название postgres метода для получения списка. По дефолту [ref_table_name]_list
		pgMethod := fld.Sql.Ref + "_list"
		if m, ok := fld.Vue.Ext["pgMethod"]; ok {
			pgMethod = m
		}
		return fmt.Sprintf(`<comp-fld-ref-search pgMethod="%s" label="%s" :item='item.%s' :ext='%s' @update="v=> item.%s = v.id" />`, pgMethod, nameRu, ajaxSelectTitle, extJsonStr, name)
	case types.FldVueTypeSelect, types.FldVueTypeMultipleSelect:
		options, err := json.Marshal(fld.Vue.Options)
		utils.CheckErr(err, fmt.Sprintf("'%s' json.Marshal(fld.Vue.Options)", fld.Name))
		multiple := ""
		if fldType == types.FldVueTypeMultipleSelect {
			multiple = "multiple"
		}
		return fmt.Sprintf(`<q-select outlined label="%s" v-model='item.%s' :options='%s' %s/>`, nameRu, name, options, multiple)
	case types.FldTypeVueComposition:
		if fld.Vue.Composition == nil {
			log.Fatal(fmt.Sprintf("fld have type '%s', but fld.Vue.Composition function is nil", types.FldTypeVueComposition))
		}
		return fld.Vue.Composition(*project, *fld.Doc)
	default:
		return fmt.Sprintf("not found vueFldTemplate for type `%s`", fldType)
	}
}
