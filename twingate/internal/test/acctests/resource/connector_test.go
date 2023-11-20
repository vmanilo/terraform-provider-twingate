package resource

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/provider/resource"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccRemoteConnectorCreate(t *testing.T) {
	t.Parallel()

	connectorName := test.RandomConnectorName()
	theResource := acctests.TerraformConnector(connectorName)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configConnector(connectorName, connectorName, test.RandomName()),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(connectorName)),
					sdk.TestCheckResourceAttrSet(theResource, attr.Name),
				),
			},
		},
	})
}

func TestAccRemoteConnectorWithCustomName(t *testing.T) {
	t.Parallel()

	connectorName := test.RandomConnectorName()
	theResource := acctests.TerraformConnector(connectorName)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configConnectorWithName(connectorName, test.RandomName(), connectorName),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(connectorName)),
					sdk.TestCheckResourceAttr(theResource, attr.Name, connectorName),
				),
			},
		},
	})
}

func TestAccRemoteConnectorImport(t *testing.T) {
	t.Parallel()

	connectorName := test.RandomConnectorName()
	theResource := acctests.TerraformConnector(connectorName)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configConnectorWithName(connectorName, test.RandomName(), connectorName),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(connectorName)),
					sdk.TestCheckResourceAttr(theResource, attr.Name, connectorName),
				),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      theResource,
				ImportStateCheck: acctests.CheckImportState(map[string]string{
					attr.Name: connectorName,
				}),
			},
		},
	})
}

func TestAccRemoteConnectorNotAllowedToChangeRemoteNetworkId(t *testing.T) {
	t.Parallel()

	terraformConnectorName := test.RandomConnectorName()
	terraformRemoteNetworkName1 := test.RandomNetworkName()
	terraformRemoteNetworkName2 := test.RandomNetworkName()
	theResource := acctests.TerraformConnector(terraformConnectorName)
	remoteNetworkName1 := test.RandomName()
	remoteNetworkName2 := test.RandomName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configConnector(terraformRemoteNetworkName1, terraformConnectorName, remoteNetworkName1),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(terraformRemoteNetworkName1)),
				),
			},
			{
				Config:      configConnector(terraformRemoteNetworkName2, terraformConnectorName, remoteNetworkName2),
				ExpectError: regexp.MustCompile(resource.ErrNotAllowChangeRemoteNetworkID.Error()),
			},
		},
	})
}

func TestAccTwingateConnectorReCreateAfterDeletion(t *testing.T) {
	t.Parallel()

	connectorName := test.RandomConnectorName()
	theResource := acctests.TerraformConnector(connectorName)
	networkName := test.RandomName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configConnector(connectorName, connectorName, networkName),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(connectorName)),
					acctests.DeleteTwingateResource(theResource, resource.TwingateConnector),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: configConnector(connectorName, connectorName, networkName),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(connectorName)),
				),
			},
		},
	})
}

func configConnector(networkTR, connectorTR, networkName string) string {
	return acctests.Nprintf(`
	%{remote_network}

	resource "twingate_connector" "%{connector_resource}" {
	  remote_network_id = twingate_remote_network.%{remote_network_resource}.id
	}
	`,
		map[string]any{
			"remote_network":          terraformResourceRemoteNetwork(networkTR, networkName),
			"connector_resource":      connectorTR,
			"remote_network_resource": networkTR,
		})
}

func configConnectorWithName(terraformResource, networkName, connectorName string) string {
	return acctests.Nprintf(`
	%{remote_network}

	resource "twingate_connector" "%{connector_resource}" {
	  remote_network_id = twingate_remote_network.%{remote_network_resource}.id
	  name  = "%{connector_name}"
	}
	`,
		map[string]any{
			"remote_network":          terraformResourceRemoteNetwork(terraformResource, networkName),
			"connector_resource":      terraformResource,
			"remote_network_resource": terraformResource,
			"connector_name":          connectorName,
		})
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
	theResource := acctests.TerraformConnector(connectorName)
	networkName := test.RandomName()
	connectorNewName := test.RandomConnectorName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configConnector(connectorName, connectorName, networkName),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(connectorName)),
					sdk.TestCheckResourceAttrSet(theResource, attr.Name),
				),
			},
			{
				Config: configConnectorWithName(connectorName, networkName, connectorNewName),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.Name, connectorNewName),
				),
			},
		},
	})
}

func TestAccRemoteConnectorCreateWithNotificationStatus(t *testing.T) {
	t.Parallel()

	connectorName := test.RandomConnectorName()
	theResource := acctests.TerraformConnector(connectorName)
	networkName := test.RandomName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configConnector(connectorName, connectorName, networkName),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(connectorName)),
					sdk.TestCheckResourceAttrSet(theResource, attr.Name),
				),
			},
			{
				// expecting no changes, as by default notifications enabled
				PlanOnly: true,
				Config:   configConnectorWithNotificationStatus(connectorName, connectorName, networkName, true),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.StatusUpdatesEnabled, "true"),
				),
			},
			{
				Config: configConnectorWithNotificationStatus(connectorName, connectorName, networkName, false),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.StatusUpdatesEnabled, "false"),
				),
			},
			{
				// expecting no changes, when user removes `status_updates_enabled` field from terraform
				PlanOnly: true,
				Config:   configConnector(connectorName, connectorName, networkName),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.StatusUpdatesEnabled, "false"),
				),
			},
		},
	})
}

func configConnectorWithNotificationStatus(terraformRemoteNetworkName, terraformConnectorName, remoteNetworkName string, notificationStatus bool) string {
	return acctests.Nprintf(`
	%{remote_network}

	resource "twingate_connector" "%{connector_resource}" {
	  remote_network_id = twingate_remote_network.%{remote_network_resource}.id
	  status_updates_enabled = %{notification_status}
	}
	`,
		map[string]any{
			"remote_network":          terraformResourceRemoteNetwork(terraformRemoteNetworkName, remoteNetworkName),
			"connector_resource":      terraformConnectorName,
			"remote_network_resource": terraformRemoteNetworkName,
			"notification_status":     notificationStatus,
		})

}
