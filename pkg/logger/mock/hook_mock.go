package mock

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type SampleHook struct{}

func (s SampleHook) Fire(e *logrus.Entry) error {
	fmt.Println("Hook success")
	return nil
}

func (s SampleHook) Levels() []logrus.Level {
	levels := []logrus.Level{logrus.InfoLevel, logrus.ErrorLevel}
	return levels
}
