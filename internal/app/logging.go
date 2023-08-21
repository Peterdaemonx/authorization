package app

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"gitlab.cmpayments.local/creditcard/platform"
	"gitlab.cmpayments.local/libraries-go/logging"
	"gitlab.cmpayments.local/libraries-go/logging/loghandling"
	"gitlab.cmpayments.local/libraries-go/logging/logprocessing"
)

func NewLogger(config Config, logName string) platform.Logger {
	var terminal logging.Handler
	if !config.Development.HumanReadableLogging {
		terminal = loghandling.NewTerminal(loghandling.FormatJson)
	} else {
		terminal = loghandling.NewTerminal(devTerminal)
	}

	var terminalMinLevel = logging.Debug
	switch strings.ToLower(config.MinLogLevel) {
	case "informational":
		terminalMinLevel = logging.Informational
	case "notice":
		terminalMinLevel = logging.Notice
	case "warning":
		terminalMinLevel = logging.Warning
	case "error":
		terminalMinLevel = logging.Error
	default:
		terminalMinLevel = logging.Debug
	}

	handler := func(ctx context.Context, record logging.Record) {
		if record.Severity.Level() <= terminalMinLevel.Level() {
			terminal(ctx, record)
		}
	}

	return logging.New(
		logName,
		// LogStateChange to terminal
		handler,
		// Set the caller for linking to the configuration file.
		logprocessing.AddCaller,
		// Mask PAN data
		logprocessing.NewPanMasker(),
		// Add stacktraces to errors and worse
		logprocessing.NewSeverityFilter(logging.Errors, logprocessing.AddStacktrace),
		// Mask the bearer token in the HTTP headers
		logprocessing.NewReplaceHttpHeaderValue("Authorization", mask),
	)
}

func mask(_ string) string {
	return `<masked>`
}

func devTerminal(record logging.Record) []byte {
	var b bytes.Buffer

	t := time.Unix(record.Timestamp, 0).UTC()

	// ANSI color green for the name.
	// ANSI color cyan for the level.
	// ANSI color red for message.
	fmt.Fprintf(
		&b,
		"{ name: \u001b[32m%s\u001B[0m, level: \u001b[36m%s\u001B[0m, message: \u001B[31m%s\u001B[0m datetime: %s  \n  data: %s \n  caller: %s:%d",
		record.Name,
		record.Severity,
		record.Message,
		t.Format("2006-01-02 15:04:05"),
		record.Data,
		record.Data["caller.file"].(string),
		record.Data["caller.line"].(int),
	)

	// Add an underlying error
	// ANSI color red for the error.
	if err, ok := record.Data["error"].(string); ok {
		b.WriteString("\n  error: " + "\u001B[31m" + err + "\u001B[0m" + " }")
	}

	// A line ends with a newline
	b.WriteString("\n")

	return b.Bytes()
}
