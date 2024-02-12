package resource

import (
	"regexp"
	"testing"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/model"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/provider/resource"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

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
