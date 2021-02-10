package yamrock

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type Children struct {
	Name    string
	Minus   string
	Type    string
	Content interface{}
	object  *object
}

type object struct {
	Name     string
	Minus    string
	Type     string
	Children []Children
}

func Gen(req map[interface{}]interface{}, pkg string) (string, error) {
	g := fmt.Sprintf("package %s\n\n\n", pkg)

	res, err := gen(req)
	if err != nil {
		return "", err
	}

	return g + res, nil
}

func gen(req map[interface{}]interface{}) (string, error) {
	data := ""

	for name, content := range req {
		Name := snakeCaseToCamelCase(fmt.Sprintf("%v", name))
		data = data + fmt.Sprintf("func New%s() %s {\n\treturn var%s\n}\n", Name, Name, Name)
		switch content.(type) {
		case map[interface{}]interface{}:
			res, err := genDataSet(name.(string), content)
			if err != nil {
				panic(err)
			}
			data = data + res
		case map[string]interface{}:
			res, err := genDataSet(name.(string), content)
			if err != nil {
				panic(err)
			}
			data = data + res
		}
	}

	return data, nil
}

func genDataSet(parentName string, parentContent interface{}) (string, error) {
	gen := ""
	parent, err := prepareDataSetRecursive(parentName, parentContent)
	if err != nil {
		return "", err
	}

	for i := len(parent); i != 0; i-- {
		res, err := generateDataSetContent(parent[i-1])
		if err != nil {
			return "", err
		}
		gen = gen + res
	}
	return gen, nil
}

func prepareDataSetRecursive(Name interface{}, content interface{}) ([]object, error) {
	name := Name.(string)
	ret := []object{}
	parent := object{}

	switch content := content.(type) {
	case map[string]interface{}:
		parent = object{
			Name:  strings.Title(name),
			Minus: name,
			Type:  strings.Title(name),
		}

		var object *object
		for childName, child := range content {
			childType, isobject := getType(childName, child)
			if isobject {
				tmp, err := prepareDataSetRecursive(childName, child)
				if err != nil {
					panic(err)
				}
				ret = append(ret, tmp...)
			}
			parent.Children = append(parent.Children, Children{
				Name:    snakeCaseToCamelCase(childName),
				Minus:   childName,
				Type:    childType,
				Content: child,
				object:  object,
			})
		}
		ret = append(ret, parent)
	}
	return ret, nil
}

func getType(name string, target interface{}) (string, bool) {
	switch target.(type) {
	case map[string]interface{}:
		return snakeCaseToCamelCase(name), true
	case map[interface{}]interface{}:
		return snakeCaseToCamelCase(name), true
	default:
		return reflect.TypeOf(target).String(), false
	}
}

func generateDataSetContent(obj object) (string, error) {
	data := fmt.Sprintf("// %s\n", obj.Name)

	// interface
	data = data + fmt.Sprintf("type %s interface {\n", obj.Name)
	for _, child := range obj.Children {
		data = data + fmt.Sprintf("\t%s() %s\n", child.Name, child.Type)
	}
	data = data + fmt.Sprintf("}\n")

	// struct
	data = data + fmt.Sprintf("type obj%s struct {\n", obj.Name)
	for _, child := range obj.Children {
		data = data + fmt.Sprintf("\tvar%-10s %-30s `yaml:\"%s\"`\n", child.Name, child.Type, child.Minus)
	}
	data = data + fmt.Sprintf("}\n")

	// var local
	data = data + fmt.Sprintf("var var%s = &obj%s{\n", obj.Name, obj.Name)
	for _, child := range obj.Children {
		data = data + fmt.Sprintf("\tvar%s: %s\n", child.Name, fmtVarType(child))
	}
	data = data + fmt.Sprintf("}\n")

	// func
	for _, child := range obj.Children {
		data = data + fmt.Sprintf("func (m *obj%s) %s() %s {\n\treturn m.var%s\n}\n", obj.Name, child.Name, child.Type, child.Name)
	}

	return data + "\n", nil
}

func fmtVarType(child Children) string {
	switch content := child.Content.(type) {
	case string:
		return fmt.Sprintf("\"%s\",", content)
	case []string:
		val := " []string{\n"
		for _, s := range content {
			val = val + fmt.Sprintf("\t\t\"%s\",\n", s)
		}
		val = val + "},"
		return val
	case []interface{}:
		val := " []interface{}{\n"
		for _, s := range content {
			val = val + fmt.Sprintf("\t\t\"%s\",\n", s)
		}
		val = val + "\t},"
		return val
	case int:
		return strconv.FormatInt(int64(content), 10) + ","
	case int8:
		return strconv.FormatInt(int64(content), 10) + ","
	case int16:
		return strconv.FormatInt(int64(content), 10) + ","
	case int32:
		return strconv.FormatInt(int64(content), 10) + ","
	case float32:
		return fmt.Sprintf("%f", content) + ","
	case float64:
		return fmt.Sprintf("%f", content) + ","
	case int64:
		return strconv.FormatInt(content, 10) + ","
	case map[string]interface{}:
		return fmt.Sprintf("var%s,", child.Name)
	default:
		log.Println(child.Name, "is type : ", reflect.TypeOf(child.Content))
	}
	return "error"
}
