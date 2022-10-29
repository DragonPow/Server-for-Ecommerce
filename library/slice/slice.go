package slice

// Map Convert []U to []V with func transform
// Example:
//	slices.Map(products, func(p Product) []int64 {
//		return p.ID
//	})
func Map[U any, V any](sources []U, transfer func(u U) V) []V {
	result := make([]V, len(sources))
	for i, el := range sources {
		result[i] = transfer(el)
	}
	return result
}

// KeyBy transforms a slice or an array of structs to a map based on a pivot callback
// Example:
//	slices.KeyBy(products, func(p Product) (int64, string) {
//		return p.ID, p.TkSku
//	})
func KeyBy[U any, K comparable, V any](sources []U, transform func(u U) (K, V)) map[K]V {
	m := make(map[K]V)
	for _, u := range sources {
		k, v := transform(u)
		m[k] = v
	}
	return m
}

// Uniq returns a duplicate-free version of an array, in which only the first occurrence of each element is kept.
// The order of result values is determined by the order they occur in the array.
func Uniq[T comparable](slices []T) []T {
	result := make([]T, 0, len(slices))
	seen := make(map[T]struct{}, len(slices))

	for _, item := range slices {
		if _, ok := seen[item]; ok {
			continue
		}

		seen[item] = struct{}{}
		result = append(result, item)
	}

	return result
}

// Contains returns true if predicate function return true.
func Contains[T any](slice []T, predicate func(T) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}

	return false
}

// Diff Get []T in sources but not in slices
// Example:
//	slices.Diff(skus, products, func(p Product) string {
//		return p.TkSku
//	})
func Diff[T comparable, E any](sources []T, slices []E, f func(v E) T) []T {
	var diffStr []T
	var defaultValue T
	m := map[T]int{}

	for _, doc := range sources {
		m[doc] = 1
	}
	for _, s := range slices {
		if k := f(s); k != defaultValue {
			m[k]++
		}
	}

	for mKey, mVal := range m {
		if mVal == 1 {
			diffStr = append(diffStr, mKey)
		}
	}
	return diffStr
}

func DiffV2[T any, E any, V comparable](slice1 []T, slice2 []E, f1 func(v1 T) V, f2 func(v2 E) V) []V {
	var diffStr []V
	m1 := make(map[V]bool)
	m2 := make(map[V]bool)

	for _, s := range slice1 {
		k := f1(s)
		m1[k] = true
	}
	for _, s := range slice2 {
		k := f2(s)
		m2[k] = true
	}

	for value := range m1 {
		if _, ok := m2[value]; !ok {
			diffStr = append(diffStr, value)
		} else {
			delete(m2, value)
		}
	}
	for value := range m2 {
		if _, ok := m1[value]; !ok {
			diffStr = append(diffStr, value)
		}
	}
	return diffStr
}
