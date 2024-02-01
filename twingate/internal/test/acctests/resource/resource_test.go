package resource

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/model"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/provider/resource"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/stretchr/testify/assert"
)

var (
	tcpPolicy                  = attr.PathAttr(attr.Protocols, attr.TCP, attr.Policy)
	udpPolicy                  = attr.PathAttr(attr.Protocols, attr.UDP, attr.Policy)
	firstTCPPort               = attr.FirstAttr(attr.Protocols, attr.TCP, attr.Ports)
	firstUDPPort               = attr.FirstAttr(attr.Protocols, attr.UDP, attr.Ports)
	tcpPortsLen                = attr.LenAttr(attr.Protocols, attr.TCP, attr.Ports)
	udpPortsLen                = attr.LenAttr(attr.Protocols, attr.UDP, attr.Ports)
	accessGroupIdsLen          = attr.Len(attr.Access, attr.GroupIDs)
	accessServiceAccountIdsLen = attr.Len(attr.Access, attr.ServiceAccountIDs)
)

func TestAccTwingateResourceCreate(t *testing.T) {
	t.Parallel()

	resourceName := test.RandomResourceName()
	theResource := acctests.TerraformResource(resourceName)
	remoteNetworkName := test.RandomName()
	name := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceBasic(resourceName, remoteNetworkName, name),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckNoResourceAttr(theResource, accessGroupIdsLen),
					sdk.TestCheckResourceAttr(acctests.TerraformRemoteNetwork(resourceName), attr.Name, remoteNetworkName),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name),
					sdk.TestCheckResourceAttr(theResource, attr.Address, "acc-test.com"),
				),
			},
		},
	})
}

func TestAccTwingateResourceUpdateProtocols(t *testing.T) {
	t.Parallel()

	resourceName := test.RandomResourceName()
	theResource := acctests.TerraformResource(resourceName)
	remoteNetworkName := test.RandomName()
	name := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceBasic(resourceName, remoteNetworkName, name),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
			{
				Config: configResourceWithSimpleProtocols(resourceName, remoteNetworkName, name),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
			{
				Config: configResourceBasic(resourceName, remoteNetworkName, name),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
		},
	})
}

func configResourceBasic(terraformResource, networkName, name string) string {
	return acctests.Nprintf(`
	resource "twingate_remote_network" "${network_resource}" {
	  name = "${network_name}"
	}
	resource "twingate_resource" "${resource_resource}" {
	  name = "${resource_name}"
	  address = "acc-test.com"
	  remote_network_id = twingate_remote_network.${network_resource}.id
	}
	`,
		map[string]any{
			"network_resource":  terraformResource,
			"network_name":      networkName,
			"resource_resource": terraformResource,
			"resource_name":     name,
		})
}

func configResourceWithSimpleProtocols(terraformResource, networkName, name string) string {
	return acctests.Nprintf(`
	resource "twingate_remote_network" "${network_resource}" {
	  name = "${network_name}"
	}
	resource "twingate_resource" "${resource_resource}" {
	  name = "${resource_name}"
	  address = "acc-test.com"
	  remote_network_id = twingate_remote_network.${network_resource}.id

	  protocols = {
        allow_icmp = true
        tcp = {
            policy = "DENY_ALL"
        }
        udp = {
            policy = "DENY_ALL"
        }
      }
	}
	`,
		map[string]any{
			"network_resource":  terraformResource,
			"network_name":      networkName,
			"resource_resource": terraformResource,
			"resource_name":     name,
		})
}

func TestAccTwingateResourceCreateWithProtocolsAndGroups(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	networkName := test.RandomName()
	group1 := test.RandomGroupName()
	group2 := test.RandomGroupName()
	name := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithProtocolsAndGroups(terraformResource, networkName, group1, group2, name),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Address, "new-acc-test.com"),
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "2"),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyRestricted),
					sdk.TestCheckResourceAttr(theResource, firstTCPPort, "80"),
				),
			},
		},
	})
}

func configResourceWithProtocolsAndGroups(terraformResource, networkName, groupName1, groupName2, resourceName string) string {
	return acctests.Nprintf(`
	resource "twingate_remote_network" "${network_resource}" {
	  name = "${network_name}"
	}

    resource "twingate_group" "g21" {
      name = "${group_1}"
    }

    resource "twingate_group" "g22" {
      name = "${group_2}"
    }

	resource "twingate_resource" "${resource_resource}" {
	  name = "${resource_name}"
	  address = "new-acc-test.com"
	  remote_network_id = twingate_remote_network.${network_resource}.id

      protocols = {
		allow_icmp = true
        tcp = {
			policy = "${tcp_policy}"
            ports = ["80", "82-83"]
        }
		udp = {
 			policy = "${udp_policy}"
		}
      }

      access {
		group_ids = [twingate_group.g21.id, twingate_group.g22.id]
      }
	}
	`,
		map[string]any{
			"network_resource":  terraformResource,
			"network_name":      networkName,
			"group_1":           groupName1,
			"group_2":           groupName2,
			"resource_resource": terraformResource,
			"resource_name":     resourceName,
			"tcp_policy":        model.PolicyRestricted,
			"udp_policy":        model.PolicyAllowAll,
		})
}

func TestAccTwingateResourceFullCreationFlow(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	networkName := test.RandomName()
	groupName := test.RandomGroupName()
	name := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configCompleteResource(terraformResource, networkName, groupName, name),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(acctests.TerraformRemoteNetwork(terraformResource), attr.Name, networkName),
					sdk.TestCheckResourceAttr(theResource, attr.Name, name),
					sdk.TestMatchResourceAttr(acctests.TerraformConnectorTokens(terraformResource), attr.AccessToken, regexp.MustCompile(".+")),
				),
			},
		},
	})
}

func configCompleteResource(terraformResource, networkName, groupName, resourceName string) string {
	return acctests.Nprintf(`
    resource "twingate_remote_network" "${network_resource}" {
      name = "${network_name}"
    }
	
    resource "twingate_connector" "${connector_resource}" {
      remote_network_id = twingate_remote_network.${network_resource}.id
    }

    resource "twingate_connector_tokens" "${connector_tokens_resource}" {
      connector_id = twingate_connector.${connector_resource}.id
    }

    resource "twingate_connector" "${connector_resource}-2" {
      remote_network_id = twingate_remote_network.${network_resource}.id
    }
	
    resource "twingate_connector_tokens" "${connector_tokens_resource}-2" {
      connector_id = twingate_connector.${connector_resource}-2.id
    }

    resource "twingate_group" "${group_resource}" {
      name = "${group_name}"
    }

    resource "twingate_resource" "${resource_resource}" {
      name = "${resource_name}"
      address = "acc-test.com"
      remote_network_id = twingate_remote_network.${network_resource}.id

      protocols = {
        allow_icmp = true
        tcp = {
            policy = "${tcp_policy}"
            ports = ["3306"]
        }
        udp = {
            policy = "${udp_policy}"
        }
      }

      access {
        group_ids = [twingate_group.${group_resource}.id]
      }
    }
    `,
		map[string]any{
			"network_resource":          terraformResource,
			"network_name":              networkName,
			"connector_resource":        terraformResource,
			"connector_tokens_resource": terraformResource,
			"group_resource":            terraformResource,
			"group_name":                groupName,
			"resource_resource":         terraformResource,
			"resource_name":             resourceName,
			"tcp_policy":                model.PolicyRestricted,
			"udp_policy":                model.PolicyAllowAll,
		})
}

func TestAccTwingateResourceWithInvalidGroupId(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	networkName := test.RandomResourceName()
	name := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		Steps: []sdk.TestStep{
			{
				Config:      configResourceWithInvalidGroupId(terraformResource, networkName, name),
				ExpectError: regexp.MustCompile("failed to create resource: Field 'groupIds' Unable to parse global ID"),
			},
		},
	})
}

func configResourceWithInvalidGroupId(terraformResource, networkName, resourceName string) string {
	return acctests.Nprintf(`
	resource "twingate_remote_network" "${network_resource}" {
	  name = "%s"
	}

	resource "twingate_resource" "${resource_resource}" {
	  name = "${resource_name}"
	  address = "acc-test.com"
	  access {
	    group_ids = ["foo", "bar"]
	  }
	  remote_network_id = twingate_remote_network.${network_resource}.id
	}
	`,
		map[string]any{
			"network_resource":  terraformResource,
			"network_name":      networkName,
			"resource_resource": terraformResource,
			"resource_name":     resourceName,
		})
}

func TestAccTwingateResourceWithTcpDenyAllPolicy(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	resourceName := test.RandomResourceName()
	networkName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithPolicy(terraformResource, networkName, resourceName, model.PolicyDenyAll, model.PolicyAllowAll),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyDenyAll),
				),
			},
			// expecting no changes - empty plan
			{
				Config:   configResourceWithPolicy(terraformResource, networkName, resourceName, model.PolicyDenyAll, model.PolicyAllowAll),
				PlanOnly: true,
			},
		},
	})
}

