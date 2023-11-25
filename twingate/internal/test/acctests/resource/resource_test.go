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
	"github.com/stretchr/testify/assert"
)

var (
	tcpPolicy                  = attr.Path(attr.Protocols, attr.TCP, attr.Policy)
	udpPolicy                  = attr.Path(attr.Protocols, attr.UDP, attr.Policy)
	firstTCPPort               = attr.First(attr.Protocols, attr.TCP, attr.Ports)
	firstUDPPort               = attr.First(attr.Protocols, attr.UDP, attr.Ports)
	tcpPortsLen                = attr.Len(attr.Protocols, attr.TCP, attr.Ports)
	udpPortsLen                = attr.Len(attr.Protocols, attr.UDP, attr.Ports)
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
				Config: configResourceWithNetwork(resourceName, remoteNetworkName, name),
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
				Config: configResourceWithNetwork(resourceName, remoteNetworkName, name),
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
				Config: configResourceWithNetwork(resourceName, remoteNetworkName, name),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
		},
	})
}

func configResourceWithNetwork(terraformResource, networkName, name string) string {
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

	  protocols {
	    allow_icmp = true
	    tcp  {
	      policy = "DENY_ALL"
	    }
	    udp {
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

      protocols {
		allow_icmp = true
        tcp  {
			policy = "${tcp_policy}"
            ports = ["80", "82-83"]
        }
		udp {
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

      protocols {
        allow_icmp = true
        tcp  {
            policy = "${tcp_policy}"
            ports = ["3306"]
        }
        udp {
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
	groupName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithPolicy(terraformResource, networkName, groupName, resourceName, model.PolicyDenyAll, model.PolicyAllowAll),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyDenyAll),
				),
			},
			// expecting no changes - empty plan
			{
				Config:   configResourceWithPolicy(terraformResource, networkName, groupName, resourceName, model.PolicyDenyAll, model.PolicyAllowAll),
				PlanOnly: true,
			},
		},
	})
}

func configResourceWithPolicy(terraformResource, networkName, groupName, resourceName, tcpPolicy, udpPolicy string) string {
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
      protocols {
        allow_icmp = true
        tcp {
          policy = "${tcp_policy}"
        }
        udp {
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

func TestAccTwingateResourceWithUdpDenyAllPolicy(t *testing.T) {
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
				Config: configResourceWithPolicy(terraformResource, remoteNetworkName, groupName, resourceName, model.PolicyAllowAll, model.PolicyDenyAll),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, udpPolicy, model.PolicyDenyAll),
				),
			},
			// expecting no changes - empty plan
			{
				Config:   configResourceWithPolicy(terraformResource, remoteNetworkName, groupName, resourceName, model.PolicyAllowAll, model.PolicyDenyAll),
				PlanOnly: true,
			},
		},
	})
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
      protocols {
        allow_icmp = true
        tcp {
          policy = "${tcp_policy}"
          ports = []
        }
        udp {
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
      protocols {
        allow_icmp = true
        tcp {
          policy = "${tcp_policy}"
          ports = [${port_range}]
        }
        udp {
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
					sdk.TestCheckResourceAttr(theResource, firstTCPPort, "82"),
					sdk.TestCheckResourceAttr(theResource, firstUDPPort, "82"),
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
				Config: configResourceWithNetwork(terraformResource, remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					acctests.DeactivateTwingateResource(theResource),
					acctests.WaitTestFunc(),
					acctests.CheckTwingateResourceActiveState(theResource, false),
				),
			},
			{
				Config: configResourceWithNetwork(terraformResource, remoteNetworkName, resourceName),
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
				Config: configResourceWithNetwork(terraformResource, remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					acctests.DeleteTwingateResource(theResource, resource.TwingateResource),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: configResourceWithNetwork(terraformResource, remoteNetworkName, resourceName),
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
				}),
			},
		},
	})
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
      protocols {
        allow_icmp = true
        tcp {
          policy = "${tcp_policy}"
          ports = ["80", "82-83"]
        }
        udp {
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
	serviceAccountConfig := configServiceAccount(terraformResource, test.RandomName())
	serviceAccountID := acctests.TerraformServiceAccount(terraformResource) + ".id"

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configResourceWithGroupsAndServiceAccount(terraformResource, test.RandomName(), test.RandomResourceName(), groups, groupsID, serviceAccountConfig, serviceAccountID),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "1"),
					sdk.TestCheckResourceAttr(theResource, accessServiceAccountIdsLen, "1"),
				),
			},
		},
	})
}

