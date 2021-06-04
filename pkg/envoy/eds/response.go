package eds

import (
	xds_discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"

	"github.com/openservicemesh/osm/pkg/catalog"
	"github.com/openservicemesh/osm/pkg/certificate"
	"github.com/openservicemesh/osm/pkg/configurator"
	"github.com/openservicemesh/osm/pkg/endpoint"
	"github.com/openservicemesh/osm/pkg/envoy"
	"github.com/openservicemesh/osm/pkg/service"
)

// NewResponse creates a new Endpoint Discovery Response.
func NewResponse(meshCatalog catalog.MeshCataloger, proxy *envoy.Proxy, _ *xds_discovery.DiscoveryRequest, _ configurator.Configurator, _ certificate.Manager) (*xds_discovery.DiscoveryResponse, error) {
	proxyIdentity, err := catalog.GetServiceAccountFromProxyCertificate(proxy.GetCertificateCommonName())
	if err != nil {
		log.Error().Err(err).Msgf("Error looking up proxy identity for proxy with SerialNumber=%s on Pod with UID=%s", proxy.GetCertificateSerialNumber(), proxy.GetPodUID())
		return nil, err
	}

<<<<<<< HEAD
	allTrafficPolicies, err := catalog.ListTrafficPolicies(proxyServiceName)
	//log.Debug().Msgf("EDS svc %s allTrafficPolicies %+v", proxyServiceName, allTrafficPolicies)

=======
	allowedEndpoints, err := getEndpointsForProxy(meshCatalog, proxyIdentity)
>>>>>>> 3d923b3f2d72006f6cdaad056938c492c364196d
	if err != nil {
		log.Error().Err(err).Msgf("Error looking up endpoints for proxy with SerialNumber=%s on Pod with UID=%s", proxy.GetCertificateSerialNumber(), proxy.GetPodUID())
		return nil, err
	}

<<<<<<< HEAD
	outboundServicesEndpoints := make(map[service.MeshServicePort][]endpoint.Endpoint)
	for _, trafficPolicy := range allTrafficPolicies {
		isSourceService := trafficPolicy.Source.Equals(proxyServiceName)
		if isSourceService {
			destService := trafficPolicy.Destination.GetMeshService()
			serviceEndpoints, err := catalog.ListEndpointsForService(destService)
			log.Trace().Msgf("EDS: proxy:%s, serviceEndpoints:%+v", proxyServiceName, serviceEndpoints)
			if err != nil {
				log.Error().Err(err).Msgf("Failed listing endpoints for proxy %s", proxyServiceName)
				return nil, err
			}
			destServicePort := trafficPolicy.Destination
			if destServicePort.Port == 0  {
				outboundServicesEndpoints[destServicePort] = serviceEndpoints
				continue
			}
			// if port specified, filter based on port
			filteredEndpoints := make([]endpoint.Endpoint, 0)
			for _, endpoint := range serviceEndpoints {
				if int(endpoint.Port) != destServicePort.Port {
					continue
				}
				filteredEndpoints = append(filteredEndpoints, endpoint)
			}
			outboundServicesEndpoints[destServicePort] = filteredEndpoints
		}
	}

	log.Trace().Msgf("Outbound service endpoints for proxy %s: %v", proxyServiceName, outboundServicesEndpoints)

	var protos []*any.Any
	for svc, endpoints := range outboundServicesEndpoints {
		if catalog.GetWitesandCataloger().IsWSEdgePodService(svc) {
			loadAssignments := cla.NewWSEdgePodClusterLoadAssignment(catalog, svc)
			for _, loadAssignment := range *loadAssignments {
				proto, err := ptypes.MarshalAny(loadAssignment)
				if err != nil {
					log.Error().Err(err).Msgf("Error marshalling EDS payload for proxy %s: %+v", proxyServiceName, loadAssignment)
					continue
				}
				protos = append(protos, proto)
			}
			continue
		} else if catalog.GetWitesandCataloger().IsWSUnicastService(svc.Name) {
			loadAssignments := cla.NewWSUnicastClusterLoadAssignment(catalog, svc)
			for _, loadAssignment := range *loadAssignments {
				proto, err := ptypes.MarshalAny(loadAssignment)
				if err != nil {
					log.Error().Err(err).Msgf("Error marshalling EDS payload for proxy %s: %+v", proxyServiceName, loadAssignment)
					continue
				}
				protos = append(protos, proto)
			}
			// fall thru for default CLAs
		}
		loadAssignment := cla.NewClusterLoadAssignment(svc, endpoints)
=======
	var protos []*any.Any
	for svc, endpoints := range allowedEndpoints {
		loadAssignment := newClusterLoadAssignment(svc, endpoints)
>>>>>>> 3d923b3f2d72006f6cdaad056938c492c364196d
		proto, err := ptypes.MarshalAny(loadAssignment)
		if err != nil {
			log.Error().Err(err).Msgf("Error marshalling EDS payload for proxy with SerialNumber=%s on Pod with UID=%s", proxy.GetCertificateSerialNumber(), proxy.GetPodUID())
			continue
		}
		protos = append(protos, proto)
	}

	log.Debug().Msgf("EDS url:%s protos: %+v", string(envoy.TypeEDS), protos)
	resp := &xds_discovery.DiscoveryResponse{
		Resources: protos,
		TypeUrl:   string(envoy.TypeEDS),
	}
	return resp, nil
}

// getEndpointsForProxy returns only those service endpoints that belong to the allowed outbound service accounts for the proxy
func getEndpointsForProxy(meshCatalog catalog.MeshCataloger, proxyIdentity service.K8sServiceAccount) (map[service.MeshService][]endpoint.Endpoint, error) {
	allowedServicesEndpoints := make(map[service.MeshService][]endpoint.Endpoint)

	for _, dstSvc := range meshCatalog.ListAllowedOutboundServicesForIdentity(proxyIdentity) {
		endpoints, err := meshCatalog.ListAllowedEndpointsForService(proxyIdentity, dstSvc)
		if err != nil {
			log.Error().Err(err).Msgf("Failed listing allowed endpoints for service %s for proxy identity %s", dstSvc, proxyIdentity)
			continue
		}
		allowedServicesEndpoints[dstSvc] = endpoints
	}
	log.Trace().Msgf("Allowed outbound service endpoints for proxy with identity %s: %v", proxyIdentity, allowedServicesEndpoints)
	return allowedServicesEndpoints, nil
}
