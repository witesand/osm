// Package main implements the main entrypoint for osm-controller and utility routines to
// bootstrap the various internal components of osm-controller.
// osm-controller is the core control plane componenent in OSM responsible for progamming sidecar proxies.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"ws/osm/pkg/certificate"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/openservicemesh/osm/pkg/catalog"
	"github.com/openservicemesh/osm/pkg/certificate/providers"
	"github.com/openservicemesh/osm/pkg/configurator"
	"github.com/openservicemesh/osm/pkg/constants"
	"github.com/openservicemesh/osm/pkg/debugger"
	"github.com/openservicemesh/osm/pkg/endpoint"
	"github.com/openservicemesh/osm/pkg/endpoint/providers/kube"
	"github.com/openservicemesh/osm/pkg/endpoint/providers/remote"
	"github.com/openservicemesh/osm/pkg/envoy/ads"
	"github.com/openservicemesh/osm/pkg/featureflags"
	"github.com/openservicemesh/osm/pkg/health"
	"github.com/openservicemesh/osm/pkg/httpserver"
	"github.com/openservicemesh/osm/pkg/ingress"
	k8s "github.com/openservicemesh/osm/pkg/kubernetes"
	"github.com/openservicemesh/osm/pkg/kubernetes/events"
	"github.com/openservicemesh/osm/pkg/logger"
	"github.com/openservicemesh/osm/pkg/metricsstore"
	"github.com/openservicemesh/osm/pkg/signals"
	"github.com/openservicemesh/osm/pkg/smi"
	"github.com/openservicemesh/osm/pkg/version"
	"github.com/openservicemesh/osm/pkg/witesand"
)

const (
	xdsServerCertificateCommonName = "ads"
)

var (
//<<<<<<< HEAD
//	verbosity            string
//	meshName             string // An ID that uniquely identifies an OSM instance
//	enableRemoteCluster  bool
//	clusterId            string
//	kubeConfigFile       string
//	osmControllerName    string
//	osmNamespace         string
//	webhookConfigName    string
//	caBundleSecretName   string
//	osmConfigMapName     string
//	metricsAddr          string
//	enableLeaderElection bool
//
//	injectorConfig injector.Config
//=======
	verbosity          string
	meshName           string // An ID that uniquely identifies an OSM instance
	kubeConfigFile     string
	osmNamespace       string
	webhookConfigName  string
	caBundleSecretName string
	osmConfigMapName   string

	certProviderKind string

	tresorOptions      providers.TresorOptions
	vaultOptions       providers.VaultOptions
	certManagerOptions providers.CertManagerOptions
//>>>>>>> 3d923b3f2d72006f6cdaad056938c492c364196d

	remoteProvider             *remote.Client
	witesandCatalog            *witesand.WitesandCatalog

	// feature flag options
	optionalFeatures featureflags.OptionalFeatures
	m                *catalog.MeshCatalog

	scheme = runtime.NewScheme()
)

var (
	flags = pflag.NewFlagSet(`osm-controller`, pflag.ExitOnError)
	port  = flags.Int("port", constants.OSMControllerPort, "Aggregated Discovery Service port number.")
	log   = logger.New("osm-controller/main")
)

