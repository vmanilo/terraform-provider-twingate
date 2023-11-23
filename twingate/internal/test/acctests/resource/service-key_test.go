package resource

import (
	"errors"
	"regexp"
	"testing"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/model"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/provider/resource"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
	sdk "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var ErrEmptyValue = errors.New("empty value")

func createServiceKey(resourceName, serviceAccountName string) string {
	return acctests.Nprintf(`
	${service_account}

	resource "twingate_service_account_key" "${service_account_key_resource}" {
	  service_account_id = twingate_service_account.${service_account_resource}.id
	}
	`,
		map[string]any{
			"service_account":              createServiceAccount(resourceName, serviceAccountName),
			"service_account_key_resource": resourceName,
			"service_account_resource":     resourceName,
		})
}

func createServiceKeyWithName(resourceName, serviceAccountName, serviceKeyName string) string {
	return acctests.Nprintf(`
	${service_account}

	resource "twingate_service_account_key" "${service_account_key_resource}" {
	  service_account_id = twingate_service_account.${service_account_resource}.id
	  name = "${name}"
	}
	`,
		map[string]any{
			"service_account":              createServiceAccount(resourceName, serviceAccountName),
			"service_account_key_resource": resourceName,
			"service_account_resource":     resourceName,
			"name":                         serviceKeyName,
		})
}

func createServiceKeyWithExpiration(resourceName, serviceAccountName string, expirationTime int) string {
	return acctests.Nprintf(`
	${service_account}

	resource "twingate_service_account_key" "${service_account_key_resource}" {
	  service_account_id = twingate_service_account.${service_account_resource}.id
	  expiration_time = ${expiration_time}
	}
	`,
		map[string]any{
			"service_account":              createServiceAccount(resourceName, serviceAccountName),
			"service_account_key_resource": resourceName,
			"service_account_resource":     resourceName,
			"expiration_time":              expirationTime,
		})
}

func nonEmptyValue(value string) error {
	if value != "" {
		return nil
	}

	return ErrEmptyValue
}

func TestAccTwingateServiceKeyCreateUpdate(t *testing.T) {
	t.Parallel()

	serviceAccountName := test.RandomName()
	terraformResourceName := test.TerraformRandName("test_key")
	serviceAccount := acctests.TerraformServiceAccount(terraformResourceName)
	serviceKey := acctests.TerraformServiceKey(terraformResourceName)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createServiceKey(terraformResourceName, serviceAccountName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceAccount),
					sdk.TestCheckResourceAttr(serviceAccount, attr.Name, serviceAccountName),
					acctests.CheckTwingateResourceExists(serviceKey),
					sdk.TestCheckResourceAttrWith(serviceKey, attr.Token, nonEmptyValue),
				),
			},
			{
				Config: createServiceKey(terraformResourceName, serviceAccountName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceAccount),
					sdk.TestCheckResourceAttr(serviceAccount, attr.Name, serviceAccountName),
					acctests.CheckTwingateResourceExists(serviceKey),
					sdk.TestCheckResourceAttrWith(serviceKey, attr.Token, nonEmptyValue),
				),
			},
		},
	})
}

func TestAccTwingateServiceKeyCreateUpdateWithName(t *testing.T) {
	t.Parallel()

	serviceAccountName := test.RandomName()
	terraformResourceName := test.TerraformRandName("test_key")
	serviceAccount := acctests.TerraformServiceAccount(terraformResourceName)
	serviceKey := acctests.TerraformServiceKey(terraformResourceName)
	name1 := test.RandomName()
	name2 := test.RandomName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createServiceKeyWithName(terraformResourceName, serviceAccountName, name1),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceAccount),
					sdk.TestCheckResourceAttr(serviceAccount, attr.Name, serviceAccountName),
					acctests.CheckTwingateResourceExists(serviceKey),
					sdk.TestCheckResourceAttr(serviceKey, attr.Name, name1),
					sdk.TestCheckResourceAttrWith(serviceKey, attr.Token, nonEmptyValue),
				),
			},
			{
				Config: createServiceKeyWithName(terraformResourceName, serviceAccountName, name2),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceAccount),
					sdk.TestCheckResourceAttr(serviceAccount, attr.Name, serviceAccountName),
					acctests.CheckTwingateResourceExists(serviceKey),
					sdk.TestCheckResourceAttr(serviceKey, attr.Name, name2),
					sdk.TestCheckResourceAttrWith(serviceKey, attr.Token, nonEmptyValue),
					acctests.WaitTestFunc(),
				),
			},
		},
	})
}

