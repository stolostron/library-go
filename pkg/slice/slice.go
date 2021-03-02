// Copyright Contributors to the Open Cluster Management project

package slice

// AppendIfDNE append stringToAppend to stringSlice if stringToAppend does not already exist in stringSlice
func AppendIfDNE(stringSlice []string, stringToAppend string) []string {
	toAppend := true

	for _, slice := range stringSlice {
		if slice == stringToAppend {
			toAppend = false
		}
	}

	if toAppend {
		stringSlice = append(stringSlice, stringToAppend)
	}

	return stringSlice
}

// RemoveFromStringSlice takes a string[] and remove all stringToRemove
func RemoveFromStringSlice(stringSlice []string, stringToRemove string) []string {
	for i, slice := range stringSlice {
		if slice == stringToRemove {
			stringSlice = append(stringSlice[0:i], stringSlice[i+1:]...)
			return RemoveFromStringSlice(stringSlice, stringToRemove)
		}
	}

	return stringSlice
}

// UniqueStringSlice takes a string[] and remove the duplicate value
func UniqueStringSlice(stringSlice []string) []string {
	keys := make(map[string]bool)
	uniqueStringSlice := []string{}

	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true

			uniqueStringSlice = append(uniqueStringSlice, entry)
		}
	}

	return uniqueStringSlice
}
