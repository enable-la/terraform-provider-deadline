resource "deadline_farm" "test" {
  display_name = "test"
  description  = "this is a test farm"
}

resource "deadline_associate_member_to_farm" "test" {
  farm_id           = deadline_farm.test.id
  member_id         = "test"
  identity_store_id = "example_identity_store"
  membership_level  = "VIEWER"
  principal_type    = "USER"
}