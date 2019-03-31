package metrics

const (
	counterMetricType      = "counter"
	gaugeMetricType        = "gauge"
	summaryMetricType      = "summary"
	histogramMetricType    = "histogram"
	counterVecMetricType   = "counterVec"
	gaugeVecMetricType     = "gaugeVec"
	summaryVecMetricType   = "summaryVec"
	histogramVecMetricType = "histogramVec"
)

var AllMetricTypes = []string{
	counterMetricType,
	gaugeMetricType,
	summaryMetricType,
	histogramMetricType,
	counterVecMetricType,
	gaugeVecMetricType,
	summaryVecMetricType,
	histogramVecMetricType,
}

var VecMetricTypes = []string{
	counterVecMetricType,
	gaugeVecMetricType,
	summaryVecMetricType,
	histogramVecMetricType,
}

func IsVectorType(mType string) bool {
	for _, t := range VecMetricTypes {
		if t == mType {
			return true
		}
	}

	return false
}
