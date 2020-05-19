package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestGetDefaultProfileYaml(t *testing.T) {
	ob, _ := yaml.Marshal(DefaultProfile())
	fmt.Println(fmt.Sprint(string(ob)))
}
