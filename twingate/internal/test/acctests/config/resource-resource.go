package config

import (
	"fmt"
	"strings"

	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test/acctests"
)

type ResourceResource struct {
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

func NewResourceResource(remoteNetworkID string) *ResourceResource {
	return &ResourceResource{
		ResourceName:    test.RandomResourceName(),
		Name:            test.RandomName(),
		Address:         test.RandomName() + ".com",
		RemoteNetworkID: remoteNetworkID,
	}
}

func (r *ResourceResource) optionalAttributes() string {
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

	for _, access := range r.Access {
		block := access.String()
		if block != "" {
			optional = append(optional, block)
		}
	}

	if r.Protocols != nil {
		optional = append(optional, fmt.Sprintf(`protocols = %s`, r.Protocols.String()))
	}

	return strings.Join(optional, "\n")
}

func (r *ResourceResource) TerraformResource() string {
	return acctests.TerraformResource(r.ResourceName)
}

func (r *ResourceResource) String() string {
	return Nprintf(`
	resource "twingate_resource" "${terraform_resource}" {
	  name = "${name}"
	  address = "${address}"
	  remote_network_id = ${remote_network_id}

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

func (r *ResourceResource) Set(values ...any) *ResourceResource {
	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		val := values[i+1]

		switch key {
		case attr.Name:
			r.Name = val.(string)
		case attr.Address:
			r.Address = val.(string)
		case attr.RemoteNetworkID:
			r.RemoteNetworkID = val.(string)

		case attr.Alias:
			r.Alias = optionalString(val)
		case attr.SecurityPolicyID:
			r.SecurityPolicyID = optionalString(val)
		case attr.IsActive:
			r.IsActive = optionalBool(val)
		case attr.IsVisible:
			r.IsVisible = optionalBool(val)
		case attr.IsAuthoritative:
			r.IsAuthoritative = optionalBool(val)
		case attr.IsBrowserShortcutEnabled:
			r.IsBrowserShortcutEnabled = optionalBool(val)
		case attr.Protocols:
			r.Protocols = val.(*Protocols)
		case attr.Access:
			r.Access = val.([]Access)
		}
	}

	return r
}

type Protocols struct {
	AllowIcmp bool
	UDP       Protocol
	TCP       Protocol
}

func (p *Protocols) String() string {
	return Nprintf(`{
	  allow_icmp = ${allow_icmp}
	  tcp = ${tcp}
	  udp = ${udp}
	}
	`, map[string]any{
		"allow_icmp": p.AllowIcmp,
		"tcp":        p.TCP.String(),
		"udp":        p.UDP.String(),
	})
}

type Protocol struct {
	Policy         string
	Ports          []string
	ShowEmptyPorts bool
}

func (p *Protocol) String() string {
	return Nprintf(`{
	  policy = "${policy}"
	  ${optional_attributes}
	}
	`, map[string]any{
		"policy":              p.Policy,
		"optional_attributes": p.optionalAttributes(),
	})
}

func (p *Protocol) optionalAttributes() string {
	if !p.ShowEmptyPorts && len(p.Ports) == 0 {
		return ""
	}

	return fmt.Sprintf(`ports = [%s]`, toStringList(p.Ports))
}

type Access struct {
	GroupIDs          []string
	ServiceAccountIDs []string
}

func (a *Access) String() string {
	if len(a.GroupIDs) == 0 && len(a.ServiceAccountIDs) == 0 {
		return ""
	}

	return Nprintf(`
	access {
	  ${optional_attributes}
	}
	`, map[string]any{
		"optional_attributes": a.optionalAttributes(),
	})
}

func (a *Access) optionalAttributes() string {
	var optional []string

	if len(a.GroupIDs) > 0 {
		optional = append(optional, fmt.Sprintf(`group_ids = [%s]`, toList(a.GroupIDs)))
	}

	if len(a.ServiceAccountIDs) > 0 {
		optional = append(optional, fmt.Sprintf(`service_account_ids = [%s]`, toList(a.ServiceAccountIDs)))
	}

	return strings.Join(optional, "\n")
}

func toList(ids []string) string {
	return strings.Join(ids, ", ")
}

func toStringList(ids []string) string {
	if len(ids) == 0 {
		return ""
	}

	return `"` + strings.Join(ids, `", "`) + `"`
}
