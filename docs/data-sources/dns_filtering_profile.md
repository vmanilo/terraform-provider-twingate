---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "twingate_dns_filtering_profile Data Source - terraform-provider-twingate"
subcategory: ""
description: |-
  DNS filtering gives you the ability to control what websites your users can access. For more information, see Twingate's documentation https://www.twingate.com/docs/dns-filtering.
---

# twingate_dns_filtering_profile (Data Source)

DNS filtering gives you the ability to control what websites your users can access. For more information, see Twingate's [documentation](https://www.twingate.com/docs/dns-filtering).

## Example Usage

```terraform
provider "twingate" {
  api_token = "1234567890abcdef"
  network   = "mynetwork"
}

data "twingate_dns_filtering_profile" "example" {
  id = "<your dns profile's id>"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (String) The DNS filtering profile's ID.

### Read-Only

- `allowed_domains` (Block, Read-only) A block with the following attributes. (see [below for nested schema](#nestedblock--allowed_domains))
- `content_categories` (Block, Read-only) A block with the following attributes. (see [below for nested schema](#nestedblock--content_categories))
- `denied_domains` (Block, Read-only) A block with the following attributes. (see [below for nested schema](#nestedblock--denied_domains))
- `fallback_method` (String) The DNS filtering profile's fallback method. One of AUTOMATIC or STRICT.
- `groups` (Set of String) A set of group IDs that have this as their DNS filtering profile. Defaults to an empty set.
- `name` (String) The DNS filtering profile's name.
- `priority` (Number) A floating point number representing the profile's priority.
- `privacy_categories` (Block, Read-only) A block with the following attributes. (see [below for nested schema](#nestedblock--privacy_categories))
- `security_categories` (Block, Read-only) A block with the following attributes. (see [below for nested schema](#nestedblock--security_categories))

<a id="nestedblock--allowed_domains"></a>
### Nested Schema for `allowed_domains`

Read-Only:

- `domains` (Set of String) A set of allowed domains.


<a id="nestedblock--content_categories"></a>
### Nested Schema for `content_categories`

Read-Only:

- `block_adult_content` (Boolean) Whether to block adult content.
- `block_dating` (Boolean) Whether to block dating content.
- `block_gambling` (Boolean) Whether to block gambling content.
- `block_games` (Boolean) Whether to block games.
- `block_piracy` (Boolean) Whether to block piracy sites.
- `block_social_media` (Boolean) Whether to block social media.
- `block_streaming` (Boolean) Whether to block streaming content.
- `enable_safesearch` (Boolean) Whether to force safe search.
- `enable_youtube_restricted_mode` (Boolean) Whether to force YouTube to use restricted mode.


<a id="nestedblock--denied_domains"></a>
### Nested Schema for `denied_domains`

Read-Only:

- `domains` (Set of String) A set of denied domains.


<a id="nestedblock--privacy_categories"></a>
### Nested Schema for `privacy_categories`

Read-Only:

- `block_ads_and_trackers` (Boolean) Whether to block ads and trackers.
- `block_affiliate_links` (Boolean) Whether to block affiliate links.
- `block_disguised_trackers` (Boolean) Whether to block disguised third party trackers.


<a id="nestedblock--security_categories"></a>
### Nested Schema for `security_categories`

Read-Only:

- `block_cryptojacking` (Boolean) Whether to block cryptojacking sites.
- `block_dns_rebinding` (Boolean) Blocks public DNS entries from returning private IP addresses.
- `block_domain_generation_algorithms` (Boolean) Blocks DGA domains.
- `block_idn_homoglyph` (Boolean) Whether to block homoglyph attacks.
- `block_newly_registered_domains` (Boolean) Blocks newly registered domains.
- `block_parked_domains` (Boolean) Block parked domains.
- `block_typosquatting` (Boolean) Blocks typosquatted domains.
- `enable_google_safe_browsing` (Boolean) Whether to use Google Safe browsing lists to block content.
- `enable_threat_intelligence_feeds` (Boolean) Whether to filter content using threat intelligence feeds.