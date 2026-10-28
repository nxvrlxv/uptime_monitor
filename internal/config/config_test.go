package config

import (
	"reflect"
	"testing"
)

func TestParseTargets(t *testing.T) {
	tests := []struct {
		Name           string
		FilePath       string
		ExpectedOutput []Target
		ExpectedErr    string
	}{
		{
			Name:     "Good output",
			FilePath: "/Users/enjoyer/go_edu/hosts_monitor/targets/targets_test1.json",
			ExpectedOutput: []Target{
				{
					ID:  1,
					URL: "https://go.dev",
				},
				{
					ID:  2,
					URL: "https://httpbin.org/status/500",
				},
				{
					ID:  3,
					URL: "https://fake-target.invalid",
				},
				{
					ID:  4,
					URL: "http://localhost:59998",
				},
			},
			ExpectedErr: "",
		},
		{
			Name:           "No file",
			FilePath:       "",
			ExpectedOutput: nil,
			ExpectedErr:    "can't open file: open : no such file or directory",
		},
		{
			Name:           "empty file",
			FilePath:       "/Users/enjoyer/go_edu/hosts_monitor/targets/targets_test2.json",
			ExpectedOutput: nil,
			ExpectedErr:    "can't parse data: EOF",
		},
		{
			Name:           "garbage in the beggining of file",
			FilePath:       "/Users/enjoyer/go_edu/hosts_monitor/targets/targets_test3.json",
			ExpectedOutput: nil,
			ExpectedErr:    "can't parse file: invalid character 'h' looking for beginning of value",
		},
	}

	for _, ts := range tests {
		targets, err := ParseTargets(ts.FilePath)
		if err == nil && ts.ExpectedErr != "" {
			t.Errorf("Error, succes, expected error: %v", ts.ExpectedErr)
		}
		if err != nil && err.Error() != ts.ExpectedErr {
			t.Errorf("Wrong err: expected: %v, got: %v", ts.ExpectedErr, err.Error())
		}

		if !reflect.DeepEqual(targets, ts.ExpectedOutput) {
			t.Errorf("Wrong answer, expected: %v, got: %v", ts.ExpectedOutput, targets)
		}
	}

}
