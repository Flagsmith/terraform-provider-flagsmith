package flagsmith

// Difference returns a slice of 64-bit integers containing the elements of a that are not present in b.
// If a or b is nil, they are treated as empty slices.
func Difference(a, b *[]int64) []int64 {
	var result []int64
	if a == nil {
		a = &[]int64{}
	}
	if b == nil {
		b = &[]int64{}
	}
	for _, value := range *a {
		found := false
		for _, otherValue := range *b {
			if value == otherValue {
				found = true
				break
			}
		}
		if !found {
			result = append(result, value)
		}
	}
	return result
}
