// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/serenize/snaker"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

var (
	reXenRefType    = regexp.MustCompile("^(.+?) ref$")
	reXenSetType    = regexp.MustCompile("^(.+?) set$")
	reXenRecordType = regexp.MustCompile("^(.+?) record$")
	reXenEnumType   = regexp.MustCompile("^enum (.+)$")
	reXenMapType    = regexp.MustCompile("^\\((.+?) -> (.+?)\\) map$")
)

func goTypeForXenType(xenType string) (goType string, err error) {
	var match []string
	if xenType == "bool" {
		goType = "bool"
	} else if xenType == "int" {
		goType = "int"
	} else if xenType == "float" {
		goType = "float64"
	} else if xenType == "string" {
		goType = "string"
	} else if xenType == "datetime" {
		goType = "time.Time"
	} else if match = reXenSetType.FindStringSubmatch(xenType); match != nil {
		var goItemType string
		goItemType, err = goTypeForXenType(match[1])
		if err != nil {
			return
		}
		goType = "[]" + goItemType
	} else if match = reXenRefType.FindStringSubmatch(xenType); match != nil {
		goType = snaker.SnakeToCamel(match[1]) + "Ref"
	} else if match = reXenRecordType.FindStringSubmatch(xenType); match != nil {
		goType = snaker.SnakeToCamel(match[1]) + "Record"
	} else if match = reXenEnumType.FindStringSubmatch(xenType); match != nil {
		goType = snaker.SnakeToCamel(match[1])
	} else if match = reXenMapType.FindStringSubmatch(xenType); match != nil {
		var goKeyType string
		goKeyType, err = goTypeForXenType(match[1])
		if err != nil {
			return
		}
		var goValueType string
		goValueType, err = goTypeForXenType(match[2])
		if err != nil {
			return
		}
		goType = "map[" + goKeyType + "]" + goValueType
	} else {
		err = fmt.Errorf("Unsupported XenAPI type: %s", xenType)
	}
	return
}

func funcPartialForXenType(xenType string) (partial string, err error) {
	var match []string
	if xenType == "bool" {
		partial = "Bool"
	} else if xenType == "int" {
		partial = "Int"
	} else if xenType == "float" {
		partial = "Float"
	} else if xenType == "string" {
		partial = "String"
	} else if xenType == "datetime" {
		partial = "Time"
	} else if match = reXenSetType.FindStringSubmatch(xenType); match != nil {
		var itemPartial string
		itemPartial, err = funcPartialForXenType(match[1])
		if err != nil {
			return
		}
		partial = itemPartial + "Set"
	} else if match = reXenRefType.FindStringSubmatch(xenType); match != nil {
		partial = snaker.SnakeToCamel(match[1]) + "Ref"
	} else if match = reXenRecordType.FindStringSubmatch(xenType); match != nil {
		partial = snaker.SnakeToCamel(match[1]) + "Record"
	} else if match = reXenEnumType.FindStringSubmatch(xenType); match != nil {
		partial = "Enum" + snaker.SnakeToCamel(match[1])
	} else if match = reXenMapType.FindStringSubmatch(xenType); match != nil {
		var keyPartial string
		keyPartial, err = funcPartialForXenType(match[1])
		if err != nil {
			return
		}
		var valuePartial string
		valuePartial, err = funcPartialForXenType(match[2])
		if err != nil {
			return
		}
		partial = keyPartial + "To" + valuePartial + "Map"
	} else {
		err = fmt.Errorf("Unsupported XenAPI type: %s", xenType)
	}
	return
}

func convertXenTypeFuncName(xenType string, direction string) (funcName string, err error) {
	funcPartial, err := funcPartialForXenType(xenType)
	if err != nil {
		return
	}

	funcName = "convert" + funcPartial + direction
	return
}

var reBeginningOfLine = regexp.MustCompile("(?m)^")

func formatGoDoc(input string) string {
	return reBeginningOfLine.ReplaceAllString(input, "// ")
}

func exportedGoIdentifier(input string) string {
	input = strings.Replace(input, "-", "_", -1)
	return snaker.SnakeToCamel(input)
}

