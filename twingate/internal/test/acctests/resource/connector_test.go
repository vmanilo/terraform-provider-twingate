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

	network := NewRemoteNetwork()
	connector := NewConnector(network.TerraformResourceID())
	theResource := connector.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(network, connector),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, network.TerraformResource()),
					sdk.TestCheckResourceAttrSet(theResource, attr.Name),
				),
			},
		},
	})
}

func TestAccRemoteConnectorWithCustomName(t *testing.T) {
	t.Parallel()

	network := NewRemoteNetwork()
	connector := NewConnector(network.TerraformResourceID()).Set(attr.Name, test.RandomConnectorName())
	theResource := connector.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(network, connector),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, network.TerraformResource()),
					sdk.TestCheckResourceAttr(theResource, attr.Name, *connector.Name),
				),
			},
		},
	})
}

func TestAccRemoteConnectorImport(t *testing.T) {
	t.Parallel()

	network := NewRemoteNetwork()
	connector := NewConnector(network.TerraformResourceID()).Set(attr.Name, test.RandomConnectorName())
	theResource := connector.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(network, connector),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, network.TerraformResource()),
					sdk.TestCheckResourceAttr(theResource, attr.Name, *connector.Name),
				),
			},
			{
				ImportState:       true,
				ImportStateVerify: true,
				ResourceName:      theResource,
				ImportStateCheck: acctests.CheckImportState(map[string]string{
					attr.Name: *connector.Name,
				}),
			},
		},
	})
}

func TestAccRemoteConnectorNotAllowedToChangeRemoteNetworkId(t *testing.T) {
	t.Parallel()

	network1 := NewRemoteNetwork()
	network2 := NewRemoteNetwork()
	connector := NewConnector(network1.TerraformResourceID())
	theResource := connector.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(network1, network2, connector),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, network1.TerraformResource()),
				),
			},
			{
				Config:      configBuilder(network1, network2, connector.Set(attr.RemoteNetworkID, network2.TerraformResourceID())),
				ExpectError: regexp.MustCompile(resource.ErrNotAllowChangeRemoteNetworkID.Error()),
			},
		},
	})
}

func TestAccTwingateConnectorReCreateAfterDeletion(t *testing.T) {
	t.Parallel()

	network := NewRemoteNetwork()
	connector := NewConnector(network.TerraformResourceID())
	theResource := connector.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(network, connector),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, network.TerraformResource()),
					acctests.DeleteTwingateResource(theResource, resource.TwingateConnector),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: configBuilder(network, connector),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, network.TerraformResource()),
				),
			},
		},
	})
}

func configConnector(networkTR, connectorTR, networkName string) string {
	return acctests.Nprintf(`
	${remote_network}

	resource "twingate_connector" "${connector_resource}" {
	  remote_network_id = twingate_remote_network.${remote_network_resource}.id
	}
	`,
		map[string]any{
			"remote_network":          configRemoteNetwork(networkTR, networkName),
			"connector_resource":      connectorTR,
			"remote_network_resource": networkTR,
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

	network := NewRemoteNetwork()
	connector := NewConnector(network.TerraformResourceID())
	theResource := connector.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(network, connector),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, network.TerraformResource()),
					sdk.TestCheckResourceAttrSet(theResource, attr.Name),
				),
			},
			{
				Config: configBuilder(network, connector.Set(attr.Name, connectorName)),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.Name, connectorName),
				),
			},
		},
	})
}

func TestAccRemoteConnectorCreateWithNotificationStatus(t *testing.T) {
	t.Parallel()

	network := NewRemoteNetwork()
	connector := NewConnector(network.TerraformResourceID())
	theResource := connector.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(network, connector),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorSetWithRemoteNetwork(theResource, network.TerraformResource()),
					sdk.TestCheckResourceAttrSet(theResource, attr.Name),
				),
			},
			{
				// expecting no changes, as by default notifications enabled
				PlanOnly: true,
				Config:   configBuilder(network, connector.Set(attr.StatusUpdatesEnabled, true)),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.StatusUpdatesEnabled, "true"),
				),
			},
			{
				Config: configBuilder(network, connector.Set(attr.StatusUpdatesEnabled, false)),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.StatusUpdatesEnabled, "false"),
				),
			},
			{
				// expecting no changes, when user removes `status_updates_enabled` field from terraform
				PlanOnly: true,
				Config:   configBuilder(network, connector.Set(attr.StatusUpdatesEnabled, nil)),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.StatusUpdatesEnabled, "false"),
				),
			},
		},
	})
}

func TestAccRemoteConnectorCreateWithNotificationStatusFalse(t *testing.T) {
	t.Parallel()

	network := NewRemoteNetwork()
	connector := NewConnector(network.TerraformResourceID())
	theResource := connector.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorAndRemoteNetworkDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(network, connector.Set(attr.StatusUpdatesEnabled, false)),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.StatusUpdatesEnabled, "false"),
				),
			},
		},
	})
}
