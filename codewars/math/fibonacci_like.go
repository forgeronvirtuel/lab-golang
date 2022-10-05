package math

func Tribonacci(init [3]float64, n int) []float64 {
	switch n {
	case 0:
		return []float64{}
	case 1:
		return []float64{init[0]}
	case 2:
		return []float64{init[0], init[1]}
	case 3:
		return []float64{init[0], init[1], init[2]}
	default:
		var res = make([]float64, n)
		res[0] = init[0]
		res[1] = init[1]
		res[2] = init[2]
		for i := 3; i < n; i++ {
			res[i] = res[i-1] + res[i-2] + res[i-3]
		}
		return res
	}
}
