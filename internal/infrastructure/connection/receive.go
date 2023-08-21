package connection

import (
	"context"
	"encoding/binary"
	"errors"
	"net"
	"sync"

	"gitlab.cmpayments.local/libraries-go/logging"
)

func (c *connection) receive(wg *sync.WaitGroup) {
	defer wg.Done()

	ctx := c.initializeLoggingCtx(context.TODO())

	for c.selectReceiverBehavior()(ctx) {
	}
}

func (c *connection) selectReceiverBehavior() func(context.Context) bool {
	if c.receiverConnected {
		// c.logger.Debug(ctx, "RECEIVER: selected connected behavior")
		return c.connectedReceiverBehavior
	}
	// c.logger.Debug(ctx, "RECEIVER: selected disconnected behavior")
	return c.disconnectedReceiverBehavior
}

func (c *connection) connectedReceiverBehavior(ctx context.Context) bool {
	select {
	case <-c.shutdownSignal:
		return c.onReceiverShutdown(ctx)
	case <-c.receiverDisconnectedSignal:
		return c.onReceiverDisconnect(ctx)
	default:
		return c.onReceiverConnectedIdle(ctx)
	}
}

func (c *connection) disconnectedReceiverBehavior(ctx context.Context) bool {
	select {
	case <-c.shutdownSignal:
		return c.onReceiverShutdown(ctx)
	case <-c.receiverConnectedSignal:
		return c.onReceiverConnected(ctx)
	}
}

func (c *connection) onReceiverShutdown(ctx context.Context) bool {
	c.logger.Debug(ctx, "RECEIVER: shutdown")
	return false
}

func (c *connection) onReceiverDisconnect(ctx context.Context) bool {
	c.logger.Debug(ctx, "RECEIVER: disconnected")
	c.receiverConnected = false
	return true
}

func (c *connection) onReceiverConnected(ctx context.Context) bool {
	c.logger.Debug(ctx, "RECEIVER: connected")
	c.receiverConnected = true
	return true
}

func (c *connection) onReceiverConnectedIdle(ctx context.Context) bool {
	lengthBuffer := make([]byte, 2)
	err := readWithTimeout(c.conn, c.keepAliveDelay, lengthBuffer)
	if err != nil {
		var netError net.Error
		switch {
		case errors.As(err, &netError) && netError.Timeout():
			c.logger.Debug(ctx, "RECEIVER - sending keep alive")
			if !c.trySend(sendCommand{context: ctx, packet: []byte{0x00, 0x00}, notifyError: false}) {
				c.logger.Warning(ctx, "RECEIVER - try send would block")
			}
		default:
			c.logger.Error(
				logging.ContextWithError(ctx, err),
				"RECEIVER - failed to receive packet length")
			if !c.tryNotifyError(err) {
				c.logger.Warning(ctx, "RECEIVER - notify error would block")
			}
			if !c.tryDisconnect() {
				c.logger.Warning(ctx, "RECEIVER - notify disconnection would block")
			}
			c.receiverConnected = false
		}
		return true
	}

	length := int(binary.BigEndian.Uint16(lengthBuffer))
	if length == 0 {
		c.logger.Debug(ctx, "RECEIVER - received zero probe, continuing")
		return true
	}

	// Visa Message Length Header (VMLH) is actually two bytes longer than what they say in the message, for unknown reasons
	length += c.msgLengthSurplus

	payload := make([]byte, length)
	err = readWithTimeout(c.conn, c.readTimeout, payload)
	if err != nil {
		c.logger.Error(
			logging.ContextWithError(ctx, err),
			"RECEIVER - failed to receive packet payload")
		c.logger.Debug(ctx, "RECEIVER - exiting")
		if !c.tryNotifyError(err) {
			c.logger.Warning(ctx, "RECEIVER - notify error would block")
		}
		if !c.tryDisconnect() {
			c.logger.Warning(ctx, "RECEIVER - notify disconnection would block")
		}
		return true
	}
	c.logger.Debug(ctx, "RECEIVER - packet received")

	if !c.tryNotifyReceived(append(lengthBuffer, payload...)) {
		c.logger.Warning(ctx, "RECEIVER - notify receive would block")
	}

	return true
}
