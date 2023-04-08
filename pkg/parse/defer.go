package parse

// func Defer(foundCities []string, iata map[string]map[string]string, log *zap.Logger) []string {
// 	log.Debug("deferring to cities")

// 	matched := make([]string, 0)

// 	for _, code := range iata {
// 		for _, candidate := range foundCities {
// 			check := code["city"]
// 			if check == "" {
// 				continue
// 			}

// 			split := strings.Split(strings.ToLower(check), " ")
// 			if !transform.InSlice(candidate, split) {
// 				//fmt.Println("continuing check: ", check)
// 				continue
// 			}

// 			fmt.Println("candidate: ", candidate)
// 			fmt.Println("found: ", check)
// 			matched = append(matched, check)
// 		}
// 	}

// 	return matched
// }
