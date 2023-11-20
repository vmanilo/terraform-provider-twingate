package resource

import (
	"fmt"
	"testing"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccRemoteConnectorWithTokens(t *testing.T) {
	t.Parallel()

	connectorName := test.RandomConnectorName()
	theResource := acctests.TerraformConnectorTokens(connectorName)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateConnectorTokensInvalidated,
		Steps: []sdk.TestStep{
			{
				Config: configConnectorTokens(connectorName, test.RandomName()),
				Check: acctests.ComposeTestCheckFunc(
					checkTwingateConnectorTokensSet(theResource),
				),
			},
		},
	})
}

func configConnectorTokens(terraformResource, networkName string) string {
	return acctests.Nprintf(`
	${connector}

	resource "twingate_connector_tokens" "${connector_token_resource}" {
	  connector_id = twingate_connector.${connector_resource}.id
      keepers = {
         foo = "bar"
      }
	}
	`, map[string]any{
		"connector":                configConnector(terraformResource, terraformResource, networkName),
		"connector_token_resource": terraformResource,
		"connector_resource":       terraformResource,
	})
}

func checkTwingateConnectorTokensSet(resourceName string) sdk.TestCheckFunc {
	return func(s *terraform.State) error {
		connectorTokens, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
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