func configResourceWithPolicy(terraformResource, networkName, resourceName, tcpPolicy, udpPolicy string) string {
	return acctests.Nprintf(`
    resource "twingate_remote_network" "${network_resource}" {
      name = "${network_name}"
    }

    resource "twingate_resource" "${resource_resource}" {
      name = "${resource_name}"
      address = "new-acc-test.com"
      remote_network_id = twingate_remote_network.${network_resource}.id

      protocols = {
        allow_icmp = true
        tcp = {
          policy = "${tcp_policy}"
        }
        udp = {
          policy = "${udp_policy}"
        }
      }
    }
    `,
		map[string]any{
			"network_resource":  terraformResource,
			"network_name":      networkName,
			"group_resource":    terraformResource,
			"resource_resource": terraformResource,
			"resource_name":     resourceName,
			"tcp_policy":        tcpPolicy,
			"udp_policy":        udpPolicy,
		})
}

func TestAccTwingateResourceWithUdpDenyAllPolicy(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithPolicy(terraformResource, remoteNetworkName, resourceName, model.PolicyAllowAll, model.PolicyDenyAll),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, udpPolicy, model.PolicyDenyAll),
				),
			},
			// expecting no changes - empty plan
			{
				Config:   configResourceWithPolicy(terraformResource, remoteNetworkName, resourceName, model.PolicyAllowAll, model.PolicyDenyAll),
				PlanOnly: true,
			},
		},
	})
}

func createResourceWithUdpDenyAllPolicy(networkName, groupName, resourceName string) string {
	return fmt.Sprintf(`
    resource "twingate_remote_network" "test6" {
      name = "%s"
    }

    resource "twingate_group" "g6" {
      name = "%s"
    }

    resource "twingate_resource" "test6" {
      name = "%s"
      address = "acc-test.com"
      remote_network_id = twingate_remote_network.test6.id
      access {
        group_ids = [twingate_group.g6.id]
      }
      protocols = {
        allow_icmp = true
        tcp = {
          policy = "%s"
        }
        udp = {
          policy = "%s"
        }
      }
    }
	`, networkName, groupName, resourceName, model.PolicyAllowAll, model.PolicyDenyAll)
}

func TestAccTwingateResourceWithDenyAllPolicyAndEmptyPortsList(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	groupName := test.RandomGroupName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithPolicyAndEmptyTCPPortsList(terraformResource, remoteNetworkName, groupName, resourceName, model.PolicyDenyAll, model.PolicyDenyAll),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.Name, resourceName),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyDenyAll),
					sdk.TestCheckNoResourceAttr(theResource, tcpPortsLen),
					sdk.TestCheckResourceAttr(theResource, udpPolicy, model.PolicyDenyAll),
					sdk.TestCheckNoResourceAttr(theResource, udpPortsLen),
				),
			},
		},
	})
}

func configResourceWithPolicyAndEmptyTCPPortsList(terraformResource, networkName, groupName, resourceName, tcpPolicy, udpPolicy string) string {
	return acctests.Nprintf(`
    resource "twingate_remote_network" "${network_resource}" {
      name = "${network_name}"
    }

    resource "twingate_group" "${group_resource}" {
      name = "${group_name}"
    }

    resource "twingate_resource" "${resource_resource}" {
      name = "${resource_name}"
      address = "new-acc-test.com"
      remote_network_id = twingate_remote_network.${network_resource}.id
      access {
        group_ids = [twingate_group.${group_resource}.id]
      }
      protocols = {
        allow_icmp = true
        tcp = {
          policy = "${tcp_policy}"
          ports = []
        }
        udp = {
          policy = "${udp_policy}"
        }
      }
    }
    `,
		map[string]any{
			"network_resource":  terraformResource,
			"network_name":      networkName,
			"group_resource":    terraformResource,
			"group_name":        groupName,
			"resource_resource": terraformResource,
			"resource_name":     resourceName,
			"tcp_policy":        tcpPolicy,
			"udp_policy":        udpPolicy,
		})
}

func TestAccTwingateResourceWithInvalidPortRange(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	expectedError := regexp.MustCompile("failed to parse protocols port range")

	genConfig := func(portRange string) string {
		return configResourceWithPortRange(terraformResource, remoteNetworkName, resourceName, portRange)
	}

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		Steps: []sdk.TestStep{
			{
				Config:      genConfig(`""`),
				ExpectError: expectedError,
			},
			{
				Config:      genConfig(`" "`),
				ExpectError: expectedError,
			},
			{
				Config:      genConfig(`"foo"`),
				ExpectError: expectedError,
			},
			{
				Config:      genConfig(`"80-"`),
				ExpectError: expectedError,
			},
			{
				Config:      genConfig(`"-80"`),
				ExpectError: expectedError,
			},
			{
				Config:      genConfig(`"80-90-100"`),
				ExpectError: expectedError,
			},
			{
				Config:      genConfig(`"80-70"`),
				ExpectError: expectedError,
			},
			{
				Config:      genConfig(`"0-65536"`),
				ExpectError: expectedError,
			},
		},
	})
}

func createResourceWithRestrictedPolicyAndPortRange(networkName, resourceName, portRange string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "test8" {
	  name = "%s"
	}

	resource "twingate_resource" "test8" {
	  name = "%s"
	  address = "new-acc-test.com"
	  remote_network_id = twingate_remote_network.test8.id
	  protocols = {
	    allow_icmp = true
	    tcp = {
	      policy = "%s"
	      ports = [%s]
	    }
	    udp = {
	      policy = "%s"
	    }
	  }
	}
	`, networkName, resourceName, model.PolicyRestricted, portRange, model.PolicyAllowAll)
}

func TestAccTwingateResourcePortReorderingCreatesNoChanges(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithPortRange(terraformResource, remoteNetworkName, resourceName, `"80", "82-83"`),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, firstTCPPort, "80"),
					sdk.TestCheckResourceAttr(theResource, firstUDPPort, "80"),
				),
			},
			// no changes
			{
				Config:   configResourceWithPortRange(terraformResource, remoteNetworkName, resourceName, `"82-83", "80"`),
				PlanOnly: true,
			},
			// no changes
			{
				Config:   configResourceWithPortRange(terraformResource, remoteNetworkName, resourceName, `"82", "83", "80"`),
				PlanOnly: true,
			},
			// new changes applied
			{
				Config: configResourceWithPortRange(terraformResource, remoteNetworkName, resourceName, `"70", "82-83"`),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, firstTCPPort, "70"),
					sdk.TestCheckResourceAttr(theResource, firstUDPPort, "70"),
				),
			},
		},
	})
}

func configResourceWithPortRange(terraformResource, networkName, resourceName, portRange string) string {
	return acctests.Nprintf(`
    resource "twingate_remote_network" "${network_resource}" {
      name = "${network_name}"
    }

    resource "twingate_resource" "${resource_resource}" {
      name = "${resource_name}"
      address = "new-acc-test.com"
      remote_network_id = twingate_remote_network.${network_resource}.id
      protocols = {
        allow_icmp = true
        tcp = {
          policy = "${tcp_policy}"
          ports = [${port_range}]
        }
        udp = {
          policy = "${udp_policy}"
          ports = [${port_range}]
        }
      }
    }
    `,
		map[string]any{
			"network_resource":  terraformResource,
			"network_name":      networkName,
			"resource_resource": terraformResource,
			"resource_name":     resourceName,
			"tcp_policy":        model.PolicyRestricted,
			"udp_policy":        model.PolicyRestricted,
			"port_range":        portRange,
		})
}

func TestAccTwingateResourcePortsRepresentationChanged(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithPortRange(terraformResource, remoteNetworkName, resourceName, `"82", "83", "80"`),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "3"),
				),
			},
		},
	})
}

func TestAccTwingateResourcePortsNotChanged(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithPortRange(terraformResource, remoteNetworkName, resourceName, `"82", "83", "80"`),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "3"),
				),
			},
			{
				PlanOnly: true,
				Config:   configResourceWithPortRange(terraformResource, remoteNetworkName, resourceName, `"80", "82-83"`),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "2"),
				),
			},
		},
	})
}

func TestAccTwingateResourcePortReorderingNoChanges(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithPortRange(terraformResource, remoteNetworkName, resourceName, `"82", "83", "80"`),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, firstTCPPort, "80"),
					sdk.TestCheckResourceAttr(theResource, firstUDPPort, "80"),
				),
			},
			// no changes
			{
				Config:   configResourceWithPortRange(terraformResource, remoteNetworkName, resourceName, `"82-83", "80"`),
				PlanOnly: true,
			},
			// no changes
			{
				Config:   configResourceWithPortRange(terraformResource, remoteNetworkName, resourceName, `"82-83", "80"`),
				PlanOnly: true,
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, udpPortsLen, "2"),
				),
			},
			// new changes applied
			{
				Config: configResourceWithPortRange(terraformResource, remoteNetworkName, resourceName, `"70", "82-83"`),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, firstTCPPort, "70"),
					sdk.TestCheckResourceAttr(theResource, firstUDPPort, "70"),
				),
			},
		},
	})
}

func TestAccTwingateResourceSetActiveStateOnUpdate(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceBasic(terraformResource, remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					acctests.DeactivateTwingateResource(theResource),
					acctests.WaitTestFunc(),
					acctests.CheckTwingateResourceActiveState(theResource, false),
				),
				// provider noticed drift and tried to change it to true
				ExpectNonEmptyPlan: true,
			},
			{
				Config: configResourceBasic(terraformResource, remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceActiveState(theResource, true),
				),
			},
		},
	})
}

func TestAccTwingateResourceReCreateAfterDeletion(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceBasic(terraformResource, remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					acctests.DeleteTwingateResource(theResource, resource.TwingateResource),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: configResourceBasic(terraformResource, remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
		},
	})
}

func TestAccTwingateResourceImport(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithPortRange(terraformResource, remoteNetworkName, resourceName, `"80", "82-83"`),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
			{
				ImportState:  true,
				ResourceName: theResource,
				ImportStateCheck: acctests.CheckImportState(map[string]string{
					attr.Address: "new-acc-test.com",
					tcpPolicy:    model.PolicyRestricted,
					tcpPortsLen:  "2",
					firstTCPPort: "80",
					udpPolicy:    model.PolicyRestricted,
				}),
			},
		},
	})
}

func createResource12(networkName, groupName1, groupName2, resourceName string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "test12" {
	  name = "%s"
	}

    resource "twingate_group" "g121" {
      name = "%s"
    }

    resource "twingate_group" "g122" {
      name = "%s"
    }

	resource "twingate_resource" "test12" {
	  name = "%s"
	  address = "acc-test.com.12"
	  remote_network_id = twingate_remote_network.test12.id
	  access {
	    group_ids = [twingate_group.g121.id, twingate_group.g122.id]
      }
      protocols = {
		allow_icmp = true
        tcp = {
			policy = "%s"
            ports = ["80", "82-83"]
        }
		udp = {
 			policy = "%s"
		}
      }
	}
	`, networkName, groupName1, groupName2, resourceName, model.PolicyRestricted, model.PolicyAllowAll)
}

