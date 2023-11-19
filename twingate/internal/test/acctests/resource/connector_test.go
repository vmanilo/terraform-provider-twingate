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

var testRegexp = regexp.MustCompile(test.Prefix() + ".*")

func TestAccRemoteConnectorCreate(t *testing.T) {
	t.Parallel()

	terraformResourceName := test.RandomConnectorName()
	theResource := acctests.TerraformConnector(terraformResourceName)
	remoteNetworkName := test.RandomName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: terraformResourceTwingateConnector(terraformResourceName, terraformResourceName, remoteNetworkName),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(terraformResourceName)),
					sdk.TestCheckResourceAttrSet(theResource, attr.Name),
				),
			},
		},
	})
}

func TestAccRemoteConnectorWithCustomName(t *testing.T) {
	t.Parallel()

	terraformResourceName := test.RandomConnectorName()
	theResource := acctests.TerraformConnector(terraformResourceName)
	remoteNetworkName := test.RandomName()
	connectorName := test.RandomConnectorName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: terraformResourceTwingateConnectorWithName(terraformResourceName, remoteNetworkName, connectorName),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(terraformResourceName)),
					sdk.TestMatchResourceAttr(theResource, attr.Name, regexp.MustCompile(connectorName)),
				),
			},
		},
	})
}

func TestAccRemoteConnectorImport(t *testing.T) {
	t.Parallel()

	terraformResourceName := test.RandomConnectorName()
	theResource := acctests.TerraformConnector(terraformResourceName)
	remoteNetworkName := test.RandomName()
	connectorName := test.RandomConnectorName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: terraformResourceTwingateConnectorWithName(terraformResourceName, remoteNetworkName, connectorName),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(terraformResourceName)),
					sdk.TestMatchResourceAttr(theResource, attr.Name, testRegexp),
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
				Config: terraformResourceTwingateConnector(terraformRemoteNetworkName1, terraformConnectorName, remoteNetworkName1),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(terraformRemoteNetworkName1)),
				),
			},
			{
				Config:      terraformResourceTwingateConnector(terraformRemoteNetworkName2, terraformConnectorName, remoteNetworkName2),
				ExpectError: regexp.MustCompile(resource.ErrNotAllowChangeRemoteNetworkID.Error()),
			},
		},
	})
}

func TestAccTwingateConnectorReCreateAfterDeletion(t *testing.T) {
	t.Parallel()

	terraformResourceName := test.RandomConnectorName()
	theResource := acctests.TerraformConnector(terraformResourceName)
	remoteNetworkName := test.RandomName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: terraformResourceTwingateConnector(terraformResourceName, terraformResourceName, remoteNetworkName),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(terraformResourceName)),
					acctests.DeleteTwingateResource(theResource, resource.TwingateConnector),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: terraformResourceTwingateConnector(terraformResourceName, terraformResourceName, remoteNetworkName),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(terraformResourceName)),
				),
			},
		},
	})
}

func terraformResourceTwingateConnector(terraformRemoteNetworkName, terraformConnectorName, remoteNetworkName string) string {
	return acctests.Nprintf(`
	%{remote_network}

	resource "twingate_connector" "%{connector_resource}" {
	  remote_network_id = twingate_remote_network.%{remote_network_resource}.id
	}
	`,
		map[string]any{
			"remote_network":          terraformResourceRemoteNetwork(terraformRemoteNetworkName, remoteNetworkName),
			"connector_resource":      terraformConnectorName,
			"remote_network_resource": terraformRemoteNetworkName,
		})
}

func terraformResourceTwingateConnectorWithName(terraformResourceName, remoteNetworkName, connectorName string) string {
	return acctests.Nprintf(`
	%{remote_network}

	resource "twingate_connector" "%{connector_resource}" {
	  remote_network_id = twingate_remote_network.%{remote_network_resource}.id
	  name  = "%{connector_name}"
	}
	`,
		map[string]any{
			"remote_network":          terraformResourceRemoteNetwork(terraformResourceName, remoteNetworkName),
			"connector_resource":      terraformResourceName,
			"remote_network_resource": terraformResourceName,
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

	terraformResourceName := test.RandomConnectorName()
	theResource := acctests.TerraformConnector(terraformResourceName)
	remoteNetworkName := test.RandomName()
	connectorName := test.RandomConnectorName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: terraformResourceTwingateConnector(terraformResourceName, terraformResourceName, remoteNetworkName),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(terraformResourceName)),
					sdk.TestCheckResourceAttrSet(theResource, attr.Name),
				),
			},
			{
				Config: terraformResourceTwingateConnectorWithName(terraformResourceName, remoteNetworkName, connectorName),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.Name, connectorName),
				),
			},
		},
	})
}

func TestAccRemoteConnectorCreateWithNotificationStatus(t *testing.T) {
	t.Parallel()

	terraformResourceName := test.RandomConnectorName()
	theResource := acctests.TerraformConnector(terraformResourceName)
	remoteNetworkName := test.RandomName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: terraformResourceTwingateConnector(terraformResourceName, terraformResourceName, remoteNetworkName),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, acctests.TerraformRemoteNetwork(terraformResourceName)),
					sdk.TestCheckResourceAttrSet(theResource, attr.Name),
				),
			},
			{
				// expecting no changes, as by default notifications enabled
				PlanOnly: true,
				Config:   terraformResourceTwingateConnectorWithNotificationStatus(terraformResourceName, terraformResourceName, remoteNetworkName, true),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.StatusUpdatesEnabled, "true"),
				),
			},
			{
				Config: terraformResourceTwingateConnectorWithNotificationStatus(terraformResourceName, terraformResourceName, remoteNetworkName, false),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.StatusUpdatesEnabled, "false"),
				),
			},
			{
				// expecting no changes, when user removes `status_updates_enabled` field from terraform
				PlanOnly: true,
				Config:   terraformResourceTwingateConnector(terraformResourceName, terraformResourceName, remoteNetworkName),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.StatusUpdatesEnabled, "false"),
				),
			},
		},
	})
}

func terraformResourceTwingateConnectorWithNotificationStatus(terraformRemoteNetworkName, terraformConnectorName, remoteNetworkName string, notificationStatus bool) string {
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
