package config

import (
	"fmt"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/provider/resource"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

type ResourceUser struct {
	ProtoResource
}

func NewResourceUser(values ...any) Resource {
	res := &ResourceUser{
		ProtoResource: ProtoResource{
			Name:     acctest.RandomWithPrefix("user"),
			Type:     resource.TwingateUser,
			Required: make(map[string]Attribute),
			Optional: make(map[string]Attribute),
		},
	}

	return res.Set(append([]any{
		attr.Email, test.RandomEmail(),
		attr.IsActive, true,
		attr.SendInvite, false,
	}, values...)...)
}

func (r *ResourceUser) Set(values ...any) Resource {
	if len(values)%2 != 0 {
		panic("Set requires key-value pairs")
	}

	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		val := values[i+1]

		switch key {
		case attr.Email:
			r.Required[key] = NewStringAttribute(key, val)
		case attr.SendInvite:
			r.Required[key] = NewAttribute(key, fmt.Sprintf("%v", val.(bool)))
		case attr.IsActive:
			r.Required[key] = NewAttribute(key, fmt.Sprintf("%v", val.(bool)))

		case attr.FirstName:
			r.Optional[key] = NewStringAttribute(key, val)
		case attr.LastName:
			r.Optional[key] = NewStringAttribute(key, val)
		case attr.Role:
			r.Optional[key] = NewStringAttribute(key, val)
		}

	}

	return r
}
