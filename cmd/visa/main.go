package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gitlab.cmpayments.local/creditcard/authorization/internal/data"
	"gitlab.cmpayments.local/creditcard/authorization/internal/entity"
	"gitlab.cmpayments.local/creditcard/authorization/internal/processing/scheme/visa"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/iso8583"
	"gitlab.cmpayments.local/creditcard/authorization/pkg/visa/base1"
	"gitlab.cmpayments.local/creditcard/platform/currencycode"
)

type config struct {
	host string
	port int
}

var (
	loglevel                         string
	version                          = "0.0.1"
	cfg                              config
	signon, echo, authorize, signoff bool
	sourceId                         string
	logError, logVerbose, logDebug   Logger = NopLogger, NopLogger, NopLogger
)

func main() {
	// the default address is the HA proxy address.
	// address of the actual EAS and its ports is documented on https://cmcom.atlassian.net/wiki/spaces/CA/pages/951287809/Visa+authorization
	//flag.StringVar(&cfg.host, "host", "10.18.27.52", "EAS server host")
	flag.StringVar(&cfg.host, "host", "mip.test.cmpayments.local", "EAS server host")
	flag.IntVar(&cfg.port, "port", 10106, "EAS server port")
	flag.StringVar(&loglevel, "log", "vvv", "log level:\n\tv = only errors are displayed,\n\tvv = opening and closing connection information text is displayed,\n\tvvv = response result is displayed")
	flag.BoolVar(&signon, "signon", false, "If a Signon should be sent")
	flag.BoolVar(&echo, "echo", true, "If an Echo should be sent")
	flag.BoolVar(&authorize, "authorize", false, "If a Authorization should be sent")
	flag.BoolVar(&signoff, "signoff", false, "If a Signoff should be sent")
	flag.StringVar(&sourceId, "sourceid", "111106", "VISA Source Station ID")

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

	//resolve tcp address
	logVerbose("resolve tcp address %s:%d\n", cfg.host, cfg.port)
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

	//set connection properties
	err = c.SetKeepAlive(true)
	if err != nil {
		logError("failed to set keepalive: %s", err.Error())
		os.Exit(1)
	}
	err = c.SetKeepAlivePeriod(30 * time.Second)
	if err != nil {
		logError("failed to set keepalive period: %s", err.Error())
		os.Exit(1)
	}
	err = c.SetDeadline(time.Now().Add(15 * time.Second))
	if err != nil {
		logError("failed to set read/write deadline: %s", err.Error())
		os.Exit(1)
	}

	//time.Sleep(10 * time.Second)

	if signon {
		//construct sign-on message
		logVerbose("construct sign-on message")
		signOnMsg := &visa.Message{
			SourceStationID: sourceId,
			Mti:             iso8583.NewMti("0800"),
			Fields: base1.Fields{
				F007_TransmissionDateTime:             base1.F007FromTime(time.Now()),
				F011_SystemTraceAuditNumber:           "100001",
				F070_NetworkManagementInformationCode: "071",
			},
		}

		signOnReq, err := visa.NewRequest(signOnMsg).Packet()
		if err != nil {
			logError("error constructing message: %s", err.Error())
			os.Exit(1)
		}

		logVerbose("sign-on message: %x", signOnReq)

		//send sign-on message
		logVerbose("send sign-on message")
		_, err = c.Write(signOnReq)
		if err != nil {
			logError("error writing message: %s", err.Error())
			os.Exit(1)
		}

		//receive a response
		logVerbose("waiting for sign-on response")
		buffer := make([]byte, 1024)
		mLen, err := c.Read(buffer)
		if err != nil {
			logError("error reading: %s", err.Error())
			os.Exit(1)
		}

		received := visa.NewResponse(buffer, nil)
		visaMsg := received.(visa.Response)

		logDebug("sign-on response: %+v", visaMsg.Message())
		logDebug("sign-on response received: %s", hex.EncodeToString(buffer[:mLen]))
	}

	if echo {
		//construct echo message
		logVerbose("construct echo message")
		echo := visa.Echo(sourceId)

		echoReq, err := visa.NewRequest(echo).Packet()
		if err != nil {
			logError("error constructing message: %s", err.Error())
			os.Exit(1)
		}

		logVerbose("echo message: %x", echoReq)

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

		received := visa.NewResponse(echoBuffer, nil)
		visaMsg := received.(visa.Response)

		logDebug("echo response: %+v", visaMsg.Message())
		logDebug("echo response received: %s", hex.EncodeToString(echoBuffer[:mLen]))
	}

	if authorize {
		hours, minutes, _ := time.Now().Clock()
		stanStr := fmt.Sprintf("00%02d%02d", hours, minutes)
		stan, _ := strconv.Atoi(stanStr)
		a := &entity.Authorization{
			ID:                       uuid.UUID{},
			LogID:                    uuid.UUID{},
			Amount:                   123,
			Currency:                 currencycode.Must(currencycode.EUR),
			CustomerReference:        "1234",
			Source:                   "moto",
			LocalTransactionDateTime: data.LocalTransactionDateTime(time.Now()),
			Status:                   "",
			Stan:                     stan,
			InstitutionID:            fmt.Sprintf("%06d", 20814),
			ProcessingDate:           time.Now(),
			CreatedAt:                time.Now(),
			Card: entity.Card{
				Number: "4761340000000035",
				Holder: "Test",
				Expiry: entity.Expiry{
					Year:  "23",
					Month: "12",
				},
			},
			CardAcceptor: entity.CardAcceptor{
				CategoryCode: "0742",
				ID:           "1234",
				Name:         "VTS Test",
				Address: entity.CardAcceptorAddress{
					PostalCode:  "4899AL",
					City:        "Breda",
					CountryCode: "NLD",
				},
			},
		}

		visa.AuthorizationSchemeData(a)

		authMsg, err := visa.MessageFromAuthorization(*a)
		if err != nil {
			logError("error defining visa message: %w", err.Error())
			os.Exit(1)
		}
		authMsg.SourceStationID = sourceId
		req := visa.NewRequest(authMsg)

		packet, err := req.Packet()
		if err != nil {
			logError("error creating auth packet: %s", err.Error())
			os.Exit(1)
		}
		_, err = c.Write(packet)
		if err != nil {
			logError("error writing auth message: %s", err.Error())
			os.Exit(1)
		}

		//receive a response
		logVerbose("waiting for auth response")
		authBuff := make([]byte, 1024)
		mLen, err := c.Read(authBuff)
		if err != nil {
			logError("error auth reading: %s", err.Error())
			os.Exit(1)
		}

		logVerbose("auth response received: %s", hex.EncodeToString(authBuff[:mLen]))

		received := visa.NewResponse(authBuff, nil)
		visaMsg := received.(visa.Response)
		logDebug("visa message F062.SF1: %#v", visaMsg.Message().Fields.F062_CustomPaymentServiceFields.SF1_AuthorizationCharacteristicsIndicator)
		logDebug("visa message F062.SF2: %#v", visaMsg.Message().Fields.F062_CustomPaymentServiceFields.SF2_TransactionIdentifier)
		if visaMsg.Error() != nil {
			logError("visa error: %s", visaMsg.Error().Error())
		}
	}

	if signoff {
		//construct sign-off message
		logVerbose("construct sign-off message")
		signOffMsg := &visa.Message{
			SourceStationID: sourceId,
			Mti:             iso8583.NewMti("0800"),
			Fields: base1.Fields{
				F007_TransmissionDateTime:             base1.F007FromTime(time.Now()),
				F011_SystemTraceAuditNumber:           "100001",
				F070_NetworkManagementInformationCode: "072",
			},
		}

		signOffReq, err := visa.NewRequest(signOffMsg).Packet()
		if err != nil {
			logError("error constructing message: %s", err.Error())
			os.Exit(1)
		}

		//send sign-off message
		logVerbose("send sign-off message")
		_, err = c.Write(signOffReq)
		if err != nil {
			logError("error writing message: %s", err.Error())
			os.Exit(1)
		}

		//receive a response
		logVerbose("waiting for response")
		signOffBuf := make([]byte, 1024)
		mLen, err := c.Read(signOffBuf)
		if err != nil {
			logError("error reading: %s", err.Error())
			os.Exit(1)
		}

		received := visa.NewResponse(signOffBuf, nil)
		visaMsg := received.(visa.Response)

		logVerbose("sign-off response: %+v", visaMsg.Message())
		logVerbose("sign-off response received: %s", hex.EncodeToString(signOffBuf[:mLen]))
	}

	os.Exit(0)
}

type Logger func(format string, a ...any)

func PrintLogger(format string, a ...any) {
	log.Printf(format, a...)
}

func NopLogger(_ string, _ ...any) {}
