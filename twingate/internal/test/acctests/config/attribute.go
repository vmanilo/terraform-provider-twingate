package config

import (
	"fmt"
	"strings"
)

type Attribute interface {
	Name() string
	Value() string
	String() string
}

func NewAttribute(name string, value any) Attribute {
	var val string

	switch value.(type) {
	case bool, int:
		val = fmt.Sprintf("%v", val)
	case string:
		val = value.(string)
	}

	return &BaseAttribute{
		name:  name,
		value: val,
	}
}

type BaseAttribute struct {
	name  string
	value string
}

func (s BaseAttribute) Name() string {
	return s.name
}

func (s BaseAttribute) Value() string {
	return s.value
}

func (s BaseAttribute) String() string {
	return s.value
}

func NewStringAttribute(name string, value any) Attribute {
	if value == nil {
		return &BaseAttribute{
			name:  name,
			value: "null",
		}
	}

	return &StringAttribute{
		BaseAttribute: BaseAttribute{
			name:  name,
			value: value.(string),
		},
	}
}

type StringAttribute struct {
	BaseAttribute
}

func (s StringAttribute) String() string {
	return fmt.Sprintf(`"%s"`, s.value)
}

type NullAttribute struct {
	BaseAttribute
}

func (s NullAttribute) String() string {
	return "null"
}

func NewSetStringAttribute(name string, values any) Attribute {
	return &SetStringAttribute{
		name:   name,
		values: values.([]string),
	}
}

type SetStringAttribute struct {
	name   string
	values []string
}

func (s SetStringAttribute) Name() string {
	return s.name
}

func (s SetStringAttribute) Value() string {
	if len(s.values) == 0 {
		return "[]"
	}

	return `["` + strings.Join(s.values, `", "`) + `"]`
}

func (s SetStringAttribute) String() string {
	return s.Value()
}

func NewSetAttribute(name string, values []string) Attribute {
	return &SetAttribute{
		name:   name,
		values: values,
	}
}

type SetAttribute struct {
	name   string
	values []string
}

func (s SetAttribute) Name() string {
	return s.name
}

func (s SetAttribute) Value() string {
	if len(s.values) == 0 {
		return "[]"
	}

	return `[` + strings.Join(s.values, `, `) + `]`
}

func (s SetAttribute) String() string {
	return s.Value()
}
