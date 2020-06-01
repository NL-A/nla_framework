package types

import (
	"fmt"
	"github.com/pepelazz/projectGenerator/utils"
	"log"
	"strings"
)

const (
	FldTypeString            = "string"
	FldTypeText              = "text"
	FldTypeInt               = "int"
	FldTypeInt64             = "int64"
	FldTypeDouble            = "double"
	FldTypeDate              = "date"
	FldTypeBool              = "bool"
	FldTypeJsonb             = "jsonb"
	FldTypeVueComposition    = "vueComposition"
	FldTypeDatetime          = "datetime"
	FldTypeTextArray         = "text[]"
	FldTypeIntArray          = "int[]"
	FldVueTypeSelect         = "select"
	FldVueTypeMultipleSelect = "multipleSelect"
	FldVueTypeTags           = "tags"
	FldVueTypeCheckbox       = "checkbox"
	FldVueTypeRadio          = "radio"
	FldVueTypeDadataAddress  = "dadataAddress"
	FldVueTypeJsonList       = "jsonList"
	FldVueTypeFiles          = "files"
	FldVueTypeImgList        = "imgList"
	FldVueTypeImg            = "img"
)

type (
	FldType struct {
		Name   string
		NameRu string
		Type   string
		Vue    FldVue
		Sql    FldSql
		Doc    *DocType // ссылка на сам документ, к которому принадлежит поле
		IntegrationData map[string]interface{} // информация по интеграции с разными системами
	}

	FldVue struct {
		Name         string
		NameRu       string
		Type         string
		RowCol       [][]int
		Class        []string
		IsRequired   bool
		IsBorderless bool
		Readonly     string
		Ext          map[string]string
		Options      []FldVueOptionsItem
		Composition  func(ProjectType, DocType) string
		Vif          string
		JsonList     FldVueJsonList
	}

	FldSql struct {
		IsSearch                 bool
		IsRequired               bool
		Ref                      string
		IsUniq                   bool
		Size                     int
		IsOptionFld              bool // признак что поле пишется не в отдельную колонку таблицы, а в json поле options
		Default                  string
		IsNotUpdatable           bool   // признак, что поле не обновляется вручную. Либо заполняется только при создании, либо обновляется триггером
		FillValueInBeforeTrigger string // строка, которая выполняется в trigger и результат, которой присваивается полю. Например new.fullname
	}

	FldVueOptionsItem struct {
		Label string      `json:"label"`
		Value interface{} `json:"value"`
	}

	FldVueJsonList struct {
		Flds []FldType
		Icon string
	}

	FldVueFilesParams struct {
		Accept      string
		MaxFileSize int64
	}

	FldVueImgParams struct {
		Accept      string
		MaxFileSize int64
		Crop        string // например, 300x400. Обрезает под данное соотношение и размер
		Width       int // обрезает максимальная ширина фото
		CanAddUrls  bool // возможность добавлять ссылки на фото, а не только загружать
	}

)

func (fld *FldType) PrintPgModel() string {
	typeStr := fmt.Sprintf(`type="%s"`, fld.Type)
	extStr := ""
	if fld.Type == "string" {
		if fld.Sql.Size > 0 {
			typeStr = fmt.Sprintf("type=\"char\",\tsize=%v", fld.Sql.Size)
		} else {
			typeStr = `type="text"`
		}
	}
	if fld.Type == FldTypeInt64 {
		typeStr = `type="int"`
	}
	if utils.CheckContainsSliceStr(fld.Type, FldTypeDate, FldTypeDatetime) {
		typeStr = `type="timestamp"`
	}
	if fld.Sql.IsRequired {
		extStr = "not null"
	}
	if len(fld.Sql.Default) > 0 {
		extStr = extStr + " default " + fld.Sql.Default
	}
	// ext может быть пустой
	ext := ""
	if len(extStr) > 0 {
		ext = fmt.Sprintf(" \text=\"%s\",", extStr)
	}
	res := fmt.Sprintf("\t{name=\"%s\",\t\t\t\t\t%s,%s\t comment=\"%s\"}", fld.Name, typeStr, ext, fld.NameRu)

	return res
}

