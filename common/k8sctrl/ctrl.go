package k8sctrl

import (
	"context"
	"net"
	"strings"

	"github.com/oschwald/maxminddb-golang"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	dnsendpoint "github.com/AbsaOSS/k8s_crd/extdns"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/external-dns/endpoint"
)

type KubeController struct {
	client      dnsendpoint.ExtDNSInterface
	controllers []cache.SharedIndexInformer
	labelFilter string
	hasSynced   bool
	resources   []*ResourceWithIndex
	indexer     cache.Indexer
	epc cache.SharedIndexInformer
}

type LookupEndpoint func(indexKey string, clientIP net.IP) (result EndpointResult)

type LookupFunc func(indexKey string, clientIP net.IP) ([]string, endpoint.TTL)

type ResourceWithIndex struct {
	Name   string
	Lookup  LookupEndpoint
	Lookup2 LookupFunc
}

const (
	defaultResyncPeriod   = 0
	endpointHostnameIndex = "endpointHostname"
)

// TODO: is new logger instance necessary
var log = clog.NewWithPlugin("k8s controller")

var OrderedResources = []*ResourceWithIndex{
	{
		Name: "DNSEndpoint",
	},
}

func NewKubeController(ctx context.Context, c *dnsendpoint.ExtDNSClient, label string) *KubeController {
	ctrl := &KubeController{
		client:      c,
		labelFilter: label,
	}
	dnsEndpoint := lookupResource("DNSEndpoint")
	if dnsEndpoint == nil {
		return ctrl
	}
	endpointController := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc:  endpointLister(ctx, ctrl.client, core.NamespaceAll, label),
			WatchFunc: endpointWatcher(ctx, ctrl.client, core.NamespaceAll, label),
		},
		&endpoint.DNSEndpoint{},
		defaultResyncPeriod,
		cache.Indexers{endpointHostnameIndex: endpointHostnameIndexFunc},
	)
	ctrl.indexer = endpointController.GetIndexer()
	ctrl.epc = endpointController
	dnsEndpoint.Lookup2 = ctrl.lookupEndpointIndex
	ctrl.controllers = append(ctrl.controllers, endpointController)
	return ctrl
}

func lookupResource(resource string) *ResourceWithIndex {
	for _, r := range OrderedResources {
		if r.Name == resource {
			return r
		}
	}
	return nil
}

func (ctrl *KubeController) UpdateResources(newResources []string) {
	ctrl.resources = []*ResourceWithIndex{}
	for _, name := range newResources {
		if resource := lookupResource(name); resource != nil {
			ctrl.resources = append(ctrl.resources, resource)
		}
	}
}

func (ctrl *KubeController) Run() {
	stopCh := make(chan struct{})
	defer close(stopCh)

	var synced []cache.InformerSynced

	for _, ctrl := range ctrl.controllers {
		go ctrl.Run(stopCh)
		synced = append(synced, ctrl.HasSynced)
	}

	if !cache.WaitForCacheSync(stopCh, synced...) {
		ctrl.hasSynced = false
	}
	log.Infof("Synced all required resources")
	ctrl.hasSynced = true

	<-stopCh
}

// HasSynced returns true if all controllers have been synced
func (ctrl *KubeController) HasSynced() bool {
	return ctrl.hasSynced
}

func endpointLister(ctx context.Context, c dnsendpoint.ExtDNSInterface, ns, label string) func(meta.ListOptions) (runtime.Object, error) {
	return func(opts meta.ListOptions) (runtime.Object, error) {
		opts.LabelSelector = label
		return c.DNSEndpoints(ns).List(ctx, opts)
	}
}

func endpointWatcher(ctx context.Context, c dnsendpoint.ExtDNSInterface, ns, label string) func(meta.ListOptions) (watch.Interface, error) {
	return func(opts meta.ListOptions) (watch.Interface, error) {
		opts.LabelSelector = label
		return c.DNSEndpoints(ns).Watch(ctx, opts)
	}
}

func endpointHostnameIndexFunc(obj interface{}) ([]string, error) {
	ep, ok := obj.(*endpoint.DNSEndpoint)
	if !ok {
		return []string{}, nil
	}

	var hostnames []string
	for _, rule := range ep.Spec.Endpoints {
		log.Infof("Adding index %s for endpoints %s", rule.DNSName, ep.Name)
		hostnames = append(hostnames, rule.DNSName)
	}
	return hostnames, nil
}

func (ctrl *KubeController) getEndpointByIndexKey(indexKey string, clientIP net.IP) (result EndpointResult) {
	log.Infof("Index key %+v", indexKey)
	objs, err := ctrl.epc.GetIndexer().ByIndex(endpointHostnameIndex, strings.ToLower(indexKey))
	if err != nil {
		return result
	}
	for _, obj := range objs {
		ep := obj.(*endpoint.DNSEndpoint)
		result = fetchEndpoint(ep, indexKey, clientIP)
		if !result.IsEmpty() {
			break
		}
	}
	return result
}

type geo struct {
	DC string `maxminddb:"datacenter"`
}

// fetchEndpoint retrieves endpoint which has DNSName equal to host
func fetchEndpoint(dnsEndpoint *endpoint.DNSEndpoint, host string, ip net.IP) (result EndpointResult) {
	if dnsEndpoint == nil {
		return result
	}
	for _, ep := range dnsEndpoint.Spec.Endpoints {
		if ep.DNSName == host {
			targets := ep.Targets
			if ep.Labels["strategy"] == "geoip" {
				targets = extractGeo(ep, ip)
			}
			result.Append(targets, ep.Labels, host, ep.RecordTTL)
			return result
		}
	}
	return result
}

func extractGeo(endpoint *endpoint.Endpoint, clientIP net.IP) (result []string) {
	db, err := maxminddb.Open("geoip.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	clientGeo := &geo{}
	err = db.Lookup(clientIP, clientGeo)
	if err != nil {
		return nil
	}

	if clientGeo.DC == "" {
		log.Infof("empty DC %+v", clientGeo)
		return result
	}

	log.Infof("clientDC: %+v", clientGeo)

	for _, ip := range endpoint.Targets {
		geoData := &geo{}
		log.Infof("processing IP %+v", ip)
		err = db.Lookup(net.ParseIP(ip), geoData)
		if err != nil {
			log.Error(err)
			continue
		}

		log.Infof("IP info: %+v", geoData.DC)
		if clientGeo.DC == geoData.DC {
			result = append(result, ip)
		}
	}
	return result
}
