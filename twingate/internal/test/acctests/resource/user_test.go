package resource

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/model"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/provider/resource"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

type User struct {
	ResourceName string
	Email        string
	FirstName    *string
	LastName     *string
	Role         *string
	SendInvite   bool
	IsActive     bool
}

func NewUser(terraformResourceName ...string) *User {
	resourceName := test.RandomUserName()
	if len(terraformResourceName) > 0 {
		resourceName = terraformResourceName[0]
	}

	return &User{
		ResourceName: resourceName,
		Email:        test.RandomEmail(),
		IsActive:     true,  // default value
		SendInvite:   false, // default value for tests
	}
}

func (u *User) optionalAttributes() string {
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

func (u *User) TerraformResource() string {
	return acctests.TerraformUser(u.ResourceName)
}

func (u *User) String() string {
	return acctests.Nprintf(`
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

func (u *User) Set(values ...any) *User {
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

func genUsers(count int, resourcePrefix ...string) []*User {
	users := make([]*User, 0, count)

	prefix := test.RandomUserName()
	if len(resourcePrefix) > 0 {
		prefix = resourcePrefix[0]
	}

	for i := 0; i < count; i++ {
		users = append(users, NewUser(fmt.Sprintf("%s_%d", prefix, i+1)))
	}

	return users
}

func TestAccTwingateUserCreateUpdate(t *testing.T) {
	t.Parallel()

	firstName := test.RandomName()
	lastName := test.RandomName()
	role := model.UserRoleSupport

	user := NewUser()
	theResource := user.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateUserDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(user),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Email, user.Email),
				),
			},
			{
				Config: configBuilder(user.Set(attr.FirstName, firstName)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Email, user.Email),
					sdk.TestCheckResourceAttr(theResource, attr.FirstName, firstName),
				),
			},
			{
				Config: configBuilder(user.Set(attr.LastName, lastName)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Email, user.Email),
					sdk.TestCheckResourceAttr(theResource, attr.FirstName, firstName),
					sdk.TestCheckResourceAttr(theResource, attr.LastName, lastName),
				),
			},
			{
				Config: configBuilder(user.Set(attr.Role, role)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Email, user.Email),
					sdk.TestCheckResourceAttr(theResource, attr.FirstName, firstName),
					sdk.TestCheckResourceAttr(theResource, attr.LastName, lastName),
					sdk.TestCheckResourceAttr(theResource, attr.Role, role),
				),
			},
		},
	})
}

func TestAccTwingateUserFullCreate(t *testing.T) {
	t.Parallel()

	user := NewUser().Set(
		attr.FirstName, test.RandomName(),
		attr.LastName, test.RandomName(),
		attr.Role, test.RandomUserRole(),
	)
	theResource := user.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateUserDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(user),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Email, user.Email),
					sdk.TestCheckResourceAttr(theResource, attr.FirstName, *user.FirstName),
					sdk.TestCheckResourceAttr(theResource, attr.LastName, *user.LastName),
					sdk.TestCheckResourceAttr(theResource, attr.Role, *user.Role),
				),
			},
		},
	})
}

func TestAccTwingateUserReCreation(t *testing.T) {
	t.Parallel()

	email1 := test.RandomEmail()
	email2 := test.RandomEmail()

	user := NewUser().Set(attr.Email, email1)
	theResource := user.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateUserDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(user),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Email, email1),
				),
			},
			{
				Config: configBuilder(user.Set(attr.Email, email2)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Email, email2),
				),
			},
		},
	})
}

func TestAccTwingateUserUpdateState(t *testing.T) {
	t.Parallel()

	user := NewUser()
	theResource := user.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateUserDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(user),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Email, user.Email),
				),
			},
			{
				Config:      configBuilder(user.Set(attr.IsActive, false)),
				ExpectError: regexp.MustCompile(`User in PENDING state`),
			},
		},
	})
}

func TestAccTwingateUserDelete(t *testing.T) {
	t.Parallel()

	user := NewUser()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateUserDestroy,
		Steps: []sdk.TestStep{
			{
				Config:  configBuilder(user),
				Destroy: true,
			},
			{
				Config: configBuilder(user),
				ConfigPlanChecks: sdk.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(user.TerraformResource(), plancheck.ResourceActionCreate),
					},
				},
			},
		},
	})
}

func TestAccTwingateUserReCreateAfterDeletion(t *testing.T) {
	t.Parallel()

	user := NewUser()
	theResource := user.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateUserDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(user),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					acctests.DeleteTwingateResource(theResource, resource.TwingateUser),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: configBuilder(user),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
		},
	})
}

func TestAccTwingateUserCreateWithUnknownRole(t *testing.T) {
	t.Parallel()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateUserDestroy,
		Steps: []sdk.TestStep{
			{
				Config:      configBuilder(NewUser().Set(attr.Role, "UnknownRole")),
				ExpectError: regexp.MustCompile(`Attribute role value must be one of`),
			},
		},
	})
}

func TestAccTwingateUserCreateWithoutEmail(t *testing.T) {
	t.Parallel()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateUserDestroy,
		Steps: []sdk.TestStep{
			{
				Config:      configUserWithoutEmail(test.RandomUserName()),
				ExpectError: regexp.MustCompile("Error: Missing required argument"),
			},
		},
	})
}

func configUserWithoutEmail(terraformResource string) string {
	return acctests.Nprintf(`
	resource "twingate_user" "${terraform_resource}" {
	  send_invite = false
	}
	`, map[string]any{
		"terraform_resource": terraformResource,
	})
}
