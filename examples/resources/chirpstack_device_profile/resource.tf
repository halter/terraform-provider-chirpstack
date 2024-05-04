resource "chirpstack_device_profile" "mydeviceprofile" {
  tenant_id                       = chirpstack_tenant.tenant.id
  name                            = "mydeviceprofile"
  description                     = "My Device Profile"
  mac_version                     = "LORAWAN_1_0_3"
  region                          = "AU915"
  region_parameters_revision      = "A"
  adr_algorithm                   = "default"
  allow_roaming                   = false
  device_status_request_frequency = 1
  device_supports_class_b         = false
  device_supports_class_c         = false
  device_supports_otaa            = true
  expected_uplink_interval        = 3600
  flush_queue_on_activate         = true
}
