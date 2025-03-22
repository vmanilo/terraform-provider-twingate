package config

import (
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/provider/resource"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

type ResourceRemoteNetwork struct {
	ProtoResource
}

func NewResourceRemoteNetwork(values ...any) Resource {
	res := &ResourceRemoteNetwork{
		ProtoResource: ProtoResource{
			Name:     acctest.RandomWithPrefix("remote_network"),
			Type:     resource.TwingateRemoteNetwork,
			Required: make(map[string]Attribute),
			Optional: make(map[string]Attribute),
		},
	}

	return res.Set(append([]any{
		attr.Name, test.RandomName(),
	}, values...)...)
}

func (r *ResourceRemoteNetwork) Set(values ...any) Resource {
	if len(values)%2 != 0 {
		panic("Set requires key-value pairs")
	}

	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		val := values[i+1]

		switch key {
		case attr.Name:
			r.Required[key] = NewStringAttribute(key, val)

		case attr.Location:
			r.Optional[key] = NewStringAttribute(key, val)
		case attr.Type:
			r.Optional[key] = NewStringAttribute(key, val)
		}

	}

	return r
}

//func (r *ResourceRemoteNetwork) Delete(attributes ...string) Resource {
//	for _, key := range attributes {
//		switch key {
//		case attr.Location, attr.Type:
//			delete(r.Optional, key)
//		}
//	}
//
//	return r
//}
