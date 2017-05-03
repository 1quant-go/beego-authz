package authz

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/hsluoyz/casbin/api"
	"net/http"
)

// NewBasicAuthorizer returns the authorizer.
// Use a casbin model file and a casbin policy file as input
func NewBasicAuthorizer(modelPath string, policyPath string) beego.FilterFunc {
	return func(ctx *context.Context) {
		e := &api.Enforcer{}
		e.InitWithFile(modelPath, policyPath)
		a := &BasicAuthorizer{enforcer: e}

		if !a.CheckPermission(ctx.Request) {
			a.RequirePermission(ctx.ResponseWriter)
		}
	}
}

// NewAuthorizer returns the authorizer.
// Use a casbin enforcer as input
func NewAuthorizer(e *api.Enforcer) beego.FilterFunc {
	return func(ctx *context.Context) {
		a := &BasicAuthorizer{enforcer: e}

		if !a.CheckPermission(ctx.Request) {
			a.RequirePermission(ctx.ResponseWriter)
		}
	}
}

// BasicAuthorizer stores the casbin handler
type BasicAuthorizer struct {
	enforcer *api.Enforcer
}

// GetUserName gets the user name from the request.
// Currently, only HTTP basic authentication is supported
func (a *BasicAuthorizer) GetUserName(r *http.Request) string {
	username, _, _ := r.BasicAuth()
	return username
}

// CheckPermission checks the user/method/path combination from the request.
// Returns true (permission granted) or false (permission forbidden)
func (a *BasicAuthorizer) CheckPermission(r *http.Request) bool {
	user := a.GetUserName(r)
	method := r.Method
	path := r.URL.Path
	return a.enforcer.Enforce(user, path, method)
}

// RequirePermission returns the 403 Forbidden to the client
func (a *BasicAuthorizer) RequirePermission(w http.ResponseWriter) {
	w.WriteHeader(403)
	w.Write([]byte("403 Forbidden\n"))
}
