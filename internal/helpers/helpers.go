package helpers

import (
	"fmt"
	"log"
	"net"
)

type CIDRReservation struct {
	CIDR string `json:"cidr"`
}

func GenerateCIDR(existingCIDRs []string, baseCIDR string, prefixSize int) (string, error) {
	_, network, err := net.ParseCIDR(baseCIDR)

	if err != nil {
		return "", fmt.Errorf("error parsing base CIDR: %v", err)
	}

	subnets, err := SplitCIDR(network, prefixSize)

	if err != nil {
		return "", fmt.Errorf("error splitting CIDR: %v", err)
	}

	// Find a subnet that doesn't overlap with existing CIDRs.
	for _, subnet := range subnets {
		if !isOverlapping(subnet, existingCIDRs) {
			return subnet.String(), nil
		}
	}

	return "", fmt.Errorf("no available CIDR found")
}

func SplitCIDR(network *net.IPNet, prefixSize int) ([]*net.IPNet, error) {
	var subnets []*net.IPNet
	basePrefix, _ := network.Mask.Size()

	if prefixSize <= basePrefix {
		return nil, fmt.Errorf("prefix size must be greater than or equal to base prefix size")
	}

	// Calculate the number of subnets
	numSubnets := 1 << (prefixSize - basePrefix)

	for i := 0; i < numSubnets; i++ {
		ip := network.IP.Mask(network.Mask)

		for j := len(ip) - 1; j >= 0; j-- {
			ip[j] += byte(i >> (8 * (len(ip) - 1 - j)))
		}

		subnet := &net.IPNet{
			IP:   ip,
			Mask: net.CIDRMask(prefixSize, 8*len(ip)),
		}

		subnets = append(subnets, subnet)
	}

	return subnets, nil
}

// isOverlapping checks if a CIDR overlaps with any CIDRs in a list.
func isOverlapping(cidr *net.IPNet, existingCIDRs []string) bool {
	for _, existingCIDR := range existingCIDRs {
		_, existingNet, err := net.ParseCIDR(existingCIDR)
		if err != nil {
			log.Printf("Skipping invalid CIDR %s: %v", existingCIDR, err)
			continue
		}
		if cidr.Contains(existingNet.IP) || existingNet.Contains(cidr.IP) {
			return true
		}
	}
	return false
}