func configResourceWithGroupsAndServiceAccount(terraformResource, networkName, resourceName string, groups, groupsID []string, serviceAccountConfig, serviceAccountID string) string {
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
      protocols {
        allow_icmp = true
        tcp {
          policy = "${tcp_policy}"
          ports = ["80", "82-83"]
        }
        udp {
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
			"service_account":    serviceAccountConfig,
			"group":              strings.Join(groups, "\n"),
			"resource_resource":  terraformResource,
			"resource_name":      resourceName,
			"tcp_policy":         model.PolicyRestricted,
			"udp_policy":         model.PolicyAllowAll,
			"service_account_id": serviceAccountID,
			"group_id":           strings.Join(groupsID, ", "),
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
      protocols {
        allow_icmp = true
        tcp {
          policy = "${tcp_policy}"
          ports = ["80", "82-83"]
        }
        udp {
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
				ExpectError: regexp.MustCompile("Error: Not enough list items"),
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
      protocols {
        allow_icmp = true
        tcp {
          policy = "${tcp_policy}"
          ports = ["80", "82-83"]
        }
        udp {
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
				ExpectError: regexp.MustCompile("Error: Not enough list items"),
			},
		},
	})
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
				ExpectError: regexp.MustCompile("Missing required argument"),
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
      protocols {
        allow_icmp = true
        tcp {
          policy = "${tcp_policy}"
          ports = ["80", "82-83"]
        }
        udp {
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
      protocols {
        allow_icmp = true
        tcp {
          policy = "${tcp_policy}"
          ports = ["80", "82-83"]
        }
        udp {
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
	const terraformResourceName = "test24"
	theResource := acctests.TerraformResource(terraformResourceName)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createSimpleResource(terraformResourceName, remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckNoResourceAttr(theResource, attr.IsVisible),
				),
			},
			{
				// expecting no changes - default value on the backend side is `true`
				PlanOnly: true,
				Config:   createResourceWithFlagIsVisible(terraformResourceName, remoteNetworkName, resourceName, true),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.IsVisible, "true"),
				),
			},
			{
				Config: createResourceWithFlagIsVisible(terraformResourceName, remoteNetworkName, resourceName, false),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.IsVisible, "false"),
				),
			},
			{
				// expecting no changes - no drift after re-applying changes
				PlanOnly: true,
				Config:   createResourceWithFlagIsVisible(terraformResourceName, remoteNetworkName, resourceName, false),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.IsVisible, "false"),
				),
			},
			{
				// expecting no changes - flag not set
				PlanOnly: true,
				Config:   createSimpleResource(terraformResourceName, remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckNoResourceAttr(theResource, attr.IsVisible),
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

func createResourceWithFlagIsVisible(terraformResourceName, networkName, resourceName string, isVisible bool) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "%s" {
	  name = "%s"
	}
	resource "twingate_resource" "%s" {
	  name = "%s"
	  address = "acc-test.com"
	  remote_network_id = twingate_remote_network.%s.id
	  is_visible = %v
	}
	`, terraformResourceName, networkName, terraformResourceName, resourceName, terraformResourceName, isVisible)
}

func TestAccTwingateCreateResourceWithFlagIsBrowserShortcutEnabled(t *testing.T) {
	const terraformResourceName = "test25"
	theResource := acctests.TerraformResource(terraformResourceName)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createSimpleResource(terraformResourceName, remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckNoResourceAttr(theResource, attr.IsBrowserShortcutEnabled),
				),
			},
			{
				// expecting no changes - default value on the backend side is `true`
				PlanOnly: true,
				Config:   createResourceWithFlagIsBrowserShortcutEnabled(terraformResourceName, remoteNetworkName, resourceName, true),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.IsBrowserShortcutEnabled, "true"),
				),
			},
			{
				Config: createResourceWithFlagIsBrowserShortcutEnabled(terraformResourceName, remoteNetworkName, resourceName, false),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.IsBrowserShortcutEnabled, "false"),
				),
			},
			{
				// expecting no changes - no drift after re-applying changes
				PlanOnly: true,
				Config:   createResourceWithFlagIsBrowserShortcutEnabled(terraformResourceName, remoteNetworkName, resourceName, false),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.IsBrowserShortcutEnabled, "false"),
				),
			},
			{
				// expecting no changes - flag not set
				PlanOnly: true,
				Config:   createSimpleResource(terraformResourceName, remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckNoResourceAttr(theResource, attr.IsBrowserShortcutEnabled),
				),
			},
		},
	})
}

func createResourceWithFlagIsBrowserShortcutEnabled(terraformResourceName, networkName, resourceName string, isBrowserShortcutEnabled bool) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "%s" {
	  name = "%s"
	}
	resource "twingate_resource" "%s" {
	  name = "%s"
	  address = "acc-test.com"
	  remote_network_id = twingate_remote_network.%s.id
	  is_browser_shortcut_enabled = %v
	}
	`, terraformResourceName, networkName, terraformResourceName, resourceName, terraformResourceName, isBrowserShortcutEnabled)
}

