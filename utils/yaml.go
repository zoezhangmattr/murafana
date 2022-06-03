package utils

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func LoadToYaml(data interface{}, fileName string) {
	yamlData, err := yaml.Marshal(data)

	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}
	filepath.Join()

	err = ioutil.WriteFile("data/"+fileName, yamlData, 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}
}
