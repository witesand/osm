package injector

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
)

func getInitContainerSpec(containerName string, containerImage string, outboundIPRangeExclusionList []string, enablePrivilegedInitContainer bool) corev1.Container {
	iptablesInitCommandsList := generateIptablesCommands(outboundIPRangeExclusionList)
	iptablesInitCommand := strings.Join(iptablesInitCommandsList, " && ")

	return corev1.Container{
//<<<<<<< HEAD
//		Name:  initContainer.Name,
//		Image: initContainer.Image,
//		ImagePullPolicy: "Always",
//=======
		Name:  containerName,
		Image: containerImage,
//>>>>>>> 3d923b3f2d72006f6cdaad056938c492c364196d
		SecurityContext: &corev1.SecurityContext{
			Privileged: &enablePrivilegedInitContainer,
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{
					"NET_ADMIN",
				},
			},
		},
//<<<<<<< HEAD
//		Env: []corev1.EnvVar{
//			{
//				Name:  "OSM_PROXY_UID",
//				Value: fmt.Sprintf("%d", constants.EnvoyUID),
//			},
//			{
//				Name:  "OSM_ENVOY_INBOUND_PORT",
//				Value: fmt.Sprintf("%d", constants.EnvoyInboundListenerPort),
//			},
//			{
//				Name:  "OSM_ENVOY_OUTBOUND_PORT",
//				Value: fmt.Sprintf("%d", constants.EnvoyOutboundListenerPort),
//			},
//			{
//				Name:  "CIDR1",
//				Value: initContainer.CIDR1,
//			},
//			{
//				Name:  "CIDR2",
//				Value: initContainer.CIDR2,
//			},
//=======
		Command: []string{"/bin/sh"},
		Args: []string{
			"-c",
			iptablesInitCommand,
//>>>>>>> 3d923b3f2d72006f6cdaad056938c492c364196d
		},
	}
}
