package validate

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/services/network/parse"
)

// VirtualNetworkID validates that the specified ID is a valid Virtual Network ID
func VirtualNetworkID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.VirtualNetworkID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}

func VirtualNetworkBgpCommunity(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	segments := strings.Split(v, ":")
	if len(segments) != 2 {
		errors = append(errors, fmt.Errorf(`invalid notation of bgp community: expected "x:y"`))
		return
	}

	asn, err := strconv.Atoi(segments[0])
	if err != nil {
		errors = append(errors, fmt.Errorf(`converting asn %q: %v`, segments[0], err))
		return
	}
	if !(asn > 0 && asn < 65535) {
		errors = append(errors, fmt.Errorf(`asn %d exceeds range: [0, 65535]`, asn))
		return
	}

	comm, err := strconv.Atoi(segments[1])
	if err != nil {
		errors = append(errors, fmt.Errorf(`converting community value %q: %v`, segments[1], err))
		return
	}
	if !(comm > 0 && comm < 65535) {
		errors = append(errors, fmt.Errorf(`community value %d exceeds range: [0, 65535]`, comm))
		return
	}
	return warnings, errors
}
