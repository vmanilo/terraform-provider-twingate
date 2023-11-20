package resource

import (
	"strings"
	"testing"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/provider/resource"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var userIdsLen = attr.Len(attr.UserIDs)

func TestAccTwingateGroupCreateUpdate(t *testing.T) {
	t.Parallel()

	groupResource := test.RandomGroupName()
	theResource := acctests.TerraformGroup(groupResource)
	name1 := test.RandomName()
	name2 := test.RandomName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateGroupDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configGroup(groupResource, name1),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name1),
				),
			},
			{
				Config: configGroup(groupResource, name2),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name2),
				),
			},
		},
	})
}

func configGroup(groupResource, name string) string {
	return acctests.Nprintf(`
	resource "twingate_group" "%{group_resource}" {
	  name = "%{group_name}"
	}
	`, map[string]interface{}{
		"group_resource": groupResource,
		"group_name":     name,
	})
}

func TestAccTwingateGroupDeleteNonExisting(t *testing.T) {
	t.Parallel()

	groupResource := test.RandomGroupName()
	theResource := acctests.TerraformGroup(groupResource)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateGroupDestroy,
		Steps: []sdk.TestStep{
			{
				Config:  configGroup(groupResource, test.RandomName()),
				Destroy: true,
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceDoesNotExists(theResource),
				),
			},
		},
	})
}

func TestAccTwingateGroupReCreateAfterDeletion(t *testing.T) {
	t.Parallel()

	groupResource := test.RandomGroupName()
	theResource := acctests.TerraformGroup(groupResource)
	groupName := test.RandomName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateGroupDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configGroup(groupResource, groupName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					acctests.DeleteTwingateResource(theResource, resource.TwingateGroup),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: configGroup(groupResource, groupName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
		},
	})
}

func TestAccTwingateGroupWithSecurityPolicy(t *testing.T) {
	t.Parallel()

	groupResource := test.RandomGroupName()
	theResource := acctests.TerraformGroup(groupResource)
	name := test.RandomName()

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
				Config: configGroup(groupResource, name),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name),
				),
			},
			{
				Config: configGroupWithSecurityPolicy(groupResource, name, testPolicy.ID),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name),
					sdk.TestCheckResourceAttr(theResource, attr.SecurityPolicyID, testPolicy.ID),
				),
			},
			{
				// expecting no changes
				PlanOnly: true,
				Config:   configGroup(groupResource, name),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name),
				),
			},
		},
	})
}

func configGroupWithSecurityPolicy(terraformResourceName, name, securityPolicyID string) string {
	return acctests.Nprintf(`
	resource "twingate_group" "%{group_resource}" {
	  name = "%{group_name}"
	  security_policy_id = "%{security_policy_id}"
	}
	`,
		map[string]interface{}{
			"group_resource":     terraformResourceName,
			"group_name":         name,
			"security_policy_id": securityPolicyID,
		})
}

func TestAccTwingateGroupUsersAuthoritativeByDefault(t *testing.T) {
	t.Parallel()

	groupResource := test.RandomGroupName()
	theResource := acctests.TerraformGroup(groupResource)
	groupName := test.RandomName()

	users, userIDs := genNewUsers("u005", 3)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateGroupDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configGroupWithUsers(groupResource, groupName, users, userIDs[:1]),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "1"),
					acctests.CheckGroupUsersLen(theResource, 1),
				),
			},
			{
				Config: configGroupWithUsers(groupResource, groupName, users, userIDs[:1]),
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
				Config: configGroupWithUsers(groupResource, groupName, users, userIDs[:1]),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "1"),
					acctests.CheckGroupUsersLen(theResource, 1),
				),
			},
			{
				// added 2 new users to the group though terraform
				Config: configGroupWithUsers(groupResource, groupName, users, userIDs[:3]),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "3"),
					acctests.CheckGroupUsersLen(theResource, 3),
				),
			},
			{
				Config: configGroupWithUsers(groupResource, groupName, users, userIDs[:3]),
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
				Config: configGroupWithUsers(groupResource, groupName, users, userIDs[:3]),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "3"),
					acctests.CheckGroupUsersLen(theResource, 3),
				),
			},
			{
				// remove 2 users from the group though terraform
				Config: configGroupWithUsers(groupResource, groupName, users, userIDs[:1]),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "1"),
					acctests.CheckGroupUsersLen(theResource, 1),
				),
			},
			{
				// expecting no drift
				Config:   configGroupWithUsersAuthoritative(groupResource, groupName, users, userIDs[:1], true),
				PlanOnly: true,
			},
			{
				Config: configGroupWithUsersAuthoritative(groupResource, groupName, users, userIDs[:2], true),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "2"),
					acctests.CheckGroupUsersLen(theResource, 2),
				),
			},
		},
	})
}

