package processors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type Sha512ProcessorTestSuite struct {
	suite.Suite
	ctx       context.Context
	processor *Sha512Processor
}

func (suite *Sha512ProcessorTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.processor = NewSha512Processor()
}

func (suite *Sha512ProcessorTestSuite) TestEncode() {
	url := "https://www.google.com"
	hash := suite.processor.Encode(suite.ctx, url)
	suite.NotEmpty(hash)
}

func TestSha512ProcessorTestSuite(t *testing.T) {
	suite.Run(t, new(Sha512ProcessorTestSuite))
}
