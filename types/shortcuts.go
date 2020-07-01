package types

import (
	"fmt"
	"github.com/pepelazz/projectGenerator/utils"
	"log"
	"strconv"
	"strings"
	"text/template"
)

// создание поля title
func GetFldTitle(params ...string) (fld FldType) {
	classStr := "col-4"
	if len(params)>0 {
		classStr= params[0]
	}
	fld = FldType {Name: "title", NameRu: "название", Type: FldTypeString, Sql: FldSql{IsRequired: true, IsUniq: true, IsSearch:true, Size:150}, Vue: FldVue{RowCol: [][]int{{1, 1}}, Class: []string{classStr}}}
	return
}

// создание поля title, которое заполняется тригером
func GetFldTitleComputed(triggerSqlString string) (fld FldType) {
	fld = FldType {Name: "title", NameRu: "название", Type: FldTypeString, Sql: FldSql{IsSearch:true, FillValueInBeforeTrigger: triggerSqlString}}
	return
}

func GetFldDouble(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	classStr := "col-4"
	if len(params)>0 {
		classStr= params[0]
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeDouble, Vue:FldVue{RowCol: rowCol, Class: []string{classStr}}}
	return
}

// создание простого поля String
func GetFldString(name, nameRu string, size int, rowCol [][]int, params ...string) (fld FldType) {
	classStr := "col-4"
	readonly := "false"
	for i, v := range params {
		if i == 0 {
			classStr = v
		} else {
			if strings.HasPrefix(v, "readonly") && strings.HasSuffix(v, "true") {
				readonly="true"
			}
		}
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeString, Vue:FldVue{RowCol: rowCol, Class: []string{classStr}, Readonly:readonly}}
	if size > 0 {
		fld.Sql.Size = size
	}
	return
}

// создание простого Date
func GetFldDate(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	classStr := "col-4"
	if len(params)>0 {
		classStr= params[0]
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeDate, Vue:FldVue{RowCol: rowCol, Class: []string{classStr}}}
	return
}

// создание простого DateTime
func GetFldDateTime(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	classStr := "col-4"
	if len(params)>0 {
		classStr= params[0]
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeDatetime, Vue:FldVue{RowCol: rowCol, Class: []string{classStr}}}
	return
}

// создание простого поля Int
func GetFldInt(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	classStr := "col-4"
	if len(params)>0 {
		classStr= params[0]
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeInt, Vue:FldVue{RowCol: rowCol, Class: []string{classStr}}}
	return
}

// создание простого поля Int64
func GetFldInt64(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	classStr := "col-4"
	if len(params)>0 {
		classStr= params[0]
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeInt64, Vue:FldVue{RowCol: rowCol, Class: []string{classStr}}}
	return
}

// создание поля UUID
func GetFldUuid(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	classStr := "col-4"
	if len(params)>0 {
		classStr= params[0]
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeUuid, Vue:FldVue{RowCol: rowCol, Class: []string{classStr}}}
	return
}

// создание простого поля Checkbox
func GetFldCheckbox(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	classStr := "col-4"
	readonly := "false"
	for i, v := range params {
		if i == 0 {
			classStr = v
		} else {
			if strings.HasPrefix(v, "readonly") && strings.HasSuffix(v, "true") {
				readonly="true"
			}
		}
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeBool, Vue:FldVue{RowCol: rowCol, Type:FldVueTypeCheckbox, Class: []string{classStr}, Readonly:readonly}}
	return
}

// создание простого поля Radio
func GetFldRadioString(name, nameRu string, rowCol [][]int, options []FldVueOptionsItem, params ...string) (fld FldType) {
	classStr := "col-4"
	readonly := "false"
	for i, v := range params {
		if i == 0 {
			classStr = v
		} else {
			if strings.HasPrefix(v, "readonly") && strings.HasSuffix(v, "true") {
				readonly="true"
			}
		}
	}

	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeString, Sql:FldSql{Size:50}, Vue:FldVue{RowCol: rowCol, Type:FldVueTypeRadio, Options: options, Class: []string{classStr}, Readonly:readonly}}
	return
}

