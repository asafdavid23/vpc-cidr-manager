package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"text/template"

	"gopkg.in/yaml.v2"
)

type CIDRReservation struct {
	CIDR string `json:"cidr"`
}

type IAMTemplateData struct {
	RoleName  string
	Principal string
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

// LoadAndRenderTemplate loads a CloudFormation template from a file, processes it with dynamic values, and returns the rendered template
func LoadAndRenderTemplate(templateFilePath string, data IAMTemplateData) (string, error) {
	// Read the template file
	cfnTemplate, err := os.ReadFile(templateFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template file: %v", err)
	}

	// Parse the template
	tmpl, err := template.New("cfnTemplate").Parse(string(cfnTemplate))
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}

	// Apply the data to the template
	var renderedTemplate bytes.Buffer
	err = tmpl.Execute(&renderedTemplate, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	return renderedTemplate.String(), nil
}

// LoadFileContent loads the content of a file as a string
// If the content is in JSON format, it converts it to YAML with proper indentation
func LoadFileContent(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %v", filePath, err)
	}

	// Attempt to parse as JSON
	var jsonData map[string]interface{}
	if json.Unmarshal(content, &jsonData) == nil {
		// Convert JSON to YAML with indentation
		yamlData, err := yaml.Marshal(jsonData)
		if err != nil {
			return "", fmt.Errorf("failed to convert JSON to YAML: %v", err)
		}

		// Return the YAML string with proper indentation
		return string(yamlData), nil
	}

	// Return as is if not JSON
	return string(content), nil
}
