package tables

type Task struct {
	TaskPk   int    `gorm:"primaryKey;column:task_pk"`
	ClientFk int    `gorm:"column:client_fk"`
	Client   Client `gorm:"foreignKey:client_fk;references:client_pk"`
}

func (t *Task) TableName() string {
	return "unknown.task"
}
