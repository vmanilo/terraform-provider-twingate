package resource

import (
	"fmt"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test/acctests/config"
	"regexp"
	"testing"

	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/provider/resource"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccRemoteConnectorCreate(t *testing.T) {
	t.Parallel()

	network := config.NewResourceRemoteNetwork()
	connector := config.NewResourceConnector(attr.RemoteNetworkID, network.TerraformResourceID())
	theResource := connector.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(network, connector),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, network.TerraformResource()),
					sdk.TestCheckResourceAttrSet(theResource, attr.Name),
					sdk.TestCheckResourceAttrSet(theResource, attr.State),
				),
			},
		},
	})
}

func TestAccRemoteConnectorUpdateNetworkName(t *testing.T) {
	t.Parallel()

	remoteNetworkName1 := test.RandomName()
	remoteNetworkName2 := test.RandomName()

	network := config.NewResourceRemoteNetwork()
	connector := config.NewResourceConnector(attr.RemoteNetworkID, network.TerraformResourceID())
	theResource := connector.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(network.Set(attr.Name, remoteNetworkName1), connector),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, network.TerraformResource()),
					sdk.TestCheckResourceAttr(network.TerraformResource(), attr.Name, remoteNetworkName1),
				),
			},
			{
				Config: config.Builder(network.Set(attr.Name, remoteNetworkName2), connector),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, network.TerraformResource()),
					sdk.TestCheckResourceAttr(network.TerraformResource(), attr.Name, remoteNetworkName2),
				),
			},
		},
	})
}

func TestAccRemoteConnectorWithCustomName(t *testing.T) {
	t.Parallel()

	connectorName := "  Some        connector   name     "
	sanitizedName := "Some connector name"

	network := config.NewResourceRemoteNetwork()
	connector := config.NewResourceConnector(attr.RemoteNetworkID, network.TerraformResourceID())
	theResource := connector.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(network, connector.Set(attr.Name, connectorName)),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, network.TerraformResource()),
					sdk.TestCheckResourceAttr(theResource, attr.Name, connectorName),
					acctests.CheckConnectorName(theResource, sanitizedName),
				),
			},
			{
				// expecting no changes
				PlanOnly: true,
				Config:   config.Builder(network, connector.Set(attr.Name, connectorName)),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, network.TerraformResource()),
					sdk.TestCheckResourceAttr(theResource, attr.Name, connectorName),
					acctests.CheckConnectorName(theResource, sanitizedName),
				),
			},
		},
	})
}

func TestAccRemoteConnectorImport(t *testing.T) {
	t.Parallel()

	connectorName := test.RandomConnectorName()

	network := config.NewResourceRemoteNetwork()
	connector := config.NewResourceConnector(attr.RemoteNetworkID, network.TerraformResourceID(), attr.Name, connectorName)
	theResource := connector.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(network, connector),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, network.TerraformResource()),
					sdk.TestMatchResourceAttr(theResource, attr.Name, regexp.MustCompile(connectorName[:len(connectorName)-3]+".*")),
				),
			},
			{
				ResourceName:      theResource,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRemoteConnectorNotAllowedToChangeRemoteNetworkId(t *testing.T) {
	t.Parallel()

	network1 := config.NewResourceRemoteNetwork()
	network2 := config.NewResourceRemoteNetwork()
	connector := config.NewResourceConnector()
	theResource := connector.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(network1, network2, connector.Set(attr.RemoteNetworkID, network1.TerraformResourceID())),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, network1.TerraformResource()),
				),
			},
			{
				Config:      config.Builder(network1, network2, connector.Set(attr.RemoteNetworkID, network2.TerraformResourceID())),
				ExpectError: regexp.MustCompile(resource.ErrNotAllowChangeRemoteNetworkID.Error()),
			},
		},
	})
}

func TestAccTwingateConnectorReCreateAfterDeletion(t *testing.T) {
	t.Parallel()

	network := config.NewResourceRemoteNetwork()
	connector := config.NewResourceConnector(attr.RemoteNetworkID, network.TerraformResourceID())
	theResource := connector.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(network, connector),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, network.TerraformResource()),
					acctests.DeleteTwingateResource(theResource, resource.TwingateConnector),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: config.Builder(network, connector),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, network.TerraformResource()),
				),
			},
		},
	})
}

