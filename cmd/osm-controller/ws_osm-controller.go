package main

import (
	"github.com/openservicemesh/osm/pkg/constants"
	"github.com/openservicemesh/osm/pkg/endpoint"
	"github.com/openservicemesh/osm/pkg/endpoint/providers/remote"
	"github.com/openservicemesh/osm/pkg/smi"
	"github.com/openservicemesh/osm/pkg/witesand"
	clientset "k8s.io/client-go/kubernetes"
)

var (
	enableRemoteCluster bool
	clusterId           string
	osmControllerName   string
	remoteProvider      *remote.Client
	witesandCatalog     *witesand.WitesandCatalog
)

func wsinit() {
	flags.BoolVar(&enableRemoteCluster, "enable-remote-cluster", false, "Enable Remote cluster")
	flags.StringVar(&clusterId, "cluster-id", "master", "Cluster Id")
	flags.StringVar(&osmControllerName, "osm-controller-name", "osm-controller", "Service name of osm-controller.")
}

func addWSCatalog(kubeClient *clientset.Clientset)  {
	log.Info().Msgf("witesand flags enableRemoteCluster:%v clusterId:%s osmcontrollername=%s", enableRemoteCluster, clusterId, osmControllerName)
	witesandCatalog = witesand.NewWitesandCatalog(kubeClient, clusterId)
}

func addWSRemoteCluster(kubeClient *clientset.Clientset, stop chan struct{}, meshSpec smi.MeshSpec, endpointsProviders *[]endpoint.Provider) error {
	log.Info().Msgf("enableRemoteCluster:%t clusterId:%s", enableRemoteCluster, clusterId)
	if enableRemoteCluster {
		remoteProvider, err := remote.NewProvider(kubeClient, witesandCatalog, clusterId, stop, meshSpec, constants.RemoteProviderName)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to initialize remote provider")
			return err
		}
		*endpointsProviders = append(*endpointsProviders, remoteProvider)
	}
	return nil
}