func internalGoIdentifier(input string) (ident string) {
	input = strings.Replace(input, "-", "_", -1)

	// The first component of the name should be all lowercase.
	_index := strings.IndexRune(input, '_')
	if _index == -1 {
		ident = strings.ToLower(input)
	} else {
		ident = strings.ToLower(input[:_index]) + snaker.SnakeToCamel(input[_index+1:])
	}

	// Rename XenAPI identifiers that conflict with Go identifiers.
	switch ident {
	case "type":
		ident = "atype"
	case "interface":
		ident = "iface"
	}

	return
}

func executeTemplateToString(templates *template.Template, name string, data interface{}) (text string, err error) {
	var buf bytes.Buffer

	err = templates.ExecuteTemplate(&buf, name, data)
	if err != nil {
		return
	}

	text = buf.String()
	return
}

type xapiLifecycle struct {
	Description string `json:"description"`
	Release     string `json:"release"`
	Transition  string `json:"transition"`
}

type xapiEnumValue struct {
	Doc  string `json:"doc"`
	Name string `json:"name"`
}

type xapiEnum struct {
	Values []*xapiEnumValue `json:"values"`
	Name   string           `json:"name"`
}

type xapiField struct {
	Default     string           `json:"default,omitempty"`
	Lifecycle   []*xapiLifecycle `json:"lifecycle"`
	Tag         string           `json:"tag"`
	Qualifier   string           `json:"qualifier"`
	Type        string           `json:"type"`
	Description string           `json:"description"`
	Name        string           `json:"name"`
}

func (field *xapiField) GoType() (string, error) {
	return goTypeForXenType(field.Type)
}

type xapiParam struct {
	Doc  string `json:"doc"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (param *xapiParam) GoType() (string, error) {
	return goTypeForXenType(param.Type)
}

type xapiResult []string

func (result *xapiResult) Type() string {
	return (*result)[0]
}

func (result *xapiResult) GoType() (string, error) {
	return goTypeForXenType(result.Type())
}

func (result *xapiResult) IsVoid() bool {
	return result.Type() == "void"
}

type xapiError struct {
	Doc  string `json:"doc"`
	Name string `json:"name"`
}

type xapiMessage struct {
	Implicit    bool             `json:"implicit"`
	Lifecycle   []*xapiLifecycle `json:"lifecycle"`
	Tag         string           `json:"tag"`
	Roles       []string         `json:"roles"`
	Errors      []*xapiError     `json:"errors"`
	Params      []*xapiParam     `json:"params"`
	Result      *xapiResult      `json:"result"`
	Description string           `json:"description"`
	Name        string           `json:"name"`
}

type xapiClass struct {
	Tag         string           `json:"tag"`
	Lifecycle   []*xapiLifecycle `json:"lifecycle"`
	Enums       []*xapiEnum      `json:"enums"`
	Messages    []*xapiMessage   `json:"messages"`
	Fields      []*xapiField     `json:"fields"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
}

const fileHeaderTemplate string = `//
// This file is generated. To change the content of this file, please do not
// apply the change to this file because it will get overwritten. Instead,
// change xenapi.go and execute 'go generate'.
//

package xenAPI

import (
	"fmt"
	"github.com/amfranz/go-xmlrpc-client"
	"reflect"
	"strconv"
	"time"
)

var _ = fmt.Errorf
var _ = xmlrpc.NewClient
var _ = reflect.TypeOf
var _ = strconv.Atoi
var _ = time.UTC
`

const enumTypeTemplate string = `
type {{ .Name|exported }} string

const ({{ range .Values }}
  {{ .Doc|godoc }}
	{{ (printf "%s_%s" $.Name .Name)|exported }} {{ $.Name|exported }} = {{ printf "%q" .Name }}{{ end }}
)
`

const recordTypeTemplate string = `
type {{ .Name|exported }}Record struct {{ "{" }}{{ range .Fields }}
  {{ .Description|godoc }}
	{{ .Name|exported }} {{ .GoType }}{{ end }}
}
`

const classTypeTemplate string = `
{{ .Description|godoc }}
type {{ .Name|exported }}Class struct {
	client *Client
}
`

