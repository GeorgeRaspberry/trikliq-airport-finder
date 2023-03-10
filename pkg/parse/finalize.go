package parse

import "trikliq-airport-finder/pkg/transform"

func Finalize(finalCandidates []string, iata map[string]map[string]string) (result map[string]any) {

	departing := make([]map[string]string, 0)
	arriving := make([]map[string]string, 0)

	for i, candidate := range finalCandidates {
		if i%2 == 0 && !transform.InSliceOfMap(departing, candidate) {
			departing = append(departing, iata[candidate])
		}

		if i%2 != 0 && !transform.InSliceOfMap(arriving, candidate) {
			arriving = append(arriving, iata[candidate])
		}
	}

	result = make(map[string]any, 0)
	result["departing"] = departing
	result["arriving"] = arriving

	return
}
