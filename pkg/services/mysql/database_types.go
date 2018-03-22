package mysql

import "github.com/Azure/open-service-broker-azure/pkg/service"

type databaseInstanceDetails struct {
	ARMDeploymentName string `json:"armDeployment"`
	DatabaseName      string `json:"database"`
}

// GetDBaseProvisionParametersSchema generates a schema for parameters used
// in database instance provisioning
func GetDBaseProvisionParametersSchema() map[string]*service.ParameterSchema {
	props := map[string]*service.ParameterSchema{}
	props["parentAlias"] = &service.ParameterSchema{
		Type: "string",
		Description: "Specifies the alias of the DBMS upon which the database " +
			"should be provisioned.",
		Required: true,
	}
	return props
}

func (d *databaseManager) SplitProvisioningParameters(
	cpp service.CombinedProvisioningParameters,
) (
	service.ProvisioningParameters,
	service.SecureProvisioningParameters,
	error,
) {
	return nil, nil, nil
}

func (d *databaseManager) SplitBindingParameters(
	params service.CombinedBindingParameters,
) (service.BindingParameters, service.SecureBindingParameters, error) {
	return nil, nil, nil
}
