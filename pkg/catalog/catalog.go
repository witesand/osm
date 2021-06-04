package catalog

import (
	"k8s.io/client-go/kubernetes"
	"time"
	"ws/osm/pkg/announcements"
	"ws/osm/pkg/kubernetes/events"
	"ws/osm/pkg/witesand"

	"github.com/openservicemesh/osm/pkg/certificate"
	"github.com/openservicemesh/osm/pkg/configurator"
	"github.com/openservicemesh/osm/pkg/endpoint"
	"github.com/openservicemesh/osm/pkg/ingress"
	k8s "github.com/openservicemesh/osm/pkg/kubernetes"
	"github.com/openservicemesh/osm/pkg/smi"
//<<<<<<< HEAD
//	"github.com/openservicemesh/osm/pkg/witesand"
//)
//
//const (
//	// this is catalog's tick rate for ticker, which triggers global proxy updates
//	// 0 disables the ticker
//	updateAtLeastEvery = 0 * time.Second
//=======
	"github.com/openservicemesh/osm/pkg/ticker"
//>>>>>>> 3d923b3f2d72006f6cdaad056938c492c364196d
)

// NewMeshCatalog creates a new service catalog
func NewMeshCatalog(kubeController k8s.Controller, kubeClient kubernetes.Interface, meshSpec smi.MeshSpec, certManager certificate.Manager, ingressMonitor ingress.Monitor, stop <-chan struct{}, cfg configurator.Configurator, wc *witesand.WitesandCatalog, endpointsProviders ...endpoint.Provider) *MeshCatalog {
	log.Info().Msg("Create a new Service MeshCatalog.")
	mc := MeshCatalog{
		endpointsProviders: endpointsProviders,
		meshSpec:           meshSpec,
		certManager:        certManager,
		ingressMonitor:     ingressMonitor,
		configurator:       cfg,

		// Kubernetes needed to determine what Services a pod that connects to XDS belongs to.
		// In multicluster scenarios this would be a map of cluster ID to Kubernetes client.
		// The certificate itself would contain the cluster ID making it easy to lookup the client in this map.
		kubeClient:     kubeClient,
		kubeController: kubeController,

		witesandCatalog: wc,
	}

	// Run release certificate handler, which listens to podDelete events
	mc.releaseCertificateHandler()

	mc.witesandHttpServerAndClient()

	go mc.dispatcher()
	ticker.InitTicker(cfg)

	return &mc
}

// GetSMISpec returns a MeshCatalog's SMI Spec
func (mc *MeshCatalog) GetSMISpec() smi.MeshSpec {
	return mc.meshSpec
}
//<<<<<<< HEAD
//
//func (mc *MeshCatalog) GetWitesandCataloger() witesand.WitesandCataloger {
//	return mc.witesandCatalog
//}
//
//func (mc *MeshCatalog) getAnnouncementChannels() []announcementChannel {
//	ticking := make(chan announcements.Announcement)
//	announcementChannels := []announcementChannel{
//		{"MeshSpec", mc.meshSpec.GetAnnouncementsChannel()},
//		{"CertManager", mc.certManager.GetAnnouncementsChannel()},
//		{"IngressMonitor", mc.ingressMonitor.GetAnnouncementsChannel()},
//		{"Ticker", ticking},
//		{"Services", mc.kubeController.GetAnnouncementsChannel(k8s.Services)},
//	}
//
//	// There could be many Endpoint Providers - iterate over all of them!
//	for _, ep := range mc.endpointsProviders {
//		annCh := announcementChannel{ep.GetID(), ep.GetAnnouncementsChannel()}
//		announcementChannels = append(announcementChannels, annCh)
//	}
//
//	if updateAtLeastEvery > 0 {
//		go func() {
//			ticker := time.NewTicker(updateAtLeastEvery)
//			for {
//				<-ticker.C
//				events.GetPubSubInstance().Publish(events.PubSubMessage{
//					AnnouncementType: announcements.ScheduleProxyBroadcast,
//					NewObj:           nil,
//					OldObj:           nil,
//				})
//			}
//		}()
//	}
//
//	return announcementChannels
//}
//
//func (mc *MeshCatalog) GetProvider(ident string) endpoint.Provider {
//	for _, ep := range mc.endpointsProviders {
//		if ep.GetID() == ident {
//			return ep
//		}
//	}
//	return nil
//}
//=======
//>>>>>>> 3d923b3f2d72006f6cdaad056938c492c364196d
