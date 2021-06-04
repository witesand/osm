package cds

import (
	mapset "github.com/deckarep/golang-set"
	xds_cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	xds_discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/golang/protobuf/ptypes"

	"github.com/openservicemesh/osm/pkg/catalog"
	"github.com/openservicemesh/osm/pkg/certificate"
	"github.com/openservicemesh/osm/pkg/configurator"
	"github.com/openservicemesh/osm/pkg/envoy"
)

// NewResponse creates a new Cluster Discovery Response.
func NewResponse(meshCatalog catalog.MeshCataloger, proxy *envoy.Proxy, _ *xds_discovery.DiscoveryRequest, cfg configurator.Configurator, _ certificate.Manager) (*xds_discovery.DiscoveryResponse, error) {
	svcList, err := meshCatalog.GetServicesFromEnvoyCertificate(proxy.GetCertificateCommonName())
	if err != nil {
		log.Error().Err(err).Msgf("Error looking up MeshService for Envoy with SerialNumber=%s on Pod with UID=%s", proxy.GetCertificateSerialNumber(), proxy.GetPodUID())
		return nil, err
	}

	var clusters []*xds_cluster.Cluster

	proxyIdentity, err := catalog.GetServiceAccountFromProxyCertificate(proxy.GetCertificateCommonName())
	if err != nil {
		log.Error().Err(err).Msgf("Error looking up proxy identity for proxy with SerialNumber=%s on Pod with UID=%s",
			proxy.GetCertificateSerialNumber(), proxy.GetPodUID())
		return nil, err
	}
	log.Debug().Msgf("svc:%s url:%s outboundServices:%+v", proxyServiceName, resp.TypeUrl, outboundServices)

	// Build remote clusters based on allowed outbound services
<<<<<<< HEAD
	for _, dstService := range outboundServices {
		if _, found := clusterFactories[dstService.String()]; found {
			// Guard against duplicates
			continue
		}

		if catalog.GetWitesandCataloger().IsWSEdgePodService(dstService) {
			getWSEdgePodUpstreamServiceCluster(catalog, dstService, proxyServiceName.GetMeshServicePort(), cfg, clusterFactories)
			continue
		} else if catalog.GetWitesandCataloger().IsWSUnicastService(dstService.Name) {
			getWSUnicastUpstreamServiceCluster(catalog, dstService, proxyServiceName.GetMeshServicePort(), cfg, clusterFactories)
			// fall thru to generate anycast cluster
		}

		remoteCluster, err := getUpstreamServiceCluster(dstService, proxyServiceName.GetMeshServicePort(), cfg)

=======
	for _, dstService := range meshCatalog.ListAllowedOutboundServicesForIdentity(proxyIdentity) {
		cluster, err := getUpstreamServiceCluster(proxyIdentity, dstService, cfg)
>>>>>>> 3d923b3f2d72006f6cdaad056938c492c364196d
		if err != nil {
			log.Error().Err(err).Msgf("Failed to construct service cluster for service %s for proxy with XDS Certificate SerialNumber=%s on Pod with UID=%s",
				dstService.Name, proxy.GetCertificateSerialNumber(), proxy.GetPodUID())
			return nil, err
		}

<<<<<<< HEAD
		if featureflags.IsBackpressureEnabled() {
			enableBackpressure(catalog, remoteCluster, dstService.GetMeshService())
		}
		log.Debug().Msgf("remoteName:%s, remoteCluster:%+v", remoteCluster.Name, remoteCluster)

		clusterFactories[remoteCluster.Name] = remoteCluster
	}

	// Create a local cluster for the service.
	// The local cluster will be used for incoming traffic.
	localClusters, err := getLocalServiceCluster(catalog, proxyServiceName)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to get local cluster config for proxy %s", proxyServiceName)
		return nil, err
	}

	count := 0
	for _, localCluster := range localClusters {
		clusterFactories[localCluster.Name] = localCluster
		log.Debug().Msgf("local:%s localCluster:%+v", localCluster.Name, localCluster)
		count++
	}

	if cfg.IsEgressEnabled() {
		// Add a pass-through cluster for egress
		passthroughCluster := getOutboundPassthroughCluster()
		clusterFactories[passthroughCluster.Name] = passthroughCluster
	}

	for _, cluster := range clusterFactories {
		log.Debug().Msgf("Proxy service %s constructed ClusterConfiguration: %+v ", proxyServiceName, cluster)
		marshalledClusters, err := ptypes.MarshalAny(cluster)
=======
		clusters = append(clusters, cluster)
	}

	// Create a local cluster for each service behind the proxy.
	// The local cluster will be used to handle incoming traffic.
	for _, proxyService := range svcList {
		localClusterName := envoy.GetLocalClusterNameForService(proxyService)
		localCluster, err := getLocalServiceCluster(meshCatalog, proxyService, localClusterName)
>>>>>>> 3d923b3f2d72006f6cdaad056938c492c364196d
		if err != nil {
			log.Error().Err(err).Msgf("Failed to get local cluster config for proxy %s", proxyService)
			return nil, err
		}
		clusters = append(clusters, localCluster)
	}

	// Add an outbound passthrough cluster for egress
	if cfg.IsEgressEnabled() {
		clusters = append(clusters, getOutboundPassthroughCluster())
	}

	// Add an inbound prometheus cluster (from Prometheus to localhost)
	if cfg.IsPrometheusScrapingEnabled() {
		clusters = append(clusters, getPrometheusCluster())
	}

	// Add an outbound tracing cluster (from localhost to tracing sink)
	if cfg.IsTracingEnabled() {
		clusters = append(clusters, getTracingCluster(cfg))
	}

	resp := &xds_discovery.DiscoveryResponse{
		TypeUrl: string(envoy.TypeCDS),
	}

	alreadyAdded := mapset.NewSet()
	for _, cluster := range clusters {
		if alreadyAdded.Contains(cluster.Name) {
			log.Error().Msgf("Found duplicate clusters with name %s; Duplicate will not be sent to Envoy with XDS Certificate SerialNumber=%s on Pod with UID=%s",
				cluster.Name, proxy.GetCertificateSerialNumber(), proxy.GetPodUID())
			continue
		}
		alreadyAdded.Add(cluster.Name)
		marshalledClusters, err := ptypes.MarshalAny(cluster)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to marshal cluster %s for Envoy with XDS Certificate SerialNumber=%s on Pod with UID=%s",
				cluster.Name, proxy.GetCertificateSerialNumber(), proxy.GetPodUID())
			return nil, err
		}
		resp.Resources = append(resp.Resources, marshalledClusters)
	}
	//log.Debug().Msgf("Proxy service %s CDS resp: %+v ", proxyServiceName, resp.Resources)

	return resp, nil
}
