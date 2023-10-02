package processors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type URLProcessorTestSuite struct {
	suite.Suite
	ctx       context.Context
	processor *URLProcessor
}

func (suite *URLProcessorTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.processor = NewURLProcessor()
}

func (suite *URLProcessorTestSuite) TestExtractDomainReturnDomainNameSuccessfully() {
	domainName := suite.processor.ExtractDomain(suite.ctx, "https://www.google.com/search?q=infracloud")
	suite.Equal("google", domainName)
}

func TestURLProcessorTestSuite(t *testing.T) {
	suite.Run(t, new(URLProcessorTestSuite))
}
