package controllers

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"

	"github.com/opensourceways/app-cla-server/dbmodels"
	"github.com/opensourceways/app-cla-server/email"
	"github.com/opensourceways/app-cla-server/models"
	"github.com/opensourceways/app-cla-server/util"
	"github.com/opensourceways/app-cla-server/worker"
)

type EmployeeSigningController struct {
	beego.Controller
}

func (this *EmployeeSigningController) Prepare() {
	if getRequestMethod(&this.Controller) == http.MethodPost {
		// sign as employee
		apiPrepare(&this.Controller, []string{PermissionIndividualSigner}, nil)
	} else {
		// get, update and delete employee
		apiPrepare(&this.Controller, []string{PermissionEmployeeManager}, nil)
	}
}

// @Title Post
// @Description sign as employee
// @Param	:cla_org_id	path 	string				true		"cla org id"
// @Param	body		body 	models.IndividualSigning	true		"body for employee signing"
// @Success 201 {int} map
// @Failure util.ErrHasSigned		"employee has signed"
// @Failure util.ErrHasNotSigned	"corp has not signed"
// @Failure util.ErrSigningUncompleted	"corp has not been enabled"
// @router /:cla_org_id [post]
func (this *EmployeeSigningController) Post() {
	var statusCode = 0
	var errCode = ""
	var reason error
	var body interface{}

	defer func() {
		sendResponse(&this.Controller, statusCode, errCode, reason, body, "sign as employee")
	}()

	claOrgID, err := fetchStringParameter(&this.Controller, ":cla_org_id")
	if err != nil {
		reason = err
		errCode = util.ErrInvalidParameter
		statusCode = 400
		return
	}

	var info models.IndividualSigning
	if err := fetchInputPayload(&this.Controller, &info); err != nil {
		reason = err
		errCode = util.ErrInvalidParameter
		statusCode = 400
		return
	}

	claOrg := &models.CLAOrg{ID: claOrgID}
	if err := claOrg.Get(); err != nil {
		reason = err
		return
	}

	corpSignedCla, corpSign, err := models.GetCorporationSigningDetail(
		claOrg.Platform, claOrg.OrgID, claOrg.RepoID, info.Email)
	if err != nil {
		reason = err
		return
	}

	if !corpSign.AdminAdded {
		reason = fmt.Errorf("the corp has not been enabled")
		errCode = util.ErrSigningUncompleted
		statusCode = 400
		return
	}

	err = (&info).Create(claOrgID, claOrg.Platform, claOrg.OrgID, claOrg.RepoID, false)
	if err != nil {
		reason = err
		return
	}
	body = "sign successfully"

	d := email.EmployeeSigning{}
	this.notifyManagers(corpSignedCla, info.Email, claOrg.OrgEmail, "Employee Signing", d)
}

// @Title GetAll
// @Description get all the employees
// @Success 200 {int} map
// @router / [get]
func (this *EmployeeSigningController) GetAll() {
	var statusCode = 0
	var errCode = ""
	var reason error
	var body interface{}

	defer func() {
		sendResponse(&this.Controller, statusCode, errCode, reason, body, "list employees")
	}()

	claOrgID, corpEmail, err := parseCorpManagerUser(&this.Controller)
	if err != nil {
		reason = err
		errCode = util.ErrUnknownToken
		statusCode = 401
		return
	}

	claOrg := &models.CLAOrg{ID: claOrgID}
	if err := claOrg.Get(); err != nil {
		reason = err
		return
	}

	opt := models.EmployeeSigningListOption{
		CLALanguage: this.GetString("cla_language"),
	}

	r, err := opt.List(corpEmail, claOrg.Platform, claOrg.OrgID, claOrg.RepoID)
	if err != nil {
		reason = err
		return
	}

	body = r
}

