// Copyright 2023 Illa Soft, Inc.
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

package oracle9i

import (
	"errors"

	"github.com/illacloud/builder-backend/src/actionruntime/common"
)

const (
	FIELD_CONTEXT = "context"
	FIELD_QUERY   = "query"
)

type Resource struct {
	Host     string `mapstructure:"host" validate:"required"`
	Port     string `mapstructure:"port" validate:"required"`
	Type     string `mapstructure:"connectionType" validate:"oneof=SID Service"`
	Name     string `mapstructure:"name"`
	SSL      bool   `mapstructure:"ssl"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type Action struct {
	Mode     string                 `mapstructure:"mode" validate:"oneof=gui sql sql-safe"`
	Opts     map[string]interface{} `mapstructure:"opts"`
	RawQuery string
	Context  map[string]interface{}
}

func (q *Action) IsSafeMode() bool {
	return q.Mode == common.MODE_SQL_SAFE
}

func (q *Action) SetRawQueryAndContext(rawTemplate map[string]interface{}) error {
	queryRaw, hit := rawTemplate[FIELD_QUERY]
	if !hit {
		return errors.New("missing query field for SetRawQueryAndContext() in query")
	}
	queryAsserted, assertPass := queryRaw.(string)
	if !assertPass {
		return errors.New("query field assert failed in SetRawQueryAndContext() method")

	}
	q.RawQuery = queryAsserted
	contextRaw, hit := rawTemplate[FIELD_CONTEXT]
	if !hit {
		return errors.New("missing context field SetRawQueryAndContext() in query")
	}
	contextAsserted, assertPass := contextRaw.(map[string]interface{})
	if !assertPass {
		return errors.New("context field assert failed in SetRawQueryAndContext() method")

	}
	q.Context = contextAsserted
	return nil
}

type SQL struct {
	Raw string `mapstructure:"raw"`
}

type GUIBulkOpts struct {
	Table   string                   `mapstructure:"table"`
	Type    string                   `mapstructure:"actionType"`
	Records []map[string]interface{} `mapstructure:"records"`
	Key     string                   `mapstructure:"primaryKey"`
}
