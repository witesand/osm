package injector

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/openservicemesh/osm/pkg/constants"
)

// iptablesRedirectionChains is the list of iptables chains created for traffic redirection via the proxy sidecar
var iptablesRedirectionChains = []string{
	// Chain to intercept inbound traffic
	"iptables -t nat -N PROXY_INBOUND",

	// Chain to redirect inbound traffic to the proxy
	"iptables -t nat -N PROXY_IN_REDIRECT",

	// Chain to intercept outbound traffic
	"iptables -t nat -N PROXY_OUTPUT",

	// Chain to redirect outbound traffic to the proxy
	"iptables -t nat -N PROXY_REDIRECT",
}

// iptablesOutboundStaticRules is the list of iptables rules related to outbound traffic interception and redirection
var iptablesOutboundStaticRules = []string{

	//#WITEDAND---START
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 22 -j ACCEPT"),  // # ssh port
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 49 -j ACCEPT"),  // # tacacs
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 69 -j ACCEPT"),  // # tftp
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 88 -j ACCEPT"),  // # nacmode
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 139 -j ACCEPT"), // # samba
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 389 -j ACCEPT"), // # radius port
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 443 -j ACCEPT"), // # aruba
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 445 -j ACCEPT"), // # samba
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 587 -j ACCEPT"), //# email port
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 636 -j ACCEPT"), // # ldaps
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 830 -j ACCEPT"), // # netconf
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 2579 -j ACCEPT"), // # kine
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 2500 -j ACCEPT"), // # osm-rest
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 4343 -j ACCEPT"), // # aruba
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 5000 -j ACCEPT"), // # devicedb
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 5432 -j ACCEPT"), // # postgres
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 5556 -j ACCEPT"), // # wsdex
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 5557 -j ACCEPT"), // # wsdex
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 7201 -j ACCEPT"), // # m3db/metricsd
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 7203 -j ACCEPT"), // # m3db/metricsd
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 7301 -j ACCEPT"), // # m3db/metricsd
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 8005 -j ACCEPT"), // # aws metricsd
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 8080 -j ACCEPT"), // # presto
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 8081 -j ACCEPT"), // # apiserver
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 8100:8110 -j ACCEPT"), // # proxyd
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 8200 -j ACCEPT"), // # valult
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 8443 -j ACCEPT"), // # apiserver
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 9000:9004 -j ACCEPT"), // # m3db/metricsd
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 9053 -j ACCEPT"), // # waves
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 9063:9064 -j ACCEPT"), // # alertdispatch
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 9158 -j ACCEPT"), // # alertruled
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 9073 -j ACCEPT"), // # identityd
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 9083 -j ACCEPT"), // # hive
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 9085 -j ACCEPT"), // # filed
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 9092 -j ACCEPT"), // # kafka
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 9131 -j ACCEPT"), // # logd-rest
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 9097 -j ACCEPT"), // # endpointd
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 9122 -j ACCEPT"), // # metricsd
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 9126 -j ACCEPT"), // # nlp rest
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 9128 -j ACCEPT"), // # historyd
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 9067 -j ACCEPT"), // # rcad
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 9200 -j ACCEPT"), // # elastic
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 9300 -j ACCEPT"), // # elastic
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 10000 -j ACCEPT"), // # radiusconfd
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 10080 -j ACCEPT"), // # byod
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 10500 -j ACCEPT"), // # deviced/proxylistener
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport 32443 -j ACCEPT"), // # sslport/apiserver
	//#WITEDAND---END


	// Redirects outbound TCP traffic hitting PROXY_REDIRECT chain to Envoy's outbound listener port
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp -j REDIRECT --to-port %d", constants.EnvoyOutboundListenerPort),

	// Traffic to the Proxy Admin port flows to the Proxy -- not redirected
	fmt.Sprintf("iptables -t nat -A PROXY_REDIRECT -p tcp --dport %d -j ACCEPT", constants.EnvoyAdminPort),



	// For outbound TCP traffic jump from OUTPUT chain to PROXY_OUTPUT chain
	"iptables -t nat -A OUTPUT -p tcp -j PROXY_OUTPUT",

	// TODO(#1266): Redirect app back calls to itself using PROXY_UID

	// Don't redirect Envoy traffic back to itself, return it to the next chain for processing
	fmt.Sprintf("iptables -t nat -A PROXY_OUTPUT -m owner --uid-owner %d -j RETURN", constants.EnvoyUID),

	// Skip localhost traffic, doesn't need to be routed via the proxy
	"iptables -t nat -A PROXY_OUTPUT -d 127.0.0.1/32 -j RETURN",

	// Redirect remaining outbound traffic to Envoy
	"iptables -t nat -A PROXY_OUTPUT -j PROXY_REDIRECT",
}