func init() {
	flags.StringVarP(&verbosity, "verbosity", "v", "info", "Set log verbosity level")
	flags.StringVar(&meshName, "mesh-name", "", "OSM mesh name")
	flags.StringVar(&kubeConfigFile, "kubeconfig", "", "Path to Kubernetes config file.")
	flags.BoolVar(&enableRemoteCluster, "enable-remote-cluster", false, "Enable Remote cluster")
	flags.StringVar(&clusterId, "cluster-id", "master", "Cluster Id")
	flags.StringVar(&osmNamespace, "osm-namespace", "", "Namespace to which OSM belongs to.")
	flags.StringVar(&osmControllerName, "osm-controller-name", "osm-controller", "Service name of osm-controller.")
	flags.StringVar(&webhookConfigName, "webhook-config-name", "", "Name of the MutatingWebhookConfiguration to be configured by osm-controller")
	flags.StringVar(&osmConfigMapName, "osm-configmap-name", "osm-config", "Name of the OSM ConfigMap")

	// Generic certificate manager/provider options
	flags.StringVar(&certProviderKind, "certificate-manager", providers.TresorKind.String(), fmt.Sprintf("Certificate manager, one of [%v]", providers.ValidCertificateProviders))
	flags.StringVar(&caBundleSecretName, "ca-bundle-secret-name", "", "Name of the Kubernetes Secret for the OSM CA bundle")

	// Vault certificate manager/provider options
	flags.StringVar(&vaultOptions.VaultProtocol, "vault-protocol", "http", "Host name of the Hashi Vault")
	flags.StringVar(&vaultOptions.VaultHost, "vault-host", "vault.default.svc.cluster.local", "Host name of the Hashi Vault")
	flags.StringVar(&vaultOptions.VaultToken, "vault-token", "", "Secret token for the the Hashi Vault")
	flags.StringVar(&vaultOptions.VaultRole, "vault-role", "openservicemesh", "Name of the Vault role dedicated to Open Service Mesh")
	flags.IntVar(&vaultOptions.VaultPort, "vault-port", 8200, "Port of the Hashi Vault")

	// Cert-manager certificate manager/provider options
	flags.StringVar(&certManagerOptions.IssuerName, "cert-manager-issuer-name", "osm-ca", "cert-manager issuer name")
	flags.StringVar(&certManagerOptions.IssuerKind, "cert-manager-issuer-kind", "Issuer", "cert-manager issuer kind")
	flags.StringVar(&certManagerOptions.IssuerGroup, "cert-manager-issuer-group", "cert-manager.io", "cert-manager issuer group")

	// feature flags
	flags.BoolVar(&optionalFeatures.WASMStats, "stats-wasm-experimental", false, "Enable a WebAssembly module that generates additional Envoy statistics.")

	_ = clientgoscheme.AddToScheme(scheme)
	_ = v1beta1.AddToScheme(scheme)
}

func main() {
	log.Info().Msgf("Starting osm-controller %s; %s; %s", version.Version, version.GitCommit, version.BuildDate)
	if err := parseFlags(); err != nil {
		log.Fatal().Err(err).Msg("Error parsing cmd line arguments")
	}

	if err := logger.SetLogLevel(verbosity); err != nil {
		log.Fatal().Err(err).Msg("Error setting log level")
	}

	if featureFlagsJSON, err := json.Marshal(featureflags.Features); err != nil {
		log.Error().Err(err).Msgf("Error marshaling feature flags struct: %+v", featureflags.Features)
	} else {
		log.Info().Msgf("Feature flags: %s", string(featureFlagsJSON))
	}

	featureflags.Initialize(optionalFeatures)
	events.GetPubSubInstance() // Just to generate the interface, single routine context

	// Initialize kube config and client
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigFile)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error creating kube config (kubeconfig=%s)", kubeConfigFile)
	}
	kubeClient := kubernetes.NewForConfigOrDie(kubeConfig)

	// Initialize the generic Kubernetes event recorder and associate it with the osm-controller pod resource
	controllerPod, err := getOSMControllerPod(kubeClient)
	if err != nil {
		log.Fatal().Msg("Error fetching osm-controller pod")
	}
	eventRecorder := events.GenericEventRecorder()
	if err := eventRecorder.Initialize(controllerPod, kubeClient, osmNamespace); err != nil {
		log.Fatal().Msg("Error initializing generic event recorder")
	}

	// This ensures CLI parameters (and dependent values) are correct.
	if err := validateCLIParams(); err != nil {
		events.GenericEventRecorder().FatalEvent(err, events.InvalidCLIParameters, "Error validating CLI parameters")
	}

	stop := signals.RegisterExitHandlers()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the default metrics store
	startMetricsStore()

	// This component will be watching the OSM ConfigMap and will make it
	// to the rest of the components.
	cfg := configurator.NewConfigurator(kubernetes.NewForConfigOrDie(kubeConfig), stop, osmNamespace, osmConfigMapName)
	configMap, err := cfg.GetConfigMap()
	if err != nil {
		log.Error().Err(err).Msgf("Error parsing ConfigMap %s", osmConfigMapName)
	}
	log.Info().Msgf("Initial ConfigMap %s: %s", osmConfigMapName, string(configMap))

	kubernetesClient, err := k8s.NewKubernetesController(kubeClient, meshName, stop)
	if err != nil {
		events.GenericEventRecorder().FatalEvent(err, events.InitializationError, "Error creating Kubernetes Controller")
	}

	meshSpec, err := smi.NewMeshSpecClient(kubeConfig, kubeClient, osmNamespace, kubernetesClient, stop)
	if err != nil {
		events.GenericEventRecorder().FatalEvent(err, events.InitializationError, "Error creating MeshSpec")
	}

	certManager, certDebugger, _, err := providers.NewCertificateProvider(kubeClient, kubeConfig, cfg, providers.Kind(certProviderKind), osmNamespace,
		caBundleSecretName, tresorOptions, vaultOptions, certManagerOptions)

	if err != nil {
		events.GenericEventRecorder().FatalEvent(err, events.InvalidCertificateManager,
			"Error fetching certificate manager of kind %s", certProviderKind)
	}

	kubeProvider, err := kube.NewProvider(kubeClient, kubernetesClient, constants.KubeProviderName, cfg)
	if err != nil {
		events.GenericEventRecorder().FatalEvent(err, events.InitializationError, "Error creating Kubernetes endpoints provider")
	}

	endpointsProviders := []endpoint.Provider{kubeProvider}

	log.Info().Msgf("enableRemoteCluster:%t clusterId:%s", enableRemoteCluster, clusterId)

	witesandCatalog = witesand.NewWitesandCatalog(kubeClient, clusterId)
	if err != nil {
		events.GenericEventRecorder().FatalEvent(err, events.InitializationError, "Error creating Witesand catalog")
	}

	if enableRemoteCluster {
		remoteProvider, err = remote.NewProvider(kubeClient, witesandCatalog, clusterId, stop, meshSpec, constants.RemoteProviderName)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to initialize remote provider")
		}
		endpointsProviders = append(endpointsProviders, remoteProvider)
	}

	ingressClient, err := ingress.NewIngressClient(kubeClient, kubernetesClient, stop, cfg)
	if err != nil {
		events.GenericEventRecorder().FatalEvent(err, events.InitializationError, "Error creating Ingress monitor client")
	}

	meshCatalog := catalog.NewMeshCatalog(
		kubernetesClient,
		kubeClient,
		meshSpec,
		certManager,
		ingressClient,
		stop,
		cfg,
		witesandCatalog,
		endpointsProviders...)

