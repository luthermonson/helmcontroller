package main

import (
	"github.com/luthermonson/helmcontroller/types/apis/some.api.group/v1"
	"github.com/rancher/norman/generator"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := generator.DefaultGenerate(v1.Schemas, "github.com/luthermonson/helmcontroller/types", false, nil); err != nil {
		logrus.Fatal(err)
	}
}
