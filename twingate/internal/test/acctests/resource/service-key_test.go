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
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

var ErrEmptyValue = errors.New("empty value")

func nonEmptyValue(value string) error {
	if value != "" {
		return nil
	}

	return ErrEmptyValue
}

func TestAccTwingateServiceKeyCreateUpdate(t *testing.T) {
	t.Parallel()

	serviceAccount := NewServiceAccount()
	serviceKey := NewServiceAccountKey(serviceAccount.TerraformResourceID())

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(serviceAccount, serviceKey),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(serviceAccount.TerraformResource(), attr.Name, serviceAccount.Name),
					sdk.TestCheckResourceAttrWith(serviceKey.TerraformResource(), attr.Token, nonEmptyValue),
				),
			},
			{
				Config: configBuilder(serviceAccount, serviceKey),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(serviceAccount.TerraformResource(), attr.Name, serviceAccount.Name),
					sdk.TestCheckResourceAttrWith(serviceKey.TerraformResource(), attr.Token, nonEmptyValue),
				),
			},
		},
	})
}

func TestAccTwingateServiceKeyCreateUpdateWithName(t *testing.T) {
	t.Parallel()

	serviceAccount := NewServiceAccount()
	serviceKey := NewServiceAccountKey(serviceAccount.TerraformResourceID())

	name1 := test.RandomName()
	name2 := test.RandomName()

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(serviceKey.Set(attr.Name, name1), serviceAccount),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(serviceAccount.TerraformResource(), attr.Name, serviceAccount.Name),
					sdk.TestCheckResourceAttr(serviceKey.TerraformResource(), attr.Name, name1),
					sdk.TestCheckResourceAttrWith(serviceKey.TerraformResource(), attr.Token, nonEmptyValue),
				),
			},
			{
				Config: configBuilder(serviceKey.Set(attr.Name, name2), serviceAccount),
				Check: acctests.ComposeTestCheckFunc(
					sdk.TestCheckResourceAttr(serviceAccount.TerraformResource(), attr.Name, serviceAccount.Name),
					sdk.TestCheckResourceAttr(serviceKey.TerraformResource(), attr.Name, name2),
					sdk.TestCheckResourceAttrWith(serviceKey.TerraformResource(), attr.Token, nonEmptyValue),
					acctests.WaitTestFunc(),
				),
			},
		},
	})
}

func TestAccTwingateServiceKeyWontReCreateAfterInactive(t *testing.T) {
	t.Parallel()

	serviceAccount := NewServiceAccount()
	serviceKey := NewServiceAccountKey(serviceAccount.TerraformResourceID())

	resourceID := new(string)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(serviceKey, serviceAccount),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceKey.TerraformResource()),
					acctests.GetTwingateResourceID(serviceKey.TerraformResource(), &resourceID),
					sdk.TestCheckResourceAttrWith(serviceKey.TerraformResource(), attr.Token, nonEmptyValue),
					acctests.RevokeTwingateServiceKey(serviceKey.TerraformResource()),
					acctests.WaitTestFunc(),
					acctests.CheckTwingateServiceKeyStatus(serviceKey.TerraformResource(), model.StatusRevoked),
				),
			},
			{
				Config: configBuilder(serviceKey, serviceAccount),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceKey.TerraformResource()),
					sdk.TestCheckResourceAttr(serviceKey.TerraformResource(), attr.IsActive, "false"),
					sdk.TestCheckResourceAttrWith(serviceKey.TerraformResource(), attr.Token, nonEmptyValue),
					sdk.TestCheckResourceAttrWith(serviceKey.TerraformResource(), attr.ID, func(value string) error {
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

	serviceAccount := NewServiceAccount()
	serviceKey := NewServiceAccountKey(serviceAccount.TerraformResourceID())

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config:  configBuilder(serviceKey, serviceAccount),
				Destroy: true,
			},
			{
				Config: configBuilder(serviceKey, serviceAccount),
				ConfigPlanChecks: sdk.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(serviceKey.TerraformResource(), plancheck.ResourceActionCreate),
					},
				},
			},
		},
	})
}

func TestAccTwingateServiceKeyReCreateAfterDeletion(t *testing.T) {
	t.Parallel()

	serviceAccount := NewServiceAccount()
	serviceKey := NewServiceAccountKey(serviceAccount.TerraformResourceID())

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(serviceKey, serviceAccount),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceKey.TerraformResource()),
					acctests.RevokeTwingateServiceKey(serviceKey.TerraformResource()),
					acctests.DeleteTwingateResource(serviceKey.TerraformResource(), resource.TwingateServiceAccountKey),
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: configBuilder(serviceKey, serviceAccount),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceKey.TerraformResource()),
					sdk.TestCheckResourceAttrWith(serviceKey.TerraformResource(), attr.Token, nonEmptyValue),
				),
			},
		},
	})
}