// @Title Update
// @Description enable/unable employee signing
// @Param	:cla_org_id	path 	string	true		"cla org id"
// @Param	:email		path 	string	true		"email"
// @Success 202 {int} map
// @router /:cla_org_id/:email [put]
func (this *EmployeeSigningController) Update() {
	var statusCode = 0
	var errCode = ""
	var reason error
	var body interface{}

	defer func() {
		sendResponse(&this.Controller, statusCode, errCode, reason, body, "enable/unable employee signing")
	}()

	if err := checkAPIStringParameter(&this.Controller, []string{":cla_org_id", ":email"}); err != nil {
		reason = err
		errCode = util.ErrInvalidParameter
		statusCode = 400
		return

	}
	employeeEmail := this.GetString(":email")
	claOrgID := this.GetString(":cla_org_id")

	orgEmail := ""
	statusCode, errCode, orgEmail, reason = this.canHandleOnEmployee(claOrgID, employeeEmail)
	if reason != nil {
		return
	}

	var info models.EmployeeSigningUdateInfo
	if err := fetchInputPayload(&this.Controller, &info); err != nil {
		reason = err
		errCode = util.ErrInvalidParameter
		statusCode = 400
		return
	}

	if err := (&info).Update(claOrgID, employeeEmail); err != nil {
		reason = err
		return
	}

	body = "enabled employee successfully"

	b := email.EmployeeNotification{}
	subject := ""
	if info.Enabled {
		b.Active = true
		subject = "Activate employee"
	} else {
		b.Inactive = true
		subject = "Inavtivate employee"
	}
	this.notifyEmployee(employeeEmail, orgEmail, subject, &b)
}

// @Title Delete
// @Description delete employee signing
// @Param	:cla_org_id	path 	string	true		"cla org id"
// @Param	:email		path 	string	true		"email"
// @Success 204 {string} delete success!
// @router /:cla_org_id/:email [delete]
func (this *EmployeeSigningController) Delete() {
	var statusCode = 0
	var errCode = ""
	var reason error
	var body string

	defer func() {
		sendResponse(&this.Controller, statusCode, errCode, reason, body, "delete employee signing")
	}()

	if err := checkAPIStringParameter(&this.Controller, []string{":cla_org_id", ":email"}); err != nil {
		reason = err
		errCode = util.ErrInvalidParameter
		statusCode = 400
		return

	}
	employeeEmail := this.GetString(":email")
	claOrgID := this.GetString(":cla_org_id")

	orgEmail := ""
	statusCode, errCode, orgEmail, reason = this.canHandleOnEmployee(claOrgID, employeeEmail)
	if reason != nil {
		return
	}

	if err := models.DeleteEmployeeSigning(claOrgID, employeeEmail); err != nil {
		reason = err
		return
	}

	body = "delete employee successfully"

	b := email.EmployeeNotification{Removing: true}
	subject := "Remove employee"
	this.notifyEmployee(employeeEmail, orgEmail, subject, &b)
}

func (this *EmployeeSigningController) canHandleOnEmployee(claOrgID, employeeEmail string) (int, string, string, error) {
	corpClaOrgID, corpEmail, err := parseCorpManagerUser(&this.Controller)
	if err != nil {
		return 401, util.ErrUnknownToken, "", err
	}

	if !isSameCorp(corpEmail, employeeEmail) {
		return 400, util.ErrNotSameCorp, "", fmt.Errorf("not same corp")
	}

	claOrg := &models.CLAOrg{ID: claOrgID}
	if err := claOrg.Get(); err != nil {
		return 0, "", "", err
	}

	corpClaOrg := &models.CLAOrg{ID: corpClaOrgID}
	if err := corpClaOrg.Get(); err != nil {
		return 0, "", "", err
	}

	if claOrg.Platform != corpClaOrg.Platform ||
		claOrg.OrgID != corpClaOrg.OrgID ||
		claOrg.RepoID != corpClaOrg.RepoID {
		return 400, util.ErrInvalidParameter, "", fmt.Errorf("not the same repo")
	}

	return 0, "", claOrg.OrgEmail, nil
}

func (this *EmployeeSigningController) notifyManagers(corpClaOrgID, employeeEmail, orgEmail, subject string, builder email.IEmailMessageBulder) {
	managers, err := models.ListCorporationManagers(corpClaOrgID, employeeEmail, dbmodels.RoleManager)
	if err != nil {
		beego.Error(err)
		return
	}

	if len(managers) == 0 {
		return
	}

	msg, err := builder.GenEmailMsg()
	if err != nil {
		beego.Error(err)
		return
	}

	to := make([]string, 0, len(managers))
	for _, item := range managers {
		if item.Role == dbmodels.RoleManager {
			to = append(to, item.Email)
		}
	}
	msg.To = to
	msg.Subject = subject

	worker.GetEmailWorker().SendSimpleMessage(orgEmail, msg)
}

func (this *EmployeeSigningController) notifyEmployee(employeeEmail, orgEmail, subject string, builder email.IEmailMessageBulder) {
	msg, err := builder.GenEmailMsg()
	if err != nil {
		beego.Error(err)
		return
	}

	msg.To = []string{employeeEmail}
	msg.Subject = subject

	worker.GetEmailWorker().SendSimpleMessage(orgEmail, msg)
}
