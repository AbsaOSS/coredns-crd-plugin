package wrr

import (
	"context"
	"fmt"
	"net"

	"github.com/AbsaOSS/k8s_crd/common/k8sctrl"
	"github.com/AbsaOSS/k8s_crd/common/netutils"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

type WeightRoundRobin struct {
	resources []*k8sctrl.ResourceWithIndex
}

const thisPlugin = "wrr"

func NewWeightRoundRobin() *WeightRoundRobin {
	return &WeightRoundRobin{
		resources: k8sctrl.OrderedResources,
	}
}

func (wrr *WeightRoundRobin) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	var clientIP net.IP
	state := request.Request{W: w, Req: r}
	clientIP = netutils.ExtractEdnsSubnet(r)
	indexKey := netutils.StripClosingDot(state.QName())

	a, _, _ := parseAnswerSection(r.Answer)

	var ep k8sctrl.LocalDNSEndpoint
	for _, resource := range wrr.resources {
		ep = resource.Lookup(indexKey, clientIP)
		if len(ep.Targets) > 0 {
			break
		}
	}

	fmt.Println(a, ep.Labels)
	return dns.RcodeSuccess, nil
}

func (wrr *WeightRoundRobin) Name() string { return thisPlugin }
