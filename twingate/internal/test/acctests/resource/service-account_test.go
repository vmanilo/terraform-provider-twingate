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

func configServiceAccount(resourceName, serviceAccountName string) string {
	return acctests.Nprintf(`
	resource "twingate_service_account" "${service_account_resource}" {
	  name = "${name}"
	}
	`,
		map[string]any{
			"service_account_resource": resourceName,
			"name":                     serviceAccountName,
		})
}

func TestAccTwingateServiceAccountCreateUpdate(t *testing.T) {
	t.Parallel()

	serviceAccount := NewServiceAccount()
	theResource := serviceAccount.TerraformResource()
	name1 := test.RandomName()
	name2 := test.RandomName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(serviceAccount.Set(attr.Name, name1)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name1),
				),
			},
			{
				Config: configBuilder(serviceAccount.Set(attr.Name, name2)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name2),
				),
			},
		},
	})
}

func TestAccTwingateServiceAccountDelete(t *testing.T) {
	t.Parallel()

	serviceAccount := NewServiceAccount()
	theResource := serviceAccount.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config:  configBuilder(serviceAccount),
				Destroy: true,
			},
			{
				Config: configBuilder(serviceAccount),
				ConfigPlanChecks: sdk.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(theResource, plancheck.ResourceActionCreate),
					},
				},
			},
		},
	})
}

func TestAccTwingateServiceAccountReCreateAfterDeletion(t *testing.T) {
	t.Parallel()

	serviceAccount := NewServiceAccount()
	theResource := serviceAccount.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(serviceAccount),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					acctests.DeleteTwingateResource(theResource, resource.TwingateServiceAccount),
					acctests.WaitTestFunc(),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: configBuilder(serviceAccount),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
		},
	})
}