const refTypeTemplate string = `
type {{ .Name|exported }}Ref string
`

const messageFuncTemplate string = `
{{ .Message.Description|godoc }}{{ if .Message.Errors }}
//
// Errors:{{ range .Message.Errors }}
//  {{ .Name }} - {{ .Doc }}{{ end }}{{ end }}
func (_class {{ .Class.Name|exported }}Class) {{ .Message.Name|exported }}({{ range $index, $param := .Message.Params }}{{ if gt $index 0 }}, {{ end }}{{ .Name|internal }} {{ .GoType }}{{ end }}) ({{ if not .Message.Result.IsVoid }}_retval {{ .Message.Result.GoType }}, {{ end }}_err error) {
	_method := "{{ .Class.Name }}.{{ .Message.Name }}"{{ range .Message.Params }}
	_{{ .Name|internal }}Arg, _err := {{ .Type|convertToXen }}(fmt.Sprintf("%s(%s)", _method, {{ printf "%q" .Name }}), {{ .Name|internal }})
	if _err != nil {
		return
	}{{ end }}
	{{ if .Message.Result.IsVoid }}_, _err = {{ else }}_result, _err :={{ end }} _class.client.APICall(_method{{ range .Message.Params }}, _{{ .Name|internal }}Arg{{ end }}){{ if not .Message.Result.IsVoid }}
	if _err != nil {
		return
	}
	_retval, _err = {{ .Message.Result.Type|convertToGo }}(_method + " -> ", _result.Value){{ end }}
	return
}
`

const clientStructTemplate string = `
type Client struct {
	rpc *xmlrpc.Client{{ range .Classes }}
	{{ .Name|exported }} {{ .Name|exported }}Class{{ end }}
}

func prepClient(rpc *xmlrpc.Client) *Client {
	var client Client
	client.rpc = rpc{{ range .Classes }}
	client.{{ .Name|exported }} = {{ .Name|exported }}Class{&client}{{ end }}
	return &client
}
`

const convertSimpleTypeToGoFuncTemplate string = `
func {{ .FuncName }}(context string, input interface{}) (value {{ .GoType }}, err error) {
	if input == nil {
		return
	}
	value, ok := input.({{ .GoType }})
	if !ok {
		err = fmt.Errorf("Failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", {{ printf "%q" .GoType }}, context, reflect.TypeOf(input), input)
	}
	return
}
`

const convertSimpleTypeToXenFuncTemplate string = `
func {{ .FuncName }}(context string, value {{ .GoType }}) ({{ .GoType }}, error) {
	return value, nil
}
`

const convertIntToGoFuncTemplate string = `
func {{ .FuncName }}(context string, input interface{}) (value int, err error) {
	strValue, ok := input.(string)
	if !ok {
		err = fmt.Errorf("Failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	} else {
  	value, err = strconv.Atoi(strValue)
	}
	return
}
`

const convertIntToXenFuncTemplate string = `
func {{ .FuncName }}(context string, value int) (string, error) {
	return strconv.Itoa(value), nil
}
`

const convertRefTypeToGoFuncTemplate string = `
func {{ .FuncName }}(context string, input interface{}) (ref {{ .GoType }}, err error) {
	value, ok := input.(string)
	if !ok {
		err = fmt.Errorf("Failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "string", context, reflect.TypeOf(input), input)
	} else {
		ref = {{ .GoType }}(value)
	}
	return
}
`

const convertRefTypeToXenFuncTemplate string = `
func {{ .FuncName }}(context string, ref {{ .GoType }}) (string, error) {
	return string(ref), nil
}
`

const convertSetTypeToGoFuncTemplate string = `
func {{ .FuncName }}(context string, input interface{}) (slice {{ .GoType }}, err error) {
	set, ok := input.([]interface{})
	if !ok {
		err = fmt.Errorf("Failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "[]interface{}", context, reflect.TypeOf(input), input)
		return
	}
	slice = make({{ .GoType }}, len(set))
	for index, item := range set {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := {{ .ItemConverter }}(itemContext, item)
		if err != nil {
			return slice, err
		}
		slice[index] = itemValue
	}
	return
}
`

