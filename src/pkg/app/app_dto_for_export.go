package app

import (
	"time"

	"github.com/google/uuid"
	"github.com/illacloud/builder-backend/internal/repository"
	"github.com/illacloud/builder-backend/src/utils/idconvertor"
)

type AppDtoForExport struct {
	ID              string                `json:"appId"` // generated by database primary key serial
	UID             uuid.UUID             `json:"uid"`
	TeamID          string                `json:"teamID"`
	Name            string                `json:"appName" validate:"required"`
	ReleaseVersion  int                   `json:"releaseVersion"`  // release version used for mark the app release version.
	MainlineVersion int                   `json:"mainlineVersion"` // mainline version keep the newest app version in database.
	Deployed        bool                  `json:"deployed"`
	Config          *repository.AppConfig `json:"config"`
	CreatedBy       string                `json:"-" `
	CreatedAt       time.Time             `json:"-"`
	UpdatedBy       string                `json:"updatedBy"`
	UpdatedAt       time.Time             `json:"updatedAt"`
	AppActivity     AppActivity           `json:"appActivity"`
}

func NewAppDtoForExport(a *AppDto) *AppDtoForExport {
	// feedback
	return &AppDtoForExport{
		ID:              idconvertor.ConvertIntToString(a.ID),
		UID:             a.UID,
		TeamID:          idconvertor.ConvertIntToString(a.TeamID),
		Name:            a.Name,
		ReleaseVersion:  a.ReleaseVersion,
		MainlineVersion: a.MainlineVersion,
		Config:          a.Config,
		CreatedBy:       idconvertor.ConvertIntToString(a.CreatedBy),
		CreatedAt:       a.CreatedAt,
		UpdatedBy:       idconvertor.ConvertIntToString(a.UpdatedBy),
		UpdatedAt:       a.UpdatedAt,
		AppActivity:     a.AppActivity,
	}
}

func (resp *AppDtoForExport) ExportForFeedback() interface{} {
	return resp
}
