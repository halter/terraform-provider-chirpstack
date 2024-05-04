resource "chirpstack_application" "application" {
  tenant_id   = chirpstack_tenant.tenant.id
  name        = "myapp"
  description = "My Chirpstack Application"
}