const convertSetTypeToXenFuncTemplate string = `
func {{ .FuncName }}(context string, slice {{ .GoType }}) (set []interface{}, err error) {
	set = make([]interface{}, len(slice))
	for index, item := range slice {
		itemContext := fmt.Sprintf("%s[%d]", context, index)
		itemValue, err := {{ .ItemConverter }}(itemContext, item)
		if err != nil {
			return set, err
		}
		set[index] = itemValue
	}
	return
}
`

const convertRecordTypeToGoFuncTemplate string = `
func {{ .FuncName }}(context string, input interface{}) (record {{ .GoType }}, err error) {
	rpcStruct, ok := input.(xmlrpc.Struct)
	if !ok {
		err = fmt.Errorf("Failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "xmlrpc.Struct", context, reflect.TypeOf(input), input)
		return
	}{{ range .Fields }}
  {{ .Name|internal }}Value, ok := rpcStruct[{{ printf "%q" .Name }}]
	if ok && {{ .Name|internal }}Value != nil {
  	record.{{ .Name|exported }}, err = {{ .Type|convertToGo }}(fmt.Sprintf("%s.%s", context, {{ printf "%q" .Name }}), {{ .Name|internal }}Value)
		if err != nil {
			return
		}
	}{{ end }}
	return
}
`

const convertRecordTypeToXenFuncTemplate string = `
func {{ .FuncName }}(context string, record {{ .GoType }}) (rpcStruct xmlrpc.Struct, err error) {{ "{\n  rpcStruct = xmlrpc.Struct{}" }}{{ range .Fields }}
  rpcStruct[{{ printf "%q" .Name }}], err = {{ .Type|convertToXen }}(fmt.Sprintf("%s.%s", context, {{ printf "%q" .Name }}), record.{{ .Name|exported }})
  if err != nil {
		return
	}{{ end }}
	return
}
`

const convertMapTypeToGoFuncTemplate string = `
func {{ .FuncName }}(context string, input interface{}) (goMap {{ .GoType }}, err error) {
	xenMap, ok := input.(xmlrpc.Struct)
	if !ok {
		err = fmt.Errorf("Failed to parse XenAPI response: expected Go type %s at %s but got Go type %s with value %v", "xmlrpc.Struct", context, reflect.TypeOf(input), input)
		return
	}
	goMap = make({{ .GoType }}, len(xenMap))
	for xenKey, xenValue := range xenMap {
		keyContext := fmt.Sprintf("%s[%s]", context, xenKey)
		goKey, err := {{ .KeyConverter }}(keyContext, xenKey)
		if err != nil {
			return goMap, err
		}
		goValue, err := {{ .ValueConverter }}(keyContext, xenValue)
		if err != nil {
			return goMap, err
		}
		goMap[goKey] = goValue
	}
	return
}
`

const convertMapTypeToXenFuncTemplate string = `
func {{ .FuncName }}(context string, goMap {{.GoType }}) (xenMap xmlrpc.Struct, err error) {
	xenMap = make(xmlrpc.Struct)
	for goKey, goValue := range goMap {
		keyContext := fmt.Sprintf("%s[%s]", context, goKey)
		xenKey, err := {{ .KeyConverter }}(keyContext, goKey)
		if err != nil {
			return xenMap, err
		}
		xenValue, err := {{ .ValueConverter }}(keyContext, goValue)
		if err != nil {
			return xenMap, err
		}
		xenMap[xenKey] = xenValue
	}
	return
}
`

const convertEnumTypeToGoFuncTemplate string = `
func {{ .FuncName }}(context string, input interface{}) (value {{ .GoType }}, err error) {
	strValue, err := {{ "string"|convertToGo }}(context, input)
	if err != nil {
		return
	}
  switch strValue {{ "{" }}{{ range .Values }}
    case {{ printf "%q" .Name }}:
      value = {{ $.GoType }}{{ .Name|exported }}{{ end }}
    default:
      err = fmt.Errorf("Unable to parse XenAPI response: got value %q for enum %s at %s, but this is not any of the known values", strValue, {{ printf "%q" .GoType }}, context)
	}
	return
}
`

