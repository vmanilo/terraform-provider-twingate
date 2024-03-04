package resource

import (
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/attr"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test"
	"github.com/Twingate/terraform-provider-twingate/twingate/internal/test/acctests"
)

type ConnectorToken struct {
	ResourceName string
	ConnectorID  string
}

func NewConnectorToken(connectorID string) *ConnectorToken {
	return &ConnectorToken{
		ResourceName: test.RandomResourceName(),
		ConnectorID:  connectorID,
	}
}

func (r *ConnectorToken) TerraformResource() string {
	return acctests.TerraformConnectorTokens(r.ResourceName)
}

func (r *ConnectorToken) String() string {
	return acctests.Nprintf(`
	resource "twingate_connector_tokens" "${terraform_resource}" {
	  connector_id = ${connector_id}
	}
	`, map[string]any{
		"terraform_resource": r.ResourceName,
		"connector_id":       r.ConnectorID,
	})
}

func (r *ConnectorToken) Set(values ...any) *ConnectorToken {
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