// создание простого поля Ref
func GetFldRef(name, nameRu, refTable string, rowCol [][]int, params ...string) (fld FldType) {
	classStr := "col-4"
	if len(params)>0 {
		classStr= params[0]
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeInt,  Sql: FldSql{Ref: refTable, IsSearch:true}, Vue:FldVue{RowCol: rowCol, Ext: map[string]string{}, Class: []string{classStr}}}
	for _, v := range params {
		// добавляем аватарку с ссылкой на выбранный документ
		if v == "isShowLink" {
			// проставляем значение pathUrl и avatar на последнем шаге, после инициализации всех документов  в методе FillVueFlds
			fld.Vue.Ext["pathUrl"] = ""
			fld.Vue.Ext["avatar"] = ""
		}
		for _, v := range params {
			if v == "isClearable" {
				fld.Vue.Ext["isClearable"] = "true"
			}
		}
	}
	return
}

// поле с кастомной композицией
func GetFldJsonbComposition(name, nameRu string, rowCol [][]int, classStr, compName string, params ...string) (fld FldType) {
	isOptionsFld := ""
	for _, v := range params {
		// IsOptionFld передаем отдельным параемтром, потому что SetIsOptionFld() срабатывает уже после того как строка с компонентой сформмирована
		if v == "IsOptionFld" {
			isOptionsFld = "options."
		}
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeJsonb,  Vue:FldVue{RowCol: rowCol, Class: []string{classStr}, Composition: func(p ProjectType, d DocType) string {
		return fmt.Sprintf("<%[1]s :fld='item.%[5]s%[2]s' :item='item' @update='item.%[5]s%[2]s = $event' label='%[3]s' %[4]s/>", compName, name, nameRu, strings.Join(params, " "), isOptionsFld)
	}}}
	return
}

// простое html поле
func GetFldSimpleHtml(rowCol [][]int, classStr, htmlStr string) (fld FldType) {
	fld = FldType{Type:FldTypeVueComposition,  Vue:FldVue{RowCol: rowCol, Class: []string{classStr}, Composition: func(p ProjectType, d DocType) string {
		return htmlStr
	}}}
	return
}

// создание простого поля Select с типом string
func GetFldSelectString(name, nameRu string, size int, rowCol [][]int, options []FldVueOptionsItem, params ...string) (fld FldType) {
	classStr := "col-4"
	readonly := "false"
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeString, Vue:FldVue{RowCol: rowCol, Type: FldVueTypeSelect, Ext: map[string]string{}, Class: []string{classStr}, Readonly:readonly, Options:options}}
	for i, v := range params {
		if i == 0 {
			classStr = v
		} else {
			if strings.HasPrefix(v, "readonly") && strings.HasSuffix(v, "true") {
				fld.Vue.Readonly = "true"
			}
			if v == "isClearable" {
				fld.Vue.Ext["isClearable"] = "true"
			}
		}
	}
	if size > 0 {
		fld.Sql.Size = size
	}
	return
}

// создание простого поля MultipleSelect с типом string
func GetFldSelectMultilple(name, nameRu string, rowCol [][]int, options []FldVueOptionsItem, params ...string) (fld FldType) {
	classStr := "col-4"
	readonly := "false"
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeTextArray, Vue:FldVue{RowCol: rowCol, Type: FldVueTypeMultipleSelect, Ext: map[string]string{}, Class: []string{classStr}, Readonly:readonly, Options:options}}
	for i, v := range params {
		if i == 0 {
			classStr = v
		} else {
			if strings.HasPrefix(v, "readonly") && strings.HasSuffix(v, "true") {
				fld.Vue.Readonly = "true"
			}
			if v == "isClearable" {
				fld.Vue.Ext["isClearable"] = "true"
			}
		}
	}

	return
}

// создание простого поля Int
func GetFldTag(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	classStr := "col-4"
	onlyExistTags := "false" // флаг для UI контрола, чтобы можно было только выбирать из существующих тэгов и нельзя было создавать новые
	for i, v := range params {
		if i == 0 {
			classStr = v
		} else {
			if strings.HasPrefix(v, "only_exist_tags") {
				onlyExistTags="true"
			}
		}
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeTextArray, Vue:FldVue{RowCol: rowCol, Type: FldVueTypeTags,  Class: []string{classStr}, Ext: map[string]string{"onlyExistTags": onlyExistTags}}}
	return
}

