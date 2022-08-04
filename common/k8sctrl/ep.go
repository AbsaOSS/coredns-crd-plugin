package k8sctrl

/*
Copyright 2022 The k8gb Contributors

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

import (
	"net"

	"github.com/oschwald/maxminddb-golang"
	"sigs.k8s.io/external-dns/endpoint"
)

type geo struct {
	DC string `maxminddb:"datacenter"`
}

type LocalDNSEndpoint struct {
	Targets []string
	TTL     endpoint.TTL
	Labels  map[string]string
	DNSName string
}

func extractLocalEndpoint(ep *endpoint.DNSEndpoint, ip net.IP, host string) (result LocalDNSEndpoint) {
	result = LocalDNSEndpoint{}
	for _, e := range ep.Spec.Endpoints {
		if e.DNSName == host {
			result.DNSName = host
			result.Labels = e.Labels
			result.TTL = e.RecordTTL
			result.Targets = e.Targets
			if e.Labels["strategy"] == "geoip" {
				targets := result.extractGeo(e, ip)
				if len(targets) > 0 {
					result.Targets = targets
				}
			}
			break
		}
	}
	return result
}

func (lep LocalDNSEndpoint) isEmpty() bool {
	return len(lep.Targets) == 0 && (len(lep.Labels) == 0) && (lep.TTL == 0)
}

func (lep LocalDNSEndpoint) extractGeo(endpoint *endpoint.Endpoint, clientIP net.IP) (result []string) {
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
