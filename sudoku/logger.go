package sudoku

import (
	"fmt"
	"strings"
)

type NamedContext interface {
	Name() string
}

type StringContext string

func (s StringContext) Name() string {
	return string(s)
}

type Logger interface {
	UpdateCell(loc CellLocation, old, new Digits)
	EnterContext(n NamedContext)
	ExitContext()
}

type voidLogger struct{}

func (v voidLogger) UpdateCell(loc CellLocation, old, new Digits) {
}

func (v voidLogger) EnterContext(n NamedContext) {
}

func (v voidLogger) ExitContext() {
}

type consoleLogger struct {
	context             []NamedContext
	printedContextLevel int
}

func NewLogger() Logger {
	return &consoleLogger{}
}

func (l *consoleLogger) UpdateCell(loc CellLocation, old, new Digits) {
	if old == new {
		return
	}
	l.printContextInfo()
	if v, isSingle := new.Single(); isSingle {
		l.printf("solved cell %v to %v", loc, v)
	} else {
		l.printf("removed candidates %s for cell %v", old&^new, loc)
	}
}

func (l *consoleLogger) EnterContext(m NamedContext) {
	l.context = append(l.context, m)
}

func (l *consoleLogger) ExitContext() {
	l.context = l.context[:len(l.context)-1]
	if l.printedContextLevel > len(l.context) {
		l.printedContextLevel--
		l.printf("}")
	}
}

func (l *consoleLogger) printContextInfo() {
	for i := l.printedContextLevel; i < len(l.context); i++ {
		l.printf("%s {", l.context[i].Name())
		l.printedContextLevel++
	}
}

func (l *consoleLogger) printf(format string, args ...interface{}) {
	fmt.Printf(strings.Repeat("  ", l.printedContextLevel)+format+"\n", args...)
}
