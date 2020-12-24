package authorize

// Adapter casbin adapter
type Adapter struct {
	ID            string   `xorm:"varchar(100) notnull pk" json:"id"`
	Name          string   `xorm:"varchar(100)" json:"name"`
	Type          string   `xorm:"varchar(100)" json:"type"`
	Param1        string   `xorm:"varchar(500)" json:"param1"`
	Param2        string   `xorm:"varchar(500)" json:"param2"`
	PolicyHeaders []string `json:"policyHeaders"`
}

// GetAdapters get all casbin adapters
func GetAdapters() []*Adapter {
	adapters := []*Adapter{}
	err := ormManager.engine.Asc("id").Find(&adapters)
	if err != nil {
		panic(err)
	}
	return adapters
}

// GetAdapter get casbin adapter
func GetAdapter(id string) *Adapter {
	adapter := Adapter{ID: id}
	existed, err := ormManager.engine.Get(&adapter)
	if err != nil {
		panic(err)
	}

	if existed {
		return &adapter
	}
	return nil
}

// NewAdapter returns new adapter
func NewAdapter() *Adapter {
	return &Adapter{
		ID:     "",
		Name:   "",
		Type:   "",
		Param1: "",
		Param2: "",
	}
}

func createAdapterTable() error {
	return ormManager.engine.Sync2(new(Adapter))
}

func dropAdapterTable() error {
	return ormManager.engine.DropTables(new(Adapter))
}

// UpdateAdapters update casbin adapters
func UpdateAdapters(adapters []*Adapter) bool {
	err := dropAdapterTable()
	if err != nil {
		panic(err)
	}

	err = createAdapterTable()
	if err != nil {
		panic(err)
	}

	affected, err := ormManager.engine.Insert(&adapters)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

// UpdateAdapter update casbin adapter
func UpdateAdapter(adapter *Adapter) bool {
	affected, err := ormManager.engine.Insert(adapter)
	if err != nil {
		panic(err)
	}
	return affected != 0
}

// DeleteAdapter delete casbin adapter
func DeleteAdapter(adapter *Adapter) bool {
	affected, err := ormManager.engine.Delete(adapter)
	if err != nil {
		panic(err)
	}
	return affected != 0
}
