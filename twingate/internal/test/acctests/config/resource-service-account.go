package config

import (
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test/acctests"
)

type ResourceServiceAccount struct {
	ResourceName string
	Name         string
}

func NewResourceServiceAccount() *ResourceServiceAccount {
	return &ResourceServiceAccount{
		ResourceName: test.RandomResourceName(),
		Name:         test.RandomServiceAccountName(),
	}
}

func (r *ResourceServiceAccount) TerraformResource() string {
	return acctests.TerraformServiceAccount(r.ResourceName)
}

func (r *ResourceServiceAccount) TerraformResourceID() string {
	return r.TerraformResource() + ".id"
}

func (r *ResourceServiceAccount) String() string {
	return Nprintf(`
	resource "twingate_service_account" "${terraform_resource}" {
	  name = "${name}"
	}
	`, map[string]any{
		"terraform_resource": r.ResourceName,
		"name":               r.Name,
	})
}

func (r *ResourceServiceAccount) Set(values ...any) *ResourceServiceAccount {
	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		val := values[i+1]

		switch key {
		case attr.Name:
			r.Name = val.(string)
		}
	}

	return r
}