const convertEnumTypeToXenFuncTemplate string = `
func {{ .FuncName }}(context string, value {{ .GoType }}) (string, error) {
	return string(value), nil
}
`

type converterFunc struct {
	name       string
	definition string
}

type apiGenerator struct {
	classes    []*xapiClass
	templates  *template.Template
	converters map[string]converterFunc
}

func newAPIGenerator() apiGenerator {
	return apiGenerator{
		converters: make(map[string]converterFunc),
	}
}

func (generator *apiGenerator) loadXenAPI(filename string) (err error) {
	xenAPI, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	return json.Unmarshal(xenAPI, &generator.classes)
}

func (generator *apiGenerator) prepTemplates() (err error) {
	generator.templates = template.New("")

	generator.templates.Funcs(template.FuncMap{
		"godoc":    formatGoDoc,
		"exported": exportedGoIdentifier,
		"internal": internalGoIdentifier,
		"convertToGo": func(xenType string) (string, error) {
			converter, err := generator.getOrCreateConverterFunc(xenType, "ToGo")
			if err != nil {
				return "", err
			}
			return converter.name, nil
		},
		"convertToXen": func(xenType string) (string, error) {
			converter, err := generator.getOrCreateConverterFunc(xenType, "ToXen")
			if err != nil {
				return "", err
			}
			return converter.name, nil
		},
	})

	templateLedger := map[string]string{
		"FileHeader":                 fileHeaderTemplate,
		"EnumType":                   enumTypeTemplate,
		"RecordType":                 recordTypeTemplate,
		"ClassType":                  classTypeTemplate,
		"RefType":                    refTypeTemplate,
		"MessageFunc":                messageFuncTemplate,
		"ClientStruct":               clientStructTemplate,
		"convertSimpleTypeToGoFunc":  convertSimpleTypeToGoFuncTemplate,
		"convertSimpleTypeToXenFunc": convertSimpleTypeToXenFuncTemplate,
		"convertIntToGoFunc":         convertIntToGoFuncTemplate,
		"convertIntToXenFunc":        convertIntToXenFuncTemplate,
		"convertRefTypeToGoFunc":     convertRefTypeToGoFuncTemplate,
		"convertRefTypeToXenFunc":    convertRefTypeToXenFuncTemplate,
		"convertSetTypeToGoFunc":     convertSetTypeToGoFuncTemplate,
		"convertSetTypeToXenFunc":    convertSetTypeToXenFuncTemplate,
		"convertRecordTypeToGoFunc":  convertRecordTypeToGoFuncTemplate,
		"convertRecordTypeToXenFunc": convertRecordTypeToXenFuncTemplate,
		"convertMapTypeToGoFunc":     convertMapTypeToGoFuncTemplate,
		"convertMapTypeToXenFunc":    convertMapTypeToXenFuncTemplate,
		"convertEnumTypeToGoFunc":    convertEnumTypeToGoFuncTemplate,
		"convertEnumTypeToXenFunc":   convertEnumTypeToXenFuncTemplate,
	}

	for name, value := range templateLedger {
		_, err = generator.templates.New(name).Parse(value)
		if err != nil {
			return
		}
	}

	return
}

func (generator *apiGenerator) buildSimpleConverterFunc(xenType string, direction string, funcName string, goType string) (string, error) {
	args := map[string]interface{}{
		"FuncName": funcName,
		"GoType":   goType,
	}

	return executeTemplateToString(generator.templates, "convertSimpleType"+direction+"Func", args)
}

func (generator *apiGenerator) buildIntConverterFunc(xenType string, direction string, funcName string) (string, error) {
	args := map[string]interface{}{
		"FuncName": funcName,
	}

	return executeTemplateToString(generator.templates, "convertInt"+direction+"Func", args)
}

func (generator *apiGenerator) buildRefConverterFunc(xenType string, direction string, funcName string, baseType string) (string, error) {
	goType, err := goTypeForXenType(xenType)
	if err != nil {
		return "", err
	}

	args := map[string]interface{}{
		"FuncName": funcName,
		"GoType":   goType,
	}

	return executeTemplateToString(generator.templates, "convertRefType"+direction+"Func", args)
}

