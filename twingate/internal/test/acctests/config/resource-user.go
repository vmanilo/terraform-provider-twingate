package config

import (
	"fmt"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test/acctests"
	"strings"
)

type ResourceUser struct {
	ResourceName string
	Email        string
	FirstName    *string
	LastName     *string
	Role         *string
	SendInvite   bool
	IsActive     bool
}

func NewResourceUser(terraformResourceName ...string) *ResourceUser {
	resourceName := test.RandomResourceName()
	if len(terraformResourceName) > 0 {
		resourceName = terraformResourceName[0]
	}

	return &ResourceUser{
		ResourceName: resourceName,
		Email:        test.RandomEmail(),
		IsActive:     true,  // default value
		SendInvite:   false, // default value for tests
	}
}

func (u *ResourceUser) optionalAttributes() string {
	var optional []string

	if u.FirstName != nil {
		optional = append(optional, fmt.Sprintf(`first_name = "%s"`, *u.FirstName))
	}

	if u.LastName != nil {
		optional = append(optional, fmt.Sprintf(`last_name = "%s"`, *u.LastName))
	}

	if u.Role != nil {
		optional = append(optional, fmt.Sprintf(`role = "%s"`, *u.Role))
	}

	return strings.Join(optional, "\n")
}

func (u *ResourceUser) TerraformResource() string {
	return acctests.TerraformUser(u.ResourceName)
}

func (u *ResourceUser) String() string {
	return Nprintf(`
	resource "twingate_user" "${terraform_resource}" {
	  email = "${email}"
	  send_invite = ${send_invite}
	  is_active = ${is_active}

	  ${optional_attributes}
	}
	`, map[string]any{
		"terraform_resource":  u.ResourceName,
		"email":               u.Email,
		"send_invite":         u.SendInvite,
		"is_active":           u.IsActive,
		"optional_attributes": u.optionalAttributes(),
	})
}

func (u *ResourceUser) Set(values ...any) *ResourceUser {
	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		val := values[i+1]

		switch key {
		case attr.Email:
			u.Email = val.(string)
		case attr.FirstName:
			u.FirstName = optionalString(val)
		case attr.LastName:
			u.LastName = optionalString(val)
		case attr.Role:
			u.Role = optionalString(val)
		case attr.SendInvite:
			u.SendInvite = val.(bool)
		case attr.IsActive:
			u.IsActive = val.(bool)
		}
	}

	return u
}

func GenUsers(count int, resourcePrefix ...string) []*ResourceUser {
	users := make([]*ResourceUser, 0, count)

	prefix := test.RandomUserName()
	if len(resourcePrefix) > 0 {
		prefix = resourcePrefix[0]
	}

	for i := 0; i < count; i++ {
		users = append(users, NewResourceUser(fmt.Sprintf("%s_%d", prefix, i+1)))
	}

	return users
}
