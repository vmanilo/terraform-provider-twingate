package config

import (
	"fmt"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test/acctests"
	"strings"
)

type ResourceConnector struct {
	ResourceName         string
	RemoteNetworkID      string
	Name                 *string
	StatusUpdatesEnabled *bool
}

func NewResourceConnector(remoteNetworkID string) *ResourceConnector {
	return &ResourceConnector{
		ResourceName:    test.RandomResourceName(),
		RemoteNetworkID: remoteNetworkID,
	}
}

func (r *ResourceConnector) optionalAttributes() string {
	var optional []string

	if r.Name != nil {
		optional = append(optional, fmt.Sprintf(`name = "%s"`, *r.Name))
	}

	if r.StatusUpdatesEnabled != nil {
		optional = append(optional, fmt.Sprintf(`status_updates_enabled = %v`, *r.StatusUpdatesEnabled))
	}

	return strings.Join(optional, "\n")
}

func (r *ResourceConnector) TerraformResource() string {
	return acctests.TerraformConnector(r.ResourceName)
}

func (r *ResourceConnector) TerraformResourceID() string {
	return r.TerraformResource() + ".id"
}

func (r *ResourceConnector) String() string {
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

func (r *ResourceConnector) Set(values ...any) *ResourceConnector {
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
