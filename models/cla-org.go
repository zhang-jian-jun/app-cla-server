package models

import (
	"time"

	"github.com/opensourceways/app-cla-server/dbmodels"
	"github.com/opensourceways/app-cla-server/util"
)

type CLAOrg struct {
	ID                   string    `json:"id"`
	Platform             string    `json:"platform"`
	OrgID                string    `json:"org_id"`
	RepoID               string    `json:"repo_id"`
	CLAID                string    `json:"cla_id"`
	CLALanguage          string    `json:"cla_language"`
	ApplyTo              string    `json:"apply_to"`
	OrgEmail             string    `json:"org_email"`
	Enabled              bool      `json:"enabled"`
	Submitter            string    `json:"submitter"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	OrgSignatureUploaded bool      `json:"org_signature_uploaded"`
}

func (this *CLAOrg) Create() error {
	this.Enabled = true

	p := dbmodels.CLAOrg{}
	if err := util.CopyBetweenStructs(this, &p); err != nil {
		return err
	}

	v, err := dbmodels.GetDB().CreateBindingBetweenCLAAndOrg(p)
	if err == nil {
		this.ID = v
	}

	return err
}

func (this CLAOrg) Delete() error {
	return dbmodels.GetDB().DeleteBindingBetweenCLAAndOrg(this.ID)
}

func (this *CLAOrg) Get() error {
	v, err := dbmodels.GetDB().GetBindingBetweenCLAAndOrg(this.ID)
	if err != nil {
		return err
	}
	return util.CopyBetweenStructs(&v, this)
}

type CLAOrgListOption dbmodels.CLAOrgListOption

func (this CLAOrgListOption) ListForSigningPage() ([]dbmodels.CLAOrg, error) {
	return dbmodels.GetDB().ListBindingForSigningPage(dbmodels.CLAOrgListOption(this))
}

func (this CLAOrgListOption) List() ([]dbmodels.CLAOrg, error) {
	return dbmodels.GetDB().ListBindingBetweenCLAAndOrg(dbmodels.CLAOrgListOption(this))
}
