package logger_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(LoggerTestSuite))
}

type LoggerTestSuite struct {
	suite.Suite
}