// создание поля-виджета со связями многие-к-многим
func GetFldLinkListWidget(linkTable string, rowCol [][]int, classStr string, opts map[string]interface{}) (fld FldType) {
	return FldType{Type: FldTypeVueComposition,  Vue: FldVue{RowCol: rowCol, Class: []string{classStr}, Composition: func(p ProjectType, d DocType) string {
		return GetVueCompLinkListWidget(p, d, linkTable, opts)
	}}}
}


// функция конвертации списка имен файлов с шаблонами в  map[string]*DocTemplate
func GetCustomTemplates(p ...string) map[string]*DocTemplate  {
	res := map[string]*DocTemplate{}
	for _, name := range p {
		res[name] = &DocTemplate{}
	}
	return res
}

// создание поля адрес с возможностью поиска через dadata
func GetFldDadataAddress(name, nameRu string, rowCol [][]int, params ...string) (fld FldType) {
	classStr := "col-4"
	if len(params)>0 {
		classStr= params[0]
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeJsonb, Vue:FldVue{RowCol: rowCol, Type: FldVueTypeDadataAddress, Class: []string{classStr}}}
	return
}

// создание поля json c редактируемым массивом элементов
func GetFldJsonList(name, nameRu string, rowCol [][]int, listParams FldVueJsonList, params ...string) (fld FldType) {
	classStr := "col-4"
	if len(params)>0 {
		classStr= params[0]
	}
	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeJsonb, Vue:FldVue{RowCol: rowCol, Type: FldVueTypeJsonList, JsonList: listParams, Class: []string{classStr}}}
	return
}

// создание поля для загрузки файлов
func GetFldFiles(name, nameRu string, rowCol [][]int, fileParams FldVueFilesParams, params ...string) (fld FldType) {
	classStr := "col-4"
	if len(params)>0 {
		classStr= params[0]
	}
	// заполняем параметры для ограничений по загрузке файлов
	ext := map[string]string{}
	if len(fileParams.Accept)>0{
		ext["accept"] = fileParams.Accept
	}
	if fileParams.MaxFileSize>0{
		ext["maxFileSize"] = strconv.FormatInt(fileParams.MaxFileSize, 10)
	}

	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeJsonb, Vue:FldVue{RowCol: rowCol, Type: FldVueTypeFiles, Ext: ext, Class: []string{classStr}}}
	return
}

// создание поля для загрузки списка изображений
func GetFldImgList(name, nameRu string, rowCol [][]int, fileParams FldVueImgParams, params ...string) (fld FldType) {
	classStr := "col-4"
	if len(params)>0 {
		classStr= params[0]
	}
	// заполняем параметры для ограничений
	ext := map[string]string{}
	if len(fileParams.Accept) > 0{
		ext["accept"] = fileParams.Accept
	}
	if fileParams.MaxFileSize > 0{
		ext["maxFileSize"] = strconv.FormatInt(fileParams.MaxFileSize, 10)
	}
	if fileParams.CanAddUrls {
		ext["canAddUrls"] = "true"
	}
	if len(fileParams.Crop) > 0 {
		// проверка что crop имеет формат 300x400
		arr := strings.Split(fileParams.Crop, "x")
		if len(arr) != 2 {
			log.Fatalf("GetFldImgList error fld: '%s' in FldVueImgParams.Crop must be such format '300x400'. You write this: %s", name, fileParams.Crop)
		}
		if _, err := strconv.Atoi(arr[0]); err != nil {
			log.Fatalf("GetFldImgList error fld: '%s' in FldVueImgParams.Crop must be such format '300x400'. %s not number", name, arr[0])
		}
		if _, err := strconv.Atoi(arr[1]); err != nil {
			log.Fatalf("GetFldImgList error fld: '%s' in FldVueImgParams.Crop must be such format '300x400'. %s not number", name, arr[1])
		}
		ext["crop"] = fileParams.Crop
	}
	if fileParams.Width > 0 {
		ext["width"] = strconv.Itoa(fileParams.Width)
	}

	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeJsonb, Vue:FldVue{RowCol: rowCol, Type: FldVueTypeImgList, Ext: ext, Class: []string{classStr}}}
	return
}