//<<<<<<< HEAD
//	// Create the sidecar-injector webhook
//	if err := injector.NewWebhook(injectorConfig, kubeClient, certManager, meshCatalog, kubernetesClient, meshName, osmControllerName, osmNamespace, webhookConfigName, stop, cfg); err != nil {
//		events.GenericEventRecorder().FatalEvent(err, events.InitializationError, "Error creating sidecar injector webhook")
//=======
	// Create the configMap validating webhook
	if err := configurator.NewValidatingWebhook(kubeClient, certManager, osmNamespace, webhookConfigName, stop); err != nil {
		events.GenericEventRecorder().FatalEvent(err, events.InitializationError, "Error creating osm-config validating webhook")
//>>>>>>> 3d923b3f2d72006f6cdaad056938c492c364196d
	}

	m = meshCatalog

	adsCert, err := certManager.IssueCertificate(xdsServerCertificateCommonName, constants.XDSCertificateValidityPeriod)
	if err != nil {
		events.GenericEventRecorder().FatalEvent(err, events.CertificateIssuanceFailure, "Error issuing XDS certificate to ADS server")
	}

	// Create and start the ADS gRPC service
	xdsServer := ads.NewADSServer(meshCatalog, cfg.IsDebugServerEnabled(), osmNamespace, cfg, certManager)
	if err := xdsServer.Start(ctx, cancel, *port, adsCert); err != nil {
		events.GenericEventRecorder().FatalEvent(err, events.InitializationError, "Error initializing ADS server")
	}

	// Initialize OSM's http service server
	httpServer := httpserver.NewHTTPServer(constants.OSMHTTPServerPort)

	// Health/Liveness probes
	funcProbes := []health.Probes{xdsServer}
	httpServer.AddHandlers(map[string]http.Handler{
		"/health/ready": health.ReadinessHandler(funcProbes, getHTTPHealthProbes()),
		"/health/alive": health.LivenessHandler(funcProbes, getHTTPHealthProbes()),
	})
	// Metrics
	httpServer.AddHandler("/metrics", metricsstore.DefaultMetricsStore.Handler())
	// Version
	httpServer.AddHandler("/version", version.GetVersionHandler())

	// Start HTTP server
	err = httpServer.Start()
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to start OSM metrics/probes HTTP server")
	}

	// Create DebugServer and start its config event listener.
	// Listener takes care to start and stop the debug server as appropriate
	debugConfig := debugger.NewDebugConfig(certDebugger, xdsServer, meshCatalog, kubeConfig, kubeClient, cfg, kubernetesClient)
	debugConfig.StartDebugServerConfigListener()

	// Wait for exit handler signal
	<-stop
	log.Info().Msgf("Stopping osm-controller %s; %s; %s", version.Version, version.GitCommit, version.BuildDate)
}

