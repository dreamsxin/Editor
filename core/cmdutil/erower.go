package cmdutil

import (
	"context"
	"io"

	"github.com/jmigpin/editor/core/toolbardata"
	"github.com/jmigpin/editor/ui"
)

// Editor Row interface
type ERower interface {
	Ed() Editorer
	Row() *ui.Row

	LoadContentClear() error
	ReloadContent() error
	SaveContent(string) error

	ToolbarData() *toolbardata.ToolbarData

	IsSpecialName() bool
	Name() string
	Filename() string
	Dir() string
	IsDir() bool
	IsRegular() bool

	Flash()

	StartExecState() context.Context
	StopExecState()
	ClearExecState(context.Context, func())

	TextAreaWriter() io.WriteCloser
	TextAreaAppendAsync(string) <-chan struct{}
	TextAreaStrHashEqual(size int, hash []byte) bool

	UpdateStateAndDuplicates()
}
