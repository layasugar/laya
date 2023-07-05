package common

const ReadmeTpl = `# {{.projectName}}

## 约定

- func返回单独结构体时, 返回该数据得指针
- laya.WebContext与laya.GrpcContext需要全局传递(ctx里面内置了记录日志与链路追踪)
- models/page 业务逻辑
- models/data 实现数据查询组装, 查询在此处完成, 尽量不要使用join(减轻数据库压力), 数据取出后, 可在该层完成组装
- models/dao 基本的请求层, 模型放置层

## 愉快编码
`
