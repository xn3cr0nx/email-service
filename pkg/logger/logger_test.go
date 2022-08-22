package logger_test

import (
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/xn3cr0nx/email-service/pkg/logger"
)

func (s *LoggerTestSuite) TestLogger() {
	logger.Setup()
	hook := test.NewLocal(logger.Log)

	logger.Info("test", "testing Info function", logger.Params{"test": "Params"})
	logger.Warn("test", "testing Warn function", logger.Params{"test": "Params"})
	logger.Debug("test", "testing Debug function", logger.Params{"test": "Params"})
	logger.Error("test", errors.New("testing Error function"), logger.Params{"error": "error"})

	s.Equal(3, len(hook.Entries))
	s.Equal(string(logrus.InfoLevel), string(hook.Entries[0].Level))
	s.Equal(string(logrus.WarnLevel), string(hook.Entries[1].Level))
	s.Equal(string(logrus.ErrorLevel), string(hook.LastEntry().Level))
	s.Equal("testing Error function", hook.LastEntry().Message)

	hook.Reset()
	s.Nil(hook.LastEntry())
}
