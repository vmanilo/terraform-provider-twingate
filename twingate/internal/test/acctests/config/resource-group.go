package config

import (
	"fmt"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/provider/resource"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

type ResourceGroup struct {
	ProtoResource
}

func NewResourceGroup(values ...any) Resource {
	res := &ResourceGroup{
		ProtoResource: ProtoResource{
			Name:     acctest.RandomWithPrefix("group"),
			Type:     resource.TwingateGroup,
			Required: make(map[string]Attribute),
			Optional: make(map[string]Attribute),
		},
	}

	return res.Set(append([]any{
		attr.Name, test.RandomName(),
	}, values...)...)
}

func (r *ResourceGroup) Set(values ...any) Resource {
	if len(values)%2 != 0 {
		panic("Set requires key-value pairs")
	}

	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		val := values[i+1]

		switch key {
		case attr.Name:
			r.Required[key] = NewStringAttribute(key, val)

		case attr.SecurityPolicyID:
			r.Optional[key] = NewStringAttribute(key, val)
		case attr.UserIDs:
			r.Optional[key] = NewSetAttribute(key, val.([]string))
		case attr.IsAuthoritative:
			r.Optional[key] = NewAttribute(key, fmt.Sprintf("%v", val.(bool)))
		}

	}

	return r
}

//func (r *ResourceGroup) Delete(attributes ...string) Resource {
//	r.ProtoResource.Delete(attributes...)
//
//	//for _, key := range attributes {
//	//	switch key {
//	//	case attr.SecurityPolicyID, attr.UserIDs, attr.IsAuthoritative:
//	//		delete(r.Optional, key)
//	//	}
//	//}
//
//	return r
//}
