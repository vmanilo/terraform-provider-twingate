package datasource

import (
	"fmt"
	"testing"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	resourcesLen     = attr.Len(attr.Resources)
	resourceNamePath = attr.Path(attr.Resources, attr.Name)
)

func TestAccDatasourceTwingateResources_basic(t *testing.T) {
	acctests.SetPageLimit(t, 1)
	networkName := test.RandomName()
	resourceName := test.RandomResourceName()
	const theDatasource = "data.twingate_resources.out_drs1"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDatasourceTwingateResources(networkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(theDatasource, resourcesLen, "2"),
					resource.TestCheckResourceAttr(theDatasource, resourceNamePath, resourceName),
				),
			},
		},
	})
}

func testDatasourceTwingateResources(networkName, resourceName string) string {
	return acctests.Nprintf(`
	resource "twingate_remote_network" "test_drs1" {
	  name = "${network_name}"
	}

	resource "twingate_resource" "test_drs1_1" {
	  name = "${resource_name}"
	  address = "acc-test.com"
	  remote_network_id = twingate_remote_network.test_drs1.id
	  protocols = {
	    allow_icmp = true
	    tcp = {
	      policy = "RESTRICTED"
	      ports = ["80-83", "85"]
	    }
	    udp = {
	      policy = "ALLOW_ALL"
	      ports = []
	    }
	  }
	}

	resource "twingate_resource" "test_drs1_2" {
	  name = "${resource_name}"
	  address = "acc-test.com"
	  remote_network_id = twingate_remote_network.test_drs1.id
	  protocols = {
	    allow_icmp = true
	    tcp = {
	      policy = "ALLOW_ALL"
	      ports = []
	    }
	    udp = {
	      policy = "ALLOW_ALL"
	      ports = []
	    }
	  }
	}

	data "twingate_resources" "out_drs1" {
	  name = "${resource_name}"

	  depends_on = [twingate_resource.test_drs1_1, twingate_resource.test_drs1_2]
	}
	`,
		map[string]any{
			"network_name":  networkName,
			"resource_name": resourceName,
		})
}

func TestAccDatasourceTwingateResources_emptyResult(t *testing.T) {
	t.Parallel()

	resourceName := test.RandomResourceName()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck: func() {
			acctests.PreCheck(t)
		},
		Steps: []resource.TestStep{
			{
				Config: testTwingateResourcesDoesNotExists(resourceName),
				Check: acctests.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.twingate_resources.out_drs2", resourcesLen, "0"),
				),
			},
		},
	})
}

func testTwingateResourcesDoesNotExists(name string) string {
	return fmt.Sprintf(`
	data "twingate_resources" "out_drs2" {
	  name = "%s"
	}

	output "my_resources_drs2" {
	  value = data.twingate_resources.out_drs2.resources
	}
	`, name)
}
