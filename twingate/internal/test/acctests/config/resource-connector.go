package config

import (
	"fmt"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/provider/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

type ResourceConnector struct {
	ProtoResource
}

func NewResourceConnector(values ...any) Resource {
	res := &ResourceConnector{
		ProtoResource: ProtoResource{
			Name:     acctest.RandomWithPrefix("connector"),
			Type:     resource.TwingateConnector,
			Required: make(map[string]Attribute),
			Optional: make(map[string]Attribute),
		},
	}

	return res.Set(values...)
}

func (r *ResourceConnector) Set(values ...any) Resource {
	if len(values)%2 != 0 {
		panic("Set requires key-value pairs")
	}

	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		val := values[i+1]

		switch key {
		case attr.Name:
			r.Required[key] = NewStringAttribute(key, val)
		case attr.RemoteNetworkID:
			r.Required[key] = NewAttribute(key, val)
		case attr.StatusUpdatesEnabled:
			r.Optional[key] = NewAttribute(key, fmt.Sprintf("%v", val.(bool)))
		}

	}

	return r
}

//func (r *ResourceConnector) Delete(attributes ...string) Resource {
//	for _, key := range attributes {
//		switch key {
//		case attr.StatusUpdatesEnabled:
//			delete(r.Optional, key)
//		}
//	}
//
//	return r
//}
