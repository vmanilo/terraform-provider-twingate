package datasource

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/utils"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	serviceAccountsLen = attr.Len(attr.ServiceAccounts)
	keyIDsLen          = attr.Len(attr.ServiceAccounts, attr.KeyIDs)
	serviceAccountName = attr.Path(attr.ServiceAccounts, attr.Name)
)

func TestAccDatasourceTwingateServicesFilterByName(t *testing.T) {
	t.Parallel()

	name := test.Prefix("orange")
	const (
		terraformResourceName = "dts_service"
		theDatasource         = "data.twingate_service_accounts.out"
	)

	config := []terraformServiceConfig{
		{
			serviceName:           name,
			terraformResourceName: test.TerraformRandName(terraformResourceName),
		},
		{
			serviceName:           test.Prefix("lemon"),
			terraformResourceName: test.TerraformRandName(terraformResourceName),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: terraformConfig(
					createServices(config),
					datasourceServices(name, config),
				),
				Check: acctests.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(theDatasource, serviceAccountsLen, "1"),
					resource.TestCheckResourceAttr(theDatasource, keyIDsLen, "1"),
					resource.TestCheckResourceAttr(theDatasource, attr.ID, "service-by-name-"+name),
				),
			},
		},
	})
}

func TestAccDatasourceTwingateServicesAll(t *testing.T) {
	t.Parallel()

	prefix := test.Prefix() + acctest.RandString(4)
	const (
		terraformResourceName = "dts_service"
		theDatasource         = "data.twingate_service_accounts.out"
	)

	config := []terraformServiceConfig{
		{
			serviceName:           prefix + "_orange",
			terraformResourceName: test.TerraformRandName(terraformResourceName),
		},
		{
			serviceName:           prefix + "_lemon",
			terraformResourceName: test.TerraformRandName(terraformResourceName),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: filterDatasourceServices(prefix, config),
				Check: acctests.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(theDatasource, attr.ID, "all-services"),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: filterDatasourceServices(prefix, config),
				Check: acctests.ComposeTestCheckFunc(
					testCheckOutputLength("my_services", 2),
				),
			},
		},
	})
}

func TestAccDatasourceTwingateServicesEmptyResult(t *testing.T) {
	t.Parallel()

	const theDatasource = "data.twingate_service_accounts.out"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: datasourceServices(test.RandomName(), nil),
				Check: acctests.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(theDatasource, serviceAccountsLen, "0"),
				),
			},
		},
	})
}

type terraformServiceConfig struct {
	terraformResourceName, serviceName string
}

func terraformConfig(resources ...string) string {
	return strings.Join(resources, "\n")
}

func datasourceServices(name string, configs []terraformServiceConfig) string {
	var dependsOn string
	ids := getTerraformServiceKeys(configs)

	if ids != "" {
		dependsOn = fmt.Sprintf("depends_on = [%s]", ids)
	}

	return fmt.Sprintf(`
	data "twingate_service_accounts" "out" {
	  name = "%s"

	  %s
	}
	`, name, dependsOn)
}

func createServices(configs []terraformServiceConfig) string {
	return strings.Join(
		utils.Map[terraformServiceConfig, string](configs, func(cfg terraformServiceConfig) string {
			return createServiceKey(cfg.terraformResourceName, cfg.serviceName)
		}),
		"\n",
	)
}

func getTerraformServiceKeys(configs []terraformServiceConfig) string {
	return strings.Join(
		utils.Map[terraformServiceConfig, string](configs, func(cfg terraformServiceConfig) string {
			return acctests.TerraformServiceKey(cfg.terraformResourceName)
		}),
		", ",
	)
}

func createServiceKey(resourceName, serviceName string) string {
	return acctests.Nprintf(`
	${service_account}

	resource "twingate_service_account_key" "${service_account_key_resource}" {
	  service_account_id = twingate_service_account.${service_account_resource}.id
	}
	`,
		map[string]any{
			"service_account":              createServiceAccount(resourceName, serviceName),
			"service_account_key_resource": resourceName,
			"service_account_resource":     resourceName,
		})
}

