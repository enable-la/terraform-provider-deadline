resource "deadline_farm" "test" {
  display_name = "test"
  description  = "this is a test farm"
}

resource "deadline_fleet" "test" {
  display_name = "test"
  farm_id      = deadline_farm.test.id
  description  = "this is a test farm"
}
