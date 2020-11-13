package validate

import (
	"fmt"

	"github.com/prusya/terraform-provider-azurerm/azurerm/internal/services/msi/parse"
)

func UserAssignedIdentityId(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if _, err := parse.UserAssignedIdentityID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
		return
	}

	return warnings, errors
}