func TestAccTwingateServiceKeyWontReCreateAfterInactive(t *testing.T) {
	t.Parallel()

	serviceAccountName := test.RandomName()
	terraformResourceName := test.TerraformRandName("test_key")
	serviceKey := acctests.TerraformServiceKey(terraformResourceName)

	resourceID := new(string)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createServiceKey(terraformResourceName, serviceAccountName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceKey),
					acctests.GetTwingateResourceID(serviceKey, &resourceID),
					sdk.TestCheckResourceAttrWith(serviceKey, attr.Token, nonEmptyValue),
					acctests.RevokeTwingateServiceKey(serviceKey),
					acctests.WaitTestFunc(),
					acctests.CheckTwingateServiceKeyStatus(serviceKey, model.StatusRevoked),
				),
			},
			{
				Config: createServiceKey(terraformResourceName, serviceAccountName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceKey),
					sdk.TestCheckResourceAttr(serviceKey, attr.IsActive, "false"),
					sdk.TestCheckResourceAttrWith(serviceKey, attr.Token, nonEmptyValue),
					sdk.TestCheckResourceAttrWith(serviceKey, attr.ID, func(value string) error {
						if *resourceID == "" {
							return errors.New("failed to fetch resource id")
						}

						if value != *resourceID {
							return errors.New("resource was re-created")
						}

						return nil
					}),
				),
			},
		},
	})
}

func TestAccTwingateServiceKeyDelete(t *testing.T) {
	t.Parallel()

	serviceAccountName := test.RandomName()
	terraformResourceName := test.TerraformRandName("test_key")
	serviceKey := acctests.TerraformServiceKey(terraformResourceName)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config:  createServiceKey(terraformResourceName, serviceAccountName),
				Destroy: true,
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceDoesNotExists(serviceKey),
				),
			},
		},
	})
}

func TestAccTwingateServiceKeyReCreateAfterDeletion(t *testing.T) {
	t.Parallel()

	serviceAccountName := test.RandomName()
	terraformResourceName := test.TerraformRandName("test_key")
	serviceKey := acctests.TerraformServiceKey(terraformResourceName)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createServiceKey(terraformResourceName, serviceAccountName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceKey),
					acctests.RevokeTwingateServiceKey(serviceKey),
					acctests.DeleteTwingateResource(serviceKey, resource.TwingateServiceAccountKey),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: createServiceKey(terraformResourceName, serviceAccountName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceKey),
					sdk.TestCheckResourceAttrWith(serviceKey, attr.Token, nonEmptyValue),
				),
			},
		},
	})
}

func TestAccTwingateServiceKeyCreateWithInvalidExpiration(t *testing.T) {
	t.Parallel()

	serviceAccountName := test.RandomName()
	terraformResourceName := test.TerraformRandName("test_key")

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config:      createServiceKeyWithExpiration(terraformResourceName, serviceAccountName, -1),
				ExpectError: regexp.MustCompile(resource.ErrInvalidExpirationTime.Error()),
			},
			{
				Config:      createServiceKeyWithExpiration(terraformResourceName, serviceAccountName, 366),
				ExpectError: regexp.MustCompile(resource.ErrInvalidExpirationTime.Error()),
			},
		},
	})
}

func TestAccTwingateServiceKeyCreateWithExpiration(t *testing.T) {
	t.Parallel()

	serviceAccountName := test.RandomName()
	terraformResourceName := test.TerraformRandName("test_key")
	serviceAccount := acctests.TerraformServiceAccount(terraformResourceName)
	serviceKey := acctests.TerraformServiceKey(terraformResourceName)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createServiceKeyWithExpiration(terraformResourceName, serviceAccountName, 365),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceAccount),
					sdk.TestCheckResourceAttr(serviceAccount, attr.Name, serviceAccountName),
					acctests.CheckTwingateResourceExists(serviceKey),
					sdk.TestCheckResourceAttr(serviceKey, attr.IsActive, "true"),
					sdk.TestCheckResourceAttrWith(serviceKey, attr.Token, nonEmptyValue),
				),
			},
		},
	})
}