func genNewGroups(resourcePrefix string, count int) ([]string, []string) {
	groups := make([]string, 0, count)
	groupsID := make([]string, 0, count)

	for i := 0; i < count; i++ {
		resourceName := fmt.Sprintf("%s_%d", resourcePrefix, i+1)
		groups = append(groups, configGroup(resourceName, test.RandomName()))
		groupsID = append(groupsID, acctests.TerraformGroup(resourceName)+".id")
	}

	return groups, groupsID
}

func getResourceNameFromID(resourceID string) string {
	idx := strings.LastIndex(resourceID, ".id")
	if idx == -1 {
		return ""
	}

	return resourceID[:idx]
}

func genNewServiceAccounts(resourcePrefix string, count int) ([]string, []string) {
	serviceAccounts := make([]string, 0, count)
	serviceAccountIDs := make([]string, 0, count)

	for i := 0; i < count; i++ {
		resourceName := fmt.Sprintf("%s_%d", resourcePrefix, i+1)
		serviceAccounts = append(serviceAccounts, configServiceAccount(resourceName, test.RandomName()))
		serviceAccountIDs = append(serviceAccountIDs, acctests.TerraformServiceAccount(resourceName)+".id")
	}

	return serviceAccounts, serviceAccountIDs
}

func TestAccTwingateResourceAddAccessServiceAccounts(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	serviceAccountConfig := configServiceAccount(terraformResource, test.RandomName())
	serviceAccountID := acctests.TerraformServiceAccount(terraformResource) + ".id"

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithServiceAccount(terraformResource, test.RandomName(), test.RandomResourceName(), serviceAccountConfig, serviceAccountID),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, accessServiceAccountIdsLen, "1"),
				),
			},
		},
	})
}

func configResourceWithServiceAccount(terraformResource, networkName, resourceName, serviceAccountConfig, serviceAccountID string) string {
	return acctests.Nprintf(`
    resource "twingate_remote_network" "${network_resource}" {
      name = "${network_name}"
    }

    ${service_account}

    resource "twingate_resource" "${resource_resource}" {
      name = "${resource_name}"
      address = "acc-test.com"
      remote_network_id = twingate_remote_network.${network_resource}.id
      protocols = {
        allow_icmp = true
        tcp = {
          policy = "${tcp_policy}"
          ports = ["80", "82-83"]
        }
        udp = {
          policy = "${udp_policy}"
        }
      }

	  access {
	    service_account_ids = [${service_account_id}]
	  }
    }
    `,
		map[string]any{
			"network_resource":   terraformResource,
			"network_name":       networkName,
			"service_account":    serviceAccountConfig,
			"resource_resource":  terraformResource,
			"resource_name":      resourceName,
			"tcp_policy":         model.PolicyRestricted,
			"udp_policy":         model.PolicyAllowAll,
			"service_account_id": serviceAccountID,
		})
}

func TestAccTwingateResourceAddAccessGroupsAndServiceAccounts(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)

	groups, groupsID := genNewGroups(terraformResource, 1)
	serviceAccountConfig := []string{configServiceAccount(terraformResource, test.RandomName())}
	serviceAccountID := []string{acctests.TerraformServiceAccount(terraformResource) + ".id"}

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithGroupsAndServiceAccounts(terraformResource, test.RandomName(), test.RandomResourceName(), groups, groupsID, serviceAccountConfig, serviceAccountID),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "1"),
					sdk.TestCheckResourceAttr(theResource, accessServiceAccountIdsLen, "1"),
				),
			},
		},
	})
}

func configResourceWithGroupsAndServiceAccounts(terraformResource, networkName, resourceName string, groups, groupIDs, serviceAccounts, serviceAccountIDs []string) string {
	return acctests.Nprintf(`
    resource "twingate_remote_network" "${network_resource}" {
      name = "${network_name}"
    }

    ${group}

    ${service_account}

    resource "twingate_resource" "${resource_resource}" {
      name = "${resource_name}"
      address = "acc-test.com"
      remote_network_id = twingate_remote_network.${network_resource}.id
      protocols = {
        allow_icmp = true
        tcp = {
          policy = "${tcp_policy}"
          ports = ["80", "82-83"]
        }
        udp = {
          policy = "${udp_policy}"
        }
      }

	  access {
	    group_ids = [${group_id}]
	    service_account_ids = [${service_account_id}]
	  }
    }
    `,
		map[string]any{
			"network_resource":   terraformResource,
			"network_name":       networkName,
			"service_account":    strings.Join(serviceAccounts, "\n"),
			"group":              strings.Join(groups, "\n"),
			"resource_resource":  terraformResource,
			"resource_name":      resourceName,
			"tcp_policy":         model.PolicyRestricted,
			"udp_policy":         model.PolicyAllowAll,
			"service_account_id": strings.Join(serviceAccountIDs, ", "),
			"group_id":           strings.Join(groupIDs, ", "),
		})
}

func TestAccTwingateResourceAccessServiceAccountsNotAuthoritative(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	serviceAccounts, serviceAccountIDs := genNewServiceAccounts(terraformResource, 3)
	serviceAccountResource := getResourceNameFromID(serviceAccountIDs[2])

	config := func(serviceAccountIDs []string) string {
		return configResourceWithServiceAccountsAndAuthoritativeFlag(terraformResource, remoteNetworkName, resourceName, serviceAccounts, serviceAccountIDs, false)
	}

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config(serviceAccountIDs[:1]),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, accessServiceAccountIdsLen, "1"),
					acctests.WaitTestFunc(),
					// added a new service account to the resource using API
					acctests.AddResourceServiceAccount(theResource, serviceAccountResource),
					acctests.WaitTestFunc(),
					acctests.CheckResourceServiceAccountsLen(theResource, 2),
				),
			},
			{
				// expecting no drift - empty plan
				Config:   config(serviceAccountIDs[:1]),
				PlanOnly: true,
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, accessServiceAccountIdsLen, "1"),
					acctests.CheckResourceServiceAccountsLen(theResource, 2),
				),
			},
			{
				// added a new service account to the resource using terraform
				Config: config(serviceAccountIDs[:2]),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, accessServiceAccountIdsLen, "2"),
					acctests.CheckResourceServiceAccountsLen(theResource, 3),
				),
			},
			{
				// remove one service account from the resource using terraform
				Config: config(serviceAccountIDs[:1]),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, accessServiceAccountIdsLen, "1"),
					acctests.CheckResourceServiceAccountsLen(theResource, 2),
				),
			},
			{
				// expecting no drift - empty plan
				Config:   config(serviceAccountIDs[:1]),
				PlanOnly: true,
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, accessServiceAccountIdsLen, "1"),
					acctests.CheckResourceServiceAccountsLen(theResource, 2),
					// delete service account from the resource using API
					acctests.DeleteResourceServiceAccount(theResource, serviceAccountResource),
					acctests.WaitTestFunc(),
					acctests.CheckResourceServiceAccountsLen(theResource, 1),
				),
			},
			{
				// expecting no drift - empty plan
				Config:   config(serviceAccountIDs[:1]),
				PlanOnly: true,
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, accessServiceAccountIdsLen, "1"),
					acctests.CheckResourceServiceAccountsLen(theResource, 1),
				),
			},
		},
	})
}

