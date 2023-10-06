package internal

import (
	"testing"

	tester_utils_testing "github.com/codecrafters-io/tester-utils/testing"
)

type stageYAML struct {
	Slug  string `yaml:"slug"`
	Title string `yaml:"name"`
}

type courseYAML struct {
	Stages []stageYAML `yaml:"stages"`
}

func TestStagesMatchYAML(t *testing.T) {
	tester_utils_testing.ValidateTesterDefinitionAgainstYAML(t, testerDefinition, "test_helpers/course_definition.yml")
}