func TestAccTwingateServiceKeyReCreateAfterChangingExpirationTime(t *testing.T) {
	t.Parallel()

	serviceAccountName := test.RandomName()
	terraformResourceName := test.TerraformRandName("test_key")
	serviceKey := acctests.TerraformServiceKey(terraformResourceName)

	resourceID := new(string)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createServiceKeyWithExpiration(terraformResourceName, serviceAccountName, 1),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceKey),
					acctests.GetTwingateResourceID(serviceKey, &resourceID),
					sdk.TestCheckResourceAttrWith(serviceKey, attr.Token, nonEmptyValue),
				),
			},
			{
				Config: createServiceKeyWithExpiration(terraformResourceName, serviceAccountName, 2),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceKey),
					sdk.TestCheckResourceAttrWith(serviceKey, attr.ID, func(value string) error {
						if *resourceID == "" {
							return errors.New("failed to fetch resource id")
						}

						if value == *resourceID {
							return errors.New("resource was not re-created")
						}

						return nil
					}),
				),
			},
		},
	})
}

func TestAccTwingateServiceKeyAndServiceAccountLifecycle(t *testing.T) {
	t.Parallel()

	serviceAccountName := test.RandomName()
	serviceAccountNameV2 := test.RandomName()
	terraformServiceAccountName := test.TerraformRandName("test_acc")
	terraformServiceAccountNameV2 := test.TerraformRandName("test_acc_v2")
	terraformServiceAccountKeyName := test.TerraformRandName("test_key")
	serviceAccount := acctests.TerraformServiceAccount(terraformServiceAccountName)
	serviceAccountV2 := acctests.TerraformServiceAccount(terraformServiceAccountNameV2)
	serviceKey := acctests.TerraformServiceKey(terraformServiceAccountKeyName)

	serviceKeyResourceID := new(string)
	serviceAccountResourceID := new(string)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config: createServiceKeyV1(terraformServiceAccountName, serviceAccountName, terraformServiceAccountNameV2, serviceAccountNameV2, terraformServiceAccountKeyName, terraformServiceAccountName),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceAccount),
					sdk.TestCheckResourceAttr(serviceAccount, attr.Name, serviceAccountName),
					acctests.CheckTwingateResourceExists(serviceKey),
					sdk.TestCheckResourceAttrWith(serviceKey, attr.Token, nonEmptyValue),
					acctests.GetTwingateResourceID(serviceKey, &serviceKeyResourceID),
					acctests.GetTwingateResourceID(serviceKey, &serviceAccountResourceID),
				),
			},
			{
				Config: createServiceKeyV1(terraformServiceAccountName, serviceAccountName, terraformServiceAccountNameV2, serviceAccountNameV2, terraformServiceAccountKeyName, terraformServiceAccountNameV2),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceAccountV2),
					sdk.TestCheckResourceAttr(serviceAccountV2, attr.Name, serviceAccountNameV2),
					acctests.CheckTwingateResourceExists(serviceKey),
					sdk.TestCheckResourceAttrWith(serviceKey, attr.Token, nonEmptyValue),

					// test resources were re-created
					sdk.TestCheckResourceAttrWith(serviceKey, attr.ID, func(value string) error {
						if *serviceKeyResourceID == "" {
							return errors.New("failed to fetch service_key resource id")
						}

						if value == *serviceKeyResourceID {
							return errors.New("service_key resource was not re-created")
						}

						return nil
					}),

					sdk.TestCheckResourceAttrWith(serviceAccountV2, attr.ID, func(value string) error {
						if *serviceAccountResourceID == "" {
							return errors.New("failed to fetch service_account resource id")
						}

						if value == *serviceAccountResourceID {
							return errors.New("service_account resource was not re-created")
						}

						return nil
					}),
				),
			},
		},
	})
}

func createServiceKeyV1(terraformServiceAccountName, serviceAccountName, terraformServiceAccountNameV2, serviceAccountNameV2, terraformServiceAccountKeyName, serviceAccount string) string {
	return acctests.Nprintf(`
	resource "twingate_service_account" "${service_account_resource_1}" {
	  name = "${name_1}"
	}

	resource "twingate_service_account" "${service_account_resource_2}" {
	  name = "${name_2}"
	}

	resource "twingate_service_account_key" "${service_account_key_resource}" {
	  service_account_id = twingate_service_account.${service_account_resource}.id
	}
	`,
		map[string]any{
			"service_account_resource_1":   terraformServiceAccountName,
			"name_1":                       serviceAccountName,
			"service_account_resource_2":   terraformServiceAccountNameV2,
			"name_2":                       serviceAccountNameV2,
			"service_account_key_resource": terraformServiceAccountKeyName,
			"service_account_resource":     serviceAccount,
		})
}
