## gtime时间包

#### 引入时间包

```
import "github.com/layasugar/laya/gtime"
```

#### 在gorm模型中的应用

```
// 声明模型User
type User struct {
	CreatedAt     gtime.Time `json:"created_at"`
	LastLoginTime gtime.Time `json:"last_login_time"`
}
```

#### 初始化一个时间
```
gtime.TimeFrom(t time.Time)
gtime.NewTime(t time.Time, valid bool)
```