func TestAccTwingateResourceGroupsAuthoritativeByDefault(t *testing.T) {
	t.Parallel()
	const theResource = "twingate_resource.test26"
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	groups, groupsID := genNewGroups("g26", 3)

	groupResource := getResourceNameFromID(groupsID[2])

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResource26(remoteNetworkName, resourceName, groups, groupsID[:1]),
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
				Config: createResource26(remoteNetworkName, resourceName, groups, groupsID[:1]),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "1"),
					acctests.CheckResourceGroupsLen(theResource, 1),
				),
			},
			{
				// added 2 new groups to the resource using terraform
				Config: createResource26(remoteNetworkName, resourceName, groups, groupsID),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "3"),
					acctests.CheckResourceGroupsLen(theResource, 3),
				),
			},
			{
				Config: createResource26(remoteNetworkName, resourceName, groups, groupsID),
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
				Config: createResource26(remoteNetworkName, resourceName, groups, groupsID),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckResourceGroupsLen(theResource, 3),
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "3"),
				),
			},
			{
				// remove 2 groups from the resource using terraform
				Config: createResource26(remoteNetworkName, resourceName, groups, groupsID[:1]),
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
	  
	  protocols {
	    allow_icmp = true
	    tcp {
	      policy = "%s"
	      ports = ["80", "82-83"]
	    }
	    udp {
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
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	groups, groupsID := genNewGroups("g28", 2)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config:      createResource28(remoteNetworkName, resourceName, groups, groupsID),
				ExpectError: regexp.MustCompile("Error: Unsupported argument"),
			},
		},
	})
}

func createResource28(networkName, resourceName string, groups, groupsID []string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "test28" {
	  name = "%s"
	}

	%s

	resource "twingate_resource" "test28" {
	  name = "%s"
	  address = "acc-test.com.28"
	  remote_network_id = twingate_remote_network.test28.id
	
	  group_ids = [%s]

	  protocols {
	    allow_icmp = true
	    tcp {
	      policy = "%s"
	      ports = ["80", "82-83"]
	    }
	    udp {
	      policy = "%s"
	    }
	  }

	}
	`, networkName, strings.Join(groups, "\n"), resourceName, strings.Join(groupsID, ", "), model.PolicyRestricted, model.PolicyAllowAll)
}

func TestAccTwingateResourceCreateWithAlias(t *testing.T) {
	const terraformResourceName = "test29"
	theResource := acctests.TerraformResource(terraformResourceName)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	const aliasName = "test.com"

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResource29(terraformResourceName, remoteNetworkName, resourceName, aliasName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, attr.Name, resourceName),
					sdk.TestCheckResourceAttr(theResource, attr.Alias, aliasName),
				),
			},
			{
				// alias attr commented out, means state keeps the same value without changes
				Config: createResource29WithoutAlias(terraformResourceName, remoteNetworkName, resourceName),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.Alias, aliasName),
				),
			},
			{
				// alias attr set with emtpy string
				Config: createResource29(terraformResourceName, remoteNetworkName, resourceName, ""),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, attr.Alias, ""),
				),
			},
		},
	})
}

func createResource29(terraformResourceName, networkName, resourceName, aliasName string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "%s" {
	  name = "%s"
	}
	resource "twingate_resource" "%s" {
	  name = "%s"
	  address = "acc-test.com"
	  remote_network_id = twingate_remote_network.%s.id
	  alias = "%s"
	}
	`, terraformResourceName, networkName, terraformResourceName, resourceName, terraformResourceName, aliasName)
}

