package config

type AnalyticsConfigs struct {
	topTransformationsPageSize int
	topTrafficPageSize         int
}

func analyticsConfigs() *AnalyticsConfigs {
	return &AnalyticsConfigs{
		topTransformationsPageSize: getInt(topTransformationsPageSize, true),
		topTrafficPageSize:         getInt(topTrafficPageSize, true),
	}
}

func (ac *AnalyticsConfigs) TopTransformationsPageSize() int {
	return ac.topTransformationsPageSize
}

func (ac *AnalyticsConfigs) TopTrafficPageSize() int {
	return ac.topTrafficPageSize
}
