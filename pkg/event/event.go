package event

import "github.com/douglasdgoulart/video-editor-api/pkg/editor"

type Event struct {
	Id            string               `json:"id"`
	EditorRequest editor.EditorRequest `json:"editor_request"`
}