func terraformResourceTwingateConnector(terraformRemoteNetworkName, terraformConnectorName, remoteNetworkName string) string {
	return fmt.Sprintf(`
	%s

	resource "twingate_connector" "%s" {
	  remote_network_id = twingate_remote_network.%s.id
	}
	`, terraformResourceRemoteNetwork(terraformRemoteNetworkName, remoteNetworkName), terraformConnectorName, terraformRemoteNetworkName)
}

func checkTwingateConnectorSetWithRemoteNetwork(connectorResource, remoteNetworkResource string) sdk.TestCheckFunc {
	return func(s *terraform.State) error {
		connector, ok := s.RootModule().Resources[connectorResource]
		if !ok {
			return fmt.Errorf("Not found: %s ", connectorResource)
		}

		if connector.Primary.ID == "" {
			return fmt.Errorf("No connectorID set ")
		}

		remoteNetwork, ok := s.RootModule().Resources[remoteNetworkResource]
		if !ok {
			return fmt.Errorf("Not found: %s ", remoteNetworkResource)
		}

		if connector.Primary.Attributes[attr.RemoteNetworkID] != remoteNetwork.Primary.ID {
			return fmt.Errorf("Remote Network ID not set properly in the connector ")
		}

		return nil
	}
}

func TestAccRemoteConnectorUpdateName(t *testing.T) {
	t.Parallel()

	connectorName := test.RandomConnectorName()
	network := config.NewResourceRemoteNetwork()
	connector := config.NewResourceConnector(attr.RemoteNetworkID, network.TerraformResourceID())
	theResource := connector.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(network, connector),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, network.TerraformResource()),
					sdk.TestCheckResourceAttrSet(theResource, attr.Name),
				),
			},
			{
				Config: config.Builder(network, connector.Set(attr.Name, connectorName)),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.Name, connectorName),
				),
			},
		},
	})
}

func TestAccRemoteConnectorCreateWithNotificationStatus(t *testing.T) {
	t.Parallel()

	network := config.NewResourceRemoteNetwork()
	connector := config.NewResourceConnector(attr.RemoteNetworkID, network.TerraformResourceID())
	theResource := connector.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(network, connector),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, network.TerraformResource()),
					sdk.TestCheckResourceAttrSet(theResource, attr.Name),
				),
			},
			{
				// expecting no changes, as by default notifications enabled
				PlanOnly: true,
				Config:   config.Builder(network, connector.Set(attr.StatusUpdatesEnabled, true)),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.StatusUpdatesEnabled, "true"),
				),
			},
			{
				Config: config.Builder(network, connector.Set(attr.StatusUpdatesEnabled, false)),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.StatusUpdatesEnabled, "false"),
				),
			},
			{
				// expecting no changes, when user removes `status_updates_enabled` field from terraform
				PlanOnly: true,
				Config:   config.Builder(network, connector.Delete(attr.StatusUpdatesEnabled)),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.StatusUpdatesEnabled, "false"),
				),
			},
		},
	})
}

func TestAccRemoteConnectorCreateWithNotificationStatusFalse(t *testing.T) {
	t.Parallel()

	network := config.NewResourceRemoteNetwork()
	connector := config.NewResourceConnector(attr.RemoteNetworkID, network.TerraformResourceID(), attr.StatusUpdatesEnabled, false)
	theResource := connector.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(network, connector),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.StatusUpdatesEnabled, "false"),
				),
			},
		},
	})
}

func TestAccRemoteConnectorNameMustBeAtLeast3CharactersLong(t *testing.T) {
	t.Parallel()

	network := config.NewResourceRemoteNetwork()
	connector := config.NewResourceConnector(attr.RemoteNetworkID, network.TerraformResourceID())
	theResource := connector.TerraformResource()

	expectedError := regexp.MustCompile("Attribute name must be at least 3 characters long")

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config:      config.Builder(network, connector.Set(attr.Name, "")),
				ExpectError: expectedError,
			},
			{
				Config:      config.Builder(network, connector.Set(attr.Name, "   ab    ")),
				ExpectError: expectedError,
			},
			{
				Config: config.Builder(network, connector.Set(attr.Name, "   a    b    ")),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckConnectorName(theResource, "a b"),
				),
			},
		},
	})

}
