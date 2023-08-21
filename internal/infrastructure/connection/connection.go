package connection

import (
	"context"
	"net"
	"sync"
	"time"

	"gitlab.cmpayments.local/libraries-go/logging"

	"gitlab.cmpayments.local/creditcard/platform"
)

const (
	connectionIdLoggingKey      = "connection-id"
	poolNameLoggingKey          = "pool-name"
	connectionAddressLoggingKey = "address"
	packetLoggingKey            = "packet"
)

type sendCommand struct {
	context     context.Context
	packet      []byte
	notifyError bool
}

type connection struct {
	address                     string
	dialTimeout                 time.Duration
	redialDelay                 time.Duration
	keepAliveDelay              time.Duration
	readTimeout                 time.Duration
	writeTimeout                time.Duration
	responseTimeout             time.Duration
	id                          string
	name                        string
	conn                        net.Conn
	shutdownSignal              chan struct{}
	receiverDisconnectedSignal  chan struct{}
	receiverConnectedSignal     chan struct{}
	attendantDisconnectedSignal chan struct{}
	attendantConnectedSignal    chan struct{}
	senderDisconnectedSignal    chan struct{}
	senderConnectedSignal       chan struct{}
	sendCmdSignal               chan sendCommand
	disconnectCmdSignal         chan struct{}
	errorSignal                 chan error
	jobSignal                   chan *job
	receivedSignal              chan []byte
	isConnected                 bool
	receivedFactory             ReceivedFactoryFunc
	logger                      platform.Logger

	attendantConnected bool
	pendingJob         *job

	receiverConnected bool

	senderConnected bool

	msgLengthSurplus int
}

func newConnection(
	id, name, address string,
	dialTimeout, redialDelay, keepAliveDelay, readTimeout, writeTimeout, responseTimeout time.Duration,
	shutdown chan struct{}, jobSignal chan *job, receivedFactory ReceivedFactoryFunc, logger platform.Logger,
	msgLengthSurplus int) *connection {
	return &connection{
		id:                          id,
		name:                        name,
		address:                     address,
		dialTimeout:                 dialTimeout,
		redialDelay:                 redialDelay,
		keepAliveDelay:              keepAliveDelay,
		readTimeout:                 readTimeout,
		writeTimeout:                writeTimeout,
		responseTimeout:             responseTimeout,
		errorSignal:                 make(chan error),
		receiverDisconnectedSignal:  make(chan struct{}),
		receiverConnectedSignal:     make(chan struct{}),
		attendantDisconnectedSignal: make(chan struct{}),
		attendantConnectedSignal:    make(chan struct{}),
		senderDisconnectedSignal:    make(chan struct{}),
		senderConnectedSignal:       make(chan struct{}),
		disconnectCmdSignal:         make(chan struct{}),
		shutdownSignal:              shutdown,
		sendCmdSignal:               make(chan sendCommand),
		jobSignal:                   jobSignal,
		receivedSignal:              make(chan []byte),
		receivedFactory:             receivedFactory,
		logger:                      logger,
		msgLengthSurplus:            msgLengthSurplus,
	}
}

func (c *connection) initializeLoggingCtx(ctx context.Context) context.Context {
	loggingCtx := logging.ContextWithValue(ctx, connectionIdLoggingKey, c.id)
	loggingCtx = logging.ContextWithValue(loggingCtx, poolNameLoggingKey, c.name)
	return logging.ContextWithValue(loggingCtx, connectionAddressLoggingKey, c.address)
}

func (c *connection) tryDisconnect() bool {
	return trySignal(c.disconnectCmdSignal)
}

func (c *connection) tryNotifyError(err error) bool {
	select {
	case c.errorSignal <- err:
		return true
	default:
		return false
	}
}

func (c *connection) tryNotifyReceived(payload []byte) bool {
	select {
	case c.receivedSignal <- payload:
		return true
	default:
		return false
	}
}

func (c *connection) trySend(cmd sendCommand) bool {
	select {
	case c.sendCmdSignal <- cmd:
		return true
	default:
		return false
	}
}

func (c *connection) start(stopWG *sync.WaitGroup) {
	go c.run(stopWG)
}

func (c *connection) run(stopWG *sync.WaitGroup) {
	defer stopWG.Done()
	var wg sync.WaitGroup

	wg.Add(4)
	go c.attend(&wg)
	go c.receive(&wg)
	go c.send(&wg)
	go c.connect(&wg)
	<-c.shutdownSignal
	wg.Wait()
}

func trySignal(signal chan struct{}) bool {
	select {
	case signal <- struct{}{}:
		return true
	default:
		return false
	}
}

func readWithTimeout(conn net.Conn, readTimeout time.Duration, buffer []byte) error {
	err := conn.SetReadDeadline(time.Now().Add(readTimeout))
	if err != nil {
		return err
	}
	_, err = conn.Read(buffer)
	return err
}

func writeWithTimeout(conn net.Conn, writeTimeout time.Duration, buffer []byte) error {
	err := conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	if err != nil {
		return err
	}
	_, err = conn.Write(buffer)
	return err
}
