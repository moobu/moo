package cmd

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/moobu/moo/internal/cli"
)

const (
	defaultGatewayPort  = 80
	defaultServerPort   = 11451
	defaultServerAddr   = "localhost:11451"
	defaultServerPreset = "local"
	defaultNamespace    = "default"
)

var cmd = &cli.Cmd{
	Name:     "moo",
	About:    "Moo is a serverless engine",
	Version:  "v0.0.1",
	Wildcard: true,
}

func RunCtx(c context.Context) error {
	cmd.Init()
	return cmd.RunCtx(c)
}

func listen(c cli.Ctx, uds bool) (net.Listener, error) {
	// we use TCP if the flag uds is not set,
	// otherwise use the UNIX doamin socket
	network := "tcp"
	address := fmt.Sprintf(":%d", c.Int("port"))
	if uds {
		address = filepath.Join(os.TempDir(), "moo.sock")
		network = "unix"
	}

	listener, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}

	// see if we need to wrap a TLS listener
	if !c.Bool("secure") {
		return listener, nil
	}

	// TODO: generate the TLS config on our own if not provided
	cert, key := c.String("cert"), c.String("key")
	if len(cert) == 0 || len(key) == 0 {
		return nil, errors.New("certificates not provided")
	}

	certificate, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ClientAuth:   tls.RequireAnyClientCert,
	}
	return tls.NewListener(listener, config), nil
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
