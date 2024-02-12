package resource

import (
	"testing"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/model"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/provider/resource"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func configRemoteNetwork(networkResource, name string) string {
	return acctests.Nprintf(`
	resource "twingate_remote_network" "${network_resource}" {
	  name = "${name}"
	}
	`,
		map[string]any{
			"network_resource": networkResource,
			"name":             name,
		})
}

func TestAccTwingateRemoteNetworkCreate(t *testing.T) {
	t.Parallel()

	network := NewRemoteNetwork()
	theResource := network.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(network.Set(attr.Location, model.LocationAzure)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, network.Name),
					sdk.TestCheckResourceAttr(theResource, attr.Location, model.LocationAzure),
				),
			},
		},
	})
}

func TestAccTwingateRemoteNetworkUpdate(t *testing.T) {
	t.Parallel()

	name1 := test.RandomName()
	name2 := test.RandomName()

	network := NewRemoteNetwork()
	theResource := network.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(network.Set(attr.Name, name1)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name1),
					sdk.TestCheckResourceAttr(theResource, attr.Location, model.LocationOther),
				),
			},
			{
				Config: configBuilder(network.Set(attr.Name, name2, attr.Location, model.LocationAWS)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name2),
					sdk.TestCheckResourceAttr(theResource, attr.Location, model.LocationAWS),
				),
			},
		},
	})
}

func TestAccTwingateRemoteNetworkDeleteNonExisting(t *testing.T) {
	t.Parallel()

	network := NewRemoteNetwork()
	theResource := network.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config:  configBuilder(network),
				Destroy: true,
			},
			{
				Config: configBuilder(network),
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

	network := NewRemoteNetwork()
	theResource := network.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(network),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					acctests.DeleteTwingateResource(theResource, resource.TwingateRemoteNetwork),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: configBuilder(network),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
		},
	})
}

func TestAccTwingateRemoteNetworkUpdateWithTheSameName(t *testing.T) {
	t.Parallel()

	network := NewRemoteNetwork()
	theResource := network.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(network),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, network.Name),
					sdk.TestCheckResourceAttr(theResource, attr.Location, model.LocationOther),
				),
			},
			{
				Config: configBuilder(network.Set(attr.Location, model.LocationAWS)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, network.Name),
					sdk.TestCheckResourceAttr(theResource, attr.Location, model.LocationAWS),
				),
			},
		},
	})
}
