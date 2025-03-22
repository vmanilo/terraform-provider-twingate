package config

import (
	"fmt"
	"strings"
)

type Resource interface {
	ResourceName() string
	ResourceType() string
	TerraformResource() string
	TerraformResourceID() string
	String() string
	Set(values ...any) Resource
	Delete(attributes ...string) Resource
	SetResourceName(name string)
}

type ProtoResource struct {
	Name string
	Type string

	Required map[string]Attribute
	Optional map[string]Attribute
}

func (p ProtoResource) ResourceName() string {
	return p.Name
}

func (p ProtoResource) SetResourceName(name string) {
	p.Name = name
}

func (p ProtoResource) ResourceType() string {
	return p.Type
}

func (p ProtoResource) TerraformResource() string {
	return fmt.Sprintf("%s.%s", p.Type, p.Name)
}

func (p ProtoResource) TerraformResourceID() string {
	return fmt.Sprintf("%s.id", p.TerraformResource())
}

func (p ProtoResource) String() string {
	return Nprintf(`
	resource "${resource_type}" "${resource_name}" {
	  ${attributes}
	}
	`, map[string]any{
		"resource_type": p.Type,
		"resource_name": p.Name,
		"attributes":    p.RenderAttributes(),
	})
}

func (p ProtoResource) RenderAttributes() string {
	var attributes []string

	const padding = "\t  "

	for key, val := range p.Required {
		if len(attributes) == 0 {
			attributes = append(attributes, fmt.Sprintf(`%s = %s`, key, val.String()))
		} else {
			attributes = append(attributes, fmt.Sprintf(`%v%s = %s`, padding, key, val.String()))
		}
	}

	for key, val := range p.Optional {
		attributes = append(attributes, fmt.Sprintf(`%v%s = %s`, padding, key, val.String()))
	}

	return strings.Join(attributes, "\n")
}

func (p ProtoResource) Set(values ...any) Resource {
	panic("not implemented")
}

func (p ProtoResource) Delete(attributes ...string) Resource {
	for _, key := range attributes {
		delete(p.Optional, key)
	}

	return p
}

type ProtoDatasource struct {
	ProtoResource
}

func (d ProtoDatasource) TerraformResource() string {
	return fmt.Sprintf("data.%s.%s", d.Type, d.Name)
}

func (d ProtoDatasource) String() string {
	return Nprintf(`
	datasource "${resource_type}" "${resource_name}" {
	  ${attributes}
	}
	`, map[string]any{
		"resource_type": d.Type,
		"resource_name": d.Name,
		"attributes":    d.RenderAttributes(),
	})
}