func (generator *apiGenerator) buildSetConverterFunc(xenType string, direction string, funcName string, itemType string) (string, error) {
	goType, err := goTypeForXenType(xenType)
	if err != nil {
		return "", err
	}

	itemConverter, err := generator.getOrCreateConverterFunc(itemType, direction)
	if err != nil {
		return "", err
	}

	args := map[string]interface{}{
		"FuncName":      funcName,
		"GoType":        goType,
		"ItemConverter": itemConverter.name,
	}

	return executeTemplateToString(generator.templates, "convertSetType"+direction+"Func", args)
}

func (generator *apiGenerator) buildRecordConverterFunc(xenType string, direction string, funcName string, itemType string) (string, error) {
	goType, err := goTypeForXenType(xenType)
	if err != nil {
		return "", err
	}

	var fields []*xapiField
	for _, class := range generator.classes {
		if class.Name+" record" == xenType {
			fields = class.Fields
			break
		}
	}
	if len(fields) == 0 {
		return "", fmt.Errorf("Unable to find definition for XenAPI %s", xenType)
	}

	args := map[string]interface{}{
		"FuncName": funcName,
		"GoType":   goType,
		"Fields":   fields,
	}

	return executeTemplateToString(generator.templates, "convertRecordType"+direction+"Func", args)
}

func (generator *apiGenerator) buildMapConverterFunc(xenType string, direction string, funcName string, keyType string, valueType string) (string, error) {
	goType, err := goTypeForXenType(xenType)
	if err != nil {
		return "", err
	}

	keyConverter, err := generator.getOrCreateConverterFunc(keyType, direction)
	if err != nil {
		return "", err
	}

	valueConverter, err := generator.getOrCreateConverterFunc(valueType, direction)
	if err != nil {
		return "", err
	}

	args := map[string]interface{}{
		"FuncName":       funcName,
		"GoType":         goType,
		"KeyConverter":   keyConverter.name,
		"ValueConverter": valueConverter.name,
	}

	return executeTemplateToString(generator.templates, "convertMapType"+direction+"Func", args)
}

func (generator *apiGenerator) buildEnumConverterFunc(xenType string, direction string, funcName string, enumType string) (string, error) {
	goType, err := goTypeForXenType(xenType)
	if err != nil {
		return "", err
	}

	var values []*xapiEnumValue
classLoop:
	for _, class := range generator.classes {
		for _, enum := range class.Enums {
			if "enum "+enum.Name == xenType {
				values = enum.Values
				break classLoop
			}
		}
	}
	if len(values) == 0 {
		return "", fmt.Errorf("Unable to find definition for XenAPI %s", xenType)
	}

	args := map[string]interface{}{
		"FuncName": funcName,
		"GoType":   goType,
		"Values":   values,
	}

	return executeTemplateToString(generator.templates, "convertEnumType"+direction+"Func", args)
}

func (generator *apiGenerator) buildConverterFunc(xenType string, direction string) (converter converterFunc, err error) {
	funcName, err := convertXenTypeFuncName(xenType, direction)
	if err != nil {
		return
	}

	var funcDefinition string
	if xenType == "string" {
		funcDefinition, err = generator.buildSimpleConverterFunc(xenType, direction, funcName, "string")
	} else if xenType == "bool" {
		funcDefinition, err = generator.buildSimpleConverterFunc(xenType, direction, funcName, "bool")
	} else if xenType == "int" {
		funcDefinition, err = generator.buildIntConverterFunc(xenType, direction, funcName)
	} else if xenType == "float" {
		funcDefinition, err = generator.buildSimpleConverterFunc(xenType, direction, funcName, "float64")
	} else if xenType == "datetime" {
		funcDefinition, err = generator.buildSimpleConverterFunc(xenType, direction, funcName, "time.Time")
	} else if match := reXenRefType.FindStringSubmatch(xenType); match != nil {
		funcDefinition, err = generator.buildRefConverterFunc(xenType, direction, funcName, match[1])
	} else if match := reXenSetType.FindStringSubmatch(xenType); match != nil {
		funcDefinition, err = generator.buildSetConverterFunc(xenType, direction, funcName, match[1])
	} else if match := reXenRecordType.FindStringSubmatch(xenType); match != nil {
		funcDefinition, err = generator.buildRecordConverterFunc(xenType, direction, funcName, match[1])
	} else if match := reXenMapType.FindStringSubmatch(xenType); match != nil {
		funcDefinition, err = generator.buildMapConverterFunc(xenType, direction, funcName, match[1], match[2])
	} else if match := reXenEnumType.FindStringSubmatch(xenType); match != nil {
		funcDefinition, err = generator.buildEnumConverterFunc(xenType, direction, funcName, match[1])
	} else {
		err = fmt.Errorf("Unable to build type conversion function for XenAPI: unsupported type %q", xenType)
	}
	if err != nil {
		return
	}

	converter = converterFunc{funcName, funcDefinition}
	return
}

