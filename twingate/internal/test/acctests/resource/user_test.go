package resource

import (
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test/acctests/config"
	"regexp"
	"testing"

	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/model"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/provider/resource"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccTwingateUserCreateUpdate(t *testing.T) {
	t.Parallel()

	email := test.RandomEmail()
	firstName := test.RandomName()
	lastName := test.RandomName()
	role := model.UserRoleSupport
	user := config.NewResourceUser(attr.Email, email)
	theResource := user.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateUserDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(user),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Email, email),
				),
			},
			{
				Config: config.Builder(user.Set(attr.FirstName, firstName)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Email, email),
					sdk.TestCheckResourceAttr(theResource, attr.FirstName, firstName),
				),
			},
			{
				Config: config.Builder(user.Set(attr.LastName, lastName)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Email, email),
					sdk.TestCheckResourceAttr(theResource, attr.FirstName, firstName),
					sdk.TestCheckResourceAttr(theResource, attr.LastName, lastName),
				),
			},
			{
				Config: config.Builder(user.Set(attr.Role, role)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Email, email),
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

	email := test.RandomEmail()
	firstName := test.RandomName()
	lastName := test.RandomName()
	role := test.RandomUserRole()
	user := config.NewResourceUser(attr.Email, email, attr.FirstName, firstName, attr.LastName, lastName, attr.Role, role)
	theResource := user.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateUserDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(user),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Email, email),
					sdk.TestCheckResourceAttr(theResource, attr.FirstName, firstName),
					sdk.TestCheckResourceAttr(theResource, attr.LastName, lastName),
					sdk.TestCheckResourceAttr(theResource, attr.Role, role),
				),
			},
		},
	})
}

func TestAccTwingateUserReCreation(t *testing.T) {
	t.Parallel()

	email1 := test.RandomEmail()
	email2 := test.RandomEmail()
	user := config.NewResourceUser()
	theResource := user.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateUserDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(user.Set(attr.Email, email1)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Email, email1),
				),
			},
			{
				Config: config.Builder(user.Set(attr.Email, email2)),
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

	email := test.RandomEmail()
	user := config.NewResourceUser(attr.Email, email)
	theResource := user.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateUserDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(user),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Email, email),
				),
			},
			{
				Config:      config.Builder(user.Set(attr.IsActive, false)),
				ExpectError: regexp.MustCompile(`User in PENDING state`),
			},
		},
	})
}

func TestAccTwingateUserDelete(t *testing.T) {
	t.Parallel()

	user := config.NewResourceUser()
	theResource := user.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateUserDestroy,
		Steps: []sdk.TestStep{
			{
				Config:  config.Builder(user),
				Destroy: true,
			},
			{
				Config: config.Builder(user),
				ConfigPlanChecks: sdk.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(theResource, plancheck.ResourceActionCreate),
					},
				},
			},
		},
	})
}

func TestAccTwingateUserReCreateAfterDeletion(t *testing.T) {
	t.Parallel()

	user := config.NewResourceUser()
	theResource := user.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateUserDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(user),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					acctests.DeleteTwingateResource(theResource, resource.TwingateUser),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: config.Builder(user),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
		},
	})
}

func TestAccTwingateUserCreateWithUnknownRole(t *testing.T) {
	t.Parallel()

	user := config.NewResourceUser(attr.Role, "UnknownRole")

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateUserDestroy,
		Steps: []sdk.TestStep{
			{
				Config:      config.Builder(user),
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
				Config: `
					resource "twingate_user" "invalid_user" {
					  send_invite = false
					}
				`,
				ExpectError: regexp.MustCompile("Error: Missing required argument"),
			},
		},
	})
}
