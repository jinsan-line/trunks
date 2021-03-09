package trunks

import (
	"context"

	"google.golang.org/grpc/metadata"
)

// BurnOpt ...
type BurnOpt func(*Burner)

// WithMetadata ...
func WithMetadata(md metadata.MD) BurnOpt {
	return func(b *Burner) {
		b.ctx = metadata.NewOutgoingContext(context.Background(), md)
	}
}

// WithWorkers ...
func WithWorkers(num uint64) BurnOpt {
	return func(b *Burner) { b.workers = num }
}

// WithMaxWorkers ...
func WithMaxWorkers(num uint64) BurnOpt {
	return func(b *Burner) { b.maxWorkers = num }
}

// WithNumConnPerHost ...
func WithNumConnPerHost(num uint64) BurnOpt {
	return func(b *Burner) { b.numConnPerHost = num }
}

// WithMaxRecvSize ...
func WithMaxRecvSize(s int) BurnOpt {
	return func(b *Burner) {
		b.maxRecv = s
	}
}

// WithMaxSendSize ...
func WithMaxSendSize(s int) BurnOpt {
	return func(b *Burner) {
		b.maxSend = s
	}
}

// WithLooping ...
func WithLooping(yesno bool) BurnOpt {
	return func(b *Burner) { b.loop = yesno }
}

// WithLoop : Deprecating; use WithLooping(bool)
func WithLoop() BurnOpt {
	return func(b *Burner) { b.loop = true }
}

// WithDumpFile ...
func WithDumpFile(fileName string) BurnOpt {
	return func(b *Burner) {
		b.dumpFile = fileName
		if fileName != "" {
			b.dump = true
		}
	}
}
