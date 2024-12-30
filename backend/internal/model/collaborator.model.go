package model

type CollaboratorModel struct {
	ID            uint `gorm:"primarykey"`
	CodeSessionId uint
	CodeSession   CodeSessionModel `gorm:"foreignKey:CodeSessionId"`
	UserId        uint
	User          UserModel `gorm:"foreignKey:UserId"`
}