// создание поля для загрузки одного
func GetFldImg(name, nameRu string, rowCol [][]int, fileParams FldVueImgParams, params ...string) (fld FldType) {
	classStr := "col-4"
	if len(params)>0 {
		classStr= params[0]
	}
	// заполняем параметры для ограничений
	ext := map[string]string{}
	if len(fileParams.Accept) > 0{
		ext["accept"] = fileParams.Accept
	}
	if fileParams.MaxFileSize > 0{
		ext["maxFileSize"] = strconv.FormatInt(fileParams.MaxFileSize, 10)
	}
	if fileParams.CanAddUrls {
		ext["canAddUrls"] = "true"
	}
	if len(fileParams.Crop) > 0 {
		// проверка что crop имеет формат 300x400
		arr := strings.Split(fileParams.Crop, "x")
		if len(arr) != 2 {
			log.Fatalf("GetFldImgList error fld: '%s' in FldVueImgParams.Crop must be such format '300x400'. You write this: %s", name, fileParams.Crop)
		}
		if _, err := strconv.Atoi(arr[0]); err != nil {
			log.Fatalf("GetFldImgList error fld: '%s' in FldVueImgParams.Crop must be such format '300x400'. %s not number", name, arr[0])
		}
		if _, err := strconv.Atoi(arr[1]); err != nil {
			log.Fatalf("GetFldImgList error fld: '%s' in FldVueImgParams.Crop must be such format '300x400'. %s not number", name, arr[1])
		}
		ext["crop"] = fileParams.Crop
	}
	if fileParams.Width > 0 {
		ext["width"] = strconv.Itoa(fileParams.Width)
	}

	fld = FldType{Name:name, NameRu:nameRu, Type:FldTypeString, Sql:FldSql{Size:500}, Vue:FldVue{RowCol: rowCol, Type: FldVueTypeImg, Ext: ext, Class: []string{classStr}}}
	return
}

// добавление для таба функциональности счетчика
// добавляется миксин, чтобы в основном табе при открытии загружался список, длина которого и является счетчиком
func (vt VueTab) AddCounter(d *DocType, tabName, pgMethod, pgParams string)  VueTab {
	tabName = utils.UpperCaseFirst(tabName)
	if d.Vue.Mixins == nil {
		d.Vue.Mixins = map[string][]VueMixin{}
	}
	if d.Vue.Mixins["docItemWithTabs"] == nil {
		d.Vue.Mixins["docItemWithTabs"] = []VueMixin{}
	}
	d.Vue.Mixins["docItemWithTabs"] = append(d.Vue.Mixins["docItemWithTabs"], VueMixin{"tabCounter"+tabName, "./mixins/tabCounter"+tabName+".js"})
	sourcePath := "../../../pepelazz/projectGenerator/templates/webClient/doc/mixins/tabCounter.js"
	funcMap := template.FuncMap{
		"VarName": func() string {return "tabCounter"+tabName},
		"PgMethod": func() string {return pgMethod},
		"PgParams": func() string {return pgParams},
	}

	docRouteName := d.Name
	if len(d.Vue.Path) > 0 {
		docRouteName = d.Vue.Path
	}
	distPath := fmt.Sprintf("../src/webClient/src/app/components/%s/mixins", docRouteName)
	d.Templates["webClient_mixin_tabCounter"+tabName+".js"] = &DocTemplate{
		Source: sourcePath,
		DistPath: distPath,
		FuncMap: funcMap,
		DistFilename: "tabCounter"+tabName+".js",
	}
	// добавляем параметры в html разметку таба
	vt.HtmlParams = vt.HtmlParams + " @updateCount='v => tabCounter"+tabName+" = v'"
	vt.HtmlInner = vt.HtmlInner + " <q-badge v-if='tabCounter"+tabName+">0' color='red' floating>{{tabCounter"+tabName+"}}</q-badge>"
	return vt
}
