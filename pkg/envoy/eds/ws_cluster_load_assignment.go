package eds

import (
	"github.com/openservicemesh/osm/pkg/catalog"
	"github.com/openservicemesh/osm/pkg/envoy"
	"github.com/openservicemesh/osm/pkg/service"
	"github.com/openservicemesh/osm/pkg/witesand"
	"strconv"
	xds_core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	xds_endpoint "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
)

func getSingleEndpointCLA(clusterName string, podIP string, servicePort int) *xds_endpoint.ClusterLoadAssignment {
	cla := xds_endpoint.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []*xds_endpoint.LocalityLbEndpoints{
			{
				Locality: &xds_core.Locality{
					Zone: zone,
				},
				LbEndpoints: []*xds_endpoint.LbEndpoint{},
			},
		},
	}

	log.Trace().Msgf("[EDS][getCLASingleEndpoint] Adding Endpoint: Cluster=%s, IP=%s, Port=%d", clusterName, podIP, servicePort)
	lbEpt := xds_endpoint.LbEndpoint{
		HostIdentifier: &xds_endpoint.LbEndpoint_Endpoint{
			Endpoint: &xds_endpoint.Endpoint{
				Address: envoy.GetAddress(podIP, uint32(servicePort)),
			},
		},
	}
	cla.Endpoints[0].LbEndpoints = append(cla.Endpoints[0].LbEndpoints, &lbEpt)
	log.Debug().Msgf("[EDS] Constructed ClusterLoadAssignment: %+v", cla)
	return &cla
}
func NewWSEdgePodClusterLoadAssignment(catalog catalog.MeshCataloger, serviceName service.MeshService) *[]*xds_endpoint.ClusterLoadAssignment {
	log.Trace().Msgf("[EDS][NewWSEdgePodClusterLoadAssignment] Adding Endpoints")
	getMultiEndpointsCLA := func(atopMap witesand.ApigroupToPodIPMap, clusterName string, servicePort int) *xds_endpoint.ClusterLoadAssignment {
		cla := xds_endpoint.ClusterLoadAssignment{
			ClusterName: clusterName,
			Endpoints: []*xds_endpoint.LocalityLbEndpoints{
				{
					Locality: &xds_core.Locality{
						Zone: zone,
					},
					LbEndpoints: []*xds_endpoint.LbEndpoint{},
				},
			},
		}

		for _, podIP := range atopMap.PodIPs {
			log.Trace().Msgf("[EDS][NewWSEdgePodClusterLoadAssignment] Adding Endpoint: Cluster=%s, Services=%s, IP=%+v, Port=%d", atopMap.Apigroup, serviceName.String(), podIP, servicePort)
			lbEpt := xds_endpoint.LbEndpoint{
				HostIdentifier: &xds_endpoint.LbEndpoint_Endpoint{
					Endpoint: &xds_endpoint.Endpoint{
						Address: envoy.GetAddress(podIP, uint32(servicePort)),
					},
				},
			}
			cla.Endpoints[0].LbEndpoints = append(cla.Endpoints[0].LbEndpoints, &lbEpt)
		}
		return &cla
	}

	var clas []*xds_endpoint.ClusterLoadAssignment
	wscatalog := catalog.GetWitesandCataloger()

	serviceEndpoints, err := catalog.GetResolvableServiceEndpoints(serviceName)
	if err != nil {
		log.Info().Msgf("getWSUnicastUpstreamServiceCluster err %+v", err)
		return nil
	}

	atopMaps, _ := wscatalog.ListApigroupToPodIPs()

	for _, endpoint := range serviceEndpoints {
		for _, atopMap := range atopMaps {
			clusterName := atopMap.Apigroup + ":" + strconv.Itoa(int(endpoint.Port))
			cla := getMultiEndpointsCLA(atopMap, clusterName, int(endpoint.Port))
			clas = append(clas, cla)

			clusterName = atopMap.Apigroup + witesand.DeviceHashSuffix + ":" + strconv.Itoa(int(endpoint.Port))
			cla = getMultiEndpointsCLA(atopMap, clusterName, int(endpoint.Port))
			clas = append(clas, cla)
		}
	}

	// Skipping this as it is done via unicast cla
	//pods, _ := wscatalog.ListAllEdgePodIPs()
	//for podName, podIP := range pods.PodToIPMap {
	//	clusterName := podName + ":" + strconv.Itoa(servicePort)
	//	cla := getSingleEndpointCLA(clusterName, podIP, servicePort)
	//	clas = append(clas, cla)
	//}
	log.Trace().Msgf("[EDS][NewWSEdgePodClusterLoadAssignment] Constructed ClusterLoadAssignment: %+v", clas)
	return &clas
}

func NewWSUnicastClusterLoadAssignment(catalog catalog.MeshCataloger, serviceName service.MeshService) *[]*xds_endpoint.ClusterLoadAssignment {
	log.Trace().Msgf("[EDS][NewWSUnicastClusterLoadAssignment] Adding Endpoints for Service:%+v", serviceName)
	//servicePort := serviceName.Port
	serviceEndpoints, err := catalog.GetResolvableServiceEndpoints(serviceName)
	if err != nil {
		log.Info().Msgf("getWSUnicastUpstreamServiceCluster err %+v", err)
		return nil
	}

	var clas []*xds_endpoint.ClusterLoadAssignment
	for _, endpoint := range serviceEndpoints {
		//fred
		//if int(endpoint.Port) != servicePort {
		//	// skip non-interesting ports
		//	continue
		//}
		clusterName := endpoint.PodName + ":" + strconv.Itoa(int(endpoint.Port))
		cla := getSingleEndpointCLA(clusterName, endpoint.IP.String(), int(endpoint.Port))
		clas = append(clas, cla)
	}
	log.Trace().Msgf("[EDS][NewWSUnicastClusterLoadAssignment] Constructed ClusterLoadAssignment: %+v", clas)
	return &clas
}
