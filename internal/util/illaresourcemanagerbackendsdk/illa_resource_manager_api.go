package illaresourcemanagerbackendsdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/illacloud/builder-backend/internal/tokenvalidator"
	"github.com/illacloud/builder-backend/internal/util/config"
)

const (
	BASEURL = "http://127.0.0.1:8008/api/v1"
	// api route part
	GET_AI_AGENT_INTERNAL_API             = "/api/v1/aiAgent/%d"
	RUN_AI_AGENT_INTERNAL_API             = "/api/v1/aiAgent/%d/run"
	DELETE_TEAM_ALL_AI_AGENT_INTERNAL_API = "/api/v1/teams/%d/aiAgent/all"
)

const (
	PRODUCT_TYPE_AIAGENTS = "aiAgents"
	PRODUCT_TYPE_APPS     = "apps"
	PRODUCT_TYPE_HUBS     = "hubs"
)

type IllaResourceManagerRestAPI struct {
	Config *config.Config
	Debug  bool `json:"-"`
}

func NewIllaResourceManagerRestAPI() (*IllaResourceManagerRestAPI, error) {
	return &IllaResourceManagerRestAPI{
		Config: config.GetInstance(),
	}, nil
}

func (r *IllaResourceManagerRestAPI) CloseDebug() {
	r.Debug = false
}

func (r *IllaResourceManagerRestAPI) OpenDebug() {
	r.Debug = true
}

func (r *IllaResourceManagerRestAPI) GetAIAgent(aiAgentID int) (map[string]interface{}, error) {
	client := resty.New()
	tokenValidator := tokenvalidator.NewRequestTokenValidator()
	uri := r.Config.GetIllaResourceManagerInternalRestAPI() + fmt.Sprintf(GET_AI_AGENT_INTERNAL_API, aiAgentID)
	resp, errInPost := client.R().
		SetHeader("Request-Token", tokenValidator.GenerateValidateToken(strconv.Itoa(aiAgentID))).
		Get(uri)
	if r.Debug {
		log.Printf("[IllaResourceManagerRestAPI.GetAiAgent()]  uri: %+v \n", uri)
		log.Printf("[IllaResourceManagerRestAPI.GetAiAgent()]  response: %+v, err: %+v \n", resp, errInPost)
		log.Printf("[IllaResourceManagerRestAPI.GetAiAgent()]  resp.StatusCode(): %+v \n", resp.StatusCode())
	}
	if errInPost != nil {
		return nil, errInPost
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return nil, errors.New(resp.String())
	}

	var aiAgent map[string]interface{}
	errInUnMarshal := json.Unmarshal([]byte(resp.String()), &aiAgent)
	if errInUnMarshal != nil {
		return nil, errInUnMarshal
	}
	return aiAgent, nil
}

func (r *IllaResourceManagerRestAPI) RunAIAgent(req map[string]interface{}) (*RunAIAgentResult, error) {
	reqInstance, errInNewReq := NewRunAIAgentRequest(req)
	if errInNewReq != nil {
		return nil, errInNewReq
	}
	client := resty.New()
	uri := r.Config.GetIllaResourceManagerInternalRestAPI() + fmt.Sprintf(RUN_AI_AGENT_INTERNAL_API, reqInstance.ExportAIAgentIDInInt())
	if r.Debug {
		log.Printf("[IllaResourceManagerRestAPI.RunAiAgent()]  uri: %+v \n", uri)
	}
	tokenValidator := tokenvalidator.NewRequestTokenValidator()
	resp, errInPost := client.R().
		SetHeader("Request-Token", tokenValidator.GenerateValidateToken(strconv.Itoa(reqInstance.ExportAIAgentIDInInt()))).
		SetBody(req).
		Post(uri)
	if r.Debug {
		log.Printf("[IllaResourceManagerRestAPI.RunAiAgent()]  response: %+v, err: %+v \n", resp, errInPost)
		log.Printf("[IllaResourceManagerRestAPI.RunAiAgent()]  resp.StatusCode(): %+v \n", resp.StatusCode())
	}
	if errInPost != nil {
		return nil, errInPost
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return nil, errors.New(resp.String())
	}

	runAIAgentResult := NewRunAIAgentResult()
	errInUnMarshal := json.Unmarshal([]byte(resp.String()), &runAIAgentResult)
	if errInUnMarshal != nil {
		return nil, errInUnMarshal
	}
	return runAIAgentResult, nil
}

func (r *IllaResourceManagerRestAPI) DeleteTeamAllAIAgent(teamID int) error {
	client := resty.New()
	tokenValidator := tokenvalidator.NewRequestTokenValidator()
	uri := r.Config.GetIllaResourceManagerInternalRestAPI() + fmt.Sprintf(DELETE_TEAM_ALL_AI_AGENT_INTERNAL_API, teamID)
	resp, errInDelete := client.R().
		SetHeader("Request-Token", tokenValidator.GenerateValidateToken(strconv.Itoa(teamID))).
		Delete(uri)
	if r.Debug {
		log.Printf("[IllaResourceManagerRestAPI.DeleteTeamAllAiAgent()]  uri: %+v \n", uri)
		log.Printf("[IllaResourceManagerRestAPI.DeleteTeamAllAiAgent()]  response: %+v, err: %+v \n", resp, errInDelete)
		log.Printf("[IllaResourceManagerRestAPI.DeleteTeamAllAiAgent()]  resp.StatusCode(): %+v \n", resp.StatusCode())
	}
	if errInDelete != nil {
		return errInDelete
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return errors.New(resp.String())
	}

	return nil
}
