package runtime

import (
	"io"
)

type Context struct {
	Table  *SymbolTable
	Input  io.Reader
	Output io.Writer
}
