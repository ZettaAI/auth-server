package authorize

// Enforcer casbin enforcer
type Enforcer struct {
	ID      string `xorm:"varchar(100) notnull pk" json:"id"`
	Name    string `xorm:"varchar(100)" json:"name"`
	Model   string `xorm:"varchar(100)" json:"model"`
	Adapter string `xorm:"varchar(100)" json:"adapter"`
}

// GetEnforcers get all casbin enforcers from db
func GetEnforcers() []*Enforcer {
	enforcers := []*Enforcer{}
	err := ormManager.engine.Asc("id").Find(&enforcers)
	if err != nil {
		panic(err)
	}

	return enforcers
}

// GetEnforcer get casbin enforcer from db
func GetEnforcer(id string) *Enforcer {
	enforcer := Enforcer{ID: id}
	existed, err := ormManager.engine.Get(&enforcer)
	if err != nil {
		panic(err)
	}

	if existed {
		return &enforcer
	}
	return nil
}

// NewEnforcer create new enforcer object
func NewEnforcer() *Enforcer {
	return &Enforcer{
		ID:      "",
		Name:    "",
		Model:   "",
		Adapter: "",
	}
}

func createEnforcerTable() error {
	return ormManager.engine.Sync2(new(Enforcer))
}

func dropEnforcerTable() error {
	return ormManager.engine.DropTables(new(Enforcer))
}

// UpdateEnforcers update enforcers in db
func UpdateEnforcers(enforcers []*Enforcer) bool {
	err := dropEnforcerTable()
	if err != nil {
		panic(err)
	}

	err = createEnforcerTable()
	if err != nil {
		panic(err)
	}

	affected, err := ormManager.engine.Insert(&enforcers)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

// UpdateEnforcer update single enforcer
func UpdateEnforcer(enforcer *Enforcer) bool {
	affected, err := ormManager.engine.Insert(enforcer)
	if err != nil {
		panic(err)
	}
	return affected != 0
}

// DeleteEnforcer delete given enforcer from db
func DeleteEnforcer(enforcer *Enforcer) bool {
	affected, err := ormManager.engine.Delete(enforcer)
	if err != nil {
		panic(err)
	}
	return affected != 0
}
