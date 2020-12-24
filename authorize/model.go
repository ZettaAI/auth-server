package authorize

// Model casbin model entity in db
type Model struct {
	ID   string `xorm:"varchar(100) notnull pk" json:"id"`
	Name string `xorm:"varchar(100)" json:"name"`
	Type string `xorm:"varchar(100)" json:"type"`
	Text string `xorm:"varchar(5000)" json:"text"`
}

// GetModels get all models from db
func GetModels() []*Model {
	models := []*Model{}
	err := ormManager.engine.Asc("id").Find(&models)
	if err != nil {
		panic(err)
	}

	return models
}

// GetModel get model from db
func GetModel(id string) *Model {
	model := Model{ID: id}
	existed, err := ormManager.engine.Get(&model)
	if err != nil {
		panic(err)
	}

	if existed {
		return &model
	}
	return nil
}

// NewModel create new model object
func NewModel() *Model {
	return &Model{
		ID:   "",
		Name: "",
		Type: "",
		Text: "",
	}
}

func createModelTable() error {
	return ormManager.engine.Sync2(new(Model))
}

func dropModelTable() error {
	return ormManager.engine.DropTables(new(Model))
}

// UpdateModels update all models in db
func UpdateModels(models []*Model) bool {
	err := dropModelTable()
	if err != nil {
		panic(err)
	}

	err = createModelTable()
	if err != nil {
		panic(err)
	}

	affected, err := ormManager.engine.Insert(&models)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

// UpdateModel update given model in db
func UpdateModel(model *Model) bool {
	affected, err := ormManager.engine.Insert(model)
	if err != nil {
		panic(err)
	}
	return affected != 0
}

// DeleteModel delete given model in db
func DeleteModel(model *Model) bool {
	affected, err := ormManager.engine.Delete(model)
	if err != nil {
		panic(err)
	}
	return affected != 0
}
