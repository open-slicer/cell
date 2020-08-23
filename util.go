package main

// someZero checks if a single value is zeroed. Supports: string, int.
func someZero(vals ...interface{}) bool {
	if vals[0] == nil {
		return true
	}
	check := func(val interface{}) bool {
		return true
	}

	switch vals[0].(type) {
	case string:
		check = func(val interface{}) bool {
			if val == "" {
				return true
			}
			return false
		}
	case int:
		check = func(val interface{}) bool {
			if val == 0 {
				return true
			}
			return false
		}
	}

	for _, val := range vals {
		if val == nil {
			return true
		}
		if check(val) {
			return true
		}
	}
	return false
}
