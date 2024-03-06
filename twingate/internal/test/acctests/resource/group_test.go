package resource

import (
	"testing"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/provider/resource"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

var userIdsLen = attr.Len(attr.UserIDs)

func TestAccTwingateGroupCreateUpdate(t *testing.T) {
	t.Parallel()

	name1 := test.RandomName()
	name2 := test.RandomName()

	group := NewGroup()
	theResource := group.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateGroupDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(group.Set(attr.Name, name1)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name1),
				),
			},
			{
				Config: configBuilder(group.Set(attr.Name, name2)),
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

	group := NewGroup()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateGroupDestroy,
		Steps: []sdk.TestStep{
			{
				Config:  configBuilder(group),
				Destroy: true,
			},
			{
				Config: configBuilder(group),
				ConfigPlanChecks: sdk.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(group.TerraformResource(), plancheck.ResourceActionCreate),
					},
				},
			},
		},
	})
}

func TestAccTwingateGroupReCreateAfterDeletion(t *testing.T) {
	t.Parallel()

	group := NewGroup()
	theResource := group.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateGroupDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(group),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					acctests.DeleteTwingateResource(theResource, resource.TwingateGroup),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: configBuilder(group),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
		},
	})
}

func TestAccTwingateGroupWithSecurityPolicy(t *testing.T) {
	t.Parallel()

	group := NewGroup()
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
				Config: configBuilder(group),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, group.Name),
				),
			},
			{
				Config: configBuilder(group.Set(attr.SecurityPolicyID, testPolicy.ID)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, group.Name),
					sdk.TestCheckResourceAttr(theResource, attr.SecurityPolicyID, testPolicy.ID),
				),
			},
			{
				// expecting no changes
				PlanOnly: true,
				Config:   configBuilder(group.Set(attr.SecurityPolicyID, nil)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, group.Name),
				),
			},
		},
	})
}

func TestAccTwingateGroupUsersAuthoritativeByDefault(t *testing.T) {
	t.Parallel()

	users := genUsers(3)
	userIDs := collectResourceIDs(users...)

	group := NewGroup()
	theResource := group.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateGroupDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(users, group.Set(attr.UserIDs, userIDs[:1])),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "1"),
					acctests.CheckGroupUsersLen(theResource, 1),
				),
			},
			{
				Config: configBuilder(users, group.Set(attr.UserIDs, userIDs[:1])),
				Check: acctests.ComposeTestCheckFunc(
					// added new user to the group though API
					acctests.AddGroupUser(theResource, group.Name, userIDs[1]),
					acctests.WaitTestFunc(),
					acctests.CheckGroupUsersLen(theResource, 2),
				),
				// expecting drift - terraform going to remove unknown user
				ExpectNonEmptyPlan: true,
			},
			{
				Config: configBuilder(users, group.Set(attr.UserIDs, userIDs[:1])),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "1"),
					acctests.CheckGroupUsersLen(theResource, 1),
				),
			},
			{
				// added 2 new users to the group though terraform
				Config: configBuilder(users, group.Set(attr.UserIDs, userIDs[:3])),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "3"),
					acctests.CheckGroupUsersLen(theResource, 3),
				),
			},
			{
				Config: configBuilder(users, group.Set(attr.UserIDs, userIDs[:3])),
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
				Config: configBuilder(users, group.Set(attr.UserIDs, userIDs[:3])),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "3"),
					acctests.CheckGroupUsersLen(theResource, 3),
				),
			},
			{
				// remove 2 users from the group though terraform
				Config: configBuilder(users, group.Set(attr.UserIDs, userIDs[:1])),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "1"),
					acctests.CheckGroupUsersLen(theResource, 1),
				),
			},
			{
				// expecting no drift
				Config:   configBuilder(users, group.Set(attr.UserIDs, userIDs[:1], attr.IsAuthoritative, true)),
				PlanOnly: true,
			},
			{
				Config: configBuilder(users, group.Set(attr.UserIDs, userIDs[:2], attr.IsAuthoritative, true)),
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

	users := genUsers(3)
	userIDs := collectResourceIDs(users...)

	group := NewGroup()
	theResource := group.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateGroupDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(users, group.Set(attr.UserIDs, userIDs[:1], attr.IsAuthoritative, false)),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "1"),
					acctests.CheckGroupUsersLen(theResource, 1),
				),
			},
			{
				Config: configBuilder(users, group.Set(attr.UserIDs, userIDs[:1], attr.IsAuthoritative, false)),
				Check: acctests.ComposeTestCheckFunc(
					// added new user to the group though API
					acctests.AddGroupUser(theResource, group.Name, userIDs[2]),
					acctests.WaitTestFunc(),
					acctests.CheckGroupUsersLen(theResource, 2),
				),
			},
			{
				// added new user to the group though terraform
				Config: configBuilder(users, group.Set(attr.UserIDs, userIDs[:2], attr.IsAuthoritative, false)),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "2"),
					acctests.CheckGroupUsersLen(theResource, 3),
				),
			},
			{
				// remove one user from the group though terraform
				Config: configBuilder(users, group.Set(attr.UserIDs, userIDs[:1], attr.IsAuthoritative, false)),
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
				Config:   configBuilder(users, group.Set(attr.UserIDs, userIDs[:1], attr.IsAuthoritative, false)),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTwingateGroupUsersCursor(t *testing.T) {
	acctests.SetPageLimit(t, 1)

	users := genUsers(3)
	userIDs := collectResourceIDs(users...)

	group := NewGroup()
	theResource := group.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateGroupDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(users, group.Set(attr.UserIDs, userIDs)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckGroupUsersLen(theResource, len(users)),
				),
			},
			{
				Config: configBuilder(users[:2], group.Set(attr.UserIDs, userIDs[:2])),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckGroupUsersLen(theResource, 2),
				),
			},
		},
	})
}
