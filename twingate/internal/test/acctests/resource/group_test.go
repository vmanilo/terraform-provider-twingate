package resource

import (
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test/acctests/config"
	"testing"

	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/provider/resource"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

var (
	userIdsLen = attr.Len(attr.UserIDs)
)

func TestAccTwingateGroupCreateUpdate(t *testing.T) {
	t.Parallel()

	name1 := test.RandomName()
	name2 := test.RandomName()

	group := config.NewResourceGroup()
	theResource := group.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateGroupDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(group.Set(attr.Name, name1)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name1),
				),
			},
			{
				Config: config.Builder(group.Set(attr.Name, name2)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name2),
				),
			},
		},
	})
}

func TestAccTwingateGroupDeleteNonExisting(t *testing.T) {
	t.Parallel()

	group := config.NewResourceGroup()
	theResource := group.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateGroupDestroy,
		Steps: []sdk.TestStep{
			{
				Config:  config.Builder(group),
				Destroy: true,
			},
			{
				Config: config.Builder(group),
				ConfigPlanChecks: sdk.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(theResource, plancheck.ResourceActionCreate),
					},
				},
			},
		},
	})
}

func TestAccTwingateGroupReCreateAfterDeletion(t *testing.T) {
	t.Parallel()

	group := config.NewResourceGroup()
	theResource := group.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateGroupDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(group),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					acctests.DeleteTwingateResource(theResource, resource.TwingateGroup),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: config.Builder(group),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
		},
	})
}

func TestAccTwingateGroupWithSecurityPolicy(t *testing.T) {
	t.Parallel()

	name := test.RandomName()
	group := config.NewResourceGroup(attr.Name, name)
	theResource := group.TerraformResource()

	securityPolicies, err := acctests.ListSecurityPolicies()
	if err != nil {
		t.Skip("can't run test:", err)
	}

	testPolicy := securityPolicies[0]

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateGroupDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(group),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name),
				),
			},
			{
				Config: config.Builder(group.Set(attr.SecurityPolicyID, testPolicy.ID)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name),
					sdk.TestCheckResourceAttr(theResource, attr.SecurityPolicyID, testPolicy.ID),
				),
			},
			{
				// expecting no changes
				PlanOnly: true,
				Config:   config.Builder(group.Delete(attr.SecurityPolicyID)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name),
				),
			},
		},
	})
}

func TestAccTwingateGroupUsersAuthoritativeByDefault(t *testing.T) {
	t.Parallel()

	groupName := test.RandomName()
	group := config.NewResourceGroup(attr.Name, groupName)
	theResource := group.TerraformResource()

	users := config.GenUsers(5)
	userIDs := config.ResourceIDs(users)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateGroupDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(users, group.Set(attr.UserIDs, userIDs[:1])),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "1"),
					acctests.CheckGroupUsersLen(theResource, 1),
				),
			},
			{
				Config: config.Builder(users, group.Set(attr.UserIDs, userIDs[:1])),
				Check: acctests.ComposeTestCheckFunc(
					// added new user to the group though API
					acctests.AddGroupUser(theResource, groupName, userIDs[1]),
					acctests.WaitTestFunc(),
					acctests.CheckGroupUsersLen(theResource, 2),
				),
				// expecting drift - terraform going to remove unknown user
				ExpectNonEmptyPlan: true,
			},
			{
				Config: config.Builder(users, group.Set(attr.UserIDs, userIDs[:1])),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "1"),
					acctests.CheckGroupUsersLen(theResource, 1),
				),
			},
			{
				// added 2 new users to the group though terraform
				Config: config.Builder(users, group.Set(attr.UserIDs, userIDs[:3])),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "3"),
					acctests.CheckGroupUsersLen(theResource, 3),
				),
			},
			{
				Config: config.Builder(users, group.Set(attr.UserIDs, userIDs[:3])),
				Check: acctests.ComposeTestCheckFunc(
					// delete one user from the group though API
					acctests.DeleteGroupUser(theResource, userIDs[2]),
					acctests.WaitTestFunc(),
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "3"),
					acctests.CheckGroupUsersLen(theResource, 2),
				),
				// expecting drift - terraform going to restore deleted user
				ExpectNonEmptyPlan: true,
			},
			{
				Config: config.Builder(users, group.Set(attr.UserIDs, userIDs[:3])),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "3"),
					acctests.CheckGroupUsersLen(theResource, 3),
				),
			},
			{
				// remove 2 users from the group though terraform
				Config: config.Builder(users, group.Set(attr.UserIDs, userIDs[:1])),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "1"),
					acctests.CheckGroupUsersLen(theResource, 1),
				),
			},
			{
				// expecting no drift
				Config:   config.Builder(users, group.Set(attr.UserIDs, userIDs[:1], attr.IsAuthoritative, true)),
				PlanOnly: true,
			},
			{
				Config: config.Builder(users, group.Set(attr.UserIDs, userIDs[:2], attr.IsAuthoritative, true)),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "2"),
					acctests.CheckGroupUsersLen(theResource, 2),
				),
			},
		},
	})
}

func TestAccTwingateGroupUsersNotAuthoritative(t *testing.T) {
	t.Parallel()

	groupName := test.RandomName()
	group := config.NewResourceGroup(attr.Name, groupName, attr.IsAuthoritative, false)
	theResource := group.TerraformResource()

	users := config.GenUsers(3)
	userIDs := config.ResourceIDs(users)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateGroupDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(users, group.Set(attr.UserIDs, userIDs[:1])),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "1"),
					acctests.CheckGroupUsersLen(theResource, 1),
				),
			},
			{
				Config: config.Builder(users, group.Set(attr.UserIDs, userIDs[:1])),
				Check: acctests.ComposeTestCheckFunc(
					// added new user to the group though API
					acctests.AddGroupUser(theResource, groupName, userIDs[2]),
					acctests.WaitTestFunc(),
					acctests.CheckGroupUsersLen(theResource, 2),
				),
			},
			{
				// added new user to the group though terraform
				Config: config.Builder(users, group.Set(attr.UserIDs, userIDs[:2])),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "2"),
					acctests.CheckGroupUsersLen(theResource, 3),
				),
			},
			{
				// remove one user from the group though terraform
				Config: config.Builder(users, group.Set(attr.UserIDs, userIDs[:1])),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "1"),
					acctests.CheckGroupUsersLen(theResource, 2),
					// remove one user from the group though API
					acctests.DeleteGroupUser(theResource, userIDs[2]),
					acctests.WaitTestFunc(),
					acctests.CheckGroupUsersLen(theResource, 1),
				),
			},
			{
				// expecting no drift - empty plan
				Config:   config.Builder(users, group.Set(attr.UserIDs, userIDs[:1])),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTwingateGroupUsersCursor(t *testing.T) {
	t.Parallel()

	group := config.NewResourceGroup()
	theResource := group.TerraformResource()

	users := config.GenUsers(3)
	userIDs := config.ResourceIDs(users)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateGroupDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(users, group.Set(attr.UserIDs, userIDs)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckGroupUsersLen(theResource, len(users)),
				),
			},
			{
				Config: config.Builder(users, group.Set(attr.UserIDs, userIDs[:2])),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckGroupUsersLen(theResource, 2),
				),
			},
		},
	})
}
