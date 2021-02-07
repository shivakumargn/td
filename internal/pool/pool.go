// Package pool contains Telegram connections pool.
package pool

import (
	"context"
	"sync"

	"go.uber.org/atomic"
	"go.uber.org/zap"
	"golang.org/x/xerrors"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/internal/tdsync"
)

// DCOptions is a Telegram data center connections pool options.
type DCOptions struct {
	// Logger is instance of zap.Logger. No logs by default.
	Logger *zap.Logger
	// MTProto options for connections.
	// Opened connection limit to the DC.
	MaxOpenConnections int64
}

// DC represents connection pool to one data center.
type DC struct {
	id int

	// Connection constructor.
	newConn func() Conn

	// Wrappers for external world, like logs or PRNG.
	log *zap.Logger // immutable

	// DC context. Will be canceled by Run on exit.
	ctx    context.Context    // immutable
	cancel context.CancelFunc // immutable

	// Connections supervisor.
	grp *tdsync.Supervisor
	// Free connections.
	free []*poolConn
	// Total connections.
	total int64
	// Connection id monotonic counter.
	nextConn atomic.Int64
	freeReq  *reqMap
	// DC mutex.
	mu sync.Mutex

	// Limit of connections.
	max int64 // immutable

	// Signal connection for cases when all connections are dead, but all requests waiting for
	// free connection in 3rd acquire case.
	stuck *tdsync.ResetReady

	// Requests wait group.
	ongoing sync.WaitGroup

	closed atomic.Bool
}

// NewDC creates new uninitialized DC.
func NewDC(ctx context.Context, id int, newConn func() Conn, opts DCOptions) *DC {
	ctx, cancel := context.WithCancel(ctx)

	return &DC{
		id:      id,
		newConn: newConn,
		log:     opts.Logger,
		ctx:     ctx,
		cancel:  cancel,
		grp:     tdsync.NewSupervisor(ctx),
		freeReq: newReqMap(),
		max:     opts.MaxOpenConnections,
		stuck:   tdsync.NewResetReady(),
	}
}

func (c *DC) createConnection(id int64) *poolConn {
	conn := &poolConn{
		Conn:    c.newConn(),
		id:      id,
		dc:      c,
		deleted: atomic.NewBool(false),
		dead:    tdsync.NewReady(),
	}

	c.grp.Go(func(groupCtx context.Context) (err error) {
		defer c.dead(conn, err)
		return conn.Run(groupCtx)
	})

	return conn
}

func (c *DC) dead(r *poolConn, deadErr error) {
	if r.deleted.Swap(true) {
		return // Already deleted.
	}

	c.stuck.Reset()
	c.mu.Lock()
	defer c.mu.Unlock()
	r.dead.Signal()

	c.total--
	remaining := c.total
	if remaining < 0 {
		panic("unreachable: remaining can'be less than zero")
	}

	idx := -1
	for i, conn := range c.free {
		// Search connection by pointer.
		if conn.id == r.id {
			idx = i
		}
	}

	if idx >= 0 {
		// Delete by index from slice tricks.
		copy(c.free[idx:], c.free[idx+1:])
		// Delete reference to prevent resource leaking.
		c.free[len(c.free)-1] = nil
		c.free = c.free[:len(c.free)-1]
	}

	c.log.Debug("Connection died",
		zap.Int64("remaining", remaining),
		zap.Int64("conn_id", r.id),
		zap.Error(deadErr),
	)
}

func (c *DC) pop() (r *poolConn, ok bool) {
	l := len(c.free)
	if l > 0 {
		r, c.free = c.free[l-1], c.free[:l-1]

		return r, true
	}

	return
}

func (c *DC) release(r *poolConn) {
	if r == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.freeReq.transfer(r) {
		c.log.Debug("Transfer connection to requester", zap.Int64("conn_id", r.id))
		return
	}
	c.log.Debug("Connection released", zap.Int64("conn_id", r.id))
	c.free = append(c.free, r)
}

var errDCIsClosed = xerrors.New("DC is closed")

func (c *DC) acquire(ctx context.Context) (r *poolConn, err error) { // nolint:gocyclo
retry:
	c.mu.Lock()
	// 1st case: have free connections.
	if r, ok := c.pop(); ok {
		c.mu.Unlock()
		select {
		case <-r.Dead():
			c.dead(r, nil)
			goto retry
		default:
		}
		c.log.Debug("Re-using free connection", zap.Int64("conn_id", r.id))
		return r, nil
	}

	// 2nd case: no free connections, but can create one.
	// c.max < 1 means unlimited
	if c.max < 1 || c.total < c.max {
		c.total++
		c.mu.Unlock()

		id := c.nextConn.Inc()
		c.log.Debug("Creating new connection",
			zap.Int64("conn_id", id),
		)
		conn := c.createConnection(id)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-c.ctx.Done():
			return nil, xerrors.Errorf("DC closed: %w", c.ctx.Err())
		case <-conn.Ready():
			return conn, nil
		case <-conn.Dead():
			c.dead(conn, nil)
			goto retry
		}
	}

	// 3rd case: no free connections, can't create yet one, wait for free.
	key, ch := c.freeReq.request()
	c.mu.Unlock()
	c.log.Debug("Waiting for free connect", zap.Int64("request_id", int64(key)))

	select {
	case conn := <-ch:
		c.log.Debug("Got connection for request",
			zap.Int64("conn_id", conn.id),
			zap.Int64("request_id", int64(key)),
		)
		return conn, nil
	case <-c.stuck.Ready():
		c.log.Debug("Some connection dead, try to create new connection, cancel waiting")

		c.freeReq.delete(key)
		select {
		default:
		case conn, ok := <-ch:
			if ok && conn != nil {
				return conn, nil
			}
		}

		goto retry
	case <-ctx.Done():
		err = ctx.Err()
	case <-c.ctx.Done():
		err = xerrors.Errorf("DC closed: %w", c.ctx.Err())
	}

	// Executed only if at least one of context is Done.
	c.freeReq.delete(key)
	select {
	default:
	case conn, ok := <-ch:
		if ok && conn != nil {
			c.release(conn)
		}
	}

	return nil, err
}

// InvokeRaw sends MTProto request using one of pool connection.
func (c *DC) InvokeRaw(ctx context.Context, input bin.Encoder, output bin.Decoder) error {
	if c.closed.Load() {
		return errDCIsClosed
	}

	c.ongoing.Add(1)
	defer c.ongoing.Done()

	for {
		conn, err := c.acquire(ctx)
		if err != nil {
			if xerrors.Is(err, ErrConnDead) {
				continue
			}
			return xerrors.Errorf("acquire connection: %w", err)
		}

		c.log.Debug("DC Invoke")
		err = conn.InvokeRaw(ctx, input, output)
		if err != nil {
			select {
			case <-conn.Dead():
				continue
			default:
			}
		}
		c.release(conn)
		if err != nil {
			c.log.Debug("DC Invoke failed", zap.Error(err))
			return xerrors.Errorf("invoke pool: %w", err)
		}

		c.log.Debug("DC Invoke complete")
		return err
	}
}

// Close waits while all ongoing requests will be done or until given context is done.
// Then, closes the DC.
func (c *DC) Close(closeCtx context.Context) error {
	if c.closed.Swap(true) {
		return xerrors.New("DC already closed")
	}
	c.log.Debug("Closing DC")
	defer c.log.Debug("DC closed")

	closed, cancel := context.WithCancel(closeCtx)
	go func() {
		c.ongoing.Wait()
		cancel()
	}()

	<-closed.Done()

	c.cancel()
	return c.grp.Wait()
}
