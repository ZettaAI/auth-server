package authorize

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	xormadapter "github.com/casbin/xorm-adapter/v2"
)

// CasbinRule casbin policy rule
type CasbinRule struct {
	PType string `xorm:"varchar(100) index not null default ''" json:"pType"`
	V0    string `xorm:"varchar(100) index not null default ''" json:"v0"`
	V1    string `xorm:"varchar(100) index not null default ''" json:"v1"`
	V2    string `xorm:"varchar(100) index not null default ''" json:"v2"`
	V3    string `xorm:"varchar(100) index not null default ''" json:"v3"`
	V4    string `xorm:"varchar(100) index not null default ''" json:"v4"`
	V5    string `xorm:"varchar(100) index not null default ''" json:"v5"`
}

func safeReturn(policy []string, i int) string {
	if len(policy) > i {
		return policy[i]
	}
	return ""
}

func matrixToCasbinRules(pType string, policies [][]string) []*CasbinRule {
	res := []*CasbinRule{}
	for _, policy := range policies {
		line := CasbinRule{
			PType: pType,
			V0:    safeReturn(policy, 0),
			V1:    safeReturn(policy, 1),
			V2:    safeReturn(policy, 2),
			V3:    safeReturn(policy, 3),
			V4:    safeReturn(policy, 4),
			V5:    safeReturn(policy, 5),
		}
		res = append(res, &line)
	}
	return res
}

func getEnforcer(adapter *Adapter) (*casbin.Enforcer, error) {
	a, err := xormadapter.NewAdapter(adapter.Param1, adapter.Param2, true)
	if err != nil {
		return nil, err
	}
	m, err := model.NewModelFromFile("casbin/model.conf")
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(m, a)
	return e, err
}

// GetAdapterPolicies get casbin rules of an adapter
func GetAdapterPolicies(id string) (bool, string, []*CasbinRule) {
	a := GetAdapter(id)
	e, err := getEnforcer(a)
	if err != nil {
		return false, err.Error(), nil
	}
	return true, "", matrixToCasbinRules("p", e.GetPolicy())
}

// GetAdapterGroupingPolicies get casbin group rules of an adapter
func GetAdapterGroupingPolicies(id string) (bool, string, []*CasbinRule) {
	a := GetAdapter(id)
	e, err := getEnforcer(a)
	if err != nil {
		return false, err.Error(), nil
	}
	return true, "", matrixToCasbinRules("g", e.GetGroupingPolicy())
}

// SetAdapterAllPolicies initialize adapter with its rules
func SetAdapterAllPolicies(id string, policies []*CasbinRule) (bool, string) {
	a := GetAdapter(id)
	e, err := getEnforcer(a)
	if err != nil {
		return false, err.Error()
	}

	e.ClearPolicy()
	for _, policy := range policies {
		if policy.PType == "p" {
			_, err = e.AddPolicy(
				policy.V0, policy.V1, policy.V2, policy.V3, policy.V4, policy.V5)
		} else if policy.PType == "g" {
			_, err = e.AddGroupingPolicy(
				policy.V0, policy.V1, policy.V2, policy.V3, policy.V4, policy.V5)
		}
		if err != nil {
			return false, err.Error()
		}
	}
	err = e.SavePolicy()
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

// AddAdapterPolicy add new rule
func AddAdapterPolicy(id string, policy ...string) (bool, string) {
	a := GetAdapter(id)
	e, err := getEnforcer(a)
	if err != nil {
		return false, err.Error()
	}
	_, err = e.AddPolicy(policy)
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}

// RemoveAdapterPolicy remove given rule
func RemoveAdapterPolicy(id string, policy ...string) (bool, string) {
	a := GetAdapter(id)
	e, err := getEnforcer(a)
	if err != nil {
		return false, err.Error()
	}
	_, err = e.RemovePolicy(policy)
	if err != nil {
		return false, err.Error()
	}
	return true, ""
}
