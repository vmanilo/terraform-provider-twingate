package resource

import (
	"fmt"
	"strings"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
)

type Resource struct {
	ResourceName             string
	Name                     string
	Address                  string
	RemoteNetworkID          string
	Protocols                *Protocols
	Access                   []Access
	IsActive                 *bool
	IsVisible                *bool
	IsAuthoritative          *bool
	IsBrowserShortcutEnabled *bool
	Alias                    *string
	SecurityPolicyID         *string
}

func NewResource(remoteNetworkID string) *Resource {
	return &Resource{
		ResourceName:    test.RandomResourceName(),
		Name:            test.RandomName(),
		Address:         test.RandomName() + ".com",
		RemoteNetworkID: remoteNetworkID,
	}
}

// todo
func (r *Resource) optionalAttributes() string {
	var optional []string

	if r.Alias != nil {
		optional = append(optional, fmt.Sprintf(`alias = "%s"`, *r.Alias))
	}

	if r.SecurityPolicyID != nil {
		optional = append(optional, fmt.Sprintf(`security_policy_id = "%s"`, *r.SecurityPolicyID))
	}

	if r.IsAuthoritative != nil {
		optional = append(optional, fmt.Sprintf(`is_authoritative = %v`, *r.IsAuthoritative))
	}

	if r.IsActive != nil {
		optional = append(optional, fmt.Sprintf(`is_active = %v`, *r.IsActive))
	}

	if r.IsVisible != nil {
		optional = append(optional, fmt.Sprintf(`is_visible = %v`, *r.IsVisible))
	}

	if r.IsBrowserShortcutEnabled != nil {
		optional = append(optional, fmt.Sprintf(`is_browser_shortcut_enabled = %v`, *r.IsBrowserShortcutEnabled))
	}

	return strings.Join(optional, "\n")
}

func (r *Resource) TerraformResource() string {
	return acctests.TerraformResource(r.ResourceName)
}

func (r *Resource) String() string {
	return acctests.Nprintf(`
	resource "twingate_resource" "${terraform_resource}" {
	  name = "${name}"
	  address = "${address}"
	  remote_network_id = "${remote_network_id}"

	  ${optional_attributes}
	}
	`, map[string]any{
		"terraform_resource":  r.ResourceName,
		"name":                r.Name,
		"address":             r.Address,
		"remote_network_id":   r.RemoteNetworkID,
		"optional_attributes": r.optionalAttributes(),
	})
}

func (r *Resource) Set(values ...any) *Resource {
	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		val := values[i+1]

		switch key {
		case attr.Name:
			r.Name = val.(string)
		case attr.SecurityPolicyID:
			r.SecurityPolicyID = optionalString(val)

		case attr.IsAuthoritative:
			r.IsAuthoritative = optionalBool(val)
		}
	}

	return r
}

type Protocols struct {
	AllowIcmp bool
	UDP       Protocol
	TCP       Protocol
}

type Protocol struct {
	Policy         string
	Ports          []string
	ShowEmptyPorts bool
}

type Access struct {
	GroupIDs          []string
	ServiceAccountIDs []string
}