func configResourceWithServiceAccountsAndAuthoritativeFlag(terraformResource, networkName, resourceName string, serviceAccounts, serviceAccountIDs []string, authoritative bool) string {
	return acctests.Nprintf(`
    resource "twingate_remote_network" "${network_resource}" {
      name = "${network_name}"
    }

    ${service_account}

    resource "twingate_resource" "${resource_resource}" {
      name = "${resource_name}"
      address = "acc-test.com"
      remote_network_id = twingate_remote_network.${network_resource}.id
      protocols = {
        allow_icmp = true
        tcp = {
          policy = "${tcp_policy}"
          ports = ["80", "82-83"]
        }
        udp = {
          policy = "${udp_policy}"
        }
      }

	  is_authoritative = ${authoritative}
	  access {
	    service_account_ids = [${service_account_id}]
	  }
    }
    `,
		map[string]any{
			"network_resource":   terraformResource,
			"network_name":       networkName,
			"service_account":    strings.Join(serviceAccounts, "\n"),
			"resource_resource":  terraformResource,
			"resource_name":      resourceName,
			"tcp_policy":         model.PolicyRestricted,
			"udp_policy":         model.PolicyAllowAll,
			"service_account_id": strings.Join(serviceAccountIDs, ", "),
			"authoritative":      authoritative,
		})
}

func TestAccTwingateResourceAccessServiceAccountsAuthoritative(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	serviceAccounts, serviceAccountIDs := genNewServiceAccounts(terraformResource, 3)
	serviceAccountResource := getResourceNameFromID(serviceAccountIDs[2])

	config := func(serviceAccountIDs []string) string {
		return configResourceWithServiceAccountsAndAuthoritativeFlag(terraformResource, remoteNetworkName, resourceName, serviceAccounts, serviceAccountIDs, true)
	}

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config(serviceAccountIDs[:1]),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, accessServiceAccountIdsLen, "1"),
					acctests.WaitTestFunc(),
					// added new service account to the resource using API
					acctests.AddResourceServiceAccount(theResource, serviceAccountResource),
					acctests.WaitTestFunc(),
					acctests.CheckResourceServiceAccountsLen(theResource, 2),
				),
				// expecting drift - terraform going to remove unknown service account
				ExpectNonEmptyPlan: true,
			},
			{
				Config: config(serviceAccountIDs[:1]),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, accessServiceAccountIdsLen, "1"),
					acctests.CheckResourceServiceAccountsLen(theResource, 1),
				),
			},
			{
				// added 2 new service accounts to the resource using terraform
				Config: config(serviceAccountIDs),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, accessServiceAccountIdsLen, "3"),
					acctests.CheckResourceServiceAccountsLen(theResource, 3),
				),
			},
			{
				Config: config(serviceAccountIDs),
				Check: acctests.ComposeTestCheckFunc(
					// delete one service account from the resource using API
					acctests.DeleteResourceServiceAccount(theResource, serviceAccountResource),
					acctests.WaitTestFunc(),
					acctests.CheckResourceServiceAccountsLen(theResource, 2),
					sdk.TestCheckResourceAttr(theResource, accessServiceAccountIdsLen, "3"),
				),
				// expecting drift - terraform going to restore deleted service account
				ExpectNonEmptyPlan: true,
			},
			{
				Config: config(serviceAccountIDs),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckResourceServiceAccountsLen(theResource, 3),
					sdk.TestCheckResourceAttr(theResource, accessServiceAccountIdsLen, "3"),
				),
			},
			{
				// remove 2 service accounts from the resource using terraform
				Config: config(serviceAccountIDs[:1]),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckResourceServiceAccountsLen(theResource, 1),
					sdk.TestCheckResourceAttr(theResource, accessServiceAccountIdsLen, "1"),
				),
			},
		},
	})
}

func createResource13(networkName, resourceName string, serviceAccounts, serviceAccountIDs []string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "test13" {
	  name = "%s"
	}

	%s

	resource "twingate_resource" "test13" {
	  name = "%s"
	  address = "acc-test.com.13"
	  remote_network_id = twingate_remote_network.test13.id
	  
	  protocols = {
	    allow_icmp = true
	    tcp = {
	      policy = "%s"
	      ports = ["80", "82-83"]
	    }
	    udp = {
	      policy = "%s"
	    }
	  }

	  is_authoritative = true
	  access {
	    service_account_ids = [%s]
	  }

	}
	`, networkName, strings.Join(serviceAccounts, "\n"), resourceName, model.PolicyRestricted, model.PolicyAllowAll, strings.Join(serviceAccountIDs, ", "))
}

func TestAccTwingateResourceAccessWithEmptyGroups(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config:      configResourceWithGroups(terraformResource, remoteNetworkName, resourceName, nil, nil),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value"),
			},
		},
	})
}

func configResourceWithGroups(terraformResource, networkName, resourceName string, groups, groupsID []string) string {
	return acctests.Nprintf(`
    resource "twingate_remote_network" "${network_resource}" {
      name = "${network_name}"
    }

    ${group}

    resource "twingate_resource" "${resource_resource}" {
      name = "${resource_name}"
      address = "acc-test.com"
      remote_network_id = twingate_remote_network.${network_resource}.id
      protocols = {
        allow_icmp = true
        tcp = {
          policy = "${tcp_policy}"
          ports = ["80", "82-83"]
        }
        udp = {
          policy = "${udp_policy}"
        }
      }

	  access {
	    group_ids = [${group_id}]
	  }
    }
    `,
		map[string]any{
			"network_resource":  terraformResource,
			"network_name":      networkName,
			"group":             strings.Join(groups, "\n"),
			"resource_resource": terraformResource,
			"resource_name":     resourceName,
			"tcp_policy":        model.PolicyRestricted,
			"udp_policy":        model.PolicyAllowAll,
			"group_id":          strings.Join(groupsID, ", "),
		})
}

func TestAccTwingateResourceAccessWithEmptyServiceAccounts(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config:      configResourceWithServiceAccount(terraformResource, remoteNetworkName, resourceName, "", ""),
				ExpectError: regexp.MustCompile("Error: Invalid Attribute Value"),
			},
		},
	})
}

func createResource19(networkName, resourceName string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "test19" {
	  name = "%s"
	}

	resource "twingate_resource" "test19" {
	  name = "%s"
	  address = "acc-test.com.19"
	  remote_network_id = twingate_remote_network.test19.id
	  
	  protocols = {
	    allow_icmp = true
	    tcp = {
	      policy = "%s"
	      ports = ["80", "82-83"]
	    }
	    udp = {
	      policy = "%s"
	    }
	  }

	  access {
	    service_account_ids = []
	  }

	}
	`, networkName, resourceName, model.PolicyRestricted, model.PolicyAllowAll)
}

func TestAccTwingateResourceAccessWithEmptyBlock(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config:      configResourceWithEmptyAccessBlock(terraformResource, remoteNetworkName, resourceName),
				ExpectError: regexp.MustCompile("invalid attribute combination"),
			},
		},
	})
}

func configResourceWithEmptyAccessBlock(terraformResource, networkName, resourceName string) string {
	return acctests.Nprintf(`
    resource "twingate_remote_network" "${network_resource}" {
      name = "${network_name}"
    }

    resource "twingate_resource" "${resource_resource}" {
      name = "${resource_name}"
      address = "acc-test.com"
      remote_network_id = twingate_remote_network.${network_resource}.id
      protocols = {
        allow_icmp = true
        tcp = {
          policy = "${tcp_policy}"
          ports = ["80", "82-83"]
        }
        udp = {
          policy = "${udp_policy}"
        }
      }

	  access {
	  }
    }
    `,
		map[string]any{
			"network_resource":  terraformResource,
			"network_name":      networkName,
			"resource_resource": terraformResource,
			"resource_name":     resourceName,
			"tcp_policy":        model.PolicyRestricted,
			"udp_policy":        model.PolicyAllowAll,
		})
}

func TestAccTwingateResourceAccessGroupsNotAuthoritative(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	groups, groupsID := genNewGroups(terraformResource, 3)
	groupResource := getResourceNameFromID(groupsID[2])

	config := func(groupsID []string) string {
		return configResourceWithGroupsAndAuthoritativeFlag(terraformResource, remoteNetworkName, resourceName, groups, groupsID, false)
	}

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config(groupsID[:1]),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "1"),
					acctests.WaitTestFunc(),
					// added a new group to the resource using API
					acctests.AddResourceGroup(theResource, groupResource),
					acctests.WaitTestFunc(),
					acctests.CheckResourceGroupsLen(theResource, 2),
				),
			},
			{
				// expecting no drift - empty plan
				Config:   config(groupsID[:1]),
				PlanOnly: true,
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "1"),
					acctests.CheckResourceGroupsLen(theResource, 2),
				),
			},
			{
				// added a new group to the resource using terraform
				Config: config(groupsID[:2]),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "2"),
					acctests.CheckResourceGroupsLen(theResource, 3),
				),
			},
			{
				// remove one group from the resource using terraform
				Config: config(groupsID[:1]),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "1"),
					acctests.CheckResourceGroupsLen(theResource, 2),
				),
			},
			{
				// expecting no drift - empty plan
				Config:   config(groupsID[:1]),
				PlanOnly: true,
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "1"),
					acctests.CheckResourceGroupsLen(theResource, 2),
					// remove one group from the resource using API
					acctests.DeleteResourceGroup(theResource, groupResource),
					acctests.WaitTestFunc(),
					acctests.CheckResourceGroupsLen(theResource, 1),
				),
			},
			{
				// expecting no drift - empty plan
				Config:   config(groupsID[:1]),
				PlanOnly: true,
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "1"),
					acctests.CheckResourceGroupsLen(theResource, 1),
				),
			},
		},
	})
}