// Start the metric store, register the metrics OSM will expose
func startMetricsStore() {
	metricsstore.DefaultMetricsStore.Start(
		metricsstore.DefaultMetricsStore.K8sAPIEventCounter,
		metricsstore.DefaultMetricsStore.K8sMonitoredNamespaceCount,
		metricsstore.DefaultMetricsStore.K8sMeshPodCount,
		metricsstore.DefaultMetricsStore.ProxyConnectCount,
		metricsstore.DefaultMetricsStore.ProxyConfigUpdateTime,
		metricsstore.DefaultMetricsStore.CertIssuedCount,
		metricsstore.DefaultMetricsStore.CertIssuedTime,
	)
}

// getHTTPHealthProbes returns the HTTP health probes served by OSM controller
func getHTTPHealthProbes() []health.HTTPProbe {
	// Example:
	// return []health.HTTPProbe{
	// 	{
	// 		URL: "https://127.0.0.1:<local-port>",
	// 		Protocol: health.ProtocolHTTPS,
	// 	},
	// }
	return nil
}

func parseFlags() error {
	if err := flags.Parse(os.Args); err != nil {
		return err
	}
	_ = flag.CommandLine.Parse([]string{})
	return nil
}

func joinURL(baseURL string, paths ...string) string {
	p := path.Join(paths...)
	return fmt.Sprintf("%s/%s", strings.TrimRight(baseURL, "/"), strings.TrimLeft(p, "/"))
}

// getOSMControllerPod returns the osm-controller pod.
// The pod name is inferred from the 'CONTROLLER_POD_NAME' env variable which is set during deployment.
func getOSMControllerPod(kubeClient kubernetes.Interface) (*corev1.Pod, error) {
	podName := os.Getenv("CONTROLLER_POD_NAME")
	if podName == "" {
		return nil, errors.New("CONTROLLER_POD_NAME env variable cannot be empty")
	}

	pod, err := kubeClient.CoreV1().Pods(osmNamespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		log.Error().Err(err).Msgf("Error retrieving osm-controller pod %s", podName)
		return nil, err
	}

	return pod, nil
}
//<<<<<<< HEAD
//
////Setting up k8s controller manager to reconcile OSM resources
//func createControllerManagerForOSMResources(certManager certificate.Manager) error {
//	log.Info().Msg("Setting up controller manager to reconcile OSM resources")
//	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
//		Scheme:             scheme,
//		MetricsBindAddress: metricsAddr,
//		LeaderElection:     enableLeaderElection,
//		Namespace:          osmNamespace,
//	})
//	if err != nil {
//		log.Error().Err(err).Msg("Error starting up controller manager")
//		return err
//	}
//
//	log.Info().Msg("Successfully setup controller for resource reconciliation")
//	log.Info().Msg("Setting up mutatingWebhookConfiguration reconciler")
//
//	// controller logic is implemented by reconciler
//	// Adding a reconciler for OSM's mutatingwehbookconfiguration
//	if err = (&reconciler.MutatingWebhookConfigrationReconciler{
//		Client:       mgr.GetClient(),
//		Scheme:       mgr.GetScheme(),
//		OsmWebhook:   fmt.Sprintf("osm-webhook-%s", meshName),
//		OsmNamespace: osmNamespace,
//		CertManager:  certManager,
//
//		OsmControllerName: osmControllerName,
//	}).SetupWithManager(mgr); err != nil {
//		log.Error().Err(err).Msg("Error creating reconcile controller for MutatingWebhookConfiguration")
//		return err
//	}
//
//	log.Info().Msg("starting manager")
//	go func() {
//		// mgr.Start() below will block until stopped
//		// See: https://github.com/kubernetes-sigs/controller-runtime/blob/release-0.6/pkg/manager/internal.go#L507-L514
//		if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
//			log.Error().Err(err).Msg("problem running manager")
//		}
//	}()
//
//	return nil
//}
//=======
//>>>>>>> 3d923b3f2d72006f6cdaad056938c492c364196d
