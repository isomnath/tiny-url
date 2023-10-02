package contracts

type TransformationResponse struct {
	TopTransformations []Transformation `json:"top_transformations"`
}

type Transformation struct {
	Domain string  `json:"domain"`
	Count  float64 `json:"count"`
}

type TrafficResponse struct {
	TopTraffic []Traffic `json:"top_traffic"`
}

type Traffic struct {
	Domain string  `json:"domain"`
	Count  float64 `json:"count"`
}
