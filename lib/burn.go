package trunks

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
)

func init() {
	encoding.RegisterCodec(codecIgnoreResp{})
}

// Burner is who runs the stress test against Gtargets.
type Burner struct {
	pool           *pool
	ctx            context.Context
	stopch         chan struct{}
	workers        uint64
	maxWorkers     uint64
	numConnPerHost uint64
	maxRecv        int
	maxSend        int
	loop           bool
	dump           bool
	dumpFile       string
	seqmu          sync.Mutex
	seq            uint64
	began          time.Time
}

const (
	DefaultGrpcConnections = 10000
)

// NewBurner creates a new Burner with burner options.
func NewBurner(hosts []string, opts ...BurnOpt) (*Burner, error) {
	if len(hosts) < 1 {
		return nil, fmt.Errorf("no host")
	}

	b := &Burner{
		ctx:            context.Background(),
		stopch:         make(chan struct{}),
		workers:        DefaultWorkers,
		maxWorkers:     DefaultMaxWorkers,
		numConnPerHost: DefaultGrpcConnections,
		began:          time.Now(),
	}

	for _, opt := range opts {
		opt(b)
	}

	var cos []grpc.CallOption
	if b.maxRecv > 0 {
		cos = append(cos, grpc.MaxCallRecvMsgSize(b.maxRecv))
	}
	if b.maxSend > 0 {
		cos = append(cos, grpc.MaxCallSendMsgSize(b.maxSend))
	}

	p, err := newPool(hosts, b.numConnPerHost, cos)
	if err != nil {
		return nil, err
	}

	b.pool = p
	return b, nil
}

// WaitDumpDone waits until all response dumpings are done
// and write everything into the dump file
func (b *Burner) WaitDumpDone() error {
	// dummy-proof
	if !b.dump {
		return nil
	}

	dumpers.Wait()

	f, err := os.OpenFile(b.dumpFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(memBuf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (b *Burner) Close() error {
	return b.pool.close()
}

func (b *Burner) Burn(tgt *Gtarget, p Pacer, du time.Duration, name string) <-chan *Result {
	var wg sync.WaitGroup

	workers := b.workers
	if workers > b.maxWorkers {
		workers = b.maxWorkers
	}

	results := make(chan *Result)
	ticks := make(chan struct{})
	for i := uint64(0); i < workers; i++ {
		wg.Add(1)
		go b.burn(tgt, name, &wg, ticks, results)
	}

	go func() {
		defer close(results)
		defer wg.Wait()
		defer close(ticks)

		began, count := time.Now(), uint64(0)
		for {
			elapsed := time.Since(began)
			if du > 0 && elapsed > du {
				return
			}

			wait, stop := p.Pace(elapsed, count)
			if stop {
				return
			}

			time.Sleep(wait)

			if workers < b.maxWorkers {
				select {
				case ticks <- struct{}{}:
					count++
					continue
				case <-b.stopch:
					return
				default:
					// all workers are blocked. start one more and try again
					workers++
					wg.Add(1)
					go b.burn(tgt, name, &wg, ticks, results)
				}
			}

			select {
			case ticks <- struct{}{}:
				count++
			case <-b.stopch:
				return
			}
		}
	}()

	return results
}

func (b *Burner) Stop() {
	select {
	case <-b.stopch:
		return
	default:
		close(b.stopch)
	}
}

func (b *Burner) burn(tgt *Gtarget, name string, workers *sync.WaitGroup, ticks <-chan struct{}, results chan<- *Result) {
	defer workers.Done()
	for range ticks {
		results <- b.hit(tgt, name)
	}
}

func (b *Burner) hit(tgt *Gtarget, name string) *Result {
	var (
		res = Result{Attack: name}
		err error
	)

	b.seqmu.Lock()
	res.Timestamp = b.began.Add(time.Since(b.began))
	res.Seq = b.seq
	b.seq++
	b.seqmu.Unlock()

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}

		res.Latency = time.Since(res.Timestamp)

		// count failure by cheating metrics
		// as if this is an HTTP call
		if err != nil {
			if res.Code < 100 {
				// if res.code is not set
				// consider it the fault from caller side
				res.Code = 400
			}
			res.Error = err.Error()
		} else {
			res.Code = 200
		}
	}()

	c, err := b.pool.pick()
	if err != nil {
		b.Stop()
		return &res
	}

	res.Method = tgt.MethodName
	// TODO
	// res.URL = tgt.URL

	req, err := tgt.pick(b.loop)
	if err != nil {
		b.Stop()
		return &res
	}

	// clone a response object to avoid race-condition
	resp := reflect.New(reflect.TypeOf(tgt.Response).Elem()).Interface().(proto.Message)
	err = c.Invoke(b.ctx, tgt.MethodName, req, resp)
	if err != nil {
		return &res
	}

	res.Body, err = proto.Marshal(resp)
	if err != nil {
		return &res
	}

	if b.dump {
		// start a new routine to serialize and dump
		// to avoid blocking processing
		dumpers.Add(1)
		go func(resp proto.Message) {
			defer dumpers.Done()
			bytes, err := json.Marshal(resp)
			if err == nil {
				memBufMutex.Lock()
				defer memBufMutex.Unlock()
				bytes = append(bytes, []byte{'\r', '\n'}...)
				memBuf.Write(bytes)
			}
		}(resp)
	}

	// TODO optimize performance
	// TODO optimize performance
	res.BytesOut = uint64(proto.Size(req))
	res.BytesIn = uint64(proto.Size(resp))

	return &res
}
