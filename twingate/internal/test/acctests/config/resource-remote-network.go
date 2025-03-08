package config

import (
	"fmt"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test/acctests"
	"strings"
)

type ResourceRemoteNetwork struct {
	ResourceName string
	Name         string
	Location     *string
}

func NewResourceRemoteNetwork() *ResourceRemoteNetwork {
	return &ResourceRemoteNetwork{
		ResourceName: test.RandomResourceName(),
		Name:         test.RandomNetworkName(),
	}
}

func (r *ResourceRemoteNetwork) optionalAttributes() string {
	var optional []string

	if r.Location != nil {
		optional = append(optional, fmt.Sprintf(`location = "%s"`, *r.Location))
	}

	return strings.Join(optional, "\n")
}

func (r *ResourceRemoteNetwork) TerraformResource() string {
	return acctests.TerraformRemoteNetwork(r.ResourceName)
}

func (r *ResourceRemoteNetwork) TerraformResourceID() string {
	return r.TerraformResource() + ".id"
}

func (r *ResourceRemoteNetwork) String() string {
	return Nprintf(`
	resource "twingate_remote_network" "${terraform_resource}" {
	  name = "${name}"


	  ${optional_attributes}
	}
	`, map[string]any{
		"terraform_resource":  r.ResourceName,
		"name":                r.Name,
		"optional_attributes": r.optionalAttributes(),
	})
}

func (r *ResourceRemoteNetwork) Set(values ...any) *ResourceRemoteNetwork {
	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		val := values[i+1]

		switch key {
		case attr.Name:
			r.Name = val.(string)
		case attr.Location:
			r.Location = optionalString(val)
		}
	}

	return r
}
