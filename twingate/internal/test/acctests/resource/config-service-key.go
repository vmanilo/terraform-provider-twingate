package resource

import (
	"fmt"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
	"strings"
)

type ServiceAccountKey struct {
	ResourceName     string
	ServiceAccountID string

	ExpirationTime *int
	Name           *string
}

func NewServiceAccountKey(serviceAccountID string) *ServiceAccountKey {
	return &ServiceAccountKey{
		ResourceName:     test.RandomResourceName(),
		ServiceAccountID: serviceAccountID,
	}
}

func (r *ServiceAccountKey) TerraformResource() string {
	return acctests.TerraformServiceKey(r.ResourceName)
}

func (r *ServiceAccountKey) TerraformResourceID() string {
	return r.TerraformResource() + ".id"
}

func (r *ServiceAccountKey) String() string {
	return acctests.Nprintf(`
	resource "twingate_service_account_key" "${terraform_resource}" {
	  service_account_id = ${service_account_id}

	  ${optional_attributes}
	}
	`, map[string]any{
		"terraform_resource":  r.ResourceName,
		"service_account_id":  r.ServiceAccountID,
		"optional_attributes": r.optionalAttributes(),
	})
}

func (r *ServiceAccountKey) Set(values ...any) *ServiceAccountKey {
	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		val := values[i+1]

		switch key {
		case attr.Name:
			r.Name = optionalString(val)
		case attr.ExpirationTime:
			r.ExpirationTime = optionalInt(val)
		case attr.ServiceAccountID:
			r.ServiceAccountID = val.(string)
		}
	}

	return r
}

func (r *ServiceAccountKey) optionalAttributes() string {
	var optional []string

	if r.Name != nil {
		optional = append(optional, fmt.Sprintf(`name = "%s"`, *r.Name))
	}

	if r.ExpirationTime != nil {
		optional = append(optional, fmt.Sprintf(`expiration_time = %v`, *r.ExpirationTime))
	}

	return strings.Join(optional, "\n")
}
