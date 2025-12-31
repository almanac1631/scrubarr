package utils

import (
	"fmt"
	"math"
)

func FormatBytes(size int64) string {
	if size < 0 {
		return "0 B"
	}
	if size == 0 {
		return "0 B"
	}

	units := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}

	// Logarithm base 1024 to find the correct index in the units slice
	i := int(math.Floor(math.Log(float64(size)) / math.Log(1024)))

	// Calculate the value in that unit
	value := float64(size) / math.Pow(1024, float64(i))

	// Format to 1 decimal place (e.g., "2.5 MB")
	return fmt.Sprintf("%.1f %s", value, units[i])
}
