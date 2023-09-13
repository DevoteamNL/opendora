package devlake

import (
	"devlake-go/group-sync/pkg/test"
	"testing"
)

func TestLargestTeamIdTableDriven(t *testing.T) {

	var tests = []struct {
		name  string
		input [][]string
		want  int
	}{
		{"should use 0 for empty table", [][]string{}, 0},
		{"should use 0 for table without numbers in the first column", test.ExampleCsvWithColumnHeaders([][]string{}), 0},
		{"should use 0 for table without numbers in the first column", test.ExampleCsvWithColumnHeaders(test.CreateTestDevLakeTeamsWithIds([]string{"A"})), 0},
		{"should use the largest number found in the first column", test.CreateTestDevLakeTeamsWithIds([]string{"3", "4", "5"}), 5},
		{"should use the largest number found in the first column", test.CreateTestDevLakeTeamsWithIds([]string{"50", "40", "30"}), 50},
		{"should use the largest number found in the first column", test.CreateTestDevLakeTeamsWithIds([]string{"40", "A", "30"}), 40},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ans := LargestTeamId(tt.input)
			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}
