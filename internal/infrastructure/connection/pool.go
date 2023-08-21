package connection

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gitlab.cmpayments.local/creditcard/platform"
)

type Request interface {
	Packet() ([]byte, error)
}

type Received interface {
	IsRequestResponse(request Request) bool
	PacketToSend() ([]byte, error)
	Error() error
}

type job struct {
	request        Request
	responseSignal chan Received
	context        context.Context
}

type ReceivedFactoryFunc func([]byte, error) Received

type Pool struct {
	connections     []*connection
	tickDelay       time.Duration
	stopWG          sync.WaitGroup
	shutdownSignal  chan struct{}
	jobSignal       chan *job
	receivedFactory ReceivedFactoryFunc
	logger          platform.Logger
}

func NewPool(cfg PoolConfiguration, receivedFactory ReceivedFactoryFunc, logger platform.Logger, msgLengthSurplus int) *Pool {
	tickDelay := cfg.TickDelay
	if tickDelay == 0 {
		tickDelay = time.Minute
	}

	pool := Pool{
		tickDelay:       tickDelay,
		jobSignal:       make(chan *job),
		shutdownSignal:  make(chan struct{}),
		receivedFactory: receivedFactory,
		logger:          logger,
	}

	var connections []*connection
	for i := 0; i < cfg.MaxConnections; i++ {
		connection := newConnection(fmt.Sprintf("%d", i), cfg.Name, cfg.Address, cfg.DialTimeout, cfg.RedialDelay, cfg.KeepAliveDelay, cfg.ReadTimeout, cfg.WriteTimeout, cfg.ResponseTimeout, pool.shutdownSignal, pool.jobSignal, receivedFactory, pool.logger, msgLengthSurplus)

		connections = append(connections, connection)
	}

	pool.connections = connections

	return &pool
}

func (p *Pool) Start() {
	p.stopWG.Add(len(p.connections))

	for _, connection := range p.connections {
		connection.start(&p.stopWG)
	}
}

func (p *Pool) Send(ctx context.Context, request Request) Received {
	job := job{
		request:        request,
		responseSignal: make(chan Received),
		context:        ctx,
	}
	ticker := time.NewTicker(p.tickDelay)
	defer ticker.Stop()

	select {
	case p.jobSignal <- &job:
		response := <-job.responseSignal
		return response
	case <-ticker.C:
		close(job.responseSignal)
		return p.receivedFactory(nil, ErrNoFreeConnections)
	}
}

func (p *Pool) Stop() {
	close(p.shutdownSignal)
	p.stopWG.Wait()
}