func createResource29WithoutAlias(terraformResourceName, networkName, resourceName string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "%s" {
	  name = "%s"
	}
	resource "twingate_resource" "%s" {
	  name = "%s"
	  address = "acc-test.com"
	  remote_network_id = twingate_remote_network.%s.id
	  # alias = "some.value"
	}
	`, terraformResourceName, networkName, terraformResourceName, resourceName, terraformResourceName)
}

func TestAccTwingateResourceGroupsCursor(t *testing.T) {
	acctests.SetPageLimit(t, 1)

	const terraformResourceName = "test27"
	theResource := acctests.TerraformResource(terraformResourceName)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	groups, groupsID := genNewGroups("g27", 3)
	serviceAccounts, serviceAccountIDs := genNewServiceAccounts("s27", 3)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResourceWithGroupsAndServiceAccounts(terraformResourceName, remoteNetworkName, resourceName, groups, groupsID, serviceAccounts, serviceAccountIDs),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, accessGroupIdsLen, "3"),
					sdk.TestCheckResourceAttr(theResource, accessServiceAccountIdsLen, "3"),
				),
			},
			{
				Config: createResourceWithGroupsAndServiceAccounts(terraformResourceName, remoteNetworkName, resourceName, groups, groupsID[:2], serviceAccounts, serviceAccountIDs[:2]),
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
	  
	  protocols {
	    allow_icmp = true
	    tcp {
	      policy = "%s"
	      ports = ["80", "82-83"]
	    }
	    udp {
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
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config:      createResourceWithPort(remoteNetworkName, resourceName, "0"),
				ExpectError: regexp.MustCompile("port 0 not in the range of 1-65535"),
			},
			{
				Config:      createResourceWithPort(remoteNetworkName, resourceName, "65536"),
				ExpectError: regexp.MustCompile("port 65536 not in the range"),
			},
			{
				Config:      createResourceWithPort(remoteNetworkName, resourceName, "0-10"),
				ExpectError: regexp.MustCompile("port 0 not in the range"),
			},
			{
				Config:      createResourceWithPort(remoteNetworkName, resourceName, "65535-65536"),
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
	  protocols {
		allow_icmp = true
		tcp  {
			policy = "%s"
			ports = ["%s"]
		}
		udp {
			policy = "%s"
		}
	  }
	}
	`, networkName, resourceName, model.PolicyRestricted, port, model.PolicyAllowAll)
}

func TestAccTwingateResourceUpdateWithPort(t *testing.T) {
	theResource := acctests.TerraformResource("test30")
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResourceWithPort(remoteNetworkName, resourceName, "1"),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(theResource, firstTCPPort, "1"),
				),
			},
			{
				Config:      createResourceWithPort(remoteNetworkName, resourceName, "0"),
				ExpectError: regexp.MustCompile("port 0 not in the range of 1-65535"),
			},
		},
	})
}

func TestAccTwingateResourceWithPortsFailsForAllowAllAndDenyAllPolicy(t *testing.T) {
	const terraformResourceName = "test28"
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config:      createResourceWithPorts(terraformResourceName, remoteNetworkName, resourceName, model.PolicyAllowAll),
				ExpectError: regexp.MustCompile(resource.ErrPortsWithPolicyAllowAll.Error()),
			},
			{
				Config:      createResourceWithPorts(terraformResourceName, remoteNetworkName, resourceName, model.PolicyDenyAll),
				ExpectError: regexp.MustCompile(resource.ErrPortsWithPolicyDenyAll.Error()),
			},
		},
	})
}

