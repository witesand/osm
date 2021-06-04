// Package injector implements OSM's automatic sidecar injection facility. The sidecar injector's mutating webhook
// admission controller intercepts pod creation requests to mutate the pod spec to inject the sidecar proxy.
package injector

import (
	mapset "github.com/deckarep/golang-set"
	"k8s.io/client-go/kubernetes"

	"github.com/openservicemesh/osm/pkg/certificate"
	"github.com/openservicemesh/osm/pkg/configurator"
	k8s "github.com/openservicemesh/osm/pkg/kubernetes"
	"github.com/openservicemesh/osm/pkg/logger"
)

const (
	envoyBootstrapConfigVolume = "envoy-bootstrap-config-volume"
)

var log = logger.New("sidecar-injector")

// mutatingWebhook is the type used to represent the webhook for sidecar injection
type mutatingWebhook struct {
	config         Config
	kubeClient     kubernetes.Interface
	certManager    certificate.Manager
	kubeController k8s.Controller
	osmNamespace   string
	meshName       string
	cert           certificate.Certificater
	configurator   configurator.Configurator

	// WITESAND OSM controller name can be different for each cluster
	osmControllerName string

	nonInjectNamespaces mapset.Set
}

// Config is the type used to represent the config options for the sidecar injection
type Config struct {
	// ListenPort defines the port on which the sidecar injector listens
	ListenPort int

	InitContainerImage string

	SidecarImage string
}

//<<<<<<< HEAD
//// JSONPatchOperation defines a Kubernetes JSON Patch operation
//type JSONPatchOperation struct {
//	Op    string      `json:"op"`
//	Path  string      `json:"path"`
//	Value interface{} `json:"value,omitempty"`
//}
//
//// InitContainer is the type used to represent information about the init container
//type InitContainer struct {
//	Name  string
//	Image string
//	CIDR1 string
//	CIDR2 string
//}
//
//// EnvoySidecarData is the type used to represent information about the Envoy sidecar
//type EnvoySidecarData struct {
//	Name           string
//	Image          string
//	EnvoyNodeID    string
//	EnvoyClusterID string
//}
//
//=======
//>>>>>>> 3d923b3f2d72006f6cdaad056938c492c364196d
// Context needed to compose the Envoy bootstrap YAML.
type envoyBootstrapConfigMeta struct {
	EnvoyAdminPort int
	XDSClusterName string
	RootCert       string
	Cert           string
	Key            string

	// Host and port of the Envoy xDS server
	XDSHost string
	XDSPort int

	// The bootstrap Envoy config will be affected by the liveness, readiness, startup probes set on
	// the pod this Envoy is fronting.
	OriginalHealthProbes healthProbes
}
