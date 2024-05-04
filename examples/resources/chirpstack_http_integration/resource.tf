resource "chirpstack_http_integration" "example" {
  application_id     = chirpstack_application.application.id
  encoding           = "JSON"
  event_endpoint_url = "https://example.com"
}
