package ksqldb

type Row []interface{}

func (r Row) String(idx int) string {
	if idx >= 0 && idx < len(r) {
		switch val := r[idx].(type) {
		case string:
			return val
		}
	}

	return ""
}

func (r Row) Int(idx int) int {
	if idx >= 0 && idx < len(r) {
		switch val := r[idx].(type) {
		case int:
			return val
		}
	}

	return 0
}

func (r Row) Bool(idx int) bool {
	if idx >= 0 && idx < len(r) {
		switch val := r[idx].(type) {
		case bool:
			return val
		}
	}

	return false
}

func (r Row) Float64(idx int) float64 {
	if idx >= 0 && idx < len(r) {
		switch val := r[idx].(type) {
		case float32:
			return float64(val)
		case float64:
			return val
		}
	}

	return 0
}
