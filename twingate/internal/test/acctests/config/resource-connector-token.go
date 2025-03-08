package config

import (
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/v3/twingate/internal/test/acctests"
)

type ResourceConnectorToken struct {
	ResourceName string
	ConnectorID  string
}

func NewResourceConnectorToken(connectorID string) *ResourceConnectorToken {
	return &ResourceConnectorToken{
		ResourceName: test.RandomResourceName(),
		ConnectorID:  connectorID,
	}
}

func (r *ResourceConnectorToken) TerraformResource() string {
	return acctests.TerraformConnectorTokens(r.ResourceName)
}

func (r *ResourceConnectorToken) String() string {
	return Nprintf(`
	resource "twingate_connector_tokens" "${terraform_resource}" {
	  connector_id = ${connector_id}
	}
	`, map[string]any{
		"terraform_resource": r.ResourceName,
		"connector_id":       r.ConnectorID,
	})
}

func (r *ResourceConnectorToken) Set(values ...any) *ResourceConnectorToken {
	for i := 0; i < len(values); i += 2 {
		key := values[i].(string)
		val := values[i+1]

		switch key {
		case attr.ConnectorID:
			r.ConnectorID = val.(string)
		}
	}

	return r
}
