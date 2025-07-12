package gen

import "math/rand/v2"

func minMaxIntN(min, max int) int {
	return min + rand.IntN(max-min+1)
}
