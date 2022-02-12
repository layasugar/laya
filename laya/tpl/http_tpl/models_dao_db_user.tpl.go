package http_tpl

const ModelsDaoDbUserTpl = `package db

// User 声明模型
type User struct {
	ID        uint64 {{.tagName}}gorm:"column:id" json:"id"{{.tagName}}
	Username  string {{.tagName}}gorm:"column:username" json:"username"{{.tagName}} // 用户名
	Nickname  string {{.tagName}}gorm:"column:nickname" json:"nickname"{{.tagName}} // 昵称
	Avatar    string {{.tagName}}gorm:"column:avatar" json:"avatar"{{.tagName}}     // 头像
	Password  string {{.tagName}}gorm:"column:password" json:"-"{{.tagName}}        // 密码
	Salt      string {{.tagName}}gorm:"column:salt" json:"-"{{.tagName}}            // 盐
	Mobile    string {{.tagName}}gorm:"column:mobile" json:"mobile"{{.tagName}}     // 手机号
	Status    uint8  {{.tagName}}gorm:"column:status" json:"status"{{.tagName}}     // 状态
	Channel   uint8  {{.tagName}}gorm:"column:channel" json:"channel"{{.tagName}}   // 注册渠道
	CreatedAt string {{.tagName}}gorm:"column:created_at" json:"created_at"{{.tagName}}
	UpdatedAt string {{.tagName}}gorm:"column:updated_at" json:"updated_at"{{.tagName}}
	DeletedAt uint8  {{.tagName}}gorm:"column:deleted_at" json:"deleted_at"{{.tagName}}
}

// TableName 将 User 的表名设置为 {{.tagName}}user{{.tagName}}
func (User) TableName() string {
	return "user"
}
`
