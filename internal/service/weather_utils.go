package service

import "strings"

func uniqueCities(cities []string) []string {
	seen := make(map[string]struct{}, len(cities))
	result := make([]string, 0, len(cities))

	for _, city := range cities {
		city = strings.TrimSpace(city)
		if city == "" {
			continue
		}
		key := strings.ToLower(city)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, city)
	}

	return result
}
