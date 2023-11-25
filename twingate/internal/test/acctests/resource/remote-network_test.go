package resource

import (
	"testing"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/model"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/provider/resource"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTwingateRemoteNetworkCreate(t *testing.T) {
	t.Parallel()

	networkResource := test.RandomNetworkName()
	theResource := acctests.TerraformRemoteNetwork(networkResource)
	networkName := test.RandomName()
	networkLocation := model.LocationAzure

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configRemoteNetworkWithLocation(networkResource, networkName, networkLocation),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, networkName),
					sdk.TestCheckResourceAttr(theResource, attr.Location, networkLocation),
				),
			},
		},
	})
}

func configRemoteNetworkWithLocation(networkResource, name, location string) string {
	return acctests.Nprintf(`
	resource "twingate_remote_network" "${network_resource}" {
	  name = "${name}"
	  location = "${location}"
	}
	`, map[string]any{
		"network_resource": networkResource,
		"name":             name,
		"location":         location,
	})
}

func TestAccTwingateRemoteNetworkUpdate(t *testing.T) {
	t.Parallel()

	networkResource := test.RandomNetworkName()
	theResource := acctests.TerraformRemoteNetwork(networkResource)
	name1 := test.RandomName()
	name2 := test.RandomName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configRemoteNetwork(networkResource, name1),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name1),
					sdk.TestCheckResourceAttr(theResource, attr.Location, model.LocationOther),
				),
			},
			{
				Config: configRemoteNetworkWithLocation(networkResource, name2, model.LocationAWS),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name2),
					sdk.TestCheckResourceAttr(theResource, attr.Location, model.LocationAWS),
				),
			},
		},
	})
}

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

func TestAccTwingateRemoteNetworkDeleteNonExisting(t *testing.T) {
	t.Parallel()

	networkResource := test.RandomNetworkName()
	theResource := acctests.TerraformRemoteNetwork(networkResource)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config:  configRemoteNetwork(networkResource, test.RandomName()),
				Destroy: true,
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceDoesNotExists(theResource),
				),
			},
		},
	})
}

func TestAccTwingateRemoteNetworkReCreateAfterDeletion(t *testing.T) {
	t.Parallel()

	networkResource := test.RandomNetworkName()
	theResource := acctests.TerraformRemoteNetwork(networkResource)
	networkName := test.RandomName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configRemoteNetwork(networkResource, networkName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					acctests.DeleteTwingateResource(theResource, resource.TwingateRemoteNetwork),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: configRemoteNetwork(networkResource, networkName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
		},
	})
}

func TestAccTwingateRemoteNetworkUpdateWithTheSameName(t *testing.T) {
	t.Parallel()

	networkResource := test.RandomNetworkName()
	theResource := acctests.TerraformRemoteNetwork(networkResource)
	name := test.RandomName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configRemoteNetwork(networkResource, name),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name),
					sdk.TestCheckResourceAttr(theResource, attr.Location, model.LocationOther),
				),
			},
			{
				Config: configRemoteNetworkWithLocation(networkResource, name, model.LocationAWS),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name),
					sdk.TestCheckResourceAttr(theResource, attr.Location, model.LocationAWS),
				),
			},
		},
	})
}
