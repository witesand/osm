package rds

import (
<<<<<<< HEAD
	"context"
	"strings"
	"fmt"

=======
	set "github.com/deckarep/golang-set"
>>>>>>> c614ca2db542271efd6f7b2b106b9d046dc64b90
	xds_route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	xds_discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/golang/protobuf/ptypes"

	"github.com/openservicemesh/osm/pkg/catalog"
	"github.com/openservicemesh/osm/pkg/configurator"
	"github.com/openservicemesh/osm/pkg/envoy"
	"github.com/openservicemesh/osm/pkg/envoy/route"
	"github.com/openservicemesh/osm/pkg/kubernetes"
	"github.com/openservicemesh/osm/pkg/service"
	"github.com/openservicemesh/osm/pkg/trafficpolicy"
)

// NewResponse creates a new Route Discovery Response.
func NewResponse(catalog catalog.MeshCataloger, proxy *envoy.Proxy, _ *xds_discovery.DiscoveryRequest, _ configurator.Configurator) (*xds_discovery.DiscoveryResponse, error) {
	svcList, err := catalog.GetServicesFromEnvoyCertificate(proxy.GetCommonName())
	if err != nil {
		log.Error().Err(err).Msgf("Error looking up MeshService for Envoy with CN=%q", proxy.GetCommonName())
		return nil, err
	}
	// Github Issue #1575
	proxyServiceName := svcList[0]

	log.Debug().Msgf("RDS proxyServiceName:%s proxy:%+v", proxyServiceName, proxy)

	allTrafficPolicies, err := catalog.ListTrafficPolicies(proxyServiceName)
	if err != nil {
		log.Error().Err(err).Msg("Failed listing routes")
		return nil, err
	}
	log.Debug().Msgf("trafficPolicies: %+v", allTrafficPolicies)

	resp := &xds_discovery.DiscoveryResponse{
		TypeUrl: string(envoy.TypeRDS),
	}

	var routeConfiguration []*xds_route.RouteConfiguration
	outboundRouteConfig := route.NewRouteConfigurationStub(route.OutboundRouteConfigName)
	inboundRouteConfig := route.NewRouteConfigurationStub(route.InboundRouteConfigName)
	outboundAggregatedRoutesByHostnames := make(map[string]map[string]trafficpolicy.RouteWeightedClusters)
	inboundAggregatedRoutesByHostnames := make(map[string]map[string]trafficpolicy.RouteWeightedClusters)

<<<<<<< HEAD
	for _, trafficPolicies := range allTrafficPolicies {
		isSourceService := trafficPolicies.Source.Equals(proxyServiceName)
		isDestinationService := trafficPolicies.Destination.GetMeshService().Equals(proxyServiceName)
		svc := trafficPolicies.Destination.GetMeshService()
		hostnames, err := catalog.GetHostnamesForService(svc)
=======
	for _, trafficPolicy := range allTrafficPolicies {
		isSourceService := trafficPolicy.Source.Equals(proxyServiceName)
		isDestinationService := trafficPolicy.Destination.Equals(proxyServiceName)
		svc := trafficPolicy.Destination
		hostnames, err := catalog.GetResolvableHostnamesForUpstreamService(proxyServiceName, svc)
>>>>>>> c614ca2db542271efd6f7b2b106b9d046dc64b90
		if err != nil {
			log.Error().Err(err).Msg("Failed listing domains")
			return nil, err
		}
		log.Debug().Msgf("RDS hostnames: %+v", hostnames)

		// multiple targets exist per service
		var weightedCluster service.WeightedCluster
		target := trafficPolicies.Destination
		if target.Port != 0 {
			hostnames = filterOnTargetPort(hostnames, target.Port)
			log.Debug().Msgf("RDS filtered hostnames: %+v", hostnames)
			weightedCluster, err = catalog.GetWeightedClusterForServicePort(target)
			if err != nil {
				log.Error().Err(err).Msg("Failed listing clusters")
				return nil, err
			}
		} else {

			weightedCluster, err = catalog.GetWeightedClusterForService(svc)
			if err != nil {
				log.Error().Err(err).Msg("Failed listing clusters")
				return nil, err
			}
		}
		log.Debug().Msgf("RDS weightedCluster: %+v", weightedCluster)

		for _, hostname := range hostnames {
			// All routes from a given source to destination are part of 1 traffic policy between the source and destination.
			for _, httpRoute := range trafficPolicy.HTTPRoutes {
				if isSourceService {
					aggregateRoutesByHost(outboundAggregatedRoutesByHostnames, httpRoute, weightedCluster, hostname)
				}

				if isDestinationService {
					aggregateRoutesByHost(inboundAggregatedRoutesByHostnames, httpRoute, weightedCluster, hostname)
				}
			}
		}
	}

	/* do not include ingress routes for now as iptables should take care of it
	if err = updateRoutesForIngress(proxyServiceName, catalog, inboundAggregatedRoutesByHostnames); err != nil {
		return nil, err
	}
	*/

	route.UpdateRouteConfiguration(outboundAggregatedRoutesByHostnames, outboundRouteConfig, route.OutboundRoute)
	route.UpdateRouteConfiguration(inboundAggregatedRoutesByHostnames, inboundRouteConfig, route.InboundRoute)
	routeConfiguration = append(routeConfiguration, outboundRouteConfig)
	routeConfiguration = append(routeConfiguration, inboundRouteConfig)

	log.Debug().Msgf("RDS routeConfiguration: %+v", routeConfiguration)

	for _, config := range routeConfiguration {
		marshalledRouteConfig, err := ptypes.MarshalAny(config)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to marshal route config for proxy")
			return nil, err
		}
		resp.Resources = append(resp.Resources, marshalledRouteConfig)
	}
	return resp, nil
}

