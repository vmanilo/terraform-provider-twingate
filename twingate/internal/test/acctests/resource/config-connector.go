package resource

import (
	"fmt"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
	"strings"
)

type Connector struct {
	ResourceName         string
	RemoteNetworkID      string
	Name                 *string
	StatusUpdatesEnabled *bool
}

func NewConnector(remoteNetworkID string) *Connector {
	return &Connector{
		ResourceName:    test.RandomResourceName(),
		RemoteNetworkID: remoteNetworkID,
	}
}

func (r *Connector) optionalAttributes() string {
	var optional []string

	if r.Name != nil {
		optional = append(optional, fmt.Sprintf(`name = "%s"`, *r.Name))
	}

	if r.StatusUpdatesEnabled != nil {
		optional = append(optional, fmt.Sprintf(`status_updates_enabled = %v`, *r.StatusUpdatesEnabled))
	}

	return strings.Join(optional, "\n")
}

func (r *Connector) TerraformResource() string {
	return acctests.TerraformConnector(r.ResourceName)
}

func (r *Connector) TerraformResourceID() string {
	return r.TerraformResource() + ".id"
}

func (r *Connector) String() string {
	return acctests.Nprintf(`
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

func (r *Connector) Set(values ...any) *Connector {
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
