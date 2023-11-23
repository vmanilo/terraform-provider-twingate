package datasource

import (
	"testing"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatasourceTwingateRemoteNetworks_read(t *testing.T) {
	acctests.SetPageLimit(t, 1)

	prefix := acctest.RandString(10)
	networkName1 := test.RandomName(prefix)
	networkName2 := test.RandomName(prefix)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateRemoteNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDatasourceTwingateRemoteNetworks2(networkName1, networkName2, prefix),
				Check: acctests.ComposeTestCheckFunc(
					testCheckOutputLength("test_networks", 2),
				),
			},
		},
	})
}

func testDatasourceTwingateRemoteNetworks2(networkName1, networkName2, prefix string) string {
	return acctests.Nprintf(`
	resource "twingate_remote_network" "test_drn1" {
		name = "${name_1}"
	}
	
	resource "twingate_remote_network" "test_drn2" {
		name = "${name_2}"
	}
	
	data "twingate_remote_networks" "all" {
		depends_on = [twingate_remote_network.test_drn1, twingate_remote_network.test_drn2]
	}

	output "test_networks" {
	  	value = [for n in [for net in data.twingate_remote_networks.all : net if can(net.*.name)][0] : n if length(regexall("${prefix}.*", n.name)) > 0]
	}
		`,
		map[string]any{
			"name_1": networkName1,
			"name_2": networkName2,
			"prefix": prefix,
		})
}
