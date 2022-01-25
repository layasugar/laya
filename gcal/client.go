/**
* @file client.go
* @desc
* @author chenxiaonan01 (zhanglei)
* @version 1.0
* @date 2018-09-03
 */

package cal

import (
	"fmt"
	"time"

	"gitlab.xthktech.cn/bs/gxe/cal/context"
	"gitlab.xthktech.cn/bs/gxe/cal/service"
)

type client struct {
	serv   service.Service
	isCopy bool
	err    error
}

// Client new client by serviceName
func Client(serviceName string) *client {
	serv, _ := service.GetService(serviceName)
	if serv == nil {
		return &client{
			err: fmt.Errorf("fail get service: %s", serviceName),
		}
	}

	return &client{
		serv: serv,
	}
}

// Do 执行
func (c *client) Do(request interface{}, response interface{}, converterType ConverterType) error {
	if c.err != nil {
		return c.err
	}
	ctx := context.NewContext()
	ctx.Caller = "CAL"

	c.err = calWithService(ctx, c.serv, request, response, converterType)
	return c.err
}

// GetIDC 得到IDC
func (c *client) GetIDC() string {
	if c.err != nil {
		return ""
	}
	return c.serv.GetIDC()
}

func (c *client) setPrepare() {
	if !c.isCopy {
		c.serv = c.serv.Clone()
		c.isCopy = true
	}
}

// SetIDC 设置IDC
func (c *client) SetIDC(idc string) *client {
	if c.err != nil {
		return c
	}

	c.setPrepare()

	c.err = c.serv.SetIDC(idc)
	return c
}

// GetProtocol 得到 protocol
func (c *client) GetProtocol() string {
	if c.err != nil {
		return ""
	}
	return c.serv.GetCalConf().Protocol
}

// SetProtocol 设置 Protocol
func (c *client) SetProtocol(p string) *client {
	if c.err != nil {
		return c
	}
	c.setPrepare()
	c.serv.GetCalConf().Protocol = p
	return c
}

// GetStrategy 得到 strategy
func (c *client) GetStrategy() string {
	if c.err != nil {
		return ""
	}
	return c.serv.GetCalConf().Strategy
}

// SetStrategy 设置 strategy
func (c *client) SetStrategy(s string) *client {
	if c.err != nil {
		return c
	}
	c.setPrepare()
	c.serv.GetCalConf().Strategy = s
	return c
}

// GetRetry 得到重试次数
func (c *client) GetRetry() int {
	if c.err != nil {
		return 0
	}
	return c.serv.GetCalConf().GetRetry()
}

// SetRetry 设置 retry
func (c *client) SetRetry(retry int) *client {
	if c.err != nil {
		return c
	}
	c.setPrepare()
	c.serv.GetCalConf().Retry = retry
	return c
}

// GetReuse 得到是否连接
func (c *client) GetReuse() bool {
	if c.err != nil {
		return false
	}
	return c.serv.GetCalConf().Reuse
}

// SetReuse 设置 是否复用连接
func (c *client) SetReuse(doReuse bool) *client {
	if c.err != nil {
		return c
	}
	c.setPrepare()
	c.serv.GetCalConf().Reuse = doReuse
	return c
}

// GetConnTimeOut 得到连接超时
func (c *client) GetConnTimeOut() time.Duration {
	if c.err != nil {
		return 0
	}
	return c.serv.GetCalConf().ConnTimeOut
}

// SetConnTimeOut 设置连接超时
func (c *client) SetConnTimeOut(ct time.Duration) *client {
	if c.err != nil {
		return c
	}
	c.setPrepare()
	c.serv.GetCalConf().ConnTimeOut = ct

	return c
}

// GetReadTimeOut 得到读超时
func (c *client) GetReadTimeOut() time.Duration {
	if c.err != nil {
		return 0
	}
	return c.serv.GetCalConf().ReadTimeOut
}

// SetReadTimeOut 设置读超时
func (c *client) SetReadTimeOut(ct time.Duration) *client {
	if c.err != nil {
		return c
	}
	c.setPrepare()
	c.serv.GetCalConf().ReadTimeOut = ct

	return c
}

// GetWriteTimeOut 得到写超时
func (c *client) GetWriteTimeOut() time.Duration {
	if c.err != nil {
		return 0
	}
	return c.serv.GetCalConf().WriteTimeOut
}

// SetWriteTimeOut 设置写超时
func (c *client) SetWriteTimeOut(ct time.Duration) *client {
	if c.err != nil {
		return c
	}
	c.setPrepare()
	c.serv.GetCalConf().WriteTimeOut = ct

	return c
}
