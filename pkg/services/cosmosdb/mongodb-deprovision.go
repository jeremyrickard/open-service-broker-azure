package cosmosdb

import (
	"context"

	"github.com/Azure/open-service-broker-azure/pkg/service"
)

func (m *mongoManager) GetDeprovisioner(
	service.Plan,
) (service.Deprovisioner, error) {
	return service.NewDeprovisioner(
		service.NewDeprovisioningStep("deleteARMDeployment", m.deleteARMDeployment),
		service.NewDeprovisioningStep(
			"deleteCosmosDBServer",
			m.deleteCosmosDBServer,
		),
	)
}

func (m *mongoManager) deleteARMDeployment(
	_ context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	err := deleteARMDeployment(m.armDeployer, instance)
	if err != nil {
		return nil, nil, err
	}
	return instance.Details, instance.SecureDetails, nil
}

func (m *mongoManager) deleteCosmosDBServer(
	ctx context.Context,
	instance service.Instance,
) (service.InstanceDetails, service.SecureInstanceDetails, error) {
	err := deleteCosmosDBServer(ctx, m.databaseAccountsClient, instance)
	if err != nil {
		return nil, nil, err
	}
	return instance.Details, instance.SecureDetails, nil
}
