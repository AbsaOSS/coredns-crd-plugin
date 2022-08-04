package gateway

import (
	"context"
	"net"
	"os"

	"github.com/AbsaOSS/k8s_crd/common/k8sctrl"

	"github.com/AbsaOSS/k8s_crd/common/netutils"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"sigs.k8s.io/external-dns/endpoint"
)

const thisPlugin = "gateway"

var log = clog.NewWithPlugin(thisPlugin)

const defaultSvc = "external-dns.kube-system"

// Gateway stores all runtime configuration of a plugin
type Gateway struct {
	externalAddrFunc func(request.Request) []dns.RR
	opts             Opts
}

func NewGateway(opts Opts) *Gateway {
	gw := &Gateway{
		opts: opts,
	}
	gw.externalAddrFunc = gw.selfAddress
	return gw
}

// ServeDNS implements the plugin.Handle interface.
func (gw *Gateway) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	var clientIP net.IP
	state := request.Request{W: w, Req: r}
	log.Infof("Incoming query %s", state.QName())

	qname := state.QName()
	zone := plugin.Zones(gw.opts.zones).Matches(qname)
	clientIP = netutils.ExtractEdnsSubnet(r)

	if zone == "" {
		log.Infof("Request %s has not matched any zones %v", qname, gw.opts.zones)
		return dns.RcodeSuccess, nil // plugin.NextOrFailure(gw.Name(), gw.Next, ctx, w, r)
	}
	zone = qname[len(qname)-len(zone):] // maintain case of original query
	state.Zone = zone

	// Computing keys to look up in cache
	indexKey := netutils.StripClosingDot(state.QName())

	log.Infof("Computed Index Keys %v", indexKey)

	for _, z := range gw.opts.zones {
		if state.Name() == z { // apex query
			ret, err := gw.serveApex(state)
			return ret, err
		}
		if dns.IsSubDomain(gw.opts.apex+"."+z, state.Name()) {
			// dns subdomain test for ns. and dns. queries
			ret, err := gw.serveSubApex(state)
			return ret, err
		}
	}

	var ep = k8sctrl.Resources.DNSEndpoint.Lookup(indexKey, clientIP)
	log.Debugf("Computed response addresses %v", ep.Targets)
	m := new(dns.Msg)
	m.SetReply(state.Req)

	if len(ep.Targets) == 0 {
		m.Rcode = dns.RcodeNameError
		m.Ns = []dns.RR{gw.soa(state)}
		if err := w.WriteMsg(m); err != nil {
			log.Errorf("Failed to send a response: %s", err)
		}
		return 0, nil
	}

	switch state.QType() {
	case dns.TypeA:
		m.Answer = gw.A(state, netutils.TargetToIP(ep.Targets), ep.TTL)
	case dns.TypeTXT:
		m.Answer = gw.TXT(state, ep.Targets, ep.TTL)
	default:
		m.Ns = []dns.RR{gw.soa(state)}
	}

	if len(m.Answer) == 0 {
		m.Ns = []dns.RR{gw.soa(state)}
	}

	if err := w.WriteMsg(m); err != nil {
		log.Errorf("Failed to send a response: %s", err)
	}

	return dns.RcodeSuccess, nil
}

// Name implements the Handler interface.
func (gw *Gateway) Name() string { return thisPlugin }

// A generates dns.RR for A record
func (gw *Gateway) A(state request.Request, results []net.IP, ttl endpoint.TTL) (records []dns.RR) {
	dup := make(map[string]struct{})
	if !ttl.IsConfigured() {
		ttl = endpoint.TTL(gw.opts.ttlLow)
	}
	for _, result := range results {
		if _, ok := dup[result.String()]; !ok {
			dup[result.String()] = struct{}{}
			records = append(records, &dns.A{Hdr: dns.RR_Header{Name: state.Name(), Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: uint32(ttl)}, A: result})
		}
	}
	return records
}

// TXT generates dns.RR for TXT record
func (gw *Gateway) TXT(state request.Request, results []string, ttl endpoint.TTL) (records []dns.RR) {
	if !ttl.IsConfigured() {
		ttl = endpoint.TTL(gw.opts.ttlLow)
	}
	return append(records, &dns.TXT{Hdr: dns.RR_Header{Name: state.Name(), Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: uint32(ttl)}, Txt: results})
}

func (gw *Gateway) selfAddress(state request.Request) (records []dns.RR) {
	// TODO: need to do self-index lookup for that i need
	// a) my own namespace - easy
	// b) my own serviceName - CoreDNS/k does that via localIP->Endpoint->Service
	// I don't really want to list Endpoints just for that so will fix that later
	//
	// As a workaround I'm reading an env variable (with a default)
	index := os.Getenv("EXTERNAL_SVC")
	if index == "" {
		index = defaultSvc
	}

	var ep = k8sctrl.Resources.DNSEndpoint.Lookup(index, net.ParseIP(state.IP()))
	m := new(dns.Msg)
	m.SetReply(state.Req)
	return gw.A(state, netutils.TargetToIP(ep.Targets), ep.TTL)
}
