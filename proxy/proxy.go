package proxy

import (
	"fmt"
	"io"
	"net"

	"github.com/rack-app/server/clusters"

	"github.com/rack-app/server/workers"
)

type Proxy struct {
	cluster clusters.Cluster
}

func New(c clusters.Cluster) *Proxy {
	return &Proxy{cluster: c}
}

func (p *proxy) ListenAndServe(addr string) error {
	l, err := net.Listen("tcp", addr)

	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}

		go handle(conn)
	}

}

func (p *proxy) handle(req net.Conn) {
	p.cluster.With(func(w workers.Worker) error {
		client, err := net.Dial("tcp", w.Addr())

		if err != nil {
			return err
		}

		go io.Copy(client, req)
		io.Copy(req, client)

		req.Close()
		client.Close()

		return nil
	})
}
