package config

import (
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/provider/resource"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

type ResourceServiceAccount struct {
	ProtoResource
}

func NewResourceServiceAccount(values ...any) Resource {
	res := &ResourceServiceAccount{
		ProtoResource: ProtoResource{
			Name:     acctest.RandomWithPrefix("service_account"),
			Type:     resource.TwingateServiceAccount,
			Required: make(map[string]Attribute),
			Optional: make(map[string]Attribute),
		},
	}

	return res.Set(append([]any{
		attr.Name, test.RandomName(),
	}, values...)...)
}

func (r *ResourceServiceAccount) Set(values ...any) Resource {
	if len(values)%2 != 0 {
		panic("Set requires key-value pairs")
	}

	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		val := values[i+1]

		switch key {
		case attr.Name:
			r.Required[key] = NewStringAttribute(key, val)
		}

	}

	return r
}
