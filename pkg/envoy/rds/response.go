package rds

import (
	"fmt"
	"strings"
	"ws/osm/pkg/kubernetes"
	"ws/osm/pkg/service"

	//<<<<<<< HEAD
//	"strings"
//	"fmt"
//
//	set "github.com/deckarep/golang-set"
//	xds_route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
//=======
//>>>>>>> 3d923b3f2d72006f6cdaad056938c492c364196d
	xds_discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/golang/protobuf/ptypes"

	"github.com/openservicemesh/osm/pkg/catalog"
	"github.com/openservicemesh/osm/pkg/certificate"
	"github.com/openservicemesh/osm/pkg/configurator"
	"github.com/openservicemesh/osm/pkg/envoy"
	"github.com/openservicemesh/osm/pkg/envoy/route"
	"github.com/openservicemesh/osm/pkg/trafficpolicy"
)

// NewResponse creates a new Route Discovery Response.
func NewResponse(cataloger catalog.MeshCataloger, proxy *envoy.Proxy, _ *xds_discovery.DiscoveryRequest, cfg configurator.Configurator, _ certificate.Manager) (*xds_discovery.DiscoveryResponse, error) {
	var inboundTrafficPolicies []*trafficpolicy.InboundTrafficPolicy
	var outboundTrafficPolicies []*trafficpolicy.OutboundTrafficPolicy

	proxyIdentity, err := catalog.GetServiceAccountFromProxyCertificate(proxy.GetCertificateCommonName())
	if err != nil {
		log.Error().Err(err).Msgf("Error looking up Service Account for Envoy with serial number=%q", proxy.GetCertificateSerialNumber())
		return nil, err
	}

//<<<<<<< HEAD
//
//	allTrafficPolicies, err := catalog.ListTrafficPolicies(proxyServiceName)
//	if err != nil {
//		log.Error().Err(err).Msg(fmt.Sprintf("Failed listing routes for proxyServiceName:%+v", proxyServiceName))
//		return nil, err
//	}
//	log.Debug().Msgf("RDS proxy:%+v trafficPolicies:%+v", proxy, allTrafficPolicies)
//
//	resp := &xds_discovery.DiscoveryResponse{
//		TypeUrl: string(envoy.TypeRDS),
//	}
//=======
	services, err := cataloger.GetServicesFromEnvoyCertificate(proxy.GetCertificateCommonName())
	if err != nil {
		log.Error().Err(err).Msgf("Error looking up services for Envoy with serial number=%q", proxy.GetCertificateSerialNumber())
		return nil, err
	}
//>>>>>>> 3d923b3f2d72006f6cdaad056938c492c364196d

	// Build traffic policies from  either SMI Traffic Target and Traffic Split or service discovery
	// depending on whether permissive mode is enabled or not
	inboundTrafficPolicies = cataloger.ListInboundTrafficPolicies(proxyIdentity, services)
	outboundTrafficPolicies = cataloger.ListOutboundTrafficPolicies(proxyIdentity)

//<<<<<<< HEAD
//	for _, trafficPolicy := range allTrafficPolicies {
//		isSourceService := trafficPolicy.Source.Equals(proxyServiceName)
//		isDestinationService := trafficPolicy.Destination.GetMeshService().Equals(proxyServiceName)
//		svc := trafficPolicy.Destination.GetMeshService()
//		hostnames, err := catalog.GetResolvableHostnamesForUpstreamService(proxyServiceName, svc)
//		//filter out traffic split service, reference to pkg/catalog/xds_certificates.go:74
//		if isTrafficSplitService(svc, allTrafficSplits) {
//			continue
//		}
//=======
	// Get Ingress inbound policies for the proxy
	for _, svc := range services {
		ingressInboundPolicies, err := cataloger.GetIngressPoliciesForService(svc)
//>>>>>>> 3d923b3f2d72006f6cdaad056938c492c364196d
		if err != nil {
			log.Error().Err(err).Msgf("Error looking up ingress policies for service=%s", svc.String())
			return nil, err
		}
//<<<<<<< HEAD
//		log.Debug().Msgf("RDS hostnames: %+v", hostnames)
//
//		// multiple targets exist per service
//		var weightedCluster service.WeightedCluster
//		target := trafficPolicy.Destination
//		if target.Port != 0 {
//			hostnames = filterOnTargetPort(hostnames, target.Port)
//			log.Debug().Msgf("RDS filtered hostnames: %+v", hostnames)
//			weightedCluster, err = catalog.GetWeightedClusterForServicePort(target)
//			if err != nil {
//				log.Error().Err(err).Msg("Failed listing clusters")
//				return nil, err
//			}
//		} else {
//
//			weightedCluster, err = catalog.GetWeightedClusterForService(svc)
//			if err != nil {
//				log.Error().Err(err).Msg("Failed listing clusters")
//				return nil, err
//			}
//		}
//		log.Debug().Msgf("RDS weightedCluster: %+v", weightedCluster)
//
//		// All routes from a given source to destination are part of 1 traffic policy between the source and destination.
//		for _, hostname := range hostnames {
//			for _, httpRoute := range trafficPolicy.HTTPRouteMatches {
//				if isSourceService {
//					aggregateRoutesByHost(outboundAggregatedRoutesByHostnames, httpRoute, weightedCluster, hostname, target.Port)
//				}
//
//				if isDestinationService {
//					aggregateRoutesByHost(inboundAggregatedRoutesByHostnames, httpRoute, weightedCluster, hostname, target.Port)
//				}
//			}
//		}
//	}
//
//	/* do not include ingress routes for now as iptables should take care of it
//	if err = updateRoutesForIngress(proxyServiceName, catalog, inboundAggregatedRoutesByHostnames); err != nil {
//		return nil, err
//=======
		inboundTrafficPolicies = trafficpolicy.MergeInboundPolicies(true, inboundTrafficPolicies, ingressInboundPolicies...)
	}

	routeConfiguration := route.BuildRouteConfiguration(inboundTrafficPolicies, outboundTrafficPolicies, proxy)
	resp := &xds_discovery.DiscoveryResponse{
		TypeUrl: string(envoy.TypeRDS),
//>>>>>>> 3d923b3f2d72006f6cdaad056938c492c364196d
	}
	*/

//<<<<<<< HEAD
//	route.UpdateRouteConfiguration(catalog, outboundAggregatedRoutesByHostnames, outboundRouteConfig, route.OutboundRoute)
//	route.UpdateRouteConfiguration(catalog, inboundAggregatedRoutesByHostnames, inboundRouteConfig, route.InboundRoute)
//	routeConfiguration = append(routeConfiguration, inboundRouteConfig)
//	routeConfiguration = append(routeConfiguration, outboundRouteConfig)
//
//	log.Debug().Msgf("RDS proxy: %+v routeConfiguration: %+v", proxy, routeConfiguration)
//
//=======
//>>>>>>> 3d923b3f2d72006f6cdaad056938c492c364196d
	for _, config := range routeConfiguration {
		marshalledRouteConfig, err := ptypes.MarshalAny(config)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to marshal route config for proxy")
			return nil, err
		}
		resp.Resources = append(resp.Resources, marshalledRouteConfig)
	}

