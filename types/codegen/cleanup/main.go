package main

import (
	"github.com/rancher/wrangler/pkg/cleanup"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := cleanup.Cleanup("./types"); err != nil {
		logrus.Fatal(err)
	}
}
