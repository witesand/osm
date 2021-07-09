package cds

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/openservicemesh/osm/pkg/catalog"
	"github.com/openservicemesh/osm/pkg/configurator"
	"github.com/openservicemesh/osm/pkg/envoy"
	"github.com/openservicemesh/osm/pkg/identity"
	"github.com/openservicemesh/osm/pkg/service"
	"github.com/openservicemesh/osm/pkg/witesand"
	"strconv"
	xds_cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
)

// getWSEdgePodUpstreamServiceCluster returns an Envoy Cluster corresponding to the given upstream service
func getWSEdgePodUpstreamServiceCluster(catalog catalog.MeshCataloger, downstreamIdentity identity.ServiceIdentity, upstreamSvc service.MeshService, cfg configurator.Configurator, cluster *[]*xds_cluster.Cluster) error {
	wscatalog := catalog.GetWitesandCataloger()
	apigroupClusterNames, err := wscatalog.ListApigroupClusterNames()
	if err != nil {
		return err
	}

	HTTP2ProtocolOptions, err := envoy.GetHTTP2ProtocolOptions()
	if err != nil {
		log.Info().Msgf("Gethttp2 option failed err %+v", err)
		return err
	}

	serviceEndpoints, err := catalog.GetResolvableServiceEndpoints(upstreamSvc)
	if err != nil {
		log.Info().Msgf("getWSUnicastUpstreamServiceCluster err %+v", err)
		return err
	}

	for _, endpoint := range serviceEndpoints {
		// create clusters with apigroup-names with ROUND_ROBIN
		for _, apigroupName := range apigroupClusterNames {
			clusterName := apigroupName + ":" + strconv.Itoa(int(endpoint.Port))

			remoteCluster := &xds_cluster.Cluster{
				Name:           clusterName,
				ConnectTimeout: ptypes.DurationProto(clusterConnectTimeout),
				CircuitBreakers: &xds_cluster.CircuitBreakers{
					Thresholds: makeWSThresholds(),
				},
				TypedExtensionProtocolOptions: HTTP2ProtocolOptions,
			}

			remoteCluster.ClusterDiscoveryType = &xds_cluster.Cluster_Type{Type: xds_cluster.Cluster_EDS}
			remoteCluster.EdsClusterConfig = &xds_cluster.Cluster_EdsClusterConfig{EdsConfig: envoy.GetADSConfigSource()}
			remoteCluster.LbPolicy = xds_cluster.Cluster_ROUND_ROBIN
			*cluster = append(*cluster, remoteCluster)
		}

		// create clusters with apigroup-names + "device-hash" with RING_HASH
		for _, apigroupName := range apigroupClusterNames {
			clusterName := apigroupName + witesand.DeviceHashSuffix + ":" + strconv.Itoa(int(endpoint.Port))

			remoteCluster := &xds_cluster.Cluster{
				Name:           clusterName,
				ConnectTimeout: ptypes.DurationProto(clusterConnectTimeout),
				CircuitBreakers: &xds_cluster.CircuitBreakers{
					Thresholds: makeWSThresholds(),
				},
				TypedExtensionProtocolOptions: HTTP2ProtocolOptions,
			}

			remoteCluster.ClusterDiscoveryType = &xds_cluster.Cluster_Type{Type: xds_cluster.Cluster_EDS}
			remoteCluster.EdsClusterConfig = &xds_cluster.Cluster_EdsClusterConfig{EdsConfig: envoy.GetADSConfigSource()}
			remoteCluster.LbPolicy = xds_cluster.Cluster_RING_HASH
			*cluster = append(*cluster, remoteCluster)
		}

		//This is no longer needed as it will be done as part of unicast support.
		//edgePodNames, err := wscatalog.ListAllEdgePods()
		//if err != nil {
		//	return err
		//}

		//// create clusters with pod-names with ROUND_ROBIN
		//for _, edgePodName := range edgePodNames {
		//	clusterName := edgePodName + ":" + strconv.Itoa(int(endpoint.Port))
		//
		//	remoteCluster := &xds_cluster.Cluster{
		//		Name:           clusterName,
		//		ConnectTimeout: ptypes.DurationProto(clusterConnectTimeout),
		//		CircuitBreakers: &xds_cluster.CircuitBreakers{
		//			Thresholds: makeWSThresholds(),
		//		},
		//	}
		//
		//	remoteCluster.ClusterDiscoveryType = &xds_cluster.Cluster_Type{Type: xds_cluster.Cluster_EDS}
		//	remoteCluster.EdsClusterConfig = &xds_cluster.Cluster_EdsClusterConfig{EdsConfig: envoy.GetADSConfigSource()}
		//	remoteCluster.LbPolicy = xds_cluster.Cluster_ROUND_ROBIN
		//	*cluster = append(*cluster, remoteCluster)
		//}
	}

	return nil
}

// create one cluster for each pod in the service.
// cluster name of the form "<pod-name>:<port-num>"
func getWSUnicastUpstreamServiceCluster(catalog catalog.MeshCataloger, downstreamIdentity identity.ServiceIdentity, upstreamSvc service.MeshService, cfg configurator.Configurator, cluster *[]*xds_cluster.Cluster) error {
	serviceEndpoints, err := catalog.GetResolvableServiceEndpoints(upstreamSvc)
	if err != nil {
		//log.Info().Msgf("getWSUnicastUpstreamServiceCluster err %+v", err)
		return err
	}

	HTTP2ProtocolOptions, err := envoy.GetHTTP2ProtocolOptions()
	if err != nil {
		log.Info().Msgf("Gethttp2 option failed err %+v", err)
		return err
	}

	//log.Info().Msgf("getWSUnicastUpstreamServiceCluster endpoints %+v", serviceEndpoints)
	// create clusters with pod-names
	for _, endpoint := range serviceEndpoints {
		clusterName := endpoint.PodName + ":" + strconv.Itoa(int(endpoint.Port))
		remoteCluster := &xds_cluster.Cluster{
			Name:                 clusterName,
			ConnectTimeout:       ptypes.DurationProto(clusterConnectTimeout),
			CircuitBreakers: &xds_cluster.CircuitBreakers{
				Thresholds:   makeWSThresholds(),
			},
			TypedExtensionProtocolOptions: HTTP2ProtocolOptions,
		}

		remoteCluster.ClusterDiscoveryType = &xds_cluster.Cluster_Type{Type: xds_cluster.Cluster_EDS}
		remoteCluster.EdsClusterConfig = &xds_cluster.Cluster_EdsClusterConfig{EdsConfig: envoy.GetADSConfigSource()}
		remoteCluster.LbPolicy = xds_cluster.Cluster_ROUND_ROBIN
		*cluster = append(*cluster, remoteCluster)
	}
	return nil
}

