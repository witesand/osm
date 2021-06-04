package configurator

import (
	"encoding/json"
	"fmt"
//<<<<<<< HEAD
//	"net"
//	"sort"
//=======
//>>>>>>> 3d923b3f2d72006f6cdaad056938c492c364196d
	"strings"
	"time"

	"github.com/openservicemesh/osm/pkg/constants"
)

const (
	// defaultServiceCertValidityDuration is the default validity duration for service certificates
	defaultServiceCertValidityDuration = 24 * time.Hour
)

// The functions in this file implement the configurator.Configurator interface

// GetOSMNamespace returns the namespace in which the OSM controller pod resides.
func (c *Client) GetOSMNamespace() string {
	return c.osmNamespace
}

func marshalConfigToJSON(config *osmConfig) ([]byte, error) {
	return json.MarshalIndent(config, "", "    ")
}

// GetConfigMap returns the ConfigMap in pretty JSON.
func (c *Client) GetConfigMap() ([]byte, error) {
	cm, err := marshalConfigToJSON(c.getConfigMap())
	if err != nil {
		log.Error().Err(err).Msgf("Error marshaling ConfigMap %s: %+v", c.getConfigMapCacheKey(), c.getConfigMap())
		return nil, err
	}
	return cm, nil
}

// IsPermissiveTrafficPolicyMode tells us whether the OSM Control Plane is in permissive mode,
// where all existing traffic is allowed to flow as it is,
// or it is in SMI Spec mode, in which only traffic between source/destinations
// referenced in SMI policies is allowed.
func (c *Client) IsPermissiveTrafficPolicyMode() bool {
	return c.getConfigMap().PermissiveTrafficPolicyMode
}

// IsEgressEnabled determines whether egress is globally enabled in the mesh or not.
func (c *Client) IsEgressEnabled() bool {
	return c.getConfigMap().Egress
}

// IsDebugServerEnabled determines whether osm debug HTTP server is enabled
func (c *Client) IsDebugServerEnabled() bool {
	return c.getConfigMap().EnableDebugServer
}

// IsPrometheusScrapingEnabled determines whether Prometheus is enabled for scraping metrics
func (c *Client) IsPrometheusScrapingEnabled() bool {
	return c.getConfigMap().PrometheusScraping
}

// IsTracingEnabled returns whether tracing is enabled
func (c *Client) IsTracingEnabled() bool {
	return c.getConfigMap().TracingEnable
}

// GetTracingHost is the host to which we send tracing spans
func (c *Client) GetTracingHost() string {
	tracingAddress := c.getConfigMap().TracingAddress
	if tracingAddress != "" {
		return tracingAddress
	}
	return fmt.Sprintf("%s.%s.svc.cluster.local", constants.DefaultTracingHost, c.GetOSMNamespace())
}

// GetTracingPort returns the tracing listener port
func (c *Client) GetTracingPort() uint32 {
	tracingPort := c.getConfigMap().TracingPort
	if tracingPort != 0 {
		return uint32(tracingPort)
	}
	return constants.DefaultTracingPort
}

// GetTracingEndpoint returns the listener's collector endpoint
func (c *Client) GetTracingEndpoint() string {
	tracingEndpoint := c.getConfigMap().TracingEndpoint
	if tracingEndpoint != "" {
		return tracingEndpoint
	}
	return constants.DefaultTracingEndpoint
}