func configGroupWithUsers(terraformResourceName, name string, users, usersID []string) string {
	return acctests.Nprintf(`
	%{users}

	resource "twingate_group" "%{group_resource}" {
	  name = "%{group_name}"
	  user_ids = [%{user_ids}]
	}
	`,
		map[string]interface{}{
			"users":          strings.Join(users, "\n"),
			"group_resource": terraformResourceName,
			"group_name":     name,
			"user_ids":       strings.Join(usersID, ", "),
		})
}

func configGroupWithUsersAuthoritative(terraformResourceName, name string, users, usersID []string, authoritative bool) string {
	return acctests.Nprintf(`
	%{users}

	resource "twingate_group" "%{group_resource}" {
	  name = "%{group_name}"
	  user_ids = [%{user_ids}]
	  is_authoritative = %{authoritative}
	}
	`,
		map[string]interface{}{
			"users":          strings.Join(users, "\n"),
			"group_resource": terraformResourceName,
			"group_name":     name,
			"user_ids":       strings.Join(usersID, ", "),
			"authoritative":  authoritative,
		})
}

func TestAccTwingateGroupUsersNotAuthoritative(t *testing.T) {
	t.Parallel()

	groupResource := test.RandomGroupName()
	theResource := acctests.TerraformGroup(groupResource)
	groupName := test.RandomName()

	users, userIDs := genNewUsers("u006", 3)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateGroupDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configGroupWithUsersAuthoritative(groupResource, groupName, users, userIDs[:1], false),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "1"),
					acctests.CheckGroupUsersLen(theResource, 1),
				),
			},
			{
				Config: configGroupWithUsersAuthoritative(groupResource, groupName, users, userIDs[:1], false),
				Check: acctests.ComposeTestCheckFunc(
					// added new user to the group though API
					acctests.AddGroupUser(theResource, groupName, userIDs[2]),
					acctests.WaitTestFunc(),
					acctests.CheckGroupUsersLen(theResource, 2),
				),
			},
			{
				// added new user to the group though terraform
				Config: configGroupWithUsersAuthoritative(groupResource, groupName, users, userIDs[:2], false),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, userIdsLen, "2"),
					acctests.CheckGroupUsersLen(theResource, 3),
				),
			},
			{
				// remove one user from the group though terraform
				Config: configGroupWithUsersAuthoritative(groupResource, groupName, users, userIDs[:1], false),
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
				Config:   configGroupWithUsersAuthoritative(groupResource, groupName, users, userIDs[:1], false),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTwingateGroupUsersCursor(t *testing.T) {
	acctests.SetPageLimit(t, 1)

	groupResource := test.RandomGroupName()
	theResource := acctests.TerraformGroup(groupResource)
	groupName := test.RandomName()

	users, userIDs := genNewUsers("u007", 3)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateGroupDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configGroupWithUsers(groupResource, groupName, users, userIDs),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckGroupUsersLen(theResource, len(users)),
				),
			},
			{
				Config: configGroupWithUsers(groupResource, groupName, users[:2], userIDs[:2]),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckGroupUsersLen(theResource, 2),
				),
			},
		},
	})
}
