package main

import (
	"github.com/douglasdgoulart/video-editor-api/pkg/api"
)

func main() {
	api.NewApi(":8080").Run()
}
