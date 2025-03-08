package config

import "fmt"

type wrapper struct {
	str string
}

func (w *wrapper) String() string {
	return w.str
}

func wrap(str string) fmt.Stringer {
	return &wrapper{str: str}
}
