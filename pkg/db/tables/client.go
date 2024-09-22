package tables

type Client struct {
	ClientPk int `gorm:"primaryKey;column:client_pk"`
}

func (t *Client) TableName() string {
	return "unknown.client"
}
