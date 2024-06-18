package app_test

import (
	"bytes"
	"fmt"
	stdlog "log"
	"os"
	"regexp"
	"runtime"
	"testing"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"gotest.tools/v3/assert"
	"wiki.leta.lv/dmitrijs_zigunovs/gopkgs/log"
)

func TestAppAllLoggerLevels(t *testing.T) {
	log.SetLevel(sentry.LevelDebug)
	var (
		buf = new(bytes.Buffer)
		l   = log.New(buf, "", stdlog.Lshortfile)
		msg = "example of something wrong"
		err = fmt.Errorf(msg)
	)
	for _, logFn := range []func(error){
		l.Debug,
		l.Info,
		l.Error,
		l.Warning,
		// func(err error) { l.Fatal(nil, err) },
	} {
		buf.Reset()
		logFn(err)
		assert.Assert(t,
			regexp.MustCompile("^\\[\\w+\\] app_test.go:\\d+: "+msg+`\n$`).
				MatchString(buf.String()),
			"buf.String() = %q",
			buf.String(),
		)
	}
}

func TestNilErr(t *testing.T) {
	log.Error(errors.WithStack(log.Chain(nil).Extra(log.Extra{"key": "value"})))
}

func TestAppSendError(t *testing.T) {
	log.SetLevel(sentry.LevelInfo)
	log.SetOutput(log.Stderr)
	log.SetProject("test")
	sentry.Init(sentry.ClientOptions{
		Dsn:        "http://88ac80e9067b47bfaac1e06012f32bdb@192.168.1.237:9080/3",
		ServerName: func() string { name, _ := os.Hostname(); return name }(),
		// Release:    version,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if executable, err := os.Executable(); err == nil {
				event.Tags["executable"] = executable
			} else {
				event.Extra["executable"] = err.Error()
			}
			return event
		},
	})
	defer sentry.Flush(2 * time.Second)
	log.Error(log.Chain(errors.New("test message")).Extra(log.Extra{
		"extraDataInt":    42,
		"extraDataString": "Hello world!",
	}))
}

func TestLogChainLocation(t *testing.T) {
	var (
		buf = new(bytes.Buffer)
		l   = log.New(buf, "", stdlog.Lshortfile)
		msg = "LogChainLocation"
		err = fmt.Errorf(msg)
		reg = func(callerSkip int) string {
			_, _, line, _ := runtime.Caller(callerSkip)
			return fmt.Sprintf(`^\[\w+\] app_test.go:%d: .*%s.*\n$`, line-1, msg)
		}
		check = func() assert.BoolOrComparison {
			return regexp.MustCompile(reg(2)).
				MatchString(buf.String())
		}
	)

	l.Error(err)
	assert.Assert(t, check(), "buf.String() = %q, but reg = %q", buf.String(), reg(1))
	buf.Reset()

	l.Error(log.Chain(err))
	assert.Assert(t, check(), "buf.String() = %q, but reg = %q", buf.String(), reg(1))
	buf.Reset()

	l.Error(log.ChainMsg(msg))
	assert.Assert(t, check(), "buf.String() = %q, but reg = %q", buf.String(), reg(1))
	buf.Reset()

	l.Error(log.ChainF("%s", msg))
	assert.Assert(t, check(), "buf.String() = %q, but reg = %q", buf.String(), reg(1))
	buf.Reset()

	l.Error(log.ChainMsg(msg))
	assert.Assert(t, check(), "buf.String() = %q, but reg = %q", buf.String(), reg(1))
	buf.Reset()

	l.Error(log.ChainWithMsg(err, msg))
	assert.Assert(t, check(), "buf.String() = %q, but reg = %q", buf.String(), reg(1))
	buf.Reset()

	l.Error(log.ChainWithMsgF(err, "%s", msg))
	assert.Assert(t, check(), "buf.String() = %q, but reg = %q", buf.String(), reg(1))
	buf.Reset()
}

func TestCustomStack(t *testing.T) {
	log.SetLevel(sentry.LevelError)
	log.SetFlags(16)
	var (
		buf = new(bytes.Buffer)
		l   = log.New(buf, "", stdlog.Lshortfile)
		msg = "simple error message"
	)
	buf.Reset()
	var err error = func(here bool) (err error) {
		if here {
			return log.Stack(0).ChainMsg(msg)
		}
		return
	}(true)
	l.Error(err)
	assert.Equal(t,
		buf.String(),
		fmt.Sprintf(
			"[ERROR] app_test.go:%d: %s\n",
			lineNum()-9,
			msg,
		),
		"buf.String() = %q",
		buf.String(),
	)
}

func lineNum() int {
	_, _, line, _ := runtime.Caller(1)
	return line
}
