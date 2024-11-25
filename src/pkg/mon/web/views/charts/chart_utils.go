package charts

import (
	"github.com/go-echarts/go-echarts/v2/opts"
)

//nolint:forcetypeassert // Be careful when using this function!
func diffSeries[T int | int32 | int64 | float32 | float64](series []opts.LineData) []opts.LineData {
	total := len(series) - 1
	result := make([]opts.LineData, total)
	for i := 0; i < total; i++ {
		value := series[i].Value.([]any)
		next := series[i+1].Value.([]any)
		diff := value[1].(T) - next[1].(T)
		if diff < 0 {
			value[1] = 0
		} else {
			value[1] = diff
		}
		result[i] = opts.LineData{Value: value}
	}
	return result
}
