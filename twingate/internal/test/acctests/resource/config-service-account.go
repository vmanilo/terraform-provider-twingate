package resource

import (
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
)

type ServiceAccount struct {
	ResourceName string
	Name         string
}

func NewServiceAccount() *ServiceAccount {
	return &ServiceAccount{
		ResourceName: test.RandomResourceName(),
		Name:         test.RandomServiceAccountName(),
	}
}

func (r *ServiceAccount) TerraformResource() string {
	return acctests.TerraformServiceAccount(r.ResourceName)
}

func (r *ServiceAccount) TerraformResourceID() string {
	return r.TerraformResource() + ".id"
}

func (r *ServiceAccount) String() string {
	return acctests.Nprintf(`
	resource "twingate_service_account" "${terraform_resource}" {
	  name = "${name}"
	}
	`, map[string]any{
		"terraform_resource": r.ResourceName,
		"name":               r.Name,
	})
}

func (r *ServiceAccount) Set(values ...any) *ServiceAccount {
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
