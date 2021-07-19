package route

import (
	mapset "github.com/deckarep/golang-set"
	xds_route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	xds_matcher "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/openservicemesh/osm/pkg/catalog"
	"github.com/openservicemesh/osm/pkg/constants"
	"github.com/openservicemesh/osm/pkg/trafficpolicy"
	"github.com/openservicemesh/osm/pkg/witesand"
)

func buildOutboundRoutesForWS(outRoutes []*trafficpolicy.RouteWeightedClusters) []*xds_route.Route {
	var routes []*xds_route.Route
	for _, outRoute := range outRoutes {
		emptyHeaders := map[string]string{}
		routes = append(routes, buildRouteForWSEdgePod(trafficpolicy.PathMatchRegex, constants.RegexMatchAll, constants.WildcardHTTPMethod, emptyHeaders, outRoute.WeightedClusters, outRoute.TotalClustersWeight(), outboundRoute))
	}
	return routes
}


func buildRouteForWSEdgePod(pathMatchTypeType trafficpolicy.PathMatchType, path string, method string, headersMap map[string]string, weightedClusters mapset.Set, totalWeight int, direction Direction) *xds_route.Route {
	t := &xds_route.RouteAction_HashPolicy_Header_{
		&xds_route.RouteAction_HashPolicy_Header{
			HeaderName:   witesand.WSHashHeader,
			RegexRewrite: nil,
		},
	}

	r := &xds_route.RouteAction_HashPolicy{
		PolicySpecifier: t,
		Terminal:        false,
	}

	// disable the timeouts, without this synchronous calls timeout
	routeTimeout := duration.Duration{Seconds: 0}

	route := xds_route.Route{
		Match: &xds_route.RouteMatch{
			Headers: getHeadersForRoute(method, headersMap),
		},
		Action: &xds_route.Route_Route{
			Route: &xds_route.RouteAction{
				ClusterSpecifier: &xds_route.RouteAction_ClusterHeader{
					ClusterHeader: witesand.WSClusterHeader,
				},
				HashPolicy: []*xds_route.RouteAction_HashPolicy{r},
				Timeout: &routeTimeout,
			},
		},
	}

	switch pathMatchTypeType {
	case trafficpolicy.PathMatchRegex:
		route.Match.PathSpecifier = &xds_route.RouteMatch_SafeRegex{
			SafeRegex: &xds_matcher.RegexMatcher{
				EngineType: &xds_matcher.RegexMatcher_GoogleRe2{GoogleRe2: &xds_matcher.RegexMatcher_GoogleRE2{}},
				Regex:      path,
			},
		}

	case trafficpolicy.PathMatchExact:
		route.Match.PathSpecifier = &xds_route.RouteMatch_Path{
			Path: path,
		}

	case trafficpolicy.PathMatchPrefix:
		route.Match.PathSpecifier = &xds_route.RouteMatch_Prefix{
			Prefix: path,
		}
	}

	return &route
}

func isWitesandOutboundHost(catalog catalog.MeshCataloger, host string) bool {
	return catalog.GetWitesandCataloger().IsWSUnicastService(host)
}
