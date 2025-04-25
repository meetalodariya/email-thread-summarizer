package main

import "time"

func getTimeString(t time.Time) string {
	return t.Format("01/02/06 15:04:05")
}

func addUnique(slice []string, value string) []string {
	for _, elem := range slice {
		if elem == value {
			return slice
		}
	}
	return append(slice, value)
}
