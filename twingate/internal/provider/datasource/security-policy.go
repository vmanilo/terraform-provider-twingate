package datasource

import (
	"context"

	"github.com/Twingate/terraform-provider-twingate/twingate/internal/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func SecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Security Policies are defined in the Twingate Admin Console and determine user and device authentication requirements for Resources.",
		ReadContext: readSecurityPolicy,
		Schema: map[string]*schema.Schema{
			fieldID: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Return a Security Policy by its ID. The ID for the Security Policy must be obtained from the Admin API.",
				ExactlyOneOf: []string{fieldName},
			},
			fieldName: {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Return a Security Policy that exactly matches this name.",
				ExactlyOneOf: []string{fieldID},
			},
		},
	}
}

func readSecurityPolicy(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*client.Client)

	securityPolicy, err := client.ReadSecurityPolicy(ctx, resourceData.Get(fieldID).(string), resourceData.Get(fieldName).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := resourceData.Set(fieldName, securityPolicy.Name); err != nil {
		return diag.FromErr(err)
	}

	resourceData.SetId(securityPolicy.ID)

	return nil
}