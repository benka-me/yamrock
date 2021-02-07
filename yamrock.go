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
	Object  *Object
}

type Object struct {
	Name     string
	Minus    string
	Type     string
	Children []Children
}

func Gen(req map[interface{}]interface{}) (string, error) {
	data := ""
	for parentName, parentContent := range req {
		parent, err := prepareObject(parentName, parentContent)
		if err != nil {
			return "", err
		}

		parentGen, err := parent.generateDataSet()
		if err != nil {
			return "", nil
		}
		data = data + parentGen
	}

	return data, nil
}

func prepareObject(parentName interface{}, parentContent interface{}) (Object, error) {
	parent := Object{}
	switch content := parentContent.(type) {
	case map[string]interface{}:
		parent = Object{
			Name:  strings.Title(parentName.(string)),
			Minus: parentName.(string),
			Type:  strings.Title(parentName.(string)),
		}

		for childName, child := range content {
			childType := getType(childName, child)
			parent.Children = append(parent.Children, Children{
				Name:    snakeCaseToCamelCase(childName),
				Minus:   childName,
				Type:    childType,
				Content: child,
				Object:  nil,
			})
		}
	}
	return parent, nil
}

func getType(name string, target interface{}) string {
	switch target.(type) {
	case map[string]interface{}:
		return snakeCaseToCamelCase(name)
	case map[interface{}]interface{}:
		return snakeCaseToCamelCase(name)
	default:
		return reflect.TypeOf(target).String()
	}
}

func (parent Object) generateDataSet() (string, error) {
	data := ""

	// interface
	data = data + fmt.Sprintf("type %s interface {\n", parent.Name)
	for _, child := range parent.Children {
		data = data + fmt.Sprintf("\t%s() %s\n", child.Name, child.Type)
	}
	data = data + fmt.Sprintf("}\n")

	// struct
	data = data + fmt.Sprintf("type Obj%s struct {\n", parent.Name)
	for _, child := range parent.Children {
		data = data + fmt.Sprintf("\tVar%-10s %-30s `yaml:\"%s\"`\n", child.Name, child.Type, child.Minus)
	}
	data = data + fmt.Sprintf("}\n")

	// var local
	data = data + fmt.Sprintf("var var%s = Obj%s{\n", parent.Name, parent.Name)
	for _, child := range parent.Children {
		data = data + fmt.Sprintf("\tVar%s: %s\n", child.Name, fmtVar(child))
	}
	data = data + fmt.Sprintf("}\n")

	return data, nil
}

func fmtVar(child Children) string {
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
	case int64:
		return strconv.FormatInt(content, 10) + ","
	case map[string]interface{}:
		return fmt.Sprintf("var%s,", child.Name)
	default:
		log.Println(child.Name, "is type : ", reflect.TypeOf(child.Content))
	}
	return "error"
}