func (fld *FldType) GoType() string {
	switch fld.Type {
	case FldTypeDouble:
		return "float64"
	case FldTypeIntArray:
		return "[]int"
	case FldTypeTextArray:
		return "[]string"
	case FldTypeDate, FldTypeDatetime:
		return "string"
	default:
		return fld.Type
	}
}

func (fld *FldType) PgInsertType() string {
	switch fld.Type {
	case FldTypeDouble:
		return "double precision"
	case FldTypeString:
		return "text"
	case FldTypeDate, FldTypeDatetime:
		return "timestamp"
	case FldTypeInt64:
		return "int"
	default:
		return fld.Type
	}
}

func (fld *FldType) PgUpdateType() string {
	switch fld.Type {
	case FldTypeInt, FldTypeInt64, FldTypeDouble:
		return "number"
	case FldTypeString:
		return "text"
	case FldTypeDate, FldTypeDatetime:
		return "timestamp"
	case FldTypeTextArray:
		return "jsonArrayText"
	case FldTypeIntArray:
		return "jsonArrayInt"
	default:
		return fld.Type
	}
}

// переписываем значение колонки и строки. Третье число - ширина колонки
func (fld FldType) SetRowCol(n ...int) FldType {
	if len(n) < 2 {
		docName := ""
		if fld.Doc != nil {
			docName = fld.Doc.Name
		}
		log.Fatalf("Doc: '%s' Fld: '%s' SetRowCol params must be more tan two numbers. Get %v", docName, fld.Name, n)
	}
	fld.Vue.RowCol = [][]int{{n[0], n[1]}}
	// если указано третье число, то заменяем класс, описыающий ширину колонки
	if len(n) > 2 {
		// копируем значения в новый массив, инчае возможны спецэффекты
		arr := []string{}
		for _, s := range fld.Vue.Class {
			if strings.HasPrefix(s, "col-") {
				s = fmt.Sprintf("col-%v", n[2])
			}
			arr = append(arr, s)
		}
		fld.Vue.Class = arr
	}
	return fld
}

func (fld FldType) SetIsRequired() FldType {
	fld.Sql.IsRequired = true
	fld.Vue.IsRequired = true
	return fld
}

func (fld FldType) SetIsOptionFld() FldType {
	fld.Sql.IsOptionFld = true
	return fld
}
func (fld FldType) SetIsSearch() FldType {
	fld.Sql.IsSearch = true
	return fld
}
func (fld FldType) SetDefault(s string) FldType {
	fld.Sql.Default = s
	return fld
}

func (fld FldType) SetIsNotUpdatable() FldType {
	fld.Sql.IsNotUpdatable = true
	return fld
}

func (fld FldType) SetIsBorderless() FldType {
	fld.Vue.IsBorderless = true
	return fld
}

func (fld FldType) AddClass(s string) FldType {
	if fld.Vue.Class == nil {
		fld.Vue.Class = []string{}
	}
	fld.Vue.Class = append(fld.Vue.Class, s)
	return fld
}

// передается либо true/false, либо функция вида ()=> item !== 'a'
func (fld FldType) SetReadonly(s string) FldType {
	fld.Vue.Readonly = s
	return fld
}

func (fld FldType) SetVif(s string) FldType {
	fld.Vue.Vif = s
	return fld
}

func (fld FldType) SetBitrixInfo(b BitrixFld) FldType {
	if fld.IntegrationData == nil {
		fld.IntegrationData = map[string]interface{}{}
	}
	fld.IntegrationData["bitrix"] = b
	return fld
}

func (fld FldVue) ClassPrint() string {
	if fld.Class != nil {
		return strings.Join(fld.Class, " ")
	}
	return ""
}

func (fld FldVue) ClassPrintOnlyCol() string {
	if fld.Class != nil {
		arr := []string{}
		for _, cName := range fld.Class {
			if strings.HasPrefix(cName, "col-") {
				arr = append(arr, cName)
			}
		}
		return strings.Join(arr, " ")
	}
	return ""
}
