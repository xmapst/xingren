package model

func InitDatabase() {
	DB.AutoMigrate(
		Detail{},
		Clinic{},
		Document{},
	)
}
