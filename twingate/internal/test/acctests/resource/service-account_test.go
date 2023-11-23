package resource

import (
	"testing"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/provider/resource"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func createServiceAccount(resourceName, serviceAccountName string) string {
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

	resourceName := test.RandomServiceAccountName()
	theResource := acctests.TerraformServiceAccount(resourceName)
	name1 := test.RandomName()
	name2 := test.RandomName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createServiceAccount(resourceName, name1),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name1),
				),
			},
			{
				Config: createServiceAccount(resourceName, name2),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name2),
				),
			},
		},
	})
}

func TestAccTwingateServiceAccountDeleteNonExisting(t *testing.T) {
	t.Parallel()

	resourceName := test.RandomServiceAccountName()
	theResource := acctests.TerraformServiceAccount(resourceName)
	name := test.RandomName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config:  createServiceAccount(resourceName, name),
				Destroy: true,
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceDoesNotExists(theResource),
				),
			},
		},
	})
}

func TestAccTwingateServiceAccountReCreateAfterDeletion(t *testing.T) {
	t.Parallel()

	resourceName := test.RandomServiceAccountName()
	theResource := acctests.TerraformServiceAccount(resourceName)
	name := test.RandomName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createServiceAccount(resourceName, name),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					acctests.DeleteTwingateResource(theResource, resource.TwingateServiceAccount),
					acctests.WaitTestFunc(),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: createServiceAccount(resourceName, name),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
		},
	})
}
