// Copyright 2022 The ILLA Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package state

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/illa-family/builder-backend/internal/repository"
	"github.com/illa-family/builder-backend/pkg/app"
	"github.com/illa-family/builder-backend/pkg/connector"
	"go.uber.org/zap"
)

type KVStateService interface {
	CreateKVState(version int, kvstate KVStateDto) (KVStateDto, error)
	DeleteKVState(kvstateID int) error
	UpdateKVState(version int, kvstate KVStateDto) (KVStateDto, error)
	GetKVStateByID(kvstateID int) (KVStateDto, error)
	FindKVStatesByVersion(version int) ([]KVStateDto, error)
	RunKVState(kvstate KVStateDto) (interface{}, error)
	ReleaseKVStateByApp(app *app.AppDto) error
}

type KVStateDto struct {
	ID        int       `json:"id"`
	StateType int       `json:"state_type"`
	AppRefID  int       `json:"app_ref_id"`
	Version   int       `json:"version"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy int       `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy int       `json:"updated_by"`
}

type KVStateServiceImpl struct {
	logger            *zap.SugaredLogger
	kvstateRepository repository.KVStateRepository
	kvstateRepository repository.kvstateRepository
}

func NewKVStateServiceImpl(logger *zap.SugaredLogger, kvstateRepository repository.KVStateRepository,
	kvstateRepository repository.kvstateRepository) *KVStateServiceImpl {
	return &KVStateServiceImpl{
		logger:            logger,
		kvstateRepository: kvstateRepository,
		kvstateRepository: kvstateRepository,
	}
}

func (impl *KVStateServiceImpl) CreateKVState(version int, kvstate KVStateDto) (KVStateDto, error) {
	// TODO: validate the version
	validate := validator.New()
	if err := validate.Struct(kvstate); err != nil {
		return KVStateDto{}, err
	}
	kvstate.CreatedAt = time.Now().UTC()
	kvstate.UpdatedAt = time.Now().UTC()
	if err := impl.kvstateRepository.Create(&repository.KVState{
		ID:        kvstate.ID,
		StateType: kvstate.StateType,
		AppRefID:  kvstate.AppRefID,
		Version:   kvstate.Version,
		Key:       kvstate.Key,
		Value:     kvstate.Value,
		CreatedAt: kvstate.CreatedAt,
		CreatedBy: kvstate.CreatedBy,
		UpdatedAt: kvstate.UpdatedAt,
		UpdatedBy: kvstate.UpdatedBy,
	}); err != nil {
		return KVStateDto{}, err
	}
	return kvstate, nil
}

func (impl *KVStateServiceImpl) DeleteKVState(kvstateID int) error {
	if err := impl.kvstateRepository.Delete(kvstateID); err != nil {
		return err
	}
	return nil
}

func (impl *KVStateServiceImpl) UpdateKVState(version int, kvstate KVStateDto) (KVStateDto, error) {
	validate := validator.New()
	if err := validate.Struct(kvstate); err != nil {
		return KVStateDto{}, err
	}
	kvstate.UpdatedAt = time.Now().UTC()
	if err := impl.kvstateRepository.Update(&repository.KVState{
		ID:        kvstate.ID,
		StateType: kvstate.StateType,
		AppRefID:  kvstate.AppRefID,
		Version:   kvstate.Version,
		Key:       kvstate.Key,
		Value:     kvstate.Value,
		CreatedAt: kvstate.CreatedAt,
		CreatedBy: kvstate.CreatedBy,
		UpdatedAt: kvstate.UpdatedAt,
		UpdatedBy: kvstate.UpdatedBy,
	}); err != nil {
		return KVStateDto{}, err
	}
	return kvstate, nil
}

func (impl *KVStateServiceImpl) GetKVStateByID(kvstateID int) (KVStateDto, error) {
	res, err := impl.kvstateRepository.RetrieveById(kvstateID)
	if err != nil {
		return KVStateDto{}, err
	}
	resDto := KVStateDto{
		ID:        res.ID,
		StateType: res.StateType,
		AppRefID:  res.AppRefID,
		Version:   res.Version,
		Key:       res.Key,
		Value:     res.Value,
		CreatedAt: res.CreatedAt,
		CreatedBy: res.CreatedBy,
		UpdatedAt: res.UpdatedAt,
		UpdatedBy: res.UpdatedBy,
	}
	return resDto, nil
}

func (impl *KVStateServiceImpl) FindKVStatesByVersion(version int) ([]KVStateDto, error) {
	res, err := impl.kvstateRepository.RetrieveKVStatesByVersion(version)
	if err != nil {
		return nil, err
	}
	resDtoSlice := make([]KVStateDto, 0, len(res))
	for _, value := range res {
		resDtoSlice = append(resDtoSlice, KVStateDto{
			KVStateId:       value.ID,
			ResourceId:      value.ResourceID,
			DisplayName:     value.Name,
			KVStateType:     value.Type,
			KVStateTemplate: value.KVStateTemplate,
			CreatedBy:       value.CreatedBy,
			CreatedAt:       value.CreatedAt,
			UpdatedBy:       value.UpdatedBy,
			UpdatedAt:       value.UpdatedAt,
		})
	}
	return resDtoSlice, nil
}

func (impl *KVStateServiceImpl) RunKVState(kvstate KVStateDto) (interface{}, error) {
	var kvstateFactory *Factory
	if kvstate.ResourceId != uuid.Nil {
		rsc, err := impl.kvstateRepository.RetrieveById(kvstate.ResourceId)
		if err != nil {
			return nil, err
		}
		resourceConn := &connector.Connector{
			Type:    rsc.Kind,
			Options: rsc.Options,
		}
		kvstateFactory = &Factory{
			Type:     kvstate.StateType,
			Template: kvstate.KVStateTemplate,
			Resource: resourceConn,
		}
	} else {
		kvstateFactory = &Factory{
			Type:     kvstate.StateType,
			Template: kvstate.KVStateTemplate,
			Resource: nil,
		}
	}
	kvstateAssemblyline := kvstateFactory.Build()
	if kvstateAssemblyline == nil {
		return nil, errors.New("invalid KVStateType:: unsupported type")
	}
	res, err := kvstateAssemblyline.Run()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// @todo: should this method be in a transaction?
func (impl *KVStateServiceImpl) ReleaseKVStateByApp(app *app.AppDto) error {
	// get edit version K-V state from database
	kvstates, err := impl.kvstateRepository.RetrieveAllTypeKVStatesByApp(app.ID, repository.APP_EDIT_VERSION)
	if err != nil {
		return err
	}
	// set version as minaline version
	for serial, _ := range kvstates {
		kvstates[serial].Version = app.MainlineVersion
	}
	// and put them to the database as duplicate
	for _, kvstate := range kvstates {
		if err := impl.kvstateRepository.Create(kvstate); err != nil {
			return err
		}
	}
	return nil
}