// GetMeshCIDRRanges returns a list of mesh CIDR ranges
func (c *Client) GetMeshCIDRRanges() []string {
	noSpaces := strings.ReplaceAll(c.getConfigMap().MeshCIDRRanges, " ", ",")
	commaSeparatedCIDRs := strings.Split(noSpaces, ",")

	cidrSet := make(map[string]interface{})
	for _, cidr := range commaSeparatedCIDRs {
		trimmedCIDR := strings.Trim(cidr, " ")
		if len(trimmedCIDR) == 0 {
			continue
		}

		_, _, err := net.ParseCIDR(trimmedCIDR)
		if err != nil {
			log.Error().Err(err).Msgf("Found incorrectly formatted in-mesh CIDR %s from ConfigMap %s/%s; Skipping CIDR", trimmedCIDR, c.osmNamespace, c.osmConfigMapName)
			continue
		}

		cidrSet[trimmedCIDR] = nil
	}

	var cidrs []string
	for cidr := range cidrSet {
		cidrs = append(cidrs, cidr)
	}

	sort.Strings(cidrs)

	return cidrs
}

// UseHTTPSIngress determines whether traffic between ingress and backend pods should use HTTPS protocol
func (c *Client) UseHTTPSIngress() bool {
	return c.getConfigMap().UseHTTPSIngress
}

// GetEnvoyLogLevel returns the envoy log level
func (c *Client) GetEnvoyLogLevel() string {
	logLevel := c.getConfigMap().EnvoyLogLevel
	if logLevel != "" {
		return logLevel
	}
	return constants.DefaultEnvoyLogLevel
}

// GetServiceCertValidityPeriod returns the validity duration for service certificates, and a default in case of invalid duration
func (c *Client) GetServiceCertValidityPeriod() time.Duration {
	durationStr := c.getConfigMap().ServiceCertValidityDuration
	validityDuration, err := time.ParseDuration(durationStr)
	if err != nil {
		log.Error().Err(err).Msgf("Error parsing service certificate validity duration %s=%s", serviceCertValidityDurationKey, durationStr)
		return defaultServiceCertValidityDuration
	}

	return validityDuration
}

// GetOutboundIPRangeExclusionList returns the list of IP ranges of the form x.x.x.x/y to exclude from outbound sidecar interception
func (c *Client) GetOutboundIPRangeExclusionList() []string {
	ipRangesStr := c.getConfigMap().OutboundIPRangeExclusionList
	if ipRangesStr == "" {
		return nil
	}

	exclusionList := strings.Split(ipRangesStr, ",")
	for i := range exclusionList {
		exclusionList[i] = strings.TrimSpace(exclusionList[i])
	}

	return exclusionList
}

// IsPrivilegedInitContainer returns whether init containers should be privileged
func (c *Client) IsPrivilegedInitContainer() bool {
	return c.getConfigMap().EnablePrivilegedInitContainer
}

// GetConfigResyncInterval returns the duration for resync interval.
// If error or non-parsable value, returns 0 duration
func (c *Client) GetConfigResyncInterval() time.Duration {
	resyncDuration := c.getConfigMap().ConfigResyncInterval
	duration, err := time.ParseDuration(resyncDuration)
	if err != nil {
		log.Debug().Err(err).Msgf("Error parsing config resync interval: %s", duration)
		return time.Duration(0)
	}
	return duration
}

// GetInboundExternalAuthConfig returns the External Authentication configuration for incoming traffic, if any
func (c *Client) GetInboundExternalAuthConfig() ExternAuthConfig {
	extAuthRet := ExternAuthConfig{}
	cfgMap := c.getConfigMap()

	extAuthRet.Enable = cfgMap.InboundExternAuthzEnable
	extAuthRet.Address = cfgMap.InboundExternAuthzAddress
	extAuthRet.Port = uint16(cfgMap.InboundExternAuthzPort)
	extAuthRet.StatPrefix = cfgMap.InboundExternAuthzStatPrefix
	extAuthRet.FailureModeAllow = cfgMap.InboundExternAuthzFailureModeAllow

	duration, err := time.ParseDuration(cfgMap.InboundExternAuthzTimeout)
	if err != nil {
		log.Debug().Err(err).Msgf("ExternAuthzTimeout: Not a valid duration %s. defaulting to 1s.", duration)
		duration = 1 * time.Second
	}
	extAuthRet.AuthzTimeout = duration

	return extAuthRet
}