func createServiceAccount(resourceName, serviceName string) string {
	return acctests.Nprintf(`
	resource "twingate_service_account" "${resource_name}" {
	  name = "${name}"
	}
	`,
		map[string]any{
			"resource_name": resourceName,
			"name":          serviceName,
		})
}

func filterDatasourceServices(prefix string, configs []terraformServiceConfig) string {
	return acctests.Nprintf(`
	${services}

	data "twingate_service_accounts" "out" {

	}

	output "my_services" {
	  	value = [for c in data.twingate_service_accounts.out.service_accounts : c if length(regexall("^${prefix}", c.name)) > 0]
	}
	`,
		map[string]any{
			"services": createServices(configs),
			"prefix":   prefix,
		})
}

func TestAccDatasourceTwingateServicesAllCursors(t *testing.T) {
	acctests.SetPageLimit(t, 1)
	prefix := test.Prefix() + acctest.RandString(4)
	const theDatasource = "data.twingate_service_accounts.out"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: datasourceServicesConfig(prefix),
				Check: acctests.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(theDatasource, attr.ID, "all-services"),
				),
			},
			{
				Config: datasourceServicesConfig(prefix),
				Check: acctests.ComposeTestCheckFunc(
					testCheckOutputLength("my_services", 3),
					testCheckOutputNestedLen("my_services", 0, attr.ResourceIDs, 1),
					testCheckOutputNestedLen("my_services", 0, attr.KeyIDs, 2),
				),
			},
		},
	})
}

func datasourceServicesConfig(prefix string) string {
	return acctests.Nprintf(`
    resource "twingate_service_account" "${prefix}_1" {
      name = "${prefix}-1"
    }
    
    resource "twingate_service_account" "${prefix}_2" {
      name = "${prefix}-2"
    }

    resource "twingate_service_account" "${prefix}_3" {
      name = "${prefix}-3"
    }
    
    resource "twingate_remote_network" "${prefix}_1" {
      name = "${prefix}-1"
    }
    
    resource "twingate_remote_network" "${prefix}_2" {
      name = "${prefix}-2"
    }
    
    resource "twingate_resource" "${prefix}_1" {
      name = "${prefix}-1"
      address = "acc-test.com"
      remote_network_id = twingate_remote_network.${prefix}_1.id
    
      access {
        service_account_ids = [twingate_service_account.${prefix}_1.id, twingate_service_account.${prefix}_2.id]
      }
    }
    
    resource "twingate_resource" "${prefix}_2" {
      name = "${prefix}-2"
      address = "acc-test.com"
      remote_network_id = twingate_remote_network.${prefix}_2.id
    
      access {
        service_account_ids = [twingate_service_account.${prefix}_3.id]
      }
    }
    
    resource "twingate_service_account_key" "${prefix}_1_1" {
      service_account_id = twingate_service_account.${prefix}_1.id
    }
    
    resource "twingate_service_account_key" "${prefix}_1_2" {
      service_account_id = twingate_service_account.${prefix}_1.id
    }
    
    resource "twingate_service_account_key" "${prefix}_2_1" {
      service_account_id = twingate_service_account.${prefix}_2.id
    }
    
    resource "twingate_service_account_key" "${prefix}_2_2" {
      service_account_id = twingate_service_account.${prefix}_2.id
    }
    
    resource "twingate_service_account_key" "${prefix}_3_1" {
      service_account_id = twingate_service_account.${prefix}_3.id
    }

    resource "twingate_service_account_key" "${prefix}_3_2" {
      service_account_id = twingate_service_account.${prefix}_3.id
    }
    
    data "twingate_service_accounts" "out" {
    	depends_on = [twingate_resource.${prefix}_1, twingate_resource.${prefix}_2]
    }
    
    output "my_services" {
      value = [for c in data.twingate_service_accounts.out.service_accounts : c if length(regexall("^${prefix}", c.name)) > 0]
      depends_on = [data.twingate_service_accounts.out]
    }
`, map[string]any{"prefix": prefix})
}

func TestAccDatasourceTwingateServicesWithMultipleFilters(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testDatasourceServicesWithMultipleFilters(test.RandomName()),
				ExpectError: regexp.MustCompile("Only one of name.*"),
			},
		},
	})
}