func createResourceWithPorts(name, networkName, resourceName, policy string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "%[1]s" {
	  name = "%[2]s"
	}
	resource "twingate_resource" "%[1]s" {
	  name = "%[3]s"
	  address = "acc-test-%[1]s.com"
	  remote_network_id = twingate_remote_network.%[1]s.id
	  
	  protocols {
	    allow_icmp = true
	    tcp {
	      policy = "%[4]s"
	      ports = ["80", "82-83"]
	    }
	    udp {
	      policy = "%[5]s"
	    }
	  }
	}
	`, name, networkName, resourceName, policy, model.PolicyAllowAll)
}

func TestAccTwingateResourceWithoutPortsOkForAllowAllAndDenyAllPolicy(t *testing.T) {
	const terraformResourceName = "test29"
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResourceName)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResourceWithoutPorts(terraformResourceName, remoteNetworkName, resourceName, model.PolicyAllowAll),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyAllowAll),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "0"),
				),
			},
			{
				Config: createResourceWithoutPorts(terraformResourceName, remoteNetworkName, resourceName, model.PolicyDenyAll),
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
	  
	  protocols {
	    allow_icmp = true
	    tcp {
	      policy = "%[4]s"
	    }
	    udp {
	      policy = "%[5]s"
	    }
	  }
	}
	`, name, networkName, resourceName, policy, model.PolicyAllowAll)
}

func TestAccTwingateResourceWithRestrictedPolicy(t *testing.T) {
	const terraformResourceName = "test30"
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	theResource := acctests.TerraformResource(terraformResourceName)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResourceWithPorts(terraformResourceName, remoteNetworkName, resourceName, model.PolicyRestricted),
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
	const terraformResourceName = "test31"
	theResource := acctests.TerraformResource(terraformResourceName)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResourceWithoutPorts(terraformResourceName, remoteNetworkName, resourceName, model.PolicyDenyAll),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyDenyAll),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "0"),
				),
			},
			{
				Config: createResourceWithPorts(terraformResourceName, remoteNetworkName, resourceName, model.PolicyRestricted),
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
	const terraformResourceName = "test32"
	theResource := acctests.TerraformResource(terraformResourceName)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResourceWithoutPorts(terraformResourceName, remoteNetworkName, resourceName, model.PolicyDenyAll),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyDenyAll),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "0"),
				),
			},
			{
				Config: createResourceWithoutPorts(terraformResourceName, remoteNetworkName, resourceName, model.PolicyAllowAll),
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
	const terraformResourceName = "test33"
	theResource := acctests.TerraformResource(terraformResourceName)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResourceWithPorts(terraformResourceName, remoteNetworkName, resourceName, model.PolicyRestricted),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyRestricted),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "2"),
				),
			},
			{
				Config: createResourceWithoutPorts(terraformResourceName, remoteNetworkName, resourceName, model.PolicyDenyAll),
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
	const terraformResourceName = "test34"
	theResource := acctests.TerraformResource(terraformResourceName)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResourceWithPorts(terraformResourceName, remoteNetworkName, resourceName, model.PolicyRestricted),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyRestricted),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "2"),
				),
			},
			{
				Config: createResourceWithoutPorts(terraformResourceName, remoteNetworkName, resourceName, model.PolicyAllowAll),
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
	const terraformResourceName = "test35"
	theResource := acctests.TerraformResource(terraformResourceName)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResourceWithPorts(terraformResourceName, remoteNetworkName, resourceName, model.PolicyRestricted),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyRestricted),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "2"),
				),
			},
			{
				Config:      createResourceWithPorts(terraformResourceName, remoteNetworkName, resourceName, model.PolicyAllowAll),
				ExpectError: regexp.MustCompile(resource.ErrPortsWithPolicyAllowAll.Error()),
			},
		},
	})
}

