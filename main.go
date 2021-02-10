package main

import (
	"fmt"
	"github.com/benka-me/yamrock/yamrock"
	"gopkg.in/yaml.v3"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("usage: yamrock source.yaml ./folder/destination package_name")
		return
	}
	src := os.Args[1]
	destDir := os.Args[2]
	packageName := os.Args[3]
	dest, err := os.Create(destDir + "/" + packageName + ".go")
	if err != nil {
		panic(err)
	}

	defer dest.Close()

	r, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	decode := yaml.NewDecoder(r)

	tmp := make(map[interface{}]interface{})
	err = decode.Decode(&tmp)
	if err != nil {
		panic(err)
	}

	data, err := yamrock.Gen(tmp, packageName)
	if err != nil {
		panic(err)
	}

	_, err = dest.WriteString(data)
	if err != nil {
		panic(err)
	}
}
