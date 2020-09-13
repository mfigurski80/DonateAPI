package main

func find(slice []string, v string) int {
	for i, item := range slice {
		if item == v {
			return i
		}
	}
	return -1
}

func remove(slice []string, i int) []string {
	slice[i] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}
