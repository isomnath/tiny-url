package processors

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type Base62ProcessorTestSuite struct {
	suite.Suite
	ctx       context.Context
	processor *Base62Processor
}

func (suite *Base62ProcessorTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.processor = NewBase62Processor()
}

func (suite *Base62ProcessorTestSuite) TestEncode() {
	encodedStr := suite.processor.Encode(suite.ctx, uint64(1234567890))
	suite.Equal("01LY7VK", encodedStr)
}

func TestBase62ProcessorTestSuite(t *testing.T) {
	suite.Run(t, new(Base62ProcessorTestSuite))
}
