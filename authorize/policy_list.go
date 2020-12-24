package authorize

// PolicyList casbin policy rule
type PolicyList struct {
	ID           string `xorm:"varchar(100) notnull pk" json:"id"`
	RuleType     string `xorm:"varchar(100)" json:"ruleType"`
	Tenant       string `xorm:"varchar(100)" json:"tenant"`
	User         string `xorm:"varchar(100)" json:"user"`
	ResourcePath string `xorm:"varchar(500)" json:"resourcePath"`
	Action       string `xorm:"varchar(500)" json:"action"`
	Service      string `xorm:"varchar(500)" json:"service"`
	AuthEffect   string `xorm:"varchar(500)" json:"authEffect"`
}

// GetPolicyLists get all casbin rules
func GetPolicyLists() []*PolicyList {
	policyLists := []*PolicyList{}
	err := ormManager.engine.Asc("id").Find(&policyLists)
	if err != nil {
		panic(err)
	}
	return policyLists
}

// GetPolicyList get casbin rule
func GetPolicyList(id string) *PolicyList {
	policyList := PolicyList{ID: id}
	existed, err := ormManager.engine.Get(&policyList)
	if err != nil {
		panic(err)
	}

	if existed {
		return &policyList
	}
	return nil
}

// NewPolicyList create new casbin rule
func NewPolicyList() *PolicyList {
	return &PolicyList{
		ID:           "",
		RuleType:     "",
		Tenant:       "",
		User:         "",
		ResourcePath: "",
		Action:       "",
		Service:      "",
		AuthEffect:   "",
	}
}

func createPolicyListTable() error {
	return ormManager.engine.Sync2(new(PolicyList))
}

func dropPolicyListTable() error {
	return ormManager.engine.DropTables(new(PolicyList))
}

// UpdatePolicyList update casbin rule
func UpdatePolicyList(policyList *PolicyList) bool {
	affected, err := ormManager.engine.Insert(policyList)
	if err != nil {
		panic(err)
	}
	return affected != 0
}

// DeletePolicyList delete casbin rule
func DeletePolicyList(policyList *PolicyList) bool {
	affected, err := ormManager.engine.Delete(policyList)
	if err != nil {
		panic(err)
	}
	return affected != 0
}