func configResourceWithGroupsAndAuthoritativeFlag(terraformResource, networkName, resourceName string, groups, groupIDs []string, authoritative bool) string {
	return acctests.Nprintf(`
    resource "twingate_remote_network" "${network_resource}" {
      name = "${network_name}"
    }

    ${groups}

    resource "twingate_resource" "${resource_resource}" {
      name = "${resource_name}"
      address = "acc-test.com"
      remote_network_id = twingate_remote_network.${network_resource}.id
      protocols = {
        allow_icmp = true
        tcp = {
          policy = "${tcp_policy}"
          ports = ["80", "82-83"]
        }
        udp = {
          policy = "${udp_policy}"
        }
      }

	  is_authoritative = ${authoritative}
	  access {
	    group_ids = [${group_ids}]
	  }
    }
    `,
		map[string]any{
			"network_resource":  terraformResource,
			"network_name":      networkName,
			"groups":            strings.Join(groups, "\n"),
			"resource_resource": terraformResource,
			"resource_name":     resourceName,
			"tcp_policy":        model.PolicyRestricted,
			"udp_policy":        model.PolicyAllowAll,
			"group_ids":         strings.Join(groupIDs, ", "),
			"authoritative":     authoritative,
		})
}

func TestAccTwingateResourceAccessGroupsAuthoritative(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	groups, groupsID := genNewGroups(terraformResource, 3)
	groupResource := getResourceNameFromID(groupsID[2])

	config := func(groupsID []string) string {
		return configResourceWithGroupsAndAuthoritativeFlag(terraformResource, remoteNetworkName, resourceName, groups, groupsID, true)
	}

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config(groupsID[:1]),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "1"),
					acctests.WaitTestFunc(),
					// added a new group to the resource using API
					acctests.AddResourceGroup(theResource, groupResource),
					acctests.WaitTestFunc(),
					acctests.CheckResourceGroupsLen(theResource, 2),
				),
				// expecting drift - terraform going to remove unknown group
				ExpectNonEmptyPlan: true,
			},
			{
				Config: config(groupsID[:1]),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "1"),
					acctests.CheckResourceGroupsLen(theResource, 1),
				),
			},
			{
				// added 2 new groups to the resource using terraform
				Config: config(groupsID),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "3"),
					acctests.CheckResourceGroupsLen(theResource, 3),
				),
			},
			{
				Config: config(groupsID),
				Check: acctests.ComposeTestCheckFunc(
					// delete one group from the resource using API
					acctests.DeleteResourceGroup(theResource, groupResource),
					acctests.WaitTestFunc(),
					acctests.CheckResourceGroupsLen(theResource, 2),
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "3"),
				),
				// expecting drift - terraform going to restore deleted group
				ExpectNonEmptyPlan: true,
			},
			{
				Config: config(groupsID),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckResourceGroupsLen(theResource, 3),
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "3"),
				),
			},
			{
				// remove 2 groups from the resource using terraform
				Config: config(groupsID[:1]),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "1"),
					acctests.CheckResourceGroupsLen(theResource, 1),
				),
			},
		},
	})
}

func createResource23(networkName, resourceName string, groups, groupsID []string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "test23" {
	  name = "%s"
	}

	%s

	resource "twingate_resource" "test23" {
	  name = "%s"
	  address = "acc-test.com.23"
	  remote_network_id = twingate_remote_network.test23.id
	  
	  protocols = {
	    allow_icmp = true
	    tcp = {
	      policy = "%s"
	      ports = ["80", "82-83"]
	    }
	    udp = {
	      policy = "%s"
	    }
	  }

	  is_authoritative = true
	  access {
	    group_ids = [%s]
	  }

	}
	`, networkName, strings.Join(groups, "\n"), resourceName, model.PolicyRestricted, model.PolicyAllowAll, strings.Join(groupsID, ", "))
}

func TestGetResourceNameFromID(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "twingate_resource.test.id",
			expected: "twingate_resource.test",
		},
		{
			input:    "twingate_resource.test",
			expected: "",
		},
		{
			input:    "",
			expected: "",
		},
	}

	for _, c := range cases {
		actual := getResourceNameFromID(c.input)
		assert.Equal(t, c.expected, actual)
	}
}

func TestAccTwingateCreateResourceWithFlagIsVisible(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceBasic(terraformResource, remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.IsVisible, "true"),
				),
			},
			{
				// expecting no changes - default value on the backend side is `true`
				PlanOnly: true,
				Config:   configResourceWithVisibleFlag(terraformResource, remoteNetworkName, resourceName, true),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.IsVisible, "true"),
				),
			},
			{
				Config: configResourceWithVisibleFlag(terraformResource, remoteNetworkName, resourceName, false),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.IsVisible, "false"),
				),
			},
			{
				// expecting no changes - no drift after re-applying changes
				PlanOnly: true,
				Config:   configResourceWithVisibleFlag(terraformResource, remoteNetworkName, resourceName, false),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.IsVisible, "false"),
				),
			},
			{
				Config: configResourceBasic(terraformResource, remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.IsVisible, "true"),
				),
			},
		},
	})
}

func createSimpleResource(terraformResourceName, networkName, resourceName string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "%s" {
	  name = "%s"
	}
	resource "twingate_resource" "%s" {
	  name = "%s"
	  address = "acc-test.com"
	  remote_network_id = twingate_remote_network.%s.id
	}
	`, terraformResourceName, networkName, terraformResourceName, resourceName, terraformResourceName)
}

func configResourceWithVisibleFlag(terraformResource, networkName, name string, isVisible bool) string {
	return acctests.Nprintf(`
	resource "twingate_remote_network" "${network_resource}" {
	  name = "${network_name}"
	}
	resource "twingate_resource" "${resource_resource}" {
	  name = "${resource_name}"
	  address = "acc-test.com"
	  remote_network_id = twingate_remote_network.${network_resource}.id
	  is_visible = ${visible}
	}
	`,
		map[string]any{
			"network_resource":  terraformResource,
			"network_name":      networkName,
			"resource_resource": terraformResource,
			"resource_name":     name,
			"visible":           isVisible,
		})
}

func TestAccTwingateCreateResourceWithFlagIsBrowserShortcutEnabled(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceBasic(terraformResource, remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.IsBrowserShortcutEnabled, "false"),
				),
			},
			{
				// expecting no changes - default value is `false`
				Config: configResourceWithBrowserShortcutEnabledFlag(terraformResource, remoteNetworkName, resourceName, false),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.IsBrowserShortcutEnabled, "false"),
				),
				PlanOnly: true,
			},
			{
				Config: configResourceWithBrowserShortcutEnabledFlag(terraformResource, remoteNetworkName, resourceName, true),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.IsBrowserShortcutEnabled, "true"),
				),
			},
			{
				// expecting no changes - no drift after re-applying changes
				PlanOnly: true,
				Config:   configResourceWithBrowserShortcutEnabledFlag(terraformResource, remoteNetworkName, resourceName, true),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.IsBrowserShortcutEnabled, "true"),
				),
			},
			{
				Config: configResourceBasic(terraformResource, remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.IsBrowserShortcutEnabled, "false"),
				),
			},
		},
	})
}

func configResourceWithBrowserShortcutEnabledFlag(terraformResource, networkName, name string, isBrowserShortcutEnabled bool) string {
	return acctests.Nprintf(`
	resource "twingate_remote_network" "${network_resource}" {
	  name = "${network_name}"
	}
	resource "twingate_resource" "${resource_resource}" {
	  name = "${resource_name}"
	  address = "acc-test.com"
	  remote_network_id = twingate_remote_network.${network_resource}.id
	  is_browser_shortcut_enabled = ${browser_shortcut_enabled}
	}
	`,
		map[string]any{
			"network_resource":         terraformResource,
			"network_name":             networkName,
			"resource_resource":        terraformResource,
			"resource_name":            name,
			"browser_shortcut_enabled": isBrowserShortcutEnabled,
		})
}

