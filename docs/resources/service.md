---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "twingate_service Resource - terraform-provider-twingate"
subcategory: ""
description: |-
  Services offer a way to provide programmatic, centrally-controlled, and consistent access controls. For more information, see Twingate's documentation https://www.twingate.com/docs/services.
---

# twingate_service (Resource)

Services offer a way to provide programmatic, centrally-controlled, and consistent access controls. For more information, see Twingate's [documentation](https://www.twingate.com/docs/services).

## Example Usage

```terraform
provider "twingate" {
  api_token = "1234567890abcdef"
  network   = "mynetwork"
}

resource "twingate_service" "github_actions_prod" {
  name = "Github Actions PROD"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the service account in Twingate

### Read-Only

- `id` (String) Autogenerated ID of the service account

