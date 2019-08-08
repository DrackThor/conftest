package test_test

import (
	"testing"

	"github.com/instrumenta/conftest/pkg/commands/test"
	"github.com/instrumenta/conftest/pkg/commands/test/testfakes"
	"github.com/spf13/viper"
)

func TestWarnQuery(t *testing.T) {

	tests := []struct {
		in  string
		exp bool
	}{
		{"", false},
		{"warn", true},
		{"warnXYZ", false},
		{"warn_", false},
		{"warn_x", true},
		{"warn_x_y_z", true},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			res := test.WarnQ.MatchString(tt.in)

			if tt.exp != res {
				t.Fatalf("%s recognized as `warn` query - expected: %v actual: %v", tt.in, tt.exp, res)
			}
		})
	}
}

func TestCombineConfig(t *testing.T) {
	viper.Set("namespace", "main")
	testTable := []struct {
		name              string
		combineConfigFlag bool
		policyPath        string
		fileList          []string
		shouldFail        bool
		skipTest          bool
	}{
		{
			name:              "combine-config flag exists",
			combineConfigFlag: false,
			policyPath:        "testdata/policy",
			fileList:          []string{"testdata/deployment.yaml"},
		},
		{
			name:              "given a valid policy and multiple configs",
			combineConfigFlag: false,
			policyPath:        "testdata/policy",
			fileList:          []string{"testdata/deployment+service.yaml", "testdata/deployment.yaml"},
		},
		{
			name:              "given a valid policy multiple configs and `combine-config` flag set to false",
			combineConfigFlag: false,
			policyPath:        "testdata/policy",
			fileList:          []string{"testdata/failing_alone.yaml", "testdata/deployment.yaml"},
			shouldFail:        true,
		},
		{
			name:              "given a valid policy multiple configs and `combine-config` flag set to true",
			combineConfigFlag: true,
			policyPath:        "testdata/policy",
			fileList:          []string{"testdata/failing_alone.yaml", "testdata/deployment.yaml"},
			skipTest:          true,
		},
	}

	for _, testunit := range testTable {
		t.Run(testunit.name, func(t *testing.T) {
			if testunit.skipTest {
				t.Skip("not yet implemented")
			}
			viper.Set("combine-config", testunit.combineConfigFlag)
			viper.Set("policy", testunit.policyPath)
			callCount := 0
			outputPrinter := new(testfakes.FakeOutputManager)
			cmd := test.NewTestCommand(func(int) {
				callCount += 1
			}, outputPrinter)
			cmd.Run(cmd, testunit.fileList)
			if outputPrinter.PutCallCount() != len(testunit.fileList) && !testunit.combineConfigFlag {
				t.Errorf("Output manager should print output for each file but it printed %v", outputPrinter.PutCallCount())
			}
			if outputPrinter.PutCallCount() != 1 && testunit.combineConfigFlag {
				t.Errorf("Output manager should have print once but it printed %v", outputPrinter.PutCallCount())
			}
			if testunit.shouldFail && callCount == 0 {
				t.Errorf("should have failed but we did not")
			}
			if !testunit.shouldFail && callCount > 0 {
				t.Errorf("we exited with a failure: %v", callCount)
			}
		})
	}

	t.Run("combine-config flag exists", func(t *testing.T) {
		callCount := 0
		cmd := test.NewTestCommand(func(int) {
			callCount += 1
		}, new(testfakes.FakeOutputManager))
		if cmd.Flag("combine-config") == nil {
			t.Errorf("combine-config flag should exist")
		}
	})
}
func TestFailQuery(t *testing.T) {

	tests := []struct {
		in  string
		exp bool
	}{
		{"", false},
		{"deny", true},
		{"denyXYZ", false},
		{"deny_", false},
		{"deny_x", true},
		{"deny_x_y_z", true},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			res := test.DenyQ.MatchString(tt.in)

			if tt.exp != res {
				t.Fatalf("%s recognized as `fail` query - expected: %v actual: %v", tt.in, tt.exp, res)
			}
		})
	}
}
