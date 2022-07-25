package service

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

// PluginContainer is keeps particular services in the container and executes them
// in order they were added
type PluginContainer interface {
	Add(plugin.Handler) error
	// Execute run plugin.Handler.ServeDNS in order as it was added.
	// If error occurs, it exits immediately end return error.
	// If !dns.RcodeSuccess is returned, it exits immediatelly end return error.
	Execute(ctx context.Context, w dns.ResponseWriter, msg *dns.Msg) error
}