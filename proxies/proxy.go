package proxies

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/rack-app/server/clusters"

	"github.com/rack-app/server/workers"
)

type Proxy struct {
	address string
	cluster clusters.Cluster
	wg      *sync.WaitGroup

	ctx    context.Context
	cancel context.CancelFunc
}

func New(addr string, c clusters.Cluster) *Proxy {
	ctx, cancel := context.WithCancel(context.Background())
	return &Proxy{cluster: c, wg: &sync.WaitGroup{}, ctx: ctx, cancel: cancel, address: addr}
}

func (p *Proxy) ListenAndServe() error {
	l, err := net.Listen("tcp", p.address)

	if err != nil {
		return err
	}

	go func() {
		<-p.ctx.Done()
		p.wg.Wait()
		l.Close()
	}()

	for {
		conn, err := l.Accept()

		if _, ok := err.(*net.OpError); ok {
			break
		}

		if err != nil {
			fmt.Println(err)
			continue
		}

		p.wg.Add(1)
		go p.handle(conn)
	}

	p.wg.Wait()
	return nil
}

func (p *Proxy) handle(req net.Conn) {
	defer p.wg.Done()
	defer req.Close()

	defer p.cluster.With(func(w workers.Worker) error {
		client, err := net.Dial("tcp", w.Addr())

		if err != nil {
			return err
		}

		defer client.Close()

		go io.Copy(client, req)
		io.Copy(req, client)
		return nil
	})
}

func (p *Proxy) Shutdown(ctx context.Context) error {
	p.cancel()

	c := make(chan struct{})

	go func() {
		defer close(c)
		p.wg.Wait()
		c <- struct{}{}
	}()

	select {
	case <-c:

	case <-ctx.Done():
		return errors.New("Shutdown failed")
	}

	return nil

}
