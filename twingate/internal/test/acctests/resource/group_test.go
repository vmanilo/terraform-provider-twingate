package resource

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/provider/resource"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/utils"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const attrTerraformResource = "terraform_resource"

func collectResourceIDs[T TerraformResource](resources ...T) []string {
	ids := make([]string, 0, len(resources))

	for _, res := range resources {
		ids = append(ids, res.TerraformResource()+".id")
	}

	return ids
}

type TerraformResource interface {
	TerraformResource() string
}

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
		ResourceName: test.RandomGroupName(),
		Name:         test.RandomName(),
	}
}

func optionalString(val any) *string {
	if val == nil {
		return nil
	}

	switch t := val.(type) {
	case string:
		return &t
	case *string:
		return t
	default:
		return nil
	}
}

func optionalBool(val any) *bool {
	if val == nil {
		return nil
	}

	switch t := val.(type) {
	case bool:
		return &t
	case *bool:
		return t
	default:
		return nil
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

type wrapper struct {
	str string
}

func (w *wrapper) String() string {
	return w.str
}

func wrap(str string) fmt.Stringer {
	return &wrapper{str: str}
}

func configBuilder(resources ...any) string {
	var list []fmt.Stringer

	for _, r := range resources {
		switch t := r.(type) {
		case fmt.Stringer:
			list = append(list, t)
		case []*User:
			list = append(list, utils.Map(t, func(item *User) fmt.Stringer {
				return item
			})...)
		}
	}

	buff := bytes.NewBufferString("")
	for _, item := range list {
		buff.WriteString(item.String() + "\n")
	}

	return buff.String()
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

var userIdsLen = attr.Len(attr.UserIDs)

func TestAccTwingateGroupCreateUpdate(t *testing.T) {
	t.Parallel()

	name1 := test.RandomName()
	name2 := test.RandomName()

	group := NewGroup().Set(attr.Name, name1)
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
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceDoesNotExists(group.TerraformResource()),
				),
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