func TestAccTwingateResourcePolicyTransitionAllowAllToRestricted(t *testing.T) {
	const terraformResourceName = "test36"
	theResource := acctests.TerraformResource(terraformResourceName)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResourceWithoutPorts(terraformResourceName, remoteNetworkName, resourceName, model.PolicyAllowAll),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyAllowAll),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "0"),
				),
			},
			{
				Config: createResourceWithPorts(terraformResourceName, remoteNetworkName, resourceName, model.PolicyRestricted),
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
	const terraformResourceName = "test37"
	theResource := acctests.TerraformResource(terraformResourceName)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResourceWithoutPorts(terraformResourceName, remoteNetworkName, resourceName, model.PolicyAllowAll),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyAllowAll),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "0"),
				),
			},
			{
				Config: createResourceWithoutPorts(terraformResourceName, remoteNetworkName, resourceName, model.PolicyDenyAll),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
					sdk.TestCheckResourceAttr(theResource, tcpPolicy, model.PolicyDenyAll),
					sdk.TestCheckResourceAttr(theResource, tcpPortsLen, "0"),
				),
			},
		},
	})
}

func TestAccTwingateResourceWithBrowserOption(t *testing.T) {
	const terraformResourceName = "test40"
	theResource := acctests.TerraformResource(terraformResourceName)
	remoteNetworkName := test.RandomName()
	resourceName := test.RandomResourceName()
	wildcardAddress := "*.acc-test.com"

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateResourceDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createResourceWithoutBrowserOption(terraformResourceName, remoteNetworkName, resourceName, wildcardAddress),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
			{
				Config: createResourceWithBrowserOption(terraformResourceName, remoteNetworkName, resourceName, wildcardAddress, false),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
			{
				Config:      createResourceWithBrowserOption(terraformResourceName, remoteNetworkName, resourceName, wildcardAddress, true),
				ExpectError: regexp.MustCompile(resource.ErrWildcardAddressWithEnabledShortcut.Error()),
			},
		},
	})
}

func TestAccTwingateResourceWithBrowserOptionFailOnUpdate(t *testing.T) {
	const terraformResourceName = "test41"
	theResource := acctests.TerraformResource(terraformResourceName)
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
				Config: createResourceWithoutBrowserOption(terraformResourceName, remoteNetworkName, resourceName, simpleAddress),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
			{
				Config: createResourceWithBrowserOption(terraformResourceName, remoteNetworkName, resourceName, simpleAddress, true),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
			{
				Config:      createResourceWithBrowserOption(terraformResourceName, remoteNetworkName, resourceName, wildcardAddress, true),
				ExpectError: regexp.MustCompile(resource.ErrWildcardAddressWithEnabledShortcut.Error()),
			},
		},
	})
}

func TestAccTwingateResourceWithBrowserOptionRecovered(t *testing.T) {
	const terraformResourceName = "test42"
	theResource := acctests.TerraformResource(terraformResourceName)
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
				Config: createResourceWithBrowserOption(terraformResourceName, remoteNetworkName, resourceName, simpleAddress, true),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
			{
				Config: createResourceWithoutBrowserOption(terraformResourceName, remoteNetworkName, resourceName, wildcardAddress),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(theResource),
				),
			},
		},
	})
}

func createResourceWithoutBrowserOption(name, networkName, resourceName, address string) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "%[1]s" {
	  name = "%[2]s"
	}
	resource "twingate_resource" "%[1]s" {
	  name = "%[3]s"
	  address = "%[4]s"
	  remote_network_id = twingate_remote_network.%[1]s.id
	}
	`, name, networkName, resourceName, address)
}

func createResourceWithBrowserOption(name, networkName, resourceName, address string, browserOption bool) string {
	return fmt.Sprintf(`
	resource "twingate_remote_network" "%[1]s" {
	  name = "%[2]s"
	}
	resource "twingate_resource" "%[1]s" {
	  name = "%[3]s"
	  address = "%[4]s"
	  remote_network_id = twingate_remote_network.%[1]s.id
	  is_browser_shortcut_enabled = %[5]v
	}
	`, name, networkName, resourceName, address, browserOption)
}