func TestAccTwingateResourceGroupsAuthoritativeByDefault(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	groups, groupsID := genNewGroups(terraformResource, 3)
	groupResource := getResourceNameFromID(groupsID[2])

	config := func(groupsID []string) string {
		return configResourceWithGroups(terraformResource, remoteNetworkName, resourceName, groups, groupsID)
	}

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config(groupsID[:1]),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "1"),
					acctests.WaitTestFunc(),
					// added a new group to the resource using API
					acctests.AddResourceGroup(theResource, groupResource),
					acctests.WaitTestFunc(),
					acctests.CheckResourceGroupsLen(theResource, 2),
				),
				// expecting drift - terraform going to remove unknown group
				ExpectNonEmptyPlan: true,
			},
			{
				Config: config(groupsID[:1]),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "1"),
					acctests.CheckResourceGroupsLen(theResource, 1),
				),
			},
			{
				// added 2 new groups to the resource using terraform
				Config: config(groupsID),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "3"),
					acctests.CheckResourceGroupsLen(theResource, 3),
				),
			},
			{
				Config: config(groupsID),
				Check: acctests.ComposeTestCheckFunc(
					// delete one group from the resource using API
					acctests.DeleteResourceGroup(theResource, groupResource),
					acctests.WaitTestFunc(),
					acctests.CheckResourceGroupsLen(theResource, 2),
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "3"),
				),
				// expecting drift - terraform going to restore deleted group
				ExpectNonEmptyPlan: true,
			},
			{
				Config: config(groupsID),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckResourceGroupsLen(theResource, 3),
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "3"),
				),
			},
			{
				// remove 2 groups from the resource using terraform
				Config: config(groupsID[:1]),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "1"),
					acctests.CheckResourceGroupsLen(theResource, 1),
				),
			},
		},
	})
}

func createResource26(networkName, resourceName string, groups, groupsID []string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "test26" {
	  name = "%s"
	}

	%s

	resource "twingate_resource" "test26" {
	  name = "%s"
	  address = "acc-test.com.26"
	  remote_network_id = twingate_remote_network.test26.id
	  
	  protocols = {
	    allow_icmp = true
	    tcp = {
	      policy = "%s"
	      ports = ["80", "82-83"]
	    }
	    udp = {
	      policy = "%s"
	    }
	  }

	  access {
	    group_ids = [%s]
	  }

	}
	`, networkName, strings.Join(groups, "\n"), resourceName, model.PolicyRestricted, model.PolicyAllowAll, strings.Join(groupsID, ", "))
}

func TestAccTwingateResourceDoesNotSupportOldGroups(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	groups, groupsID := genNewGroups(terraformResource, 2)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config:      configResourceWithOldGroups(terraformResource, remoteNetworkName, resourceName, groups, groupsID),
				ExpectError: regexp.MustCompile("Error: Unsupported argument"),
			},
		},
	})
}

func configResourceWithOldGroups(terraformResource, networkName, resourceName string, groups, groupsID []string) string {
	return acctests.Nprintf(`
    resource "twingate_remote_network" "${network_resource}" {
      name = "${network_name}"
    }

    ${group}

    resource "twingate_resource" "${resource_resource}" {
      name = "${resource_name}"
      address = "acc-test.com"
      remote_network_id = twingate_remote_network.${network_resource}.id
      protocols = {
        allow_icmp = true
        tcp = {
          policy = "${tcp_policy}"
          ports = ["80", "82-83"]
        }
        udp = {
          policy = "${udp_policy}"
        }
      }

	  group_ids = [${group_id}]
    }
    `,
		map[string]any{
			"network_resource":  terraformResource,
			"network_name":      networkName,
			"group":             strings.Join(groups, "\n"),
			"resource_resource": terraformResource,
			"resource_name":     resourceName,
			"tcp_policy":        model.PolicyRestricted,
			"udp_policy":        model.PolicyAllowAll,
			"group_id":          strings.Join(groupsID, ", "),
		})
}

func TestAccTwingateResourceCreateWithAlias(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	const aliasName = "test.com"

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithAlias(terraformResource, remoteNetworkName, resourceName, aliasName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, resourceName),
					sdk.TestCheckResourceAttr(theResource, attr.Alias, aliasName),
				),
			},
			{
				Config: configResourceBasic(terraformResource, remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckNoResourceAttr(theResource, attr.Alias),
				),
			},
			{
				// alias attr set with empty string
				Config: configResourceWithAlias(terraformResource, remoteNetworkName, resourceName, ""),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.Alias, ""),
				),
			},
		},
	})
}

func configResourceWithAlias(terraformResource, networkName, name, alias string) string {
	return acctests.Nprintf(`
	resource "twingate_remote_network" "${network_resource}" {
	  name = "${network_name}"
	}
	resource "twingate_resource" "${resource_resource}" {
	  name = "${resource_name}"
	  address = "acc-test.com"
	  remote_network_id = twingate_remote_network.${network_resource}.id
	  alias = "${alias}"
	}
	`,
		map[string]any{
			"network_resource":  terraformResource,
			"network_name":      networkName,
			"resource_resource": terraformResource,
			"resource_name":     name,
			"alias":             alias,
		})
}

func TestAccTwingateResourceGroupsCursor(t *testing.T) {
	acctests.SetPageLimit(t, 1)

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	groups, groupsID := genNewGroups(terraformResource, 3)
	serviceAccounts, serviceAccountIDs := genNewServiceAccounts(terraformResource, 3)

	config := func(groupsID, serviceAccountIDs []string) string {
		return configResourceWithGroupsAndServiceAccounts(terraformResource, remoteNetworkName, resourceName, groups, groupsID, serviceAccounts, serviceAccountIDs)
	}

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config(groupsID, serviceAccountIDs),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "3"),
					sdk.TestCheckResourceAttr(theResource, accessServiceAccountIdsLen, "3"),
				),
			},
			{
				Config: config(groupsID[:2], serviceAccountIDs[:2]),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "2"),
					sdk.TestCheckResourceAttr(theResource, accessServiceAccountIdsLen, "2"),
				),
			},
		},
	})
}

func createResourceWithGroupsAndServiceAccounts(name, networkName, resourceName string, groups, groupsID, serviceAccounts, serviceAccountIDs []string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "%s" {
	  name = "%s"
	}

	%s

	%s

	resource "twingate_resource" "%s" {
	  name = "%s"
	  address = "acc-test.com.26"
	  remote_network_id = twingate_remote_network.%s.id
	  
	  protocols = {
	    allow_icmp = true
	    tcp = {
	      policy = "%s"
	      ports = ["80", "82-83"]
	    }
	    udp = {
	      policy = "%s"
	    }
	  }

	  access {
	    group_ids = [%s]
	    service_account_ids = [%s]
	  }

	}
	`, name, networkName, strings.Join(groups, "\n"), strings.Join(serviceAccounts, "\n"), name, resourceName, name, model.PolicyRestricted, model.PolicyAllowAll, strings.Join(groupsID, ", "), strings.Join(serviceAccountIDs, ", "))
}

func TestAccTwingateResourceCreateWithPort(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	config := func(portRange string) string {
		return configResourceWithPortRange(terraformResource, remoteNetworkName, resourceName, portRange)
	}

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config:      config(`"0"`),
				ExpectError: regexp.MustCompile("port 0 not in the range of 1-65535"),
			},
			{
				Config:      config(`"65536"`),
				ExpectError: regexp.MustCompile("port 65536 not in the range"),
			},
			{
				Config:      config(`"0-10"`),
				ExpectError: regexp.MustCompile("port 0 not in the range"),
			},
			{
				Config:      config(`"65535-65536"`),
				ExpectError: regexp.MustCompile("port 65536 not in the[\\n\\s]+range"),
			},
		},
	})
}

func createResourceWithPort(networkName, resourceName, port string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "test30" {
	  name = "%s"
	}
	resource "twingate_resource" "test30" {
	  name = "%s"
	  address = "new-acc-test.com"
	  remote_network_id = twingate_remote_network.test30.id
	  protocols = {
		allow_icmp = true
		tcp = {
			policy = "%s"
			ports = ["%s"]
		}
		udp = {
			policy = "%s"
		}
	  }
	}
	`, networkName, resourceName, model.PolicyRestricted, port, model.PolicyAllowAll)
}

func TestAccTwingateResourceUpdateWithPort(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	config := func(portRange string) string {
		return configResourceWithPortRange(terraformResource, remoteNetworkName, resourceName, portRange)
	}

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config(`"1"`),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, firstTCPPort, "1"),
				),
			},
			{
				Config:      config(`"0"`),
				ExpectError: regexp.MustCompile("port 0 not in the range of 1-65535"),
			},
		},
	})
}

func TestAccTwingateResourceWithPortsFailsForAllowAllAndDenyAllPolicy(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	config := func(policy string) string {
		return configResourceWithPolicyAndPortRange(terraformResource, remoteNetworkName, resourceName, policy, `"80", "82-83"`)
	}

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config:      config(model.PolicyAllowAll),
				ExpectError: regexp.MustCompile(resource.ErrPortsWithPolicyAllowAll.Error()),
			},
			{
				Config:      config(model.PolicyDenyAll),
				ExpectError: regexp.MustCompile(resource.ErrPortsWithPolicyDenyAll.Error()),
			},
		},
	})
}

func configResourceWithPolicyAndPortRange(terraformResource, networkName, resourceName, policy, portRange string) string {
	return acctests.Nprintf(`
    resource "twingate_remote_network" "${network_resource}" {
      name = "${network_name}"
    }

    resource "twingate_resource" "${resource_resource}" {
      name = "${resource_name}"
      address = "new-acc-test.com"
      remote_network_id = twingate_remote_network.${network_resource}.id
      
      protocols = {
        allow_icmp = true
        tcp = {
          policy = "${policy}"
          ports = [${port_range}]
        }
        udp = {
          policy = "${policy}"
          ports = [${port_range}]
        }
      }
    }
    `,
		map[string]any{
			"network_resource":  terraformResource,
			"network_name":      networkName,
			"resource_resource": terraformResource,
			"resource_name":     resourceName,
			"policy":            policy,
			"port_range":        portRange,
		})
}

func TestAccTwingateResourceWithoutPortsOkForAllowAllAndDenyAllPolicy(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)

	config := func(policy string) string {
		return configResourceWithPolicy(terraformResource, remoteNetworkName, resourceName, policy, policy)
	}

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config(model.PolicyAllowAll),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyAllowAll),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "0"),
				),
			},
			{
				Config: config(model.PolicyDenyAll),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyDenyAll),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "0"),
				),
			},
		},
	})
}

func createResourceWithoutPorts(name, networkName, resourceName, policy string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "%[1]s" {
	  name = "%[2]s"
	}
	resource "twingate_resource" "%[1]s" {
	  name = "%[3]s"
	  address = "acc-test-%[1]s.com"
	  remote_network_id = twingate_remote_network.%[1]s.id
	  
	  protocols = {
	    allow_icmp = true
	    tcp = {
	      policy = "%[4]s"
	    }
	    udp = {
	      policy = "%[5]s"
	    }
	  }
	}
	`, name, networkName, resourceName, policy, model.PolicyAllowAll)
}

