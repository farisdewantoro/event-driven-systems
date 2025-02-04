package logger_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"testing"

	"eventdrivensystem/pkg/logger"
	"eventdrivensystem/pkg/logger/mock"

	"gotest.tools/assert"
)

type Dependency struct {
	log    logger.Logger
	writer bytes.Buffer
}

type lg struct {
	RequestID string `json:"request_id"`
	Message   string `json:"msg"`
	Level     string `json:"level"`
}

var d Dependency

func TestMain(m *testing.M) {
	var buf bytes.Buffer

	d.writer = buf
	d.log = logger.Init(logger.Options{
		Output:       logger.OutputCustom,
		CustomWriter: &d.writer,
		Formatter:    logger.FormatJSON,
		Level:        logger.LevelInfo,
		ContextFields: map[string]string{
			"request_id": "x-request-id",
		},
		Hooks: []logger.Hook{mock.SampleHook{}},
		Hook:  logger.OTEL,
	})

	// Run tests!
	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestInfoLogger(t *testing.T) {

	testCases := []struct {
		name string
		in   string
		out  struct {
			expected lg
			err      error
		}
	}{
		{
			name: "Test Info",
			in:   "test",
			out: struct {
				expected lg
				err      error
			}{
				expected: lg{
					Message: "test",
					Level:   "info",
				},
				err: nil,
			},
		},
		{
			name: "Test Info 2",
			in:   "asd",
			out: struct {
				expected lg
				err      error
			}{
				expected: lg{
					Message: "asd",
					Level:   "info",
				},
				err: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d.writer.Reset()
			d.log.Info(tc.in)

			var tLog lg
			err := json.Unmarshal(d.writer.Bytes(), &tLog)
			if err != nil {
				t.Error(err)
			}

			t.Log(tLog)

			assert.DeepEqual(t, tLog.Message, tc.out.expected.Message)
			assert.DeepEqual(t, tLog.Level, tc.out.expected.Level)
		})
	}

}

func TestInfoCtxLogger(t *testing.T) {

	type ctx struct {
		RequestID string
	}

	testCases := []struct {
		name  string
		inCtx ctx
		in    string
		out   struct {
			expected lg
			err      error
		}
	}{
		{
			name: "Test Info Context",
			inCtx: ctx{
				RequestID: "req ID",
			},
			in: "asd",
			out: struct {
				expected lg
				err      error
			}{
				expected: lg{
					RequestID: "req ID",
					Message:   "asd",
					Level:     "info",
				},
				err: nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d.writer.Reset()
			contx := context.WithValue(context.TODO(), "x-request-id", tc.inCtx.RequestID)
			d.log.InfoWithContext(contx, tc.in)

			var tLog lg
			err := json.Unmarshal(d.writer.Bytes(), &tLog)
			if err != nil {
				t.Error(err)
			}

			assert.DeepEqual(t, tLog.Message, tc.out.expected.Message)
			assert.DeepEqual(t, tLog.RequestID, tc.out.expected.RequestID)
			assert.DeepEqual(t, tLog.Level, tc.out.expected.Level)
		})
	}

}
