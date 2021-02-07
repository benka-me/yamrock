package main

import (
	"fmt"
	"github.com/benka-me/yamrock"
	"gopkg.in/yaml.v3"
	"os"
)

func main() {
	r, err := os.Open("test.yaml")
	if err != nil {
		panic(err)
	}
	decode := yaml.NewDecoder(r)

	tmp := make(map[interface{}]interface{})
	err = decode.Decode(&tmp)
	if err != nil {
		panic(err)
	}

	data, err := yamrock.Gen(tmp)
	if err != nil {
		panic(err)
	}

	fmt.Println(data)

	//tmp2 := struct {
	//	Config res.ObjConfig `yaml:"config"`
	//}{}
	//err = decode.Decode(&tmp2)
	//if err != nil {
	//	panic(err)
	//}
	//spew.Dump(tmp2)

	//config := res.GetConfig()
	//spew.Dump(config.Sub().Name())
}