func testDatasourceServicesWithMultipleFilters(name string) string {
	return fmt.Sprintf(`
	data "twingate_service_accounts" "with-multiple-filters" {
	  name_regexp = "%[1]s"
	  name_contains = "%[1]s"
	}
	`, name)
}

func TestAccDatasourceTwingateServicesFilterByPrefix(t *testing.T) {
	t.Parallel()

	const (
		terraformResourceName = "dts_service"
		theDatasource         = "data.twingate_service_accounts.out"
	)

	prefix := test.Prefix("orange")
	name := acctest.RandomWithPrefix(prefix)
	config := []terraformServiceConfig{
		{
			serviceName:           name,
			terraformResourceName: test.TerraformRandName(terraformResourceName),
		},
		{
			serviceName:           test.Prefix("lemon"),
			terraformResourceName: test.TerraformRandName(terraformResourceName),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: terraformConfig(
					createServices(config),
					datasourceServicesWithFilter(config, prefix, attr.FilterByPrefix),
				),
				Check: acctests.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(theDatasource, serviceAccountsLen, "1"),
					resource.TestCheckResourceAttr(theDatasource, serviceAccountName, name),
				),
			},
		},
	})
}

func datasourceServicesWithFilter(configs []terraformServiceConfig, name, filter string) string {
	var dependsOn string
	ids := getTerraformServiceKeys(configs)

	if ids != "" {
		dependsOn = fmt.Sprintf("depends_on = [%s]", ids)
	}

	return fmt.Sprintf(`
	data "twingate_service_accounts" "out" {
	  name%s = "%s"

	  %s
	}
	`, filter, name, dependsOn)
}

func TestAccDatasourceTwingateServicesFilterBySuffix(t *testing.T) {
	t.Parallel()

	const (
		terraformResourceName = "dts_service"
		theDatasource         = "data.twingate_service_accounts.out"
	)

	name := test.Prefix("orange")
	config := []terraformServiceConfig{
		{
			serviceName:           name,
			terraformResourceName: test.TerraformRandName(terraformResourceName),
		},
		{
			serviceName:           test.Prefix("lemon"),
			terraformResourceName: test.TerraformRandName(terraformResourceName),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: terraformConfig(
					createServices(config),
					datasourceServicesWithFilter(config, "orange", attr.FilterBySuffix),
				),
				Check: acctests.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(theDatasource, serviceAccountsLen, "1"),
					resource.TestCheckResourceAttr(theDatasource, serviceAccountName, name),
				),
			},
		},
	})
}

func TestAccDatasourceTwingateServicesFilterByContains(t *testing.T) {
	t.Parallel()

	const (
		terraformResourceName = "dts_service"
		theDatasource         = "data.twingate_service_accounts.out"
	)

	name := test.Prefix("orange")
	config := []terraformServiceConfig{
		{
			serviceName:           name,
			terraformResourceName: test.TerraformRandName(terraformResourceName),
		},
		{
			serviceName:           test.Prefix("lemon"),
			terraformResourceName: test.TerraformRandName(terraformResourceName),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: terraformConfig(
					createServices(config),
					datasourceServicesWithFilter(config, "rang", attr.FilterByContains),
				),
				Check: acctests.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(theDatasource, serviceAccountsLen, "1"),
					resource.TestCheckResourceAttr(theDatasource, serviceAccountName, name),
				),
			},
		},
	})
}

func TestAccDatasourceTwingateServicesFilterByRegexp(t *testing.T) {
	t.Parallel()

	const (
		terraformResourceName = "dts_service"
		theDatasource         = "data.twingate_service_accounts.out"
	)

	name := test.Prefix("orange")
	config := []terraformServiceConfig{
		{
			serviceName:           name,
			terraformResourceName: test.TerraformRandName(terraformResourceName),
		},
		{
			serviceName:           test.Prefix("lemon"),
			terraformResourceName: test.TerraformRandName(terraformResourceName),
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: terraformConfig(
					createServices(config),
					datasourceServicesWithFilter(config, ".*ora.*", attr.FilterByRegexp),
				),
				Check: acctests.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(theDatasource, serviceAccountsLen, "1"),
					resource.TestCheckResourceAttr(theDatasource, serviceAccountName, name),
				),
			},
		},
	})
}
