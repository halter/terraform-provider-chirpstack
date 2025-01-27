---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "chirpstack_http_integration Resource - chirpstack"
subcategory: ""
description: |-
  Http Integration resource
---

# chirpstack_http_integration (Resource)

Http Integration resource

## Example Usage

```terraform
resource "chirpstack_http_integration" "example" {
  application_id     = chirpstack_application.application.id
  encoding           = "JSON"
  event_endpoint_url = "https://example.com"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `application_id` (String) Application ID
- `encoding` (String) Http Integration encoding. JSON or PROTOBUF.
- `event_endpoint_url` (String) Http Integration URL

### Optional

- `headers` (Map of String) Http Integration headers

### Read-Only

- `id` (String) Http Integration identifier
