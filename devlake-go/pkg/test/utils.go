package test

import "strconv"

func ExampleCsvWithColumnHeaders(otherRows [][]string) [][]string {
	return append([][]string{{"Id", "Name", "Alias", "ParentId", "SortingIndex"}}, otherRows...)
}

func CreateTestDevLakeTeamsWithIds(ids []string) [][]string {
	teams := [][]string{}

	for _, id := range ids {
		teams = append(teams, []string{id, "B", "C", "D", "E"})
	}

	return teams
}

func CreateTestDevLakeTeamsWithNames(names []string) [][]string {
	teams := [][]string{}

	for i, name := range names {
		teams = append(teams, []string{strconv.Itoa(i), name, "C", "D", "E"})
	}

	return teams
}
