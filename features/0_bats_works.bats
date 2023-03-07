#!/usr/bin/env bats
load '/helpers/bats-assert/load'

@test "BATS works" {
  run echo 'hello world!'
  assert_success
}