// iptablesInboundStaticRules is the list of iptables rules related to inbound traffic interception and redirection
var iptablesInboundStaticRules = []string{
	// Redirects inbound TCP traffic hitting the PROXY_IN_REDIRECT chain to Envoy's inbound listener port
	fmt.Sprintf("iptables -t nat -A PROXY_IN_REDIRECT -p tcp -j REDIRECT --to-port %d", constants.EnvoyInboundListenerPort),

	// For inbound traffic jump from PREROUTING chain to PROXY_INBOUND chain
	"iptables -t nat -A PREROUTING -p tcp -j PROXY_INBOUND",

	// Skip metrics query traffic being directed to Envoy's inbound prometheus listener port
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport %d -j RETURN", constants.EnvoyPrometheusInboundListenerPort),
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 5432 -j RETURN"),
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 49 -j RETURN"), //  # tacacs
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 69 -j RETURN"), //  # tftp
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 88 -j RETURN"), //  # nacmode
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 139 -j RETURN"), //  # samba
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 389 -j RETURN"), //  # radius
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 443 -j RETURN"), //  # aruba
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 445 -j RETURN"), //  # samba
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 587 -j RETURN"), //  # email
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 636 -j RETURN"), //  # ldaps
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 830 -j RETURN"), //  # netconf
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 2579 -j RETURN"), //  # kine
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 2500 -j RETURN"), //  # osm-rest
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 4343 -j RETURN"), //  # aruba
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 5000 -j RETURN"), //  # devicedb
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 5432 -j RETURN"), //  # postgres
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 5556 -j RETURN"), //  # wsdex
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 5557 -j RETURN"), //  # wsdex
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 7201 -j RETURN"), //  # m3db/metricsd
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 7203 -j RETURN"), //  # m3db/metricsd
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 7301 -j RETURN"), //  # m3db/metricsd
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 8005 -j RETURN"), //  # aws metricsd
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 8080 -j RETURN"), //  # presto
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 8081 -j RETURN"), //  # apiserver
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 8100:8110 -j RETURN"), //  # proxyd
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 8200 -j RETURN"), //  # valult
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 8443 -j RETURN"), //  # apiserver
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 9000:9004 -j RETURN"), //  # m3db/metricsd
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 9053 -j RETURN"), //  # waves
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 9063:9064 -j RETURN"), //  # alertdispatch
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 9158 -j RETURN"), //  # alertruled
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 9073 -j RETURN"), //  # identityd
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 9083 -j RETURN"), //  # hive
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 9085 -j RETURN"), //  # filed
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 9092 -j RETURN"), //  # kafka
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 9131 -j RETURN"), //  # logd-rest
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 9097 -j RETURN"), //  # endpointd
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 9122 -j RETURN"), //  # metricsd
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 9126 -j RETURN"), //  # nlp rest
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 9128 -j RETURN"), //  # historyd
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 9067 -j RETURN"), //  # rcad
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 9200 -j RETURN"), //  # elastic
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 9300 -j RETURN"), //  # elastic
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 10000 -j RETURN"), // # radiusconfd
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 10080 -j RETURN"), // # byod
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 10500 -j RETURN"), // # deviced/proxylistener
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport 32443 -j RETURN"), // # sslpoort/apiserver

	// Skip inbound health probes; These ports will be explicitly handled by listeners configured on the
	// Envoy proxy IF any health probes have been configured in the Pod Spec.
	// TODO(draychev): Do not add these if no health probes have been defined (https://github.com/openservicemesh/osm/issues/2243)
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport %d -j RETURN", livenessProbePort),
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport %d -j RETURN", readinessProbePort),
	fmt.Sprintf("iptables -t nat -A PROXY_INBOUND -p tcp --dport %d -j RETURN", startupProbePort),

	// Redirect remaining inbound traffic to Envoy
	"iptables -t nat -A PROXY_INBOUND -p tcp -j PROXY_IN_REDIRECT",
}

// generateIptablesCommands generates a list of iptables commands to set up sidecar interception and redirection
func generateIptablesCommands(outboundIPRangeExclusionList []string, outboundPortExclusionList []int, inboundPortExclusionList []int) []string {
	var cmd []string

	// 1. Create redirection chains
	cmd = append(cmd, iptablesRedirectionChains...)

	// 2. Create outbound rules
	cmd = append(cmd, iptablesOutboundStaticRules...)

	// 3. Create inbound rules
	cmd = append(cmd, iptablesInboundStaticRules...)

	// 4. Create dynamic outbound ip ranges exclusion rules
	for _, cidr := range outboundIPRangeExclusionList {
		// *Note: it is important to use the insert option '-I' instead of the append option '-A' to ensure the exclusion
		// rules take precedence over the static redirection rules. Iptables rules are evaluated in order.
		rule := fmt.Sprintf("iptables -t nat -I PROXY_OUTPUT -d %s -j RETURN", cidr)
		cmd = append(cmd, rule)
	}

	// 5. Create dynamic outbound ports exclusion rules
	if len(outboundPortExclusionList) > 0 {
		var portExclusionListStr []string
		for _, port := range outboundPortExclusionList {
			portExclusionListStr = append(portExclusionListStr, strconv.Itoa(port))
		}
		outboundPortsToExclude := strings.Join(portExclusionListStr, ",")
		rule := fmt.Sprintf("iptables -t nat -I PROXY_OUTPUT -p tcp --match multiport --dports %s -j RETURN", outboundPortsToExclude)
		cmd = append(cmd, rule)
	}

	// 6. Create dynamic inbound ports exclusion rules
	if len(inboundPortExclusionList) > 0 {
		var portExclusionListStr []string
		for _, port := range inboundPortExclusionList {
			portExclusionListStr = append(portExclusionListStr, strconv.Itoa(port))
		}
		inboundPortsToExclude := strings.Join(portExclusionListStr, ",")
		rule := fmt.Sprintf("iptables -t nat -I PROXY_INBOUND -p tcp --match multiport --dports %s -j RETURN", inboundPortsToExclude)
		cmd = append(cmd, rule)
	}

	return cmd
}
