package config

import (
	"fmt"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test/acctests"
	"strings"
)

type ResourceConnectorV1 struct {
	ResourceName         string
	RemoteNetworkID      string
	Name                 *string
	StatusUpdatesEnabled *bool
}

func NewResourceConnectorV1(remoteNetworkID string) *ResourceConnectorV1 {
	return &ResourceConnectorV1{
		ResourceName:    test.RandomResourceName(),
		RemoteNetworkID: remoteNetworkID,
	}
}

//type AttributeType string
//
//const (
//	AttributeTypeString  AttributeType = "string"
//	AttributeTypeBoolean AttributeType = "boolean"
//	AttributeTypeNumber  AttributeType = "number"
//	AttributeTypeUnknown AttributeType = "unknown"
//)

type Resource interface {
	ResourceName() string
	ResourceType() string
	TerraformResource() string
	TerraformResourceID() string
	Render() string
	Set(values ...any) Resource
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

func (p ProtoResource) ResourceType() string {
	return p.Type
}

func (p ProtoResource) TerraformResource() string {
	return fmt.Sprintf("%s.%s", p.Type, p.Name)
}

func (p ProtoResource) TerraformResourceID() string {
	return fmt.Sprintf("%s.id", p.TerraformResource())
}

func (p ProtoResource) Render() string {
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

//func (p ProtoResource) Set(values ...any) *Resource {
//	//TODO implement me
//	panic("implement me")
//}

func (p ProtoResource) RenderAttributes() string {
	var attributes []string

	for key, val := range p.Required {
		attributes = append(attributes, fmt.Sprintf(`%s = %s`, key, val.String()))
	}

	for key, val := range p.Optional {
		attributes = append(attributes, fmt.Sprintf(`%s = %s`, key, val.String()))
	}

	return strings.Join(attributes, "\n")
}

type ProtoDatasource struct {
	ProtoResource
}

func (d ProtoDatasource) TerraformResource() string {
	return fmt.Sprintf("data.%s.%s", d.Type, d.Name)
}

func (d ProtoDatasource) Render() string {
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

type SimpleResource struct {
	ProtoResource
}

func NewSimpleResource(name string, values ...any) Resource {
	res := &SimpleResource{
		ProtoResource: ProtoResource{
			Name: name,
			Type: "twingate_connector",
		},
	}

	return res.Set(values...)
}

func (r *SimpleResource) Set(values ...any) Resource {
	if len(values)%2 != 0 {
		panic("Set requires key-value pairs")
	}

	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		val := values[i+1]

		switch key {
		case attr.Name:
			r.Required[key] = StringAttribute{name: key, value: val.(string)}
		}

	}

	return r
}

type Attribute interface {
	Name() string
	Value() string
	String() string
}

type StringAttribute struct {
	name  string
	value string
}

func (s StringAttribute) Name() string {
	return s.name
}

func (s StringAttribute) Value() string {
	return s.value
}

func (s StringAttribute) String() string {
	return fmt.Sprintf(`"%s"`, s.value)
}

//
//type Attribute struct {
//	Name  string
//	Value fmt.Stringer
//	//Type  AttributeType
//}
//
//func (a *Attribute) String() string {
//	return fmt.Sprintf(`%s = %s`, a.Name, a.Value.String())
//}
//
//type StringAttribute string
//
//func (s StringAttribute) String() string {
//	return fmt.Sprintf(`"%s"`, s)
//}

func (r *ResourceConnectorV1) optionalAttributes() string {
	var optional []string

	if r.Name != nil {
		optional = append(optional, fmt.Sprintf(`name = "%s"`, *r.Name))
	}

	if r.StatusUpdatesEnabled != nil {
		optional = append(optional, fmt.Sprintf(`status_updates_enabled = %v`, *r.StatusUpdatesEnabled))
	}

	return strings.Join(optional, "\n")
}

func (r *ResourceConnectorV1) TerraformResource() string {
	return acctests.TerraformConnector(r.ResourceName)
}

func (r *ResourceConnectorV1) TerraformResourceID() string {
	return r.TerraformResource() + ".id"
}

func (r *ResourceConnectorV1) String() string {
	return Nprintf(`
	resource "twingate_connector" "${terraform_resource}" {
	  remote_network_id = ${remote_network_id}

	  ${optional_attributes}
	}
	`, map[string]any{
		"terraform_resource":  r.ResourceName,
		"remote_network_id":   r.RemoteNetworkID,
		"optional_attributes": r.optionalAttributes(),
	})
}

func (r *ResourceConnectorV1) Set(values ...any) *ResourceConnectorV1 {
	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		val := values[i+1]

		switch key {
		case attr.Name:
			r.Name = optionalString(val)
		case attr.RemoteNetworkID:
			r.RemoteNetworkID = val.(string)
		case attr.StatusUpdatesEnabled:
			r.StatusUpdatesEnabled = optionalBool(val)
		}
	}

	return r
}
