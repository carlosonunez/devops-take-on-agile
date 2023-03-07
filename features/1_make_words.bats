#!/usr/bin/env bats
load '/helpers/bats-support/load'
load '/helpers/bats-assert/load'

@test "'make words' outputs a frequency map" {
  expected_frequency_map=$(cat <<-EOF
Top five things DevOps hates about Agile, sorted by frequency:

standups 500
grooming 400
story point 300
burndown 200
jira 100
EOF
)
  run `make words`
  assert_output --partial "$expected_frequency_map"
}
