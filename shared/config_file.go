package shared

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ConfigStruct struct {
	Username *string `yaml:"username"`
	Password *string `yaml:"password"`
	Room     *string `yaml:"room"`
	FullName *string `yaml:"full_name"`
}

var (
	Config = &ConfigStruct{}
)

func (config *ConfigStruct) IsConfigured() bool {
	return *config.Username != "" && *config.Password != "" && *config.FullName != "" && *config.Room != ""
}

func (config *ConfigStruct) WriteToFile(path string) {
	data, err := yaml.Marshal(config)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		panic(err)
	}
}

func ParseConfig(configFilePath string, config *ConfigStruct) {
	fileContents, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(fileContents, config)
	if err != nil {
		panic(err)
	}
}

func WriteSampleFile(path string) {
	userName := "307354_1904343"
	password := "password"
	fullName := "John Doe"
	room := "307354_name"

	sampleConfig := &ConfigStruct{
		Username: &userName,
		Password: &password,
		Room:     &room,
		FullName: &fullName,
	}

	sampleConfig.WriteToFile(path)
	fmt.Println("Sample config file written to:", path)

}
