package authorize

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	"github.com/casbin/casbin/v2/util"
	"github.com/labstack/echo"
)

// Authorize main
func Authorize(c echo.Context, email string) error {
	uri := c.Request().Header.Get("X-Forwarded-Uri")
	method := c.Request().Header.Get("X-Forwarded-Method")
	domain := c.Request().Header.Get("X-Forwarded-Prefix")

	// add forward headers for backend use
	c.Response().Header().Set("X-Forwarded-Uri", uri)
	c.Response().Header().Set("X-Forwarded-Method", method)
	c.Response().Header().Set("X-Forwarded-User", email)
	c.Response().Header().Set("X-Forwarded-Domain", domain)

	if strings.Contains(uri, "logout") {
		return c.String(http.StatusOK, "success")
	}

	authorized := enforce(email, domain, method, uri)
	if !authorized {
		// not enough permissions 403
		return c.String(
			http.StatusForbidden,
			fmt.Sprintf(
				"User %s not authorized to %s %s in domain %s", email, method, uri, domain),
		)
	}
	return c.String(
		http.StatusOK,
		fmt.Sprintf("User %s authorized to %s %s in domain %s", email, method, uri, domain),
	)
}

func enforce(email string, domain string, method string, uri string) bool {
	e, err := casbin.NewEnforcer("casbin/model.conf", "casbin/policy.csv")
	if err != nil {
		// if authorization model is missing assume every request is authorized
		log.Println("No policy, allow by default.")
		return true
	}

	rm := e.GetRoleManager()
	rm.(*defaultrolemanager.RoleManager).AddMatchingFunc("KeyMatch2", util.KeyMatch2)
	rm.(*defaultrolemanager.RoleManager).AddDomainMatchingFunc("KeyMatch2", util.KeyMatch2)

	authorized, err := e.Enforce(email, domain, uri, method)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	return authorized
}
