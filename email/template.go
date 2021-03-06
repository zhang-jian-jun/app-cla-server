package email

import (
	"fmt"
	"text/template"

	"github.com/opensourceways/app-cla-server/util"
)

const (
	TmplCorporationSigning    = "corporation signing"
	TmplIndividualSigning     = "individual signing"
	TmplEmployeeSigning       = "employee signing"
	TmplCorpSigningVerifiCode = "verificaition code"
	TmplAddingCorpAdmin       = "adding corp admin"
	TmplAddingCorpManager     = "adding corp manager"
	TmplRemovingCorpManager   = "removing corp manager"
	TmplActivatingEmployee    = "activating employee"
	TmplInactivaingEmployee   = "inactivating employee"
	TmplRemovingingEmployee   = "removing employee"
)

var msgTmpl = map[string]*template.Template{}

func initTemplate() error {
	items := map[string]string{
		TmplCorporationSigning:    "./conf/email-template/corporation-signing.tmpl",
		TmplIndividualSigning:     "./conf/email-template/individual-signing.tmpl",
		TmplEmployeeSigning:       "./conf/email-template/employee-signing.tmpl",
		TmplCorpSigningVerifiCode: "./conf/email-template/verification-code.tmpl",
		TmplAddingCorpAdmin:       "./conf/email-template/adding-corp-admin.tmpl",
		TmplAddingCorpManager:     "./conf/email-template/adding-corp-manager.tmpl",
		TmplRemovingCorpManager:   "./conf/email-template/removing-corp-manager.tmpl",
		TmplActivatingEmployee:    "./conf/email-template/activating-employee.tmpl",
		TmplInactivaingEmployee:   "./conf/email-template/inactivating-employee.tmpl",
		TmplRemovingingEmployee:   "./conf/email-template/removing-employee.tmpl",
	}

	for name, path := range items {
		tmpl, err := util.NewTemplate(name, path)
		if err != nil {
			return err
		}
		msgTmpl[name] = tmpl
	}

	return nil
}

func findTmpl(name string) *template.Template {
	v, ok := msgTmpl[name]
	if ok {
		return v
	}
	return nil
}

func genEmailMsg(tmplName string, data interface{}) (*EmailMessage, error) {
	tmpl := findTmpl(tmplName)
	if tmpl == nil {
		return nil, fmt.Errorf("Failed to generate email msg: didn't find msg template: %s", tmplName)
	}

	str, err := util.RenderTemplate(tmpl, data)
	if err != nil {
		return nil, err
	}
	return &EmailMessage{Content: str}, nil
}

type IEmailMessageBulder interface {
	// msg returned only includes content
	GenEmailMsg() (*EmailMessage, error)
}

type CorporationSigning struct{}

func (this CorporationSigning) GenEmailMsg() (*EmailMessage, error) {
	return genEmailMsg(TmplCorporationSigning, this)
}

type IndividualSigning struct{}

func (this IndividualSigning) GenEmailMsg() (*EmailMessage, error) {
	return genEmailMsg(TmplIndividualSigning, this)
}

type CorpSigningVerificationCode struct {
	Code string
}

func (this CorpSigningVerificationCode) GenEmailMsg() (*EmailMessage, error) {
	return genEmailMsg(TmplCorpSigningVerifiCode, this)
}

type AddingCorpManager struct {
	Admin    bool
	Password string
}

func (this AddingCorpManager) GenEmailMsg() (*EmailMessage, error) {
	if this.Admin {
		return genEmailMsg(TmplAddingCorpAdmin, this)
	}
	return genEmailMsg(TmplAddingCorpManager, this)
}

type RemovingCorpManager struct {
}

func (this RemovingCorpManager) GenEmailMsg() (*EmailMessage, error) {
	return genEmailMsg(TmplRemovingCorpManager, this)
}

type EmployeeSigning struct {
}

func (this EmployeeSigning) GenEmailMsg() (*EmailMessage, error) {
	return genEmailMsg(TmplEmployeeSigning, this)
}

type EmployeeNotification struct {
	Removing bool
	Active   bool
	Inactive bool
}

func (this EmployeeNotification) GenEmailMsg() (*EmailMessage, error) {
	if this.Active {
		return genEmailMsg(TmplActivatingEmployee, this)
	}

	if this.Inactive {
		return genEmailMsg(TmplInactivaingEmployee, this)
	}

	if this.Removing {
		return genEmailMsg(TmplRemovingingEmployee, this)
	}

	return nil, fmt.Errorf("do nothing")
}
