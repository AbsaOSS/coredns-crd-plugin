/*
Copyright 2021 ABSA Group Limited

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Generated by GoLic, for more details see: https://github.com/AbsaOSS/golic
*/
package gateway

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"sigs.k8s.io/external-dns/endpoint"
)

const thisPlugin = "gateway"

var log = clog.NewWithPlugin(thisPlugin)

const defaultSvc = "external-dns.kube-system"

type lookupFunc func(indexKey string, clientIP net.IP) ([]string, endpoint.TTL)

type resourceWithIndex struct {
	name   string
	lookup lookupFunc
}

var orderedResources = []*resourceWithIndex{
	{
		name: "DNSEndpoint",
	},
}

var (
	ttlLowDefault     = uint32(60)
	ttlHighDefault    = uint32(3600)
	defaultApex       = "dns"
	defaultHostmaster = "hostmaster"
)

// Gateway stores all runtime configuration of a plugin
type Gateway struct {
	Next             plugin.Handler
	Zones            []string
	Resources        []*resourceWithIndex
	ttlLow           uint32
	ttlHigh          uint32
	Controller       *KubeController
	apex             string
	hostmaster       string
	Filter           string
	Annotation       string
	ExternalAddrFunc func(request.Request) []dns.RR
}

func NewGateway() *Gateway {
	return &Gateway{
		apex:       defaultApex,
		Resources:  orderedResources,
		ttlLow:     ttlLowDefault,
		ttlHigh:    ttlHighDefault,
		hostmaster: defaultHostmaster,
	}
}

func lookupResource(resource string) *resourceWithIndex {

	for _, r := range orderedResources {
		if r.name == resource {
			return r
		}
	}
	return nil
}

func (gw *Gateway) UpdateResources(newResources []string) {

	gw.Resources = []*resourceWithIndex{}

	for _, name := range newResources {
		if resource := lookupResource(name); resource != nil {
			gw.Resources = append(gw.Resources, resource)
		}
	}
}

func extractEdnsSubnet(msg *dns.Msg) net.IP {
	edns := msg.IsEdns0()
	if edns == nil {
		return nil
	}
	for _, o := range edns.Option {
		if o.Option() == dns.EDNS0SUBNET {
			subnet := o.(*dns.EDNS0_SUBNET)
			return subnet.Address
		}
	}
	return nil
}

func targetToIP(targets []string) (ips []net.IP) {
	for _, ip := range targets {
		ips = append(ips, net.ParseIP(ip))
	}
	return
}

// ServeDNS implements the plugin.Handle interface.
func (gw *Gateway) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	var clientIP net.IP
	state := request.Request{W: w, Req: r}
	log.Infof("Incoming query %s", state.QName())

	qname := state.QName()
	zone := plugin.Zones(gw.Zones).Matches(qname)
	clientIP = extractEdnsSubnet(r)

	if zone == "" {
		log.Infof("Request %s has not matched any zones %v", qname, gw.Zones)
		return dns.RcodeSuccess, nil // plugin.NextOrFailure(gw.Name(), gw.Next, ctx, w, r)
	}
	zone = qname[len(qname)-len(zone):] // maintain case of original query
	state.Zone = zone

	// Computing keys to look up in cache
	indexKey := stripClosingDot(state.QName())

	log.Infof("Computed Index Keys %v", indexKey)

	if !gw.Controller.HasSynced() {
		// TODO maybe there's a better way to do this? e.g. return an error back to the client?
		return dns.RcodeServerFailure, plugin.Error(thisPlugin, fmt.Errorf("could not sync required resources"))
	}

	for _, z := range gw.Zones {
		if state.Name() == z { // apex query
			ret, err := gw.serveApex(state)
			return ret, err
		}
		if dns.IsSubDomain(gw.apex+"."+z, state.Name()) {
			// dns subdomain test for ns. and dns. queries
			ret, err := gw.serveSubApex(state)
			return ret, err
		}
	}

	var addrs []string
	var ttl endpoint.TTL

	// Iterate over supported resources and lookup DNS queries
	// Stop once we've found at least one match
	for _, resource := range gw.Resources {
		addrs, ttl = resource.lookup(indexKey, clientIP)
		if len(addrs) > 0 {
			break
		}
	}
	log.Debugf("Computed response addresses %v", addrs)

	m := new(dns.Msg)
	m.SetReply(state.Req)

	if len(addrs) == 0 {
		m.Rcode = dns.RcodeNameError
		m.Ns = []dns.RR{gw.soa(state)}
		if err := w.WriteMsg(m); err != nil {
			log.Errorf("Failed to send a response: %s", err)
		}
		return 0, nil
	}

	switch state.QType() {
	case dns.TypeA:
		m.Answer = gw.A(state, targetToIP(addrs), ttl)
	case dns.TypeTXT:
		m.Answer = gw.TXT(state, addrs, ttl)
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
		ttl = endpoint.TTL(gw.ttlLow)
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
		ttl = endpoint.TTL(gw.ttlLow)
	}
	return append(records, &dns.TXT{Hdr: dns.RR_Header{Name: state.Name(), Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: uint32(ttl)}, Txt: results})
}

func (gw *Gateway) SelfAddress(state request.Request) (records []dns.RR) {
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

	var addrs []string
	var ttl endpoint.TTL
	for _, resource := range gw.Resources {
		addrs, ttl = resource.lookup(index, net.ParseIP(state.IP()))
		if len(addrs) > 0 {
			break
		}
	}

	m := new(dns.Msg)
	m.SetReply(state.Req)
	return gw.A(state, targetToIP(addrs), ttl)
}

func (gw *Gateway) SetTTLLow(ttl uint32) {
	gw.ttlLow = ttl
}

func (gw *Gateway) SetTTLHigh(ttl uint32) {
	gw.ttlHigh = ttl
}

func (gw *Gateway) SetApex(apex string) {
	gw.apex = apex
}

// Strips the closing dot unless it's "."
func stripClosingDot(s string) string {
	if len(s) > 1 {
		return strings.TrimSuffix(s, ".")
	}
	return s
}
