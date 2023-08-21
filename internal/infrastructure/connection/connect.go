package connection

import (
	"context"
	"net"
	"sync"
	"time"

	"gitlab.cmpayments.local/libraries-go/logging"
)

func (c *connection) connect(wg *sync.WaitGroup) {
	defer wg.Done()

	ctx := c.initializeLoggingCtx(context.TODO())

	reconnectTicker := time.NewTicker(c.redialDelay)

	for c.selectConnectorBehavior()(ctx, reconnectTicker) {
	}
}

func (c *connection) selectConnectorBehavior() func(context.Context, *time.Ticker) bool {
	if c.isConnected {
		return c.connectedConnectorBehavior
	}
	return c.disconnectedConnectorBehavior
}

func (c *connection) connectedConnectorBehavior(ctx context.Context, reconnectTicker *time.Ticker) bool {
	select {
	case <-c.shutdownSignal:
		return c.onConnectorShutdown(ctx)
	case <-c.disconnectCmdSignal:
		return c.onConnectorDisconnect(ctx, reconnectTicker)
	}
}

func (c *connection) disconnectedConnectorBehavior(ctx context.Context, reconnectTicker *time.Ticker) bool {
	select {
	case <-c.shutdownSignal:
		return c.onConnectorShutdown(ctx)
	case <-reconnectTicker.C:
		return c.onConnectorReconnectAttempt(ctx, reconnectTicker)
	case <-c.disconnectCmdSignal:
		return c.onConnectorDisconnect(ctx, reconnectTicker)
	}
}

func (c *connection) onConnectorShutdown(ctx context.Context) bool {
	c.logger.Debug(ctx, "CONNECT: shutdown")
	if c.isConnected {
		c.isConnected = false
		c.logger.Debug(ctx, "CONNECT: closing connection")
		err := c.conn.Close()
		if err != nil {
			c.logger.Warning(
				logging.ContextWithError(ctx, err),
				"CONNECT: failed to close connection")
		}
	}
	c.isConnected = false
	return false
}

func (c *connection) onConnectorDisconnect(ctx context.Context, reconnectTicker *time.Ticker) bool {
	c.logger.Debug(ctx, "CONNECT: disconnect command")
	if c.isConnected {
		c.isConnected = false
		c.logger.Debug(ctx, "CONNECT: closing connection")
		err := c.conn.Close()
		if err != nil {
			c.logger.Warning(
				logging.ContextWithError(ctx, err),
				"CONNECT: failed to close connection")
		}
		if !trySignal(c.attendantDisconnectedSignal) {
			c.logger.Debug(ctx, "CONNECT: attender disconnect would block")
		}
		if !trySignal(c.receiverDisconnectedSignal) {
			c.logger.Debug(ctx, "CONNECT: receiver disconnect would block")
		}
		if !trySignal(c.senderDisconnectedSignal) {
			c.logger.Debug(ctx, "CONNECT: sender disconnect would block")
		}

		if reconnectTicker == nil {
			reconnectTicker = time.NewTicker(1 * time.Second) //nolint:staticcheck	// False positive
		} else {
			reconnectTicker.Reset(1 * time.Second)
		}
	}
	return true
}

func (c *connection) onConnectorReconnectAttempt(ctx context.Context, reconnectTicker *time.Ticker) bool {
	time.Sleep(c.redialDelay)
	c.logger.Debug(ctx, "CONNECT: trying to connect")
	var err error
	c.conn, err = net.DialTimeout("tcp", c.address, c.dialTimeout)
	if err != nil {
		c.logger.Error(
			logging.ContextWithError(ctx, err),
			"CONNECT: failed to connect")
		return true
	}
	c.logger.Debug(ctx, "CONNECT: connected")
	reconnectTicker.Stop()
	c.isConnected = true
	if !trySignal(c.receiverConnectedSignal) {
		c.logger.Debug(ctx, "CONNECT: receiver disconnect would block")
	}
	if !trySignal(c.senderConnectedSignal) {
		c.logger.Debug(ctx, "CONNECT: sender disconnect would block")
	}
	if !trySignal(c.attendantConnectedSignal) {
		c.logger.Debug(ctx, "CONNECT: attender disconnect would block")
	}
	return true
}
