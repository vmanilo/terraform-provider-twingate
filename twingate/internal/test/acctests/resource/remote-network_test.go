package resource

import (
	"fmt"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test/acctests/config"
	"testing"

	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/model"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/provider/resource"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestAccTwingateRemoteNetworkCreate(t *testing.T) {
	t.Parallel()

	networkName := test.RandomName()
	networkLocation := model.LocationAzure

	network := config.NewResourceRemoteNetwork(attr.Name, networkName, attr.Location, networkLocation)
	theResource := network.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(network),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, networkName),
					sdk.TestCheckResourceAttr(theResource, attr.Location, networkLocation),
				),
			},
		},
	})
}

func TestAccTwingateRemoteNetworkUpdate(t *testing.T) {
	t.Parallel()

	nameBefore := test.RandomName()
	nameAfter := test.RandomName()

	network := config.NewResourceRemoteNetwork()
	theResource := network.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(network.Set(attr.Name, nameBefore)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, nameBefore),
					sdk.TestCheckResourceAttr(theResource, attr.Location, model.LocationOther),
				),
			},
			{
				Config: config.Builder(network.Set(attr.Name, nameAfter, attr.Location, model.LocationAWS)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, nameAfter),
					sdk.TestCheckResourceAttr(theResource, attr.Location, model.LocationAWS),
				),
			},
		},
	})
}

func terraformResourceRemoteNetwork(terraformResourceName, name string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "%s" {
	  name = "%s"
	}
	`, terraformResourceName, name)
}

func TestAccTwingateRemoteNetworkDeleteNonExisting(t *testing.T) {
	t.Parallel()

	network := config.NewResourceRemoteNetwork()
	theResource := network.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config:  config.Builder(network),
				Destroy: true,
			},
			{
				Config: config.Builder(network),
				ConfigPlanChecks: sdk.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(theResource, plancheck.ResourceActionCreate),
					},
				},
			},
		},
	})
}

func TestAccTwingateRemoteNetworkReCreateAfterDeletion(t *testing.T) {
	t.Parallel()

	network := config.NewResourceRemoteNetwork()
	theResource := network.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(network),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					acctests.DeleteTwingateResource(theResource, resource.TwingateRemoteNetwork),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: config.Builder(network),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
		},
	})
}

func TestAccTwingateRemoteNetworkUpdateWithTheSameName(t *testing.T) {
	t.Parallel()

	name := test.RandomName()
	network := config.NewResourceRemoteNetwork(attr.Name, name)
	theResource := network.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(network),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name),
					sdk.TestCheckResourceAttr(theResource, attr.Location, model.LocationOther),
				),
			},
			{
				Config: config.Builder(network.Set(attr.Location, model.LocationAWS)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name),
					sdk.TestCheckResourceAttr(theResource, attr.Location, model.LocationAWS),
				),
			},
		},
	})
}

func TestAccTwingateRemoteNetworkCreateExitNode(t *testing.T) {
	t.Parallel()

	name := test.RandomName()
	network := config.NewResourceRemoteNetwork(attr.Name, name)
	theResource := network.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(network.Set(attr.Type, model.NetworkTypeExit)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name),
					sdk.TestCheckResourceAttr(theResource, attr.Type, model.NetworkTypeExit),
				),
			},
			{
				Config: config.Builder(network.Set(attr.Type, model.NetworkTypeRegular)),
				ConfigPlanChecks: sdk.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(theResource, plancheck.ResourceActionReplace),
					},
				},
			},
		},
	})
}

func TestAccTwingateRemoteNetworkUpdateOnlyName(t *testing.T) {
	t.Parallel()

	name1 := test.RandomName()
	name2 := test.RandomName()

	network := config.NewResourceRemoteNetwork()
	theResource := network.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(network.Set(attr.Name, name1)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name1),
				),
			},
			{
				Config: config.Builder(network.Set(attr.Name, name2)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name2),
				),
			},
		},
	})
}