func (generator *apiGenerator) getOrCreateConverterFunc(xenType string, direction string) (converter converterFunc, err error) {
	converterKey := xenType + direction
	converter, found := generator.converters[converterKey]
	if !found {
		converter, err = generator.buildConverterFunc(xenType, direction)
		if err != nil {
			return
		}
		generator.converters[converterKey] = converter
	}
	return
}

func (generator *apiGenerator) generateClassAPI(class *xapiClass) (err error) {
	apiFilename := fmt.Sprintf("%s_gen.go", strings.ToLower(class.Name))

	fileHandle, err := os.Create(apiFilename)
	if err != nil {
		return
	}

	defer fileHandle.Close()

	err = generator.templates.ExecuteTemplate(fileHandle, "FileHeader", nil)
	if err != nil {
		return
	}

	for _, enum := range class.Enums {
		err = generator.templates.ExecuteTemplate(fileHandle, "EnumType", enum)
		if err != nil {
			return
		}
	}

	if len(class.Fields) > 0 {
		err = generator.templates.ExecuteTemplate(fileHandle, "RecordType", class)
		if err != nil {
			return
		}
	}

	err = generator.templates.ExecuteTemplate(fileHandle, "RefType", class)
	if err != nil {
		return
	}

	err = generator.templates.ExecuteTemplate(fileHandle, "ClassType", class)
	if err != nil {
		return
	}

	for _, message := range class.Messages {

		context := map[string]interface{}{
			"Class":   class,
			"Message": message,
		}

		err = generator.templates.ExecuteTemplate(fileHandle, "MessageFunc", context)
		if err != nil {
			return
		}

	}

	return
}

func (generator *apiGenerator) generateConverters() (err error) {
	fileHandle, err := os.Create("convert_gen.go")
	if err != nil {
		return
	}

	defer fileHandle.Close()

	err = generator.templates.ExecuteTemplate(fileHandle, "FileHeader", nil)
	if err != nil {
		return
	}

	var keys []string
	for key := range generator.converters {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		converter := generator.converters[key]
		_, err = fileHandle.WriteString(converter.definition)
		if err != nil {
			return
		}
	}

	return
}

func (generator *apiGenerator) generateClient() (err error) {
	fileHandle, err := os.Create("client_gen.go")
	if err != nil {
		return
	}

	defer fileHandle.Close()

	err = generator.templates.ExecuteTemplate(fileHandle, "FileHeader", nil)
	if err != nil {
		return
	}

	err = generator.templates.ExecuteTemplate(fileHandle, "ClientStruct", map[string]interface{}{
		"Classes": generator.classes,
	})

	return
}

func (generator *apiGenerator) run() (err error) {
	err = generator.loadXenAPI("xenapi.json")
	if err != nil {
		return
	}

	err = generator.prepTemplates()
	if err != nil {
		return
	}

	for _, class := range generator.classes {
		err = generator.generateClassAPI(class)
		if err != nil {
			return
		}
	}

	err = generator.generateConverters()
	if err != nil {
		return
	}

	err = generator.generateClient()
	return
}

func main() {
	generator := newAPIGenerator()
	err := generator.run()
	if err != nil {
		panic(err)
	}
}
