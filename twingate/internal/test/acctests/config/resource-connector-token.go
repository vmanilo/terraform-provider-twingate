package config

import (
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/provider/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

type ResourceConnectorToken struct {
	ProtoResource
}

func NewResourceConnectorToken(values ...any) Resource {
	res := &ResourceConnectorToken{
		ProtoResource: ProtoResource{
			Name:     acctest.RandomWithPrefix("connector_tokens"),
			Type:     resource.TwingateConnectorTokens,
			Required: make(map[string]Attribute),
			Optional: make(map[string]Attribute),
		},
	}

	return res.Set(values...)
}

func (r *ResourceConnectorToken) Set(values ...any) Resource {
	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		val := values[i+1]

		switch key {
		case attr.ConnectorID:
			r.Required[key] = NewAttribute(key, val.(string))
		}
	}

	return r
}