func TestAccTwingateResourceWithRestrictedPolicy(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)

	config := func(policy string) string {
		return configResourceWithPolicyAndPortRange(terraformResource, remoteNetworkName, resourceName, policy, `"80", "82-83"`)
	}

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config(model.PolicyRestricted),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyRestricted),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "2"),
				),
			},
		},
	})
}

func TestAccTwingateResourcePolicyTransitionDenyAllToRestricted(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithPolicy(terraformResource, remoteNetworkName, resourceName, model.PolicyDenyAll, model.PolicyDenyAll),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyDenyAll),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "0"),
				),
			},
			{
				Config: configResourceWithPolicyAndPortRange(terraformResource, remoteNetworkName, resourceName, model.PolicyRestricted, `"80", "82-83"`),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyRestricted),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "2"),
				),
			},
		},
	})
}

func TestAccTwingateResourcePolicyTransitionDenyAllToAllowAll(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	config := func(policy string) string {
		return configResourceWithPolicy(terraformResource, remoteNetworkName, resourceName, policy, policy)
	}

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config(model.PolicyDenyAll),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyDenyAll),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "0"),
				),
			},
			{
				Config: config(model.PolicyAllowAll),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyAllowAll),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "0"),
				),
			},
		},
	})
}

func TestAccTwingateResourcePolicyTransitionRestrictedToDenyAll(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithPolicyAndPortRange(terraformResource, remoteNetworkName, resourceName, model.PolicyRestricted, `"80", "82-83"`),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyRestricted),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "2"),
				),
			},
			{
				Config: configResourceWithPolicy(terraformResource, remoteNetworkName, resourceName, model.PolicyDenyAll, model.PolicyDenyAll),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyDenyAll),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "0"),
				),
			},
		},
	})
}

func TestAccTwingateResourcePolicyTransitionRestrictedToAllowAll(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithPolicyAndPortRange(terraformResource, remoteNetworkName, resourceName, model.PolicyRestricted, `"80", "82-83"`),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyRestricted),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "2"),
				),
			},
			{
				Config: configResourceWithPolicy(terraformResource, remoteNetworkName, resourceName, model.PolicyAllowAll, model.PolicyAllowAll),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyAllowAll),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "0"),
				),
			},
		},
	})
}

func TestAccTwingateResourcePolicyTransitionRestrictedToAllowAllWithPortsShouldFail(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	config := func(policy string) string {
		return configResourceWithPolicyAndPortRange(terraformResource, remoteNetworkName, resourceName, policy, `"80", "82-83"`)
	}

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config(model.PolicyRestricted),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyRestricted),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "2"),
				),
			},
			{
				Config:      config(model.PolicyAllowAll),
				ExpectError: regexp.MustCompile(resource.ErrPortsWithPolicyAllowAll.Error()),
			},
		},
	})
}

func TestAccTwingateResourcePolicyTransitionAllowAllToRestricted(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithPolicy(terraformResource, remoteNetworkName, resourceName, model.PolicyAllowAll, model.PolicyAllowAll),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyAllowAll),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "0"),
				),
			},
			{
				Config: configResourceWithPolicyAndPortRange(terraformResource, remoteNetworkName, resourceName, model.PolicyRestricted, `"80", "82-83"`),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyRestricted),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "2"),
				),
			},
		},
	})
}

func TestAccTwingateResourcePolicyTransitionAllowAllToDenyAll(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	config := func(policy string) string {
		return configResourceWithPolicy(terraformResource, remoteNetworkName, resourceName, policy, policy)
	}

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: config(model.PolicyAllowAll),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyAllowAll),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "0"),
				),
			},
			{
				Config: config(model.PolicyDenyAll),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyDenyAll),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "0"),
				),
			},
		},
	})
}

func TestAccTwingateResourceTestCaseInsensitiveAlias(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	const aliasName = "test.com"

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithAlias(terraformResource, remoteNetworkName, resourceName, aliasName),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.Alias, aliasName),
				),
			},
			{
				// expecting no changes
				PlanOnly: true,
				Config:   configResourceWithAlias(terraformResource, remoteNetworkName, resourceName, strings.ToUpper(aliasName)),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.Alias, aliasName),
				),
			},
		},
	})
}

func TestAccTwingateResourceWithBrowserOption(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	wildcardAddress := "*.acc-test.com"

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithAddress(terraformResource, remoteNetworkName, resourceName, wildcardAddress),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
			{
				Config: configResourceWithAddressAndBrowserOption(terraformResource, remoteNetworkName, resourceName, wildcardAddress, false),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
			{
				Config:      configResourceWithAddressAndBrowserOption(terraformResource, remoteNetworkName, resourceName, wildcardAddress, true),
				ExpectError: regexp.MustCompile("Resources with a CIDR range or wildcard"),
			},
		},
	})
}

func TestAccTwingateResourceWithBrowserOptionFailOnUpdate(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	wildcardAddress := "*.acc-test.com"
	simpleAddress := "acc-test.com"

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithAddress(terraformResource, remoteNetworkName, resourceName, simpleAddress),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
			{
				Config: configResourceWithAddressAndBrowserOption(terraformResource, remoteNetworkName, resourceName, simpleAddress, true),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
			{
				Config:      configResourceWithAddressAndBrowserOption(terraformResource, remoteNetworkName, resourceName, wildcardAddress, true),
				ExpectError: regexp.MustCompile("Resources with a CIDR range or wildcard"),
			},
		},
	})
}

func TestAccTwingateResourceWithBrowserOptionRecovered(t *testing.T) {
	t.Parallel()

	terraformResource := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResource)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	wildcardAddress := "*.acc-test.com"
	simpleAddress := "acc-test.com"

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithAddressAndBrowserOption(terraformResource, remoteNetworkName, resourceName, simpleAddress, true),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
			{
				Config: configResourceWithAddress(terraformResource, remoteNetworkName, resourceName, wildcardAddress),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
		},
	})
}

func configResourceWithAddress(terraformResource, networkName, name, address string) string {
	return acctests.Nprintf(`
	resource "twingate_remote_network" "${network_resource}" {
	  name = "${network_name}"
	}
	resource "twingate_resource" "${resource_resource}" {
	  name = "${resource_name}"
	  address = "${address}"
	  remote_network_id = twingate_remote_network.${network_resource}.id
	}
	`,
		map[string]any{
			"network_resource":  terraformResource,
			"network_name":      networkName,
			"resource_resource": terraformResource,
			"resource_name":     name,
			"address":           address,
		})
}

func configResourceWithAddressAndBrowserOption(terraformResource, networkName, name, address string, browserFlag bool) string {
	return acctests.Nprintf(`
	resource "twingate_remote_network" "${network_resource}" {
	  name = "${network_name}"
	}
	resource "twingate_resource" "${resource_resource}" {
	  name = "${resource_name}"
	  address = "${address}"
	  remote_network_id = twingate_remote_network.${network_resource}.id
	  is_browser_shortcut_enabled = ${browser_flag}
	}
	`,
		map[string]any{
			"network_resource":  terraformResource,
			"network_name":      networkName,
			"resource_resource": terraformResource,
			"resource_name":     name,
			"address":           address,
			"browser_flag":      browserFlag,
		})
}

