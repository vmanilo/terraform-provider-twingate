package config

import (
	"fmt"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test/acctests"
	"strings"
)

type ResourceServiceAccountKey struct {
	ResourceName     string
	ServiceAccountID string

	ExpirationTime *int
	Name           *string
}

func NewResourceServiceAccountKey(serviceAccountID string) *ResourceServiceAccountKey {
	return &ResourceServiceAccountKey{
		ResourceName:     test.RandomResourceName(),
		ServiceAccountID: serviceAccountID,
	}
}

func (r *ResourceServiceAccountKey) TerraformResource() string {
	return acctests.TerraformServiceKey(r.ResourceName)
}

func (r *ResourceServiceAccountKey) TerraformResourceID() string {
	return r.TerraformResource() + ".id"
}

func (r *ResourceServiceAccountKey) String() string {
	return Nprintf(`
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

func (r *ResourceServiceAccountKey) Set(values ...any) *ResourceServiceAccountKey {
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

func (r *ResourceServiceAccountKey) optionalAttributes() string {
	var optional []string

	if r.Name != nil {
		optional = append(optional, fmt.Sprintf(`name = "%s"`, *r.Name))
	}

	if r.ExpirationTime != nil {
		optional = append(optional, fmt.Sprintf(`expiration_time = %v`, *r.ExpirationTime))
	}

	return strings.Join(optional, "\n")
}
