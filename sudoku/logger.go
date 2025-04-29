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

type Logger[D Digits] interface {
	UpdateCell(loc CellLocation, old, new D)
	EnterContext(n NamedContext)
	ExitContext()
}

type voidLogger[D Digits] struct{}

func (v voidLogger[D]) UpdateCell(loc CellLocation, old, new D) {
}

func (v voidLogger[D]) EnterContext(n NamedContext) {
}

func (v voidLogger[D]) ExitContext() {
}

type consoleLogger[D Digits] struct {
	context             []NamedContext
	printedContextLevel int
}

func NewLogger[D Digits]() Logger[D] {
	return &consoleLogger[D]{}
}

func (l *consoleLogger[D]) UpdateCell(loc CellLocation, before, after D) {
	if before == after {
		return
	}
	l.printContextInfo()
	if v, isSingle := after.Single(); isSingle {
		l.printf("solved cell %v to %v", loc, v)
	} else {
		//removed := before &^ after
		removed := new(D)
		l.printf("removed candidates %s for cell %v", removed, loc)
	}
}

func (l *consoleLogger[D]) EnterContext(m NamedContext) {
	l.context = append(l.context, m)
}

func (l *consoleLogger[D]) ExitContext() {
	l.context = l.context[:len(l.context)-1]
	if l.printedContextLevel > len(l.context) {
		l.printedContextLevel--
		l.printf("}")
	}
}

func (l *consoleLogger[D]) printContextInfo() {
	for i := l.printedContextLevel; i < len(l.context); i++ {
		l.printf("%s {", l.context[i].Name())
		l.printedContextLevel++
	}
}

func (l *consoleLogger[D]) printf(format string, args ...interface{}) {
	fmt.Printf(strings.Repeat("  ", l.printedContextLevel)+format+"\n", args...)
}
