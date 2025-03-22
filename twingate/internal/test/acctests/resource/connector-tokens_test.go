package resource

import (
	"fmt"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test/acctests/config"
	"testing"

	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccRemoteConnectorWithTokens(t *testing.T) {
	t.Parallel()

	network := config.NewResourceRemoteNetwork()
	connector := config.NewResourceConnector(attr.RemoteNetworkID, network.TerraformResourceID())
	connectorTokens := config.NewResourceConnectorToken(attr.ConnectorID, connector.TerraformResourceID())
	theResource := connectorTokens.TerraformResource()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorTokensInvalidated,
		Steps: []sdk.TestStep{
			{
				Config: config.Builder(network, connector, connectorTokens),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorTokensSet(theResource),
				),
			},
		},
	})
}

func checkTwingateConnectorTokensSet(connectorNameTokens string) sdk.TestCheckFunc {
	return func(s *terraform.State) error {
		connectorTokens, ok := s.RootModule().Resources[connectorNameTokens]

		if !ok {
			return fmt.Errorf("not found: %s", connectorNameTokens)
		}

		if connectorTokens.Primary.ID == "" {
			return fmt.Errorf("no connectorTokensID set")
		}

		if connectorTokens.Primary.Attributes[attr.AccessToken] == "" {
			return fmt.Errorf("no access token set")
		}

		if connectorTokens.Primary.Attributes[attr.RefreshToken] == "" {
			return fmt.Errorf("no refresh token set")
		}

		return nil
	}
}
