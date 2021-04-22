package parser

import "testing"

func TestParse(t *testing.T) {
	testSuccessDataSet := map[string][3]string{
		"foo|bar->zee": {"foo", "bar", "zee"},
		"foo||bar->zee": {"foo", "|bar", "zee"},
		"foo|bar->->zee": {"foo", "bar", "->zee"},
		"foo||bar->->zee": {"foo", "|bar", "->zee"},
	}

	for input, expected := range testSuccessDataSet {
		scope, label, instruction, err := Parse(input)
		if scope != expected[0] || label != expected[1] || instruction != expected[2] {
			t.Errorf("actual [%q, %q, %q], expected [%q, %q, %q]",
				scope, label, instruction, expected[0], expected[1], expected[2])
		} else if err != nil {
			t.Errorf("unexpected error while parsing: [%q]", input)
		}
	}

	testFailedDataSet := []string{
		"foo",
		"|bar",
		"bar->",
		"|foo->",
		"foo|bar",
		"bar->zee",
	}

	for _, input := range testFailedDataSet {
		_, _, _, err := Parse(input)
		if err == nil {
			t.Errorf("expected error while parsing: [%q]", input)
		}
	}
}

func TestGetScopeLabelPair(t *testing.T) {
	testDataSet := map[string][2]string{
		"foo|bar":  {"foo", "bar"},
		"foo||bar": {"foo", "|bar"},
	}

	for expected, input := range testDataSet {
		actual := GetScopeLabelPair(input[0], input[1])
		if actual != expected {
			t.Errorf("actual %q, expected %q", actual, expected)
		}
	}
}

func TestGetEntry(t *testing.T) {
	testDataSet := map[string][3]string{
		"foo|bar->zee":    {"foo", "bar", "zee"},
		"foo||bar->zee":   {"foo", "|bar", "zee"},
		"foo|bar->->zee":  {"foo", "bar", "->zee"},
		"foo||bar->->zee": {"foo", "|bar", "->zee"},
	}

	for expected, input := range testDataSet {
		actual := GetEntry(input[0], input[1], input[2])
		if actual != expected {
			t.Errorf("actual %q, expected %q", actual, expected)
		}
	}
}
