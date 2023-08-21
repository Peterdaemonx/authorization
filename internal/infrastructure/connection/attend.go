package connection

import (
	"context"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"gitlab.cmpayments.local/libraries-go/logging"
)

func (c *connection) attend(wg *sync.WaitGroup) {
	defer wg.Done()

	ctx := c.initializeLoggingCtx(context.TODO())

	for c.selectAttendantBehavior()(ctx) {
	}
}

func (c *connection) selectAttendantBehavior() func(context.Context) bool {
	if c.attendantConnected {
		if c.pendingJob != nil {
			return c.connectedBusyAttendantBehavior
		} else {
			return c.connectedIdleAttendantBehavior
		}
	}
	return c.disconnectedAttendantBehavior
}

func (c *connection) connectedIdleAttendantBehavior(ctx context.Context) bool {
	select {
	case <-c.shutdownSignal:
		return c.onAttendantShutdown(ctx)
	case <-c.attendantDisconnectedSignal:
		return c.onAttendantDisconnected(ctx)
	case job := <-c.jobSignal:
		return c.onAttendantJob(ctx, job)
	case packet := <-c.receivedSignal:
		return c.onAttendantReceived(ctx, packet)
	}
}

func (c *connection) connectedBusyAttendantBehavior(ctx context.Context) bool {
	select {
	case <-c.shutdownSignal:
		return c.onAttendantShutdown(ctx)
	case <-c.attendantDisconnectedSignal:
		return c.onAttendantDisconnected(ctx)
	case err := <-c.errorSignal:
		return c.onAttendantError(ctx, err)
	case packet := <-c.receivedSignal:
		return c.onAttendantReceived(ctx, packet)
	case <-c.pendingJob.context.Done():
		return c.onAttendantJobCanceled()
	case <-time.After(c.responseTimeout):
		return c.onAttendantRequestTimeout()
	}
}

func (c *connection) disconnectedAttendantBehavior(ctx context.Context) bool {
	select {
	case <-c.shutdownSignal:
		return c.onAttendantShutdown(ctx)
	case <-c.attendantConnectedSignal:
		return c.onAttendantConnected(ctx)
	}
}

func (c *connection) onAttendantShutdown(ctx context.Context) bool {
	c.logger.Debug(ctx, "ATTEND: shutdown")
	c.finishPendingJob(c.receivedFactory(nil, ErrRequestAbandoned))
	return false
}

func (c *connection) onAttendantError(ctx context.Context, err error) bool {
	var loggingCtx context.Context
	if c.pendingJob != nil {
		loggingCtx = c.initializeLoggingCtx(c.pendingJob.context)
	} else {
		loggingCtx = ctx
	}

	c.logger.Debug(
		logging.ContextWithError(loggingCtx, err),
		"ATTEND: error")
	c.finishPendingJob(c.receivedFactory(nil, err))
	return true
}

func (c *connection) onAttendantConnected(ctx context.Context) bool {
	c.logger.Debug(ctx, "ATTEND: connected")
	c.attendantConnected = true
	return true
}

func (c *connection) onAttendantDisconnected(ctx context.Context) bool {
	c.logger.Debug(ctx, "ATTEND: disconnected")
	c.attendantConnected = false
	c.finishPendingJob(c.receivedFactory(nil, ErrRequestAbandoned))
	return true
}

func (c *connection) finishPendingJob(received Received) {
	if c.pendingJob == nil {
		return
	}

	c.pendingJob.responseSignal <- received
	close(c.pendingJob.responseSignal)
	c.pendingJob = nil
}

func (c *connection) onAttendantJob(ctx context.Context, job *job) bool {
	c.logger.Debug(ctx, "ATTEND: job received")

	loggingCtx := c.initializeLoggingCtx(job.context)

	if c.pendingJob != nil {
		c.logger.Warning(loggingCtx, "ATTENDER: pending job when received another")
		c.finishPendingJob(c.receivedFactory(nil, ErrRequestAbandoned))
	}
	c.pendingJob = job

	packet, err := c.pendingJob.request.Packet()
	if err != nil {
		c.logger.Error(
			logging.ContextWithError(loggingCtx, err),
			"ATTENDER: failed to prepare packet")
		c.finishPendingJob(c.receivedFactory(nil, err))
		return true
	}

	if !c.trySend(sendCommand{context: loggingCtx, packet: packet, notifyError: true}) {
		c.logger.Warning(loggingCtx, "ATTENDER: try send would block")
	}

	return true
}

func (c *connection) onAttendantReceived(ctx context.Context, packet []byte) bool {
	var loggingCtx context.Context
	if c.pendingJob != nil {
		loggingCtx = c.initializeLoggingCtx(c.pendingJob.context)
	} else {
		loggingCtx = ctx
	}
	loggingCtx = logging.ContextWithValue(loggingCtx, packetLoggingKey, hex.EncodeToString(packet))
	c.logger.Debug(loggingCtx, "ATTEND: received packet")

	received := c.receivedFactory(packet, nil)
	if received.Error() != nil {
		c.logger.Error(
			logging.ContextWithError(loggingCtx, received.Error()),
			"ATTEND: failed to translate response payload")
		c.finishPendingJob(received)
		return true
	}

	loggingCtx = logging.ContextWithValue(loggingCtx, "message", fmt.Sprintf("%#v", received))
	c.logger.Debug(loggingCtx, "ATTEND: received message")

	if c.pendingJob != nil && received.IsRequestResponse(c.pendingJob.request) {
		c.logger.Debug(loggingCtx, "ATTEND: send response")
		c.finishPendingJob(received)
		return true
	}

	c.logger.Warning(loggingCtx, "ATTEND: received network request")

	packet, err := received.PacketToSend()
	if err != nil {
		c.logger.Error(
			logging.ContextWithError(loggingCtx, err),
			"ATTEND: failed to get packet to send")
		return true
	}

	if !c.trySend(sendCommand{context: loggingCtx, packet: packet, notifyError: false}) {
		c.logger.Warning(loggingCtx, "ATTEND: try send would block")
	}

	return true
}

func (c *connection) onAttendantJobCanceled() bool {
	loggingCtx := c.initializeLoggingCtx(c.pendingJob.context)
	c.logger.Debug(loggingCtx, "ATTEND: request canceled")
	c.finishPendingJob(c.receivedFactory(nil, ErrRequestCanceled))
	return true
}

func (c *connection) onAttendantRequestTimeout() bool {
	loggingCtx := c.initializeLoggingCtx(c.pendingJob.context)
	c.logger.Debug(loggingCtx, "ATTEND: request timeout")
	c.finishPendingJob(c.receivedFactory(nil, ErrRequestTimeout))
	return true
}
