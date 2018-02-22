package service

import (
	"encoding/json"
	"time"

	"github.com/Azure/open-service-broker-azure/pkg/crypto"
)

// Instance represents an instance of a service
type Instance struct {
	InstanceID                            string                       `json:"instanceId"` // nolint: lll
	Alias                                 string                       `json:"alias"`      // nolint: lll
	ServiceID                             string                       `json:"serviceId"`  // nolint: lll
	Service                               Service                      `json:"-"`
	PlanID                                string                       `json:"planId"` // nolint: lll
	Plan                                  Plan                         `json:"-"`
	ProvisioningParameters                ProvisioningParameters       `json:"provisioningParameters"`       // nolint: lll
	EncryptedSecureProvisioningParameters []byte                       `json:"secureProvisioningParameters"` // nolint: lll
	SecureProvisioningParameters          SecureProvisioningParameters `json:"-"`
	EncryptedUpdatingParameters           []byte                       `json:"updatingParameters"` // nolint: lll
	UpdatingParameters                    UpdatingParameters           `json:"-"`
	Status                                string                       `json:"status"`        // nolint: lll
	StatusReason                          string                       `json:"statusReason"`  // nolint: lll
	Location                              string                       `json:"location"`      // nolint: lll
	ResourceGroup                         string                       `json:"resourceGroup"` // nolint: lll
	Parent                                *Instance                    `json:"-"`
	ParentAlias                           string                       `json:"parentAlias"`   // nolint: lll
	Tags                                  map[string]string            `json:"tags"`          // nolint: lll
	Details                               InstanceDetails              `json:"details"`       // nolint: lll
	EncryptedSecureDetails                []byte                       `json:"secureDetails"` // nolint: lll
	SecureDetails                         SecureInstanceDetails        `json:"-"`
	Created                               time.Time                    `json:"created"` // nolint: lll
}

// NewInstanceFromJSON returns a new Instance unmarshalled from the provided
// JSON []byte
func NewInstanceFromJSON(
	jsonBytes []byte,
	pp ProvisioningParameters,
	spp SecureProvisioningParameters,
	up UpdatingParameters,
	dt InstanceDetails,
	sdt InstanceDetails,
	codec crypto.Codec,
) (Instance, error) {
	instance := Instance{
		ProvisioningParameters:       pp,
		SecureProvisioningParameters: spp,
		UpdatingParameters:           up,
		Details:                      dt,
		SecureDetails:                sdt,
	}
	if err := json.Unmarshal(jsonBytes, &instance); err != nil {
		return instance, err
	}
	if instance.ProvisioningParameters == nil {
		instance.ProvisioningParameters = pp
	}
	if instance.Details == nil {
		instance.Details = dt
	}
	return instance.decrypt(codec)
}

// ToJSON returns a []byte containing a JSON representation of the
// instance
func (i Instance) ToJSON(codec crypto.Codec) ([]byte, error) {
	var err error
	if i, err = i.encrypt(codec); err != nil {
		return nil, err
	}
	return json.Marshal(i)
}

func (i Instance) encrypt(codec crypto.Codec) (Instance, error) {
	var err error
	if i, err = i.encryptSecureProvisioningParameters(codec); err != nil {
		return i, err
	}
	if i, err = i.encryptUpdatingParameters(codec); err != nil {
		return i, err
	}
	return i.encryptSecureDetails(codec)
}

func (i Instance) encryptSecureProvisioningParameters(
	codec crypto.Codec,
) (Instance, error) {
	jsonBytes, err := json.Marshal(i.SecureProvisioningParameters)
	if err != nil {
		return i, err
	}
	i.EncryptedSecureProvisioningParameters, err = codec.Encrypt(jsonBytes)
	return i, err
}

func (i Instance) encryptUpdatingParameters(
	codec crypto.Codec,
) (Instance, error) {
	jsonBytes, err := json.Marshal(i.UpdatingParameters)
	if err != nil {
		return i, err
	}
	i.EncryptedUpdatingParameters, err = codec.Encrypt(jsonBytes)
	return i, err
}

func (i Instance) encryptSecureDetails(
	codec crypto.Codec,
) (Instance, error) {
	jsonBytes, err := json.Marshal(i.SecureDetails)
	if err != nil {
		return i, err
	}
	i.EncryptedSecureDetails, err = codec.Encrypt(jsonBytes)
	return i, err
}

func (i Instance) decrypt(codec crypto.Codec) (Instance, error) {
	var err error
	if i, err = i.decryptSecureProvisioningParameters(codec); err != nil {
		return i, err
	}
	if i, err = i.decryptUpdatingParameters(codec); err != nil {
		return i, err
	}
	return i.decryptSecureDetails(codec)
}

func (i Instance) decryptSecureProvisioningParameters(
	codec crypto.Codec,
) (Instance, error) {
	if len(i.EncryptedSecureProvisioningParameters) == 0 ||
		i.SecureProvisioningParameters == nil {
		return i, nil
	}
	plaintext, err := codec.Decrypt(i.EncryptedSecureProvisioningParameters)
	if err != nil {
		return i, err
	}
	return i, json.Unmarshal(plaintext, i.SecureProvisioningParameters)
}

func (i Instance) decryptUpdatingParameters(
	codec crypto.Codec,
) (Instance, error) {
	if len(i.EncryptedUpdatingParameters) == 0 ||
		i.UpdatingParameters == nil {
		return i, nil
	}
	plaintext, err := codec.Decrypt(i.EncryptedUpdatingParameters)
	if err != nil {
		return i, err
	}
	return i, json.Unmarshal(plaintext, i.UpdatingParameters)
}

func (i Instance) decryptSecureDetails(codec crypto.Codec) (Instance, error) {
	if len(i.EncryptedSecureDetails) == 0 ||
		i.SecureDetails == nil {
		return i, nil
	}
	plaintext, err := codec.Decrypt(i.EncryptedSecureDetails)
	if err != nil {
		return i, err
	}
	return i, json.Unmarshal(plaintext, i.SecureDetails)
}
