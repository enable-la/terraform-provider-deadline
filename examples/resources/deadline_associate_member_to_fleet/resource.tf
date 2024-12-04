resource "deadline_farm" "test" {
  display_name = "test"
  description  = "this is a test farm"
}

resource "deadline_fleet" "test" {
  display_name = "test"
  farm_id      = deadline_farm.test.id
  description  = "this is a test farm"
}

resource "deadline_associate_member_to_fleet" "test" {
  farm_id           = deadline_farm.test.id
  fleet_id          = deadline_fleet.test.id
  member_id         = "test"
  identity_store_id = "example_identity_store"
  membership_level  = "VIEWER"
  principal_type    = "USER"
}