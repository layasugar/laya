## 使用介绍

## 引入日志包

```
import "github.com/layatips/laya/glogs"
```

##

## 日志打印

#### 初始化

```
glogs.InitLog(appName, logType, logPath string)
```

#### 使用

```
glogs.Info(template string, args ...interface{})
glogs.InfoF(c *gin.Context, title string, template string, args ...interface{})
glogs.Warn(template string, args ...interface{})
glogs.WarnF(c *gin.Context, title string, template string, args ...interface{})
glogs.Error(template string, args ...interface{})
glogs.ErrorF(c *gin.Context, title string, template string, args ...interface{})
```

#### 备注

- 带F的方法会记录title和request_id
- logPath 配置如下"/home/logs/app/appName"到appName结束,不带最后一个斜杠

##

## gorm_logger的使用

#### 初始化db连接的时候载入logger

```
  import "gorm.io/gorm"
  import "gorm.io/gorm/logger"
    
  DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
    Logger: glogs.Default(glogs.Sugar, logger.Info),
  })
```

#### 备注

- glogs.Sugar是uber的zap包的*zap.Logger，可用自己实现的，也可用glogs包
- logger.Info是info类型，gorm的logger还提供的warn和error类型

##

## trace （zipkin的链路追踪）

#### 初始化

```
glogs.InitTrace(genv.AppName(), genv.HttpListen(), zipkinAddr, mod)
```

#### 使用

- 配合gin使用,加入已下中间件

```
// gin全局trace中间件
func SetTrace(c *gin.Context) {
	if !gutils.InSliceString(c.Request.RequestURI, gutils.IgnoreRoutes) {
		span := glogs.StartSpanR(c.Request, c.Request.RequestURI)
		if span != nil {
			span.Tag(glogs.RequestIDName, c.GetHeader(glogs.RequestIDName))
			_ = glogs.Inject(context.WithValue(context.Background(), glogs.GetSpanContextKey(), span.Context()), c.Request)
			c.Next()
			span.Finish()
		}
	}
}
```

- 配合gin使用,平级使用,span1和span2是平级

```
span1 := glogs.StartSpanFromReq(c.Request, "我是span1")
time.Sleep(time.Second)
glogs.StopSpan(span1)

span2 := glogs.StartSpanFromReq(c.Request, "我是span2")
time.Sleep(100 * time.Microsecond)
glogs.StopSpan(span2)
```

- 配合gin使用,上下级使用,span3的上级是span2,span2的上级是span1

```
span1 := glogs.StartSpanR(c.Request, "我是span1")
time.Sleep(time.Second)
glogs.StopSpan(span1)

span2 := glogs.StartSpanP(span1.Context(), "我是span2")
time.Sleep(100 * time.Microsecond)
glogs.StopSpan(span2)

span3 := glogs.StartSpanP(span2.Context(), "我是span3")
time.Sleep(200 * time.Microsecond)
glogs.StopSpan(span3)
```

- 单独使用，平级和上下级, span2和span3都是span1的子集,span2和span3是平级

```
span1 := glogs.StartSpan("我是span1")
time.Sleep(time.Second)
glogs.StopSpan(span1)

span2 := glogs.StartSpanP(span1.Context(), "我是span2")
time.Sleep(100 * time.Microsecond)
glogs.StopSpan(span2)

span3 := glogs.StartSpanP(span1.Context(), "我是span3")
time.Sleep(200 * time.Microsecond)
glogs.StopSpan(span3)
```

#### 备注

##

## 钉钉推送

#### 初始化

```
glogs.InitDing(robotKey, robotHost)
```

#### 使用

```
var alarmData = &glogs.AlarmData{
    Title:       "我是一个快乐的机器人",
    Description: "快乐的机器人",
    Content: map[string]interface{}{
        "time": time.Now(),
        "haha": "流弊机器人",
    },
}
glogs.SendDing(alarmData)
```

#### 备注

- 钉钉推送不需要开协程，方法里面已经做了处理