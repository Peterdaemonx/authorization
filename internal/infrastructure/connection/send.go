package connection

import (
	"context"
	"encoding/hex"
	"sync"

	"gitlab.cmpayments.local/libraries-go/logging"
)

func (c *connection) send(wg *sync.WaitGroup) {
	defer wg.Done()

	ctx := c.initializeLoggingCtx(context.TODO())

	for c.selectSenderBehavior()(ctx) {
	}
}

func (c *connection) selectSenderBehavior() func(context.Context) bool {
	if c.senderConnected {
		// c.logger.Debug(ctx, "SENDER: selected connected behavior")
		return c.connectedSenderBehavior
	}
	// c.logger.Debug(ctx, "SENDER: selected disconnected behavior")
	return c.disconnectedSenderBehavior
}

func (c *connection) connectedSenderBehavior(ctx context.Context) bool {
	select {
	case <-c.shutdownSignal:
		return c.onSenderShutdown(ctx)
	case <-c.senderDisconnectedSignal:
		return c.onSenderDisconnected(ctx)
	case sendCmd := <-c.sendCmdSignal:
		return c.onSenderCommand(ctx, sendCmd)
	}
}

func (c *connection) disconnectedSenderBehavior(ctx context.Context) bool {
	select {
	case <-c.shutdownSignal:
		return c.onSenderShutdown(ctx)
	case <-c.senderConnectedSignal:
		return c.onSenderConnected(ctx)
	}
}

func (c *connection) onSenderShutdown(ctx context.Context) bool {
	c.logger.Debug(ctx, "SENDER - shutdown")
	return false
}

func (c *connection) onSenderDisconnected(ctx context.Context) bool {
	c.logger.Debug(ctx, "SENDER - disconnected")
	c.senderConnected = false
	return true
}

func (c *connection) onSenderCommand(ctx context.Context, cmd sendCommand) bool {
	c.logger.Debug(ctx, "SENDER - command")
	loggingCtx := logging.ContextWithValue(cmd.context, packetLoggingKey, hex.EncodeToString(cmd.packet))
	c.logger.Debug(loggingCtx, "SENDER - send packet")
	err := writeWithTimeout(c.conn, c.writeTimeout, cmd.packet)
	if err != nil {
		c.logger.Error(
			logging.ContextWithError(loggingCtx, err),
			"SENDER - failed to send packet")
		if cmd.notifyError {
			if !c.tryNotifyError(err) {
				c.logger.Warning(ctx, "SENDER - notify error would block")
			}
		}
		if !c.tryDisconnect() {
			c.logger.Warning(ctx, "SENDER - notify disconnection would block")
		}
		c.senderConnected = false
	}
	return true
}

func (c *connection) onSenderConnected(ctx context.Context) bool {
	c.logger.Debug(ctx, "SENDER - connected")
	c.senderConnected = true
	return true
}
