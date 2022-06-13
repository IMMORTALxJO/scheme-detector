package main

func stringInArray(k string, arr []string) bool {
	for _, a := range arr {
		if a == k {
			return true
		}
	}
	return false
}