func createResourceWithSecurityPolicy(remoteNetwork, resource, policyID string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "%[1]s" {
	  name = "%[1]s"
	}
	resource "twingate_resource" "%[2]s" {
	  name = "%[2]s"
	  address = "acc-test-address.com"
	  remote_network_id = twingate_remote_network.%[1]s.id
	  security_policy_id = "%[3]s"
	}
	`, remoteNetwork, resource, policyID)
}

func createResourceWithoutSecurityPolicy(remoteNetwork, resource string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "%[1]s" {
	  name = "%[1]s"
	}
	resource "twingate_resource" "%[2]s" {
	  name = "%[2]s"
	  address = "acc-test-address.com"
	  remote_network_id = twingate_remote_network.%[1]s.id
	}
	`, remoteNetwork, resource)
}

func TestAccTwingateResourceUpdateWithDefaultProtocols(t *testing.T) {
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	theResource := acctests.TerraformResource(resourceName)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResourceWithProtocols(remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
			{
				Config: createResourceWithoutProtocols(remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
		},
	})
}

func createResourceWithProtocols(remoteNetwork, resource string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "%[1]s" {
	  name = "%[1]s"
	}
	resource "twingate_resource" "%[2]s" {
	  name = "%[2]s"
	  address = "acc-test-address.com"
	  remote_network_id = twingate_remote_network.%[1]s.id
	  protocols = {
	    allow_icmp = true
	    tcp = {
	      policy = "RESTRICTED"
	      ports = ["80-83"]
	    }
	    udp = {
	      policy = "RESTRICTED"
	      ports = ["80"]
	    }
	  }
	}
	`, remoteNetwork, resource)
}

func createResourceWithoutProtocols(remoteNetwork, resource string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "%[1]s" {
	  name = "%[1]s"
	}
	resource "twingate_resource" "%[2]s" {
	  name = "%[2]s"
	  address = "acc-test-address.com"
	  remote_network_id = twingate_remote_network.%[1]s.id
	}
	`, remoteNetwork, resource)
}

func TestAccTwingateResourceUpdatePortsFromEmptyListToNull(t *testing.T) {
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	theResource := acctests.TerraformResource(resourceName)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResourceWithEmptyArrayPorts(remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
			{
				// expect no changes
				PlanOnly: true,
				Config:   createResourceWithDefaultPorts(remoteNetworkName, resourceName),
			},
		},
	})
}

func TestAccTwingateResourceUpdatePortsFromNullToEmptyList(t *testing.T) {
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	theResource := acctests.TerraformResource(resourceName)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResourceWithDefaultPorts(remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
			{
				// expect no changes
				PlanOnly: true,
				Config:   createResourceWithEmptyArrayPorts(remoteNetworkName, resourceName),
			},
		},
	})
}

func createResourceWithDefaultPorts(remoteNetwork, resource string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "%[1]s" {
	  name = "%[1]s"
	}
	resource "twingate_resource" "%[2]s" {
	  name = "%[2]s"
	  address = "acc-test-address.com"
	  remote_network_id = twingate_remote_network.%[1]s.id
	  protocols = {
	    allow_icmp = true
	    tcp = {
	      policy = "ALLOW_ALL"
	    }
	    udp = {
	      policy = "ALLOW_ALL"
	    }
	  }
	}
	`, remoteNetwork, resource)
}

func createResourceWithEmptyArrayPorts(remoteNetwork, resource string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "%[1]s" {
	  name = "%[1]s"
	}
	resource "twingate_resource" "%[2]s" {
	  name = "%[2]s"
	  address = "acc-test-address.com"
	  remote_network_id = twingate_remote_network.%[1]s.id
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
	`, remoteNetwork, resource)
}

func TestAccTwingateResourceUpdateSecurityPolicy(t *testing.T) {
	t.Parallel()

	resourceName := test.RandomResourceName()
	theResource := acctests.TerraformResource(resourceName)
	remoteNetworkName := test.RandomName()

	defaultPolicy, testPolicy := preparePolicies(t)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResourceWithSecurityPolicy(remoteNetworkName, resourceName, defaultPolicy),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.SecurityPolicyID, defaultPolicy),
				),
			},
			{
				Config: createResourceWithSecurityPolicy(remoteNetworkName, resourceName, testPolicy),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.SecurityPolicyID, testPolicy),
				),
			},
			{
				Config: createResourceWithoutSecurityPolicy(remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.SecurityPolicyID, defaultPolicy),
				),
			},
			{
				Config: createResourceWithSecurityPolicy(remoteNetworkName, resourceName, ""),
				// no changes
				PlanOnly: true,
			},
		},
	})
}

func preparePolicies(t *testing.T) (string, string) {
	policies, err := acctests.ListSecurityPolicies()
	if err != nil {
		t.Skipf("failed to retrieve security policies: %v", err)
	}

	if len(policies) < 2 {
		t.Skip("requires at least 2 security policy for the test")
	}

	var defaultPolicy, testPolicy string
	if policies[0].Name == resource.DefaultSecurityPolicyName {
		defaultPolicy = policies[0].ID
		testPolicy = policies[1].ID
	} else {
		testPolicy = policies[0].ID
		defaultPolicy = policies[1].ID
	}

	return defaultPolicy, testPolicy
}

func TestAccTwingateResourceSetDefaultSecurityPolicyByDefault(t *testing.T) {
	t.Parallel()

	resourceName := test.RandomResourceName()
	theResource := acctests.TerraformResource(resourceName)
	remoteNetworkName := test.RandomName()

	defaultPolicy, testPolicy := preparePolicies(t)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResourceWithSecurityPolicy(remoteNetworkName, resourceName, testPolicy),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.SecurityPolicyID, testPolicy),
				),
			},
			{
				Config: createResourceWithoutSecurityPolicy(remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.SecurityPolicyID, defaultPolicy),
					acctests.CheckResourceSecurityPolicy(theResource, defaultPolicy),
					// set new policy via API
					acctests.UpdateResourceSecurityPolicy(theResource, testPolicy),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: createResourceWithSecurityPolicy(remoteNetworkName, resourceName, ""),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckResourceSecurityPolicy(theResource, defaultPolicy),
				),
			},
			{
				Config: createResourceWithoutSecurityPolicy(remoteNetworkName, resourceName),
				// no changes
				PlanOnly: true,
			},
		},
	})
}

func TestAccTwingateResourceSecurityPolicy(t *testing.T) {
	t.Parallel()

	resourceName := test.RandomResourceName()
	theResource := acctests.TerraformResource(resourceName)
	remoteNetworkName := test.RandomName()

	_, testPolicy := preparePolicies(t)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResourceWithoutSecurityPolicy(remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckNoResourceAttr(theResource, attr.SecurityPolicyID),
				),
			},
			{
				Config: createResourceWithSecurityPolicy(remoteNetworkName, resourceName, testPolicy),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.SecurityPolicyID, testPolicy),
				),
			},
		},
	})
}

func TestAccTwingateResourceCreateInactive(t *testing.T) {
	t.Parallel()

	resourceName := test.RandomResourceName()
	theResource := acctests.TerraformResource(resourceName)
	remoteNetworkName := test.RandomName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResourceWithIsActiveFlag(remoteNetworkName, resourceName, false),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.IsActive, "false"),
					acctests.CheckTwingateResourceActiveState(theResource, false),
				),
			},
		},
	})
}

func createResourceWithIsActiveFlag(networkName, resourceName string, isActive bool) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "%[1]s" {
	  name = "%[1]s"
	}
	resource "twingate_resource" "%[2]s" {
	  name = "%[2]s"
	  address = "acc-test.com"
	  remote_network_id = twingate_remote_network.%[1]s.id
	  is_active = %[3]v
	}
	`, networkName, resourceName, isActive)
}

func TestAccTwingateResourceTestInactiveFlag(t *testing.T) {
	t.Parallel()

	resourceName := test.RandomResourceName()
	theResource := acctests.TerraformResource(resourceName)
	remoteNetworkName := test.RandomName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResourceWithIsActiveFlag(remoteNetworkName, resourceName, true),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.IsActive, "true"),
				),
			},
			{
				Config: createResourceWithIsActiveFlag(remoteNetworkName, resourceName, false),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.IsActive, "false"),
					acctests.CheckTwingateResourceActiveState(theResource, false),
				),
			},
		},
	})
}

func TestAccTwingateResourceTestPlanOnDisabledResource(t *testing.T) {
	t.Parallel()

	resourceName := test.RandomResourceName()
	theResource := acctests.TerraformResource(resourceName)
	remoteNetworkName := test.RandomName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResource(remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceActiveState(theResource, true),
					acctests.DeactivateTwingateResource(theResource),
					acctests.CheckTwingateResourceActiveState(theResource, false),
				),
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: sdk.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
						plancheck.ExpectResourceAction(theResource, plancheck.ResourceActionUpdate),
						acctests.CheckResourceActiveState(theResource, false),
					},
				},
			},
		},
	})
}

func createResource(networkName, resourceName string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "%[1]s" {
	  name = "%[1]s"
	}
	resource "twingate_resource" "%[2]s" {
	  name = "%[2]s"
	  address = "acc-test.com"
	  remote_network_id = twingate_remote_network.%[1]s.id
	}
	`, networkName, resourceName)
}
