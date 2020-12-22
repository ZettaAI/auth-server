package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/casbin/casbin/v2"
	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	"github.com/casbin/casbin/v2/util"
	"github.com/labstack/echo"
)

const (
	// AuthorizeEP main auth endpoint
	AuthorizeEP = "/can-i/read/table/vnc"
)

// Authorize main
func Authorize(c echo.Context, email string) error {
	domain := c.Request().Header.Get("X-Forwarded-Prefix")
	e, err := casbin.NewEnforcer("casbin/model.conf", "casbin/policy.csv")
	if err != nil {
		log.Fatalf("error: adapter: %s", err)
	}

	rm := e.GetRoleManager()
	rm.(*defaultrolemanager.RoleManager).AddMatchingFunc("regexMatch", util.RegexMatch)
	// rm.PrintRoles()

	uri := c.Request().Header.Get("X-Forwarded-Uri")
	method := c.Request().Header.Get("X-Forwarded-Method")

	log.Print(email)
	log.Print(domain)
	log.Print(method)
	log.Print(uri)

	authorized, err := e.Enforce(email, domain, uri, method)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	if !authorized {
		// not enough permissions 403
		return c.String(
			http.StatusForbidden,
			fmt.Sprintf("User %v does not have the permission to %v %v", email, method, uri),
		)
	}
	return c.JSON(http.StatusOK, authorized)
}
