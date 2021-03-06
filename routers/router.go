// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"

	"github.com/opensourceways/app-cla-server/controllers"
)

func init() {
	ns := beego.NewNamespace("/v1",
		/*
			beego.NSNamespace("/cla",
				beego.NSInclude(
					&controllers.CLAController{},
				),
			),
		*/
		beego.NSNamespace("/cla-org",
			beego.NSInclude(
				&controllers.CLAOrgController{},
			),
		),
		beego.NSNamespace("/individual-signing",
			beego.NSInclude(
				&controllers.IndividualSigningController{},
			),
		),
		/*
			beego.NSNamespace("/employee-signing",
				beego.NSInclude(
					&controllers.EmployeeSigningController{},
				),
			),
			beego.NSNamespace("/corporation-signing",
				beego.NSInclude(
					&controllers.CorporationSigningController{},
				),
			),
			beego.NSNamespace("/corporation-manager",
				beego.NSInclude(
					&controllers.CorporationManagerController{},
				),
			),
			beego.NSNamespace("/employee-manager",
				beego.NSInclude(
					&controllers.EmployeeManagerController{},
				),
			),
			beego.NSNamespace("/email",
				beego.NSInclude(
					&controllers.EmailController{},
				),
			),
		*/
		beego.NSNamespace("/auth",
			beego.NSInclude(
				&controllers.AuthController{},
			),
		),
		/*
			beego.NSNamespace("/org-signature",
				beego.NSInclude(
					&controllers.OrgSignatureController{},
				),
			),
		*/
	)
	beego.AddNamespace(ns)
}
