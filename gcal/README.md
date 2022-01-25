# RAL
RAL是一个支持多种交互协议和打包格式的扩展包。

RAL规定了一套高度抽象的交互过程规范，将整个后端交互过程分成了交互协议和数据打包/解包两大块，可以支持一些常用的后端交互协议，标准化协议扩充的开发过程，促进代码复用。

RAL集成了负载均衡、健康检查等功能，让上游端不需要再关注这些繁琐的通用逻辑，同时实现版本可以在性能方面有更优的表现。

## 文件配置
RAL的配置文件都放在项目的`conf/cal/services/`路径下。一个典型的RAL配置文件内容如下:
```toml
Name = "aps"
# 连接超时
ConnTimeOut = 1500
# 写数据超时
WriteTimeOut = 1500
# 读数据超时
ReadTimeOut = 1500
# 请求失败后的重试次数：总请求次数 = Retry + 1
Retry = 2
# 数据头协议
Protocol = "http" # http / nshead
#Protocol = "nshead" # http / nshead
# 数据体格式
Converter = "form"  # form / mcpack1 / mapack2 / string
# use http proxy 
HTTPProxy = "" 
# use https proxy 
HTTPSProxy = ""
# 资源使用策略
# random: 纯随机
# roundrobin: 依次轮询（未实现）
# weight-random: 带权重随机（未实现）
# hash: 使用hashid按范围访问（未实现）
Strategy = "random" # random / roundrobin / weight-random / hash
[Resource.BNS]
BNSName = "group.sre-webpage-searchbox-inrouter.orp.all"
EnableSmartBNS = true
# 资源定位：手动配置 - 使用IP、端口
#[Resource.Manual]
#[[Resource.Manual.tc]]
#Host = "cp01-chenxiaonan01.epc.baidu.com"
#Port = 8220
# 资源定位 支持单个机房 单独配置 支持 BNS 和 Manual混用
[Resource.SingleIDC]
[Resource.SingleIDC.BNS]
# 一个机房支持一个BNS配置
[Resource.SingleIDC.BNS.tc]
BNSName = "group.sre-webpage-searchbox-inrouter.orp.all"
EnableSmartBNS = true
#[Resource.SingleIDC.Manual]
#[[Resource.SingleIDC.Manual.tc]]
#Host = "cp01-chenxiaonan01.epc.baidu.com"
#Port = 34231
```
* `Name`: 每个配置都有一个名字，用来与其他的配置进行区分，RAL包会根据这个名字来查找对应的配置内容。
* `ConnTimeOut`: 连接超时时间配置，单位毫秒
* `WriteTimeOut`: 写超时时间配置, 单位毫秒
* `ReadTimeOut`: 写超时时间配置，单位毫秒
* `Retry`: 请求重试次数
* `Protocol`: 网络交互的协议，支持`http`,`https`,`nshead`等
* `Converter`: 数据打包格式,支持表单(`form`), `mcpack2`, 字符串(`string`)等，这里只是一个默认的，还可以在发送RAL请求的时候指定打包格式，因为可能一个BNS对应多种打包格式，这样更加灵活。
* `Strategy`: 负载均衡策略，目前只支持`random`策略，这种策略会有IDC优先级的考虑，默认是`LocalIDC` > `DefaultIDC` > `BackupIDC` > `AllIDC`
* `Resource.BNS`: 支持BNS的方式访问服务端
* `Resource.BNS/BNSName`: 百度内部服务对应的BNS名称。
* `Resource.BNS/EnableSmartBNS`: 是否使用智能BNS。参考[BNS](/cal/cal.md#BNS)， 如果遇到BNS解析错误的问题，先确认此项是否填写正确，不知道是否是智能BNS可以与业务线OP确认。
* `Resource.Manual`: 支持手动配置BNS
* `Resource.Manula.idcmap`: idcmap是指机房对应的名称，表示这个机房下的机器配置
* `Resource.Manual.idcmap/Host`: 此idcmap下的机器的Host地址
* `Resource.Manual.idcmap/Post`: 此idcmap下的机器的端口
* `Resource.SingleIDC`: 支持单独IDC配置方式
* `Resource.SingleIDC.BNS`: 支持单独IDC配置BNS
* `Resource.SingleIDC.BNS/BNSName`: 百度内部服务对应的BNS名称。
* `Resource.SingleIDC.BNS/EnableSmartBNS`: 是否使用智能BNS。参考[BNS](/cal/cal.md#BNS)， 如果遇到BNS解析错误的问题，先确认此项是否填写正确，不知道是否是智能BNS可以与业务线OP确认。
* `Resource.SingleIDC.Manual`: 手动配置单机房BNS
* `Resource.SingleIDC.Manual.idcmap`:idcmap是指机房对应的名称，表示这个机房下的机器配置
* `Resource.SingleIDC.Manual.idcmap/Host`: 此idcmap下的机器的Host地址
* `Resource.SingleIDC.Manual.idcmap/Post`: 此idcmap下的机器的端口


## 协议与打包格式
* 目前支持的协议是`http`, `https`, `nshead`。`nshead`是公司内部使用比较广泛的一个交互协议，其本质是`tcp`协议，只不过是自定义了一个协议头。
* 打包格式目前支持`form`, `mcpack2`, `string`等方式。
一般来说`nshead + mcpack2`组成一个对进行请求。`http/https`与`form/string`等自由组合。

## cal.Cal调用
我们通过`cal.Cal`函数调用来发送请求，其定义如下:
```go
Cal(serviceName string, request interface{}, response interface{}, converterType ConverterType) error
```
下面对其参数进行分析:
* `serviceName`: 这个`conf/cal/services/`目录下cal的配置文件的`Name`字段。在初始化阶段配置文件的内容及其对应的BNS的解析已经完成，放到了一个map中，Cal函数内部会从map中取出对应的对象进行使用。
* `request`: 这个是发送的请求的数据结构，针对不同的协议其格式也不一样，下面会详细说明。
* `response`: 这里定义一个请求返回的结构体，当请求成功后会把结果赋值到这个变量中，这样就可以方便的使用返回的结果了。
* `ConverterType`: 这个是要指定请求的服务器返回的数据的打包类型，之后正确指定这个类型cal.Cal内部才会对其进行正确的解析赋值。

下面针对不同协议进行详细说明：

### HTTP/HTTPS请求 {#HTTP/HTTPS}
一个请求的事例:
```go
package main

import (
        "fmt"
        _ "github.com/layasugar/laya"
        "github.com/layasugar/laya/gcal"
)

func main() {
        header := map[string][]string{
                "Content-Type": {"application/x-www-form-urlencoded"},
        }
        body := map[string]string{
                "abc": "123123123",
                "xxx": "xxxxxxx",
        }

        request := cal.HTTPRequest{
                Header:    header,
                Method:    "POST",
                Path:      "webpage?type=user&action=home&format=json&uk=xN_v9qFPpxXG3LGD3wJoaQ",
                Body:      body,
                Converter: cal.FORMConverter,
        }

        type Res struct {
                Errno     int    `json:"errno"`
                RequestID string `json:"request_id"`
                Data      struct {
                        User struct {
                                Name string `json:"display_name"`
                        } `json:"user"`
                } `json:"data"`
        }

        type resp struct {
                Head    cal.HTTPHead
                Body    Res
                Raw     []byte
                TagHead cal.HTTPHead `cal:"head"`
                TagBody Res          `cal:"body"`
                TagRaw  []byte       `cal:"raw"`
        }
        var response = resp{}
        err := cal.Cal("aps", request, &response, cal.JSONConverter)
        if err != nil {
                fmt.Printf("%v", err)
        } else {
                fmt.Printf("Head: %v\n", response.Head)
                fmt.Printf("Body: %v\n", response.Body)
                fmt.Printf("Raw: %v\n", response.Raw)
                fmt.Printf("TagHead: %v\n", response.TagHead)
                fmt.Printf("TagBody: %v\n", response.TagBody)
                fmt.Printf("TagRaw: %v\n", response.TagRaw)
        }
}
```
我们通过`cal.Cal`函数调用发送请求。当我们通过HTTP/HTTPS协议进行交互时， 我们的请求结构体必须是`cal.HTTPRequest`类型,其定义如下:
```go
type Request struct {
    CustomHost string
    CustomPort int
    Header     map[string][]string
    Method     string
    Body       interface{}
    Path       string
    LogID      int64
    Converter  converters.ConverterType
}
* `CustomHost`和`CustomPort`分别对应自定义的host和port, 当传递这两个参数时就会优先向这个server发送请求，如果没有指定就会根据Cal的配置文件BNS发送请求。
* `Header`：发送请求的header,这里的类型与`net/http`保持了一致，便于统一处理。
* `Method`: 发送请求的方法，这里支持http的标准方法，包括:`GET`, `HEAD`, `POST`, `PUT`, `DELETE`, `CONNECT`, `OPTIONS`, `TRACE`。
* `Converters`: 发送请求是以什么样的打包形式发送请求.
* `Body`: 发送请求的body, 其格式与对应的converterType有关系，他们的关系如下:

|Converter|Body|
|-:-|-:-|
|`FORMConverter`|`map[string]string`, `net/url.Values`|
|`JSONConverter`|`string`,`[]byte`|
|`RAWConverter`|`string`,`[]byte`|

* `Path`: 要访问的url地址
* `LogID`: 透传logid, 便于串联请求，追踪日志。

`response`结构字段是固定的:
```go
type resp struct {
        Head    cal.HTTPHead
        Body    Res
        Raw     []byte
}
```
* `Head`:是返回的结果的头
* `Body`:是返回的具体内容，需要提前知道返回的结果格式并使用对应的`Converter`进行接收
* `Raw` :会返回对应的Body的`[]byte`结构，可以在cal不支持对应的结构解析时自己解析或者返回数据解析不对需要看原始数据时使用。

另外`response`还支持`tag`:
```go
type resp struct {
        TagHead cal.HTTPHead `cal:"head"`
        TagBody Res          `cal:"body"`
        TagRaw  []byte       `cal:"raw"`
}
```
支持的`tag`有`head`,`body`,`raw`与之前的字段`Head`,`Body`,`Raw`一一对应。
Cal解析的顺序是先解析数据的到字段，再解析`tag`，在冲突的情况下，后面的会覆盖前面的值。
### NSHEAD请求
请求事例:

```go
package main

import (
    "fmt"
    _ "github.com/layasugar/laya"
    "github.com/layasugar/laya/gcal"
)

func main() {

    type Con struct {
        Name string
    }
    con := Con{"Request"}
    request := cal.NSHEADRequest{
        Body:      con,
        Converter: cal.MCPACK2Converter,
    }

    type Ns struct {
        Name string
        Age  int
    }

    type resp struct {
        Body Ns
    }
    var response = resp{}

    cal.Cal("ns", request, &response, cal.MCPACK2Converter)
    fmt.Printf("%v", response.Body)
}
```
我们通过`cal.Cal`函数调用发送请求。当我们通过`NSHEAD`协议进行交互时， 我们的请求结构体必须是`cal.NSHEADRequest`类型,其定义如下:
```go
type Request struct {
    CustomHost string
    CustomPort int
    Body       interface{}
    Converter  converters.ConverterType
}
```
* `CustomHost`和`CustomPort`分别对应自定义的host和port, 当传递这两个参数时就会优先向这个server发送请求，如果没有指定就会根据Cal的配置文件BNS发送请求。
* `Converters`: 发送请求是以什么样的打包形式发送请求.
* `Body`: 发送请求的body, 其格式与对应的converterType有关系，他们的关系如下:

|Converter|Body|
|-:-|-:-|
|`MCPACK2Converter`|`interface{}`|
这里的Body其实是要访问的Server定义的结构体。

`response`结构字段是固定的,与[HTTP/HTTPS](/cal/cal.md#HTTP/HTTPS)一样，唯一不一样的是Body是一个自定义的结构提，根据Server的返回来自己定义。

