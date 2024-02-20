package resource

import (
	"fmt"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
	"strings"
)

type Group struct {
	ResourceName     string
	Name             string
	SecurityPolicyID *string

	UserIDs        []string
	userIDsEnabled bool

	IsAuthoritative *bool
}

func NewGroup() *Group {
	return &Group{
		ResourceName: test.RandomResourceName(),
		Name:         test.RandomGroupName(),
	}
}

func (g *Group) optionalAttributes() string {
	var optional []string

	if g.SecurityPolicyID != nil {
		optional = append(optional, fmt.Sprintf(`security_policy_id = "%s"`, *g.SecurityPolicyID))
	}

	if g.userIDsEnabled {
		optional = append(optional, fmt.Sprintf(`user_ids = [%s]`, strings.Join(g.UserIDs, ", ")))
	}

	if g.IsAuthoritative != nil {
		optional = append(optional, fmt.Sprintf(`is_authoritative = %v`, *g.IsAuthoritative))
	}

	return strings.Join(optional, "\n")
}

func (g *Group) TerraformResource() string {
	return acctests.TerraformGroup(g.ResourceName)
}

func (g *Group) String() string {
	return acctests.Nprintf(`
	resource "twingate_group" "${terraform_resource}" {
	  name = "${name}"

	  ${optional_attributes}
	}
	`, map[string]any{
		"terraform_resource":  g.ResourceName,
		"name":                g.Name,
		"optional_attributes": g.optionalAttributes(),
	})
}

func (g *Group) Set(values ...any) *Group {
	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		val := values[i+1]

		switch key {
		case attr.Name:
			g.Name = val.(string)
		case attr.SecurityPolicyID:
			g.SecurityPolicyID = optionalString(val)
		case attr.UserIDs:
			g.UserIDs = val.([]string)
			g.userIDsEnabled = len(g.UserIDs) > 0
		case attr.IsAuthoritative:
			g.IsAuthoritative = optionalBool(val)
		}
	}

	return g
}

func configGroup(groupResource, name string) string {
	return acctests.Nprintf(`
	resource "twingate_group" "${group_resource}" {
	  name = "${name}"
	}
	`, map[string]any{
		"group_resource": groupResource,
		"name":           name,
	})
}
