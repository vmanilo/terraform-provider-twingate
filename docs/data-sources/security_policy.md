---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "twingate_security_policy Data Source - terraform-provider-twingate"
subcategory: ""
description: |-
  Security Policies are defined in the Twingate Admin Console and determine user and device authentication requirements for Resources.
---

# twingate_security_policy (Data Source)

Security Policies are defined in the Twingate Admin Console and determine user and device authentication requirements for Resources.

## Example Usage

```terraform
data "twingate_security_policy" "foo" {
  name = "<your security policy name>"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (String) Return a Security Policy by its ID. The ID for the Security Policy must be obtained from the Admin API.
- `name` (String) Return a Security Policy that exactly matches this name.

