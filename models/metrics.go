package models

type DomainTransformationCounter struct {
	Domain          string
	Transformations float64
}

type DomainRedirectionCounter struct {
	Domain  string
	Traffic float64
}
