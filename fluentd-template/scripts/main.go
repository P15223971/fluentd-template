package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func panicError(err error) {
	if err != nil {	panic(err) }
}

type Config struct {
	LogType 			logType 			`yaml:"logType"`
	RecordTransformer 	recordTransformer	`yaml:"recordTransformer"`
	Fields				[]fields			`yaml:"fields"`
	MultiFormat 		[]formats			`yaml:"multiFormat"`
}

type logType struct {
	Source     		string 	`yaml:"source"`
	System     		string 	`yaml:"system"`
	Path       		string 	`yaml:"path"`
	PosPath    		string 	`yaml:"posPath"`
	Tag				string	`yaml:"tag"`
	TimeFormat 		string 	`yaml:"timeFormat"`
	ParseType  		string 	`yaml:"parseType"`
	FormatFirstLine string 	`yaml:"formatFirstLine"`
}

type recordTransformer struct {
	RemoveKeys 		[]string 			`yaml:"removeKeys"`
	ModifyFields 	[]fieldToModify 	`yaml:"modifyFields"`
}

type fields struct {
	Name 		string `yaml:"name"`
	Regex 		string `yaml:"regex"`
	Delimiter 	string `yaml:"delimiter"`
}

type fieldToModify struct {
	Name 	string	`yaml:"field"`
	Action	string 	`yaml:"modify"`
}

type formats struct {
	Type 	string 		`yaml:"type"`
	Fields 	[]fields 	`yaml:"fields"`
}

func (c *Config) ParseConfiguration(source string) *Config {
	f, err := ioutil.ReadFile(source)

	if err != nil {
		fmt.Errorf("Error opening file: %v", err)
	}

	if err := yaml.Unmarshal(f, c); err != nil {
		fmt.Errorf("Error parsing configuration file: %v", err)
	}

	return c
}

func ListConfigs(dir string) ([]string, error){
	var configs []string
	f, err := os.Open(dir)
	if err != nil {
		return nil, fmt.Errorf("Error opening directory: %v", err)
	}
	fileInfo, err := f.Readdir(-1)
	err = f.Close()
	if err != nil {
		return nil, err
	}
	for _, file := range fileInfo {
		configs = append(configs, file.Name())
	}
	return configs, nil
}

func GenTemplate(file string, config Config) {
	t, err := template.ParseFiles("./templates/template.yml")
	panicError(err)
	s := strings.Split(file, ".")
	file = s[0]
	f, err := os.Create(fmt.Sprintf("./fluentd-conf/%v.conf", file))
	panicError(err)

	err = t.Execute(f, config)
	panicError(err)
}

func main() {

	confList, err := ListConfigs("../log-configs/")
	if err != nil {
		fmt.Errorf("Failed to list files in directory: %v", err)
	}
	fmt.Println(confList)
	for _, config := range confList {
		var c Config
		c.ParseConfiguration("../log-configs/"+config)
		if err != nil {
			fmt.Errorf("Failed to parse config: %v", err)
		}
		GenTemplate(config, c)
	}
}
