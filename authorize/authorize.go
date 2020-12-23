package authorize

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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
	authorized := enforce(email, domain, method, uri)
	if !authorized {
		// not enough permissions 403
		return c.String(
			http.StatusForbidden,
			fmt.Sprintf(
				"User %v does not have the permission to %v %v", email, method, uri),
		)
	}
	// add forward headers for backend use
	c.Response().Header().Set("X-Forwarded-User", email)
	c.Response().Header().Set("X-Forwarded-Domain", domain)
	return c.String(http.StatusOK, strconv.FormatBool(authorized))
}

func enforce(email string, domain string, method string, uri string) bool {
	e, err := casbin.NewEnforcer("casbin/model.conf", "casbin/policy.csv")
	if err != nil {
		log.Fatalf("error: adapter: %s", err)
	}

	rm := e.GetRoleManager()
	rm.(*defaultrolemanager.RoleManager).AddMatchingFunc("regexMatch", util.RegexMatch)
	// rm.PrintRoles()

	authorized, err := e.Enforce(email, domain, uri, method)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	return authorized
}