func TestAccTwingateServiceKeyCreateWithInvalidExpiration(t *testing.T) {
	t.Parallel()

	serviceAccount := NewServiceAccount()
	serviceKey := NewServiceAccountKey(serviceAccount.TerraformResourceID())

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config:      configBuilder(serviceAccount, serviceKey.Set(attr.ExpirationTime, -1)),
				ExpectError: regexp.MustCompile(resource.ErrInvalidExpirationTime.Error()),
			},
			{
				Config:      configBuilder(serviceAccount, serviceKey.Set(attr.ExpirationTime, 366)),
				ExpectError: regexp.MustCompile(resource.ErrInvalidExpirationTime.Error()),
			},
		},
	})
}

func TestAccTwingateServiceKeyCreateWithExpiration(t *testing.T) {
	t.Parallel()

	serviceAccount := NewServiceAccount()
	serviceKey := NewServiceAccountKey(serviceAccount.TerraformResourceID())

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(serviceAccount, serviceKey.Set(attr.ExpirationTime, 365)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceAccount.TerraformResource()),
					sdk.TestCheckResourceAttr(serviceAccount.TerraformResource(), attr.Name, serviceAccount.Name),
					acctests.CheckTwingateResourceExists(serviceKey.TerraformResource()),
					sdk.TestCheckResourceAttr(serviceKey.TerraformResource(), attr.IsActive, "true"),
					sdk.TestCheckResourceAttrWith(serviceKey.TerraformResource(), attr.Token, nonEmptyValue),
				),
			},
		},
	})
}

func TestAccTwingateServiceKeyReCreateAfterChangingExpirationTime(t *testing.T) {
	t.Parallel()

	serviceAccount := NewServiceAccount()
	serviceKey := NewServiceAccountKey(serviceAccount.TerraformResourceID())

	resourceID := new(string)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(serviceAccount, serviceKey.Set(attr.ExpirationTime, 1)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceKey.TerraformResource()),
					acctests.GetTwingateResourceID(serviceKey.TerraformResource(), &resourceID),
					sdk.TestCheckResourceAttrWith(serviceKey.TerraformResource(), attr.Token, nonEmptyValue),
				),
			},
			{
				Config: configBuilder(serviceAccount, serviceKey.Set(attr.ExpirationTime, 2)),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceKey.TerraformResource()),
					sdk.TestCheckResourceAttrWith(serviceKey.TerraformResource(), attr.ID, func(value string) error {
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

	serviceAccount1 := NewServiceAccount()
	serviceAccount2 := NewServiceAccount()
	serviceKey := NewServiceAccountKey(serviceAccount1.TerraformResourceID())

	serviceKeyResourceID := new(string)
	serviceAccountResourceID := new(string)

	sdk.Test(t, sdk.TestCase{
		ProtoV6ProviderFactories: acctests.ProviderFactories,
		PreCheck:                 func() { acctests.PreCheck(t) },
		CheckDestroy:             acctests.CheckTwingateServiceAccountDestroy,
		Steps: []sdk.TestStep{
			{
				Config: configBuilder(serviceKey, serviceAccount1, serviceAccount2),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceAccount1.TerraformResource()),
					sdk.TestCheckResourceAttr(serviceAccount1.TerraformResource(), attr.Name, serviceAccount1.Name),
					acctests.CheckTwingateResourceExists(serviceKey.TerraformResource()),
					sdk.TestCheckResourceAttrWith(serviceKey.TerraformResource(), attr.Token, nonEmptyValue),
					acctests.GetTwingateResourceID(serviceKey.TerraformResource(), &serviceKeyResourceID),
					acctests.GetTwingateResourceID(serviceKey.TerraformResource(), &serviceAccountResourceID),
				),
			},
			{
				Config: configBuilder(serviceKey.Set(attr.ServiceAccountID, serviceAccount2.TerraformResourceID()), serviceAccount1, serviceAccount2),
				Check: acctests.ComposeTestCheckFunc(
					acctests.CheckTwingateResourceExists(serviceAccount2.TerraformResource()),
					sdk.TestCheckResourceAttr(serviceAccount2.TerraformResource(), attr.Name, serviceAccount2.Name),
					acctests.CheckTwingateResourceExists(serviceKey.TerraformResource()),
					sdk.TestCheckResourceAttrWith(serviceKey.TerraformResource(), attr.Token, nonEmptyValue),

					// test resources were re-created
					sdk.TestCheckResourceAttrWith(serviceKey.TerraformResource(), attr.ID, func(value string) error {
						if *serviceKeyResourceID == "" {
							return errors.New("failed to fetch service_key resource id")
						}

						if value == *serviceKeyResourceID {
							return errors.New("service_key resource was not re-created")
						}

						return nil
					}),

					sdk.TestCheckResourceAttrWith(serviceAccount2.TerraformResource(), attr.ID, func(value string) error {
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