func aggregateRoutesByHost(routesPerHost map[string]map[string]trafficpolicy.RouteWeightedClusters, routePolicy trafficpolicy.HTTPRoute, weightedCluster service.WeightedCluster, hostname string) {
	host := kubernetes.GetServiceFromHostname(hostname)
	_, exists := routesPerHost[host]
	if !exists {
		// no host found, create a new route map
		routesPerHost[host] = make(map[string]trafficpolicy.RouteWeightedClusters)
	}
	routePolicyWeightedCluster, routeFound := routesPerHost[host][routePolicy.PathRegex]
	log.Debug().Msgf("RDS aggregateRoutesByHost: routeFound:%t pathregex:%+v", routeFound, routePolicy.PathRegex)
	if routeFound {
		// add the cluster to the existing route
		routePolicyWeightedCluster.WeightedClusters.Add(weightedCluster)
		routePolicyWeightedCluster.HTTPRoute.Methods = append(routePolicyWeightedCluster.HTTPRoute.Methods, routePolicy.Methods...)
		if routePolicyWeightedCluster.HTTPRoute.Headers == nil {
			routePolicyWeightedCluster.HTTPRoute.Headers = make(map[string]string)
		}
		for headerKey, headerValue := range routePolicy.Headers {
			routePolicyWeightedCluster.HTTPRoute.Headers[headerKey] = headerValue
		}
		routePolicyWeightedCluster.Hostnames.Add(hostname)
		routesPerHost[host][routePolicy.PathRegex] = routePolicyWeightedCluster
	} else {
		// no route found, create a new route and cluster mapping on host
		routesPerHost[host][routePolicy.PathRegex] = createRoutePolicyWeightedClusters(routePolicy, weightedCluster, hostname)
	}
}

func createRoutePolicyWeightedClusters(routePolicy trafficpolicy.HTTPRoute, weightedCluster service.WeightedCluster, hostname string) trafficpolicy.RouteWeightedClusters {
	return trafficpolicy.RouteWeightedClusters{
		HTTPRoute:        routePolicy,
		WeightedClusters: set.NewSet(weightedCluster),
		Hostnames:        set.NewSet(hostname),
	}
}

// return only those hostnames whose name ends with ":<port>"
func filterOnTargetPort(hostnames string, port int) string {
	newHostnames := make([]string, 0)
	strs := strings.Split(hostnames, ",")
	toMatch := fmt.Sprintf(":%d", port)
	for _, name := range strs {
		if strings.HasSuffix(name, toMatch) {
			newHostnames = append(newHostnames, name)
		}
	}
	if len(newHostnames) == 0 {
		return joinTargetPort(hostnames, port)
	}
	return strings.Join(newHostnames, ",")
}

// join port on all hostnames
func joinTargetPort(hostnames string, port int) string {
	newHostnames := make([]string, 0)
	strs := strings.Split(hostnames, ",")
	portStr := fmt.Sprintf(":%d", port)
	for _, name := range strs {
		if !strings.Contains(name, ":") {
			newHostname := name + portStr
			newHostnames = append(newHostnames, newHostname)
		}
	}
	return strings.Join(newHostnames, ",")
}
