package eds

import (
	"github.com/openservicemesh/osm/pkg/catalog"
	"github.com/openservicemesh/osm/pkg/witesand"
	"strconv"
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

func NewWSEdgePodClusterLoadAssignment(catalog catalog.MeshCataloger, serviceName service.MeshServicePort) *[]*xds_endpoint.ClusterLoadAssignment {
	log.Trace().Msgf("[EDS][NewWSEdgePodClusterLoadAssignment] Adding Endpoints")
	servicePort := serviceName.Port
	getMultiEndpointsCLA := func(atopMap witesand.ApigroupToPodIPMap, clusterName string) *xds_endpoint.ClusterLoadAssignment {
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

	atopMaps, _ := wscatalog.ListApigroupToPodIPs()
	for _, atopMap := range atopMaps {
		clusterName := atopMap.Apigroup + ":" + strconv.Itoa(servicePort)
		cla := getMultiEndpointsCLA(atopMap, clusterName)
		clas = append(clas, cla)

		clusterName = atopMap.Apigroup + witesand.DeviceHashSuffix + ":" + strconv.Itoa(servicePort)
		cla = getMultiEndpointsCLA(atopMap, clusterName)
		clas = append(clas, cla)
	}

	pods, _ := wscatalog.ListAllEdgePodIPs()
	for podName, podIP := range pods.PodToIPMap {
		clusterName := podName + ":" + strconv.Itoa(servicePort)
		cla := getSingleEndpointCLA(clusterName, podIP, servicePort)
		clas = append(clas, cla)
	}
	log.Trace().Msgf("[EDS][NewWSEdgePodClusterLoadAssignment] Constructed ClusterLoadAssignment: %+v", clas)
	return &clas
}

func NewWSUnicastClusterLoadAssignment(catalog catalog.MeshCataloger, serviceName service.MeshServicePort) *[]*xds_endpoint.ClusterLoadAssignment {
	log.Trace().Msgf("[EDS][NewWSUnicastClusterLoadAssignment] Adding Endpoints for Service:%+v", serviceName)
	servicePort := serviceName.Port
	serviceEndpoints, err := catalog.ListEndpointsForService(serviceName.GetMeshService())
	if err != nil {
		log.Error().Msgf("[EDS][NewWSEdgePodClusterLoadAssignment] Error adding Endpoints for Service:%+v, err:%+v", serviceName, err)
		return nil
	}
	var clas []*xds_endpoint.ClusterLoadAssignment
	for _, endpoint := range serviceEndpoints {
		if int(endpoint.Port) != servicePort {
			// skip non-interesting ports
			continue
		}
		clusterName := endpoint.PodName + ":" + strconv.Itoa(servicePort)
		cla := getSingleEndpointCLA(clusterName, endpoint.IP.String(), servicePort)
		clas = append(clas, cla)
	}
	log.Trace().Msgf("[EDS][NewWSUnicastClusterLoadAssignment] Constructed ClusterLoadAssignment: %+v", clas)
	return &clas
}