//<<<<<<< HEAD
//func aggregateRoutesByHost(routesPerHost map[string]map[string]trafficpolicy.RouteWeightedClusters, routePolicy trafficpolicy.HTTPRouteMatch, weightedCluster service.WeightedCluster, hostname string, targetPort int) {
//	host := kubernetes.GetServiceFromHostname(hostname)
//	if targetPort != 0 {
//		host = fmt.Sprintf("%s:%d", host, targetPort)
//	}
//	_, exists := routesPerHost[host]
//	if !exists {
//		// no host found, create a new route map
//		routesPerHost[host] = make(map[string]trafficpolicy.RouteWeightedClusters)
//	}
//	routePolicyWeightedCluster, routeFound := routesPerHost[host][routePolicy.PathRegex]
//	//log.Debug().Msgf("RDS aggregateRoutesByHost: routeFound:%t pathregex:%+v", routeFound, routePolicy.PathRegex)
//	if routeFound {
//		// add the cluster to the existing route
//		routePolicyWeightedCluster.WeightedClusters.Add(weightedCluster)
//		routePolicyWeightedCluster.HTTPRouteMatch.Methods = append(routePolicyWeightedCluster.HTTPRouteMatch.Methods, routePolicy.Methods...)
//		if routePolicyWeightedCluster.HTTPRouteMatch.Headers == nil {
//			routePolicyWeightedCluster.HTTPRouteMatch.Headers = make(map[string]string)
//		}
//		for headerKey, headerValue := range routePolicy.Headers {
//			routePolicyWeightedCluster.HTTPRouteMatch.Headers[headerKey] = headerValue
//		}
//		routePolicyWeightedCluster.Hostnames.Add(hostname)
//		routesPerHost[host][routePolicy.PathRegex] = routePolicyWeightedCluster
//	} else {
//		// no route found, create a new route and cluster mapping on host
//		routesPerHost[host][routePolicy.PathRegex] = createRoutePolicyWeightedClusters(routePolicy, weightedCluster, hostname)
//	}
//}
//
//func createRoutePolicyWeightedClusters(routePolicy trafficpolicy.HTTPRouteMatch, weightedCluster service.WeightedCluster, hostname string) trafficpolicy.RouteWeightedClusters {
//	return trafficpolicy.RouteWeightedClusters{
//		HTTPRouteMatch:   routePolicy,
//		WeightedClusters: set.NewSet(weightedCluster),
//		Hostnames:        set.NewSet(hostname),
//	}
//=======
	return resp, nil
//>>>>>>> 3d923b3f2d72006f6cdaad056938c492c364196d
}

// return only those hostnames whose name ends with ":<port>"
func filterOnTargetPort(hostnames []string, port int) []string {
	newHostnames := make([]string, 0)
	toMatch := fmt.Sprintf(":%d", port)
	for _, name := range hostnames {
		if strings.HasSuffix(name, toMatch) {
			newHostnames = append(newHostnames, name)
		}
	}
	if len(newHostnames) == 0 {
		return joinTargetPort(hostnames, port)
	}
	return newHostnames
}

// join port on all hostnames
func joinTargetPort(hostnames []string, port int) []string {
	newHostnames := make([]string, 0)
	portStr := fmt.Sprintf(":%d", port)
	for _, name := range hostnames {
		if !strings.Contains(name, ":") {
			newHostname := name + portStr
			newHostnames = append(newHostnames, newHostname)
		}
	}
	return newHostnames
}
