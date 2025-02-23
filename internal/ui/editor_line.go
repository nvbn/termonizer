package ui

type editorLine struct {
	start         int
	end           int
	editorContent string // full editorContent of the editor
}

func (e *editorLine) String() string {
	return e.editorContent[e.start:e.end]
}
