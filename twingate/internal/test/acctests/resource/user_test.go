package resource

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/model"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/provider/resource"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTwingateUserCreateUpdate(t *testing.T) {
	t.Parallel()

	userResource := test.RandomUserName()
	theResource := acctests.TerraformUser(userResource)
	email := test.RandomEmail()
	firstName := test.RandomName()
	lastName := test.RandomName()
	role := model.UserRoleSupport

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateUserDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configUser(userResource, email),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Email, email),
				),
			},
			{
				Config: configUserWithFirstName(userResource, email, firstName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Email, email),
					sdk.TestCheckResourceAttr(theResource, attr.FirstName, firstName),
				),
			},
			{
				Config: configUserWithLastName(userResource, email, lastName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Email, email),
					sdk.TestCheckResourceAttr(theResource, attr.FirstName, firstName),
					sdk.TestCheckResourceAttr(theResource, attr.LastName, lastName),
				),
			},
			{
				Config: configUserWithRole(userResource, email, role),
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

func configUser(userResource, email string) string {
	return acctests.Nprintf(`
	resource "twingate_user" "${user_resource}" {
	  email = "${email}"
	  send_invite = false
	}
	`, map[string]any{
		"user_resource": userResource,
		"email":         email,
	})
}

func configUserWithFirstName(userResource, email, firstName string) string {
	return acctests.Nprintf(`
	resource "twingate_user" "${user_resource}" {
	  email = "${email}"
	  first_name = "${first_name}"
	  send_invite = false
	}
	`, map[string]any{
		"user_resource": userResource,
		"email":         email,
		"first_name":    firstName,
	})
}

func configUserWithLastName(userResource, email, lastName string) string {
	return acctests.Nprintf(`
	resource "twingate_user" "${user_resource}" {
	  email = "${email}"
	  last_name = "${last_name}"
	  send_invite = false
	}
	`, map[string]any{
		"user_resource": userResource,
		"email":         email,
		"last_name":     lastName,
	})
}

func configUserWithRole(userResource, email, role string) string {
	return acctests.Nprintf(`
	resource "twingate_user" "${user_resource}" {
	  email = "${email}"
	  role = "${role}"
	  send_invite = false
	}
	`, map[string]any{
		"user_resource": userResource,
		"email":         email,
		"role":          role,
	})
}

func TestAccTwingateUserFullCreate(t *testing.T) {
	t.Parallel()

	userResource := test.RandomUserName()
	theResource := acctests.TerraformUser(userResource)
	email := test.RandomEmail()
	firstName := test.RandomName()
	lastName := test.RandomName()
	role := test.RandomUserRole()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateUserDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configUserFull(userResource, email, firstName, lastName, role),
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

func configUserFull(userResource, email, firstName, lastName, role string) string {
	return acctests.Nprintf(`
	resource "twingate_user" "${user_resource}" {
	  email = "${email}"
	  first_name = "${first_name}"
	  last_name = "${last_name}"
	  role = "${role}"
	  send_invite = false
	}
	`, map[string]any{
		"user_resource": userResource,
		"email":         email,
		"first_name":    firstName,
		"last_name":     lastName,
		"role":          role,
	})
}

func TestAccTwingateUserReCreation(t *testing.T) {
	t.Parallel()

	userResource := test.RandomUserName()
	theResource := acctests.TerraformUser(userResource)
	email1 := test.RandomEmail()
	email2 := test.RandomEmail()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateUserDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configUser(userResource, email1),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Email, email1),
				),
			},
			{
				Config: configUser(userResource, email2),
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

	userResource := test.RandomUserName()
	theResource := acctests.TerraformUser(userResource)
	email := test.RandomEmail()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateUserDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configUser(userResource, email),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Email, email),
				),
			},
			{
				Config:      configUserDisabled(userResource, email),
				ExpectError: regexp.MustCompile(`User in PENDING state`),
			},
		},
	})
}

func configUserDisabled(userResource, email string) string {
	return acctests.Nprintf(`
	resource "twingate_user" "${user_resource}" {
	  email = "${email}"
	  send_invite = false
	  is_active = false
	}
	`, map[string]any{
		"user_resource": userResource,
		"email":         email,
	})
}

func TestAccTwingateUserDelete(t *testing.T) {
	t.Parallel()

	userResource := test.RandomUserName()
	theResource := acctests.TerraformUser(userResource)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateUserDestroy,
		Steps: []sdk.TestStep{
			{
				Config:  configUser(userResource, test.RandomEmail()),
				Destroy: true,
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceDoesNotExists(theResource),
				),
			},
		},
	})
}

func TestAccTwingateUserReCreateAfterDeletion(t *testing.T) {
	t.Parallel()

	userResource := test.RandomUserName()
	theResource := acctests.TerraformUser(userResource)
	email := test.RandomEmail()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateUserDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configUser(userResource, email),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					acctests.DeleteTwingateResource(theResource, resource.TwingateUser),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: configUser(userResource, email),
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
				Config:      configUserWithRole(test.RandomUserName(), test.RandomEmail(), "UnknownRole"),
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

func configUserWithoutEmail(userResource string) string {
	return acctests.Nprintf(`
	resource "twingate_user" "${user_resource}" {
	  send_invite = false
	}
	`, map[string]any{
		"user_resource": userResource,
	})
}

func genNewUsers(resourcePrefix string, count int) ([]string, []string) {
	users := make([]string, 0, count)
	userIDs := make([]string, 0, count)

	for i := 0; i < count; i++ {
		resourceName := fmt.Sprintf("%s_%d", resourcePrefix, i+1)
		users = append(users, configUser(resourceName, test.RandomEmail()))
		userIDs = append(userIDs, fmt.Sprintf("twingate_user.%s.id", resourceName))
	}

	return users, userIDs
}
