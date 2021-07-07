package configurator

// GetMeshCIDRRanges returns a list of mesh CIDR ranges
//func (c *Client) GetMeshCIDRRanges() []string {
//	noSpaces := strings.ReplaceAll(c.getMeshConfig().Spec.MeshCIDRRanges, " ", ",")
//	commaSeparatedCIDRs := strings.Split(noSpaces, ",")
//
//	cidrSet := make(map[string]interface{})
//	for _, cidr := range commaSeparatedCIDRs {
//		trimmedCIDR := strings.Trim(cidr, " ")
//		if len(trimmedCIDR) == 0 {
//			continue
//		}
//
//		_, _, err := net.ParseCIDR(trimmedCIDR)
//		if err != nil {
//			log.Error().Err(err).Msgf("Found incorrectly formatted in-mesh CIDR %s from ConfigMap %s/%s; Skipping CIDR", trimmedCIDR, c.osmNamespace, c.osmConfigMapName)
//			continue
//		}
//
//		cidrSet[trimmedCIDR] = nil
//	}
//
//	var cidrs []string
//	for cidr := range cidrSet {
//		cidrs = append(cidrs, cidr)
//	}
//
//	sort.Strings(cidrs)
//
//	return cidrs
//}
