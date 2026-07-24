package system

func Utilization(used, total uint64) float64 {
	if total == 0 {
		return 0
	}
	return clampPercentage(float64(used) / float64(total) * 100)
}

func clampPercentage(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return value
}

func uint64Pointer(value uint64) *uint64  { return &value }
func intPointer(value int) *int           { return &value }
func floatPointer(value float64) *float64 { return &value }
