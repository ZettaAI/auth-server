package authorize

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

// GetPolicyFromDB get policy from adapter
func GetPolicyFromDB(c echo.Context) error {
	id := c.Param("id")
	log.Printf("fetching policy from adapter %v", id)

	ok, err, res := GetAdapterPolicies(id)
	if err != "" {
		log.Fatal(err)
	}

	if ok {
		log.Print(res)
	}
	return c.String(http.StatusOK, fmt.Sprintf("%v", res))
}
