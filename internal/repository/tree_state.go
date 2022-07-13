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

package repository

import (
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TreeState struct {
	ID                 int       `json:"id" 					gorm:"column:id;type:bigserial"`
	StateType          int       `json:"state_type" 			gorm:"column:state_type;type:bigint"`
	ParentNodeRefID    int       `json:"parent_node_ref_id" 	gorm:"column:parent_node_ref_id;type:bigint"`
	ChildrenNodeRefIDs int       `json:"children_node_ref_ids"  gorm:"column:children_node_ref_ids;type:bigint"`
	AppRefID           int       `json:"app_ref_id" 			gorm:"column:app_ref_id;type:bigint"`
	Version            int       `json:"version" 				gorm:"column:version;type:bigint"`
	Name               string    `json:"name" 					gorm:"column:name;type:text"`
	Content            string    `json:"content"    			gorm:"column:content;type:jsonb"`
	CreatedAt          time.Time `json:"created_at" 			gorm:"column:created_at;type:timestamp"`
	CreatedBy          int       `json:"created_by" 			gorm:"column:created_by;type:bigint"`
	UpdatedAt          time.Time `json:"updated_at" 			gorm:"column:updated_at;type:timestamp"`
	UpdatedBy          int       `json:"updated_by" 			gorm:"column:updated_by;type:bigint"`
}

type TreeStateRepository interface {
	Create(treestate *TreeState) error
	Delete(treestateID int) error
	Update(treestate *TreeState) error
	RetrieveById(treestateID int) (*TreeState, error)
	RetrieveTreeStatesByVersion(versionID int) ([]*TreeState, error)
}

type TreeStateRepositoryImpl struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

func NewTreeStateRepositoryImpl(logger *zap.SugaredLogger, db *gorm.DB) *TreeStateRepositoryImpl {
	return &TreeStateRepositoryImpl{
		logger: logger,
		db:     db,
	}
}

func (impl *TreeStateRepositoryImpl) Create(treestate *TreeState) error {
	if err := impl.db.Create(treestate).Error; err != nil {
		return err
	}
	return nil
}

func (impl *TreeStateRepositoryImpl) Delete(treestateID int) error {
	if err := impl.db.Delete(&TreeState{}, treestateID).Error; err != nil {
		return err
	}
	return nil
}

func (impl *TreeStateRepositoryImpl) Update(treestate *TreeState) error {
	if err := impl.db.Model(treestate).Updates(TreeState{
		ID:                 treestate.ID,
		StateType:          treestate.StateType,
		ParentNodeRefID:    treestate.ParentNodeRefID,
		ChildrenNodeRefIDs: treestate.ChildrenNodeRefIDs,
		AppRefID:           treestate.AppRefID,
		Version:            treestate.Version,
		Name:               treestate.Name,
		Content:            treestate.Content,
		UpdatedAt:          treestate.UpdatedAt,
		UpdatedBy:          treestate.UpdatedBy,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (impl *TreeStateRepositoryImpl) RetrieveByID(treestateID int) (*TreeState, error) {
	treestate := &TreeState{}
	if err := impl.db.First(treestate, treestateID).Error; err != nil {
		return &TreeState{}, err
	}
	return treestate, nil
}

func (impl *TreeStateRepositoryImpl) RetrieveTreeStatesByVersion(version int) ([]*TreeState, error) {
	var treestates []*TreeState
	if err := impl.db.Where("version = ?", version).Find(&treestates).Error; err != nil {
		return nil, err
	}
	return treestates, nil
}

func (impl *TreeStateRepositoryImpl) RetrieveTreeStatesByName(name string) ([]*TreeState, error) {
	var treestates []*TreeState
	if err := impl.db.Where("key = ?", name).Find(&treestates).Error; err != nil {
		return nil, err
	}
	return treestates, nil
}

func (impl *TreeStateRepositoryImpl) RetrieveTreeStatesByApp(apprefid int, statetype int, version int) ([]*TreeState, error) {
	var treestates []*TreeState
	if err := impl.db.Where("app_ref_id = ? AND state_type = ? AND version = ?", apprefid, statetype, version).Find(&treestates).Error; err != nil {
		return nil, err
	}
	return treestates, nil
}
