package models

type Location struct {
	ID   int64
	Name string
}

type Badpic struct {
	ID               int64  `gorm:"PRIMARY_KEY;auto_increment"` //id
	PicPath          string //图片路径 pic_path
	LocationID       int    //违法地点外键 location_id
	LicKind          string
	ShootInstitution string
	Backup           string
	ShootTime        string //违法时间 shoot_time
	RuleType         string //违规类型 rule_type
	Rollback         bool   // 是否经过回滚 rollback
	LicPlate         string //车牌号码 lic_plate
	Direct           string //方向 direct
	Status           bool   //校对1 未校对 0 status
	Location         Location
}

type Goodpic struct {
	ID               int64  `gorm:"PRIMARY_KEY;auto_increment"` //id
	PicPath          string //图片路径 pic_path
	LocationID       int    //违法地点外键 location_id
	LicKind          string
	ShootInstitution string
	Backup           string
	ShootTime        string //违法时间 shoot_time
	RuleType         string //违规类型 rule_type
	Rollback         bool   // 是否经过回滚 rollback
	LicPlate         string //车牌号码 lic_plate
	Direct           string //方向 direct
	Status           bool   //校对1 未校对 0 status
	Location         Location
}

type User struct{
	ID	int64 `gorm:"PRIMARY_KEY;auto_increment"`
	Username string //用户名
	Password string //密码
}

func migration() {
	// 自动迁移模式
	DB.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&Location{}, &Badpic{}, &Goodpic{}，&User{})
}
