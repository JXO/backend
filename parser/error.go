package parser

import (
	"fmt"
)

type Error interface {
	Line() int
	Column() int
	Description() string
	Error() string
}

type BasicError struct {
	line        int
	column      int
	description string
}

func NewError(line, column int, description string) Error {
	return &BasicError{line, column, description}
}
func (be *BasicError) Error() string {
	return fmt.Sprintf("%d,%d: %s", be.line, be.column, be.description)
}
func (be *BasicError) Line() int           { return be.line }
func (be *BasicError) Column() int         { return be.column }
func (be *BasicError) Description() string { return be.description }
