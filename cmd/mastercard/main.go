package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/scheme/mastercard"
)

type config struct {
	host string
	port int
}

var (
	loglevel                       string
	version                        = "0.0.1"
	cfg                            config
	logError, logVerbose, logDebug Logger = NopLogger, NopLogger, NopLogger
)

func main() {
	// the default address is the HA proxy address.
	// address of the actual MIP and its ports is documented on https://cmcom.atlassian.net/wiki/spaces/CA/pages/941588492/Mastercard+Interface+Processor+MIP

	flag.StringVar(&cfg.host, "host", "mip.test.cmpayments.local", "MIP server host")
	flag.IntVar(&cfg.port, "port", 7043, "MIP server port")
	flag.StringVar(&loglevel, "log", "vvv", "log level:\n\tv = only errors are displayed,\n\tvv = opening and closing connection information text is displayed,\n\tvvv = response result is displayed")

	flag.Parse()

	switch true {
	case strings.Contains("v", loglevel):
		logError = PrintLogger
	case strings.Contains("vv", loglevel):
		logError = PrintLogger
		logVerbose = PrintLogger
	case strings.Contains("vvv", loglevel):
		logError = PrintLogger
		logVerbose = PrintLogger
		logDebug = PrintLogger
	}

	log.Printf("binary version: %s", version)

	logVerbose("resolve tcp address %s:%d", cfg.host, cfg.port)
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", cfg.host, cfg.port))
	if err != nil {
		logError("failed to resolve TCP address: %s", err.Error())
		os.Exit(1)
	}

	//establish connection
	logVerbose("establish %s connection to %s:%d\n", "tcp", cfg.host, cfg.port)
	c, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		logError("error dialing server: %s", err.Error())
		os.Exit(1)
	}
	logVerbose("connection established")
	defer c.Close()

	logVerbose("construct echo message")
	echo := mastercard.Echo()

	echoReq, err := mastercard.NewRequest(echo).Packet()
	if err != nil {
		logError("error constructing message: %w", err.Error())
		os.Exit(1)
	}

	//send echo message
	logVerbose("send echo message")
	_, err = c.Write(echoReq)
	if err != nil {
		logError("error writing message: %s", err.Error())
		os.Exit(1)
	}

	//receive an echo response
	logVerbose("waiting for response")
	echoBuffer := make([]byte, 1024)
	mLen, err := c.Read(echoBuffer)
	if err != nil {
		logError("error reading: %s", err.Error())
		os.Exit(1)
	}

	received := mastercard.NewResponse(echoBuffer, nil)
	mastercardResponse := received.(mastercard.Response)

	logDebug("echo response: %+v", mastercardResponse.Message())
	logDebug("echo response received: %s", hex.EncodeToString(echoBuffer[:mLen]))

	os.Exit(0)
}

type Logger func(format string, a ...any)

func PrintLogger(format string, a ...any) {
	log.Printf(format, a...)
}

func NopLogger(_ string, _ ...any) {}
