# 项目名称
Exchange Cache提供基于内存的缓存功能，主要用于缓存一些使用频率较高，且不会更新或者对更新时效性要求不高的数据。

## 快速开始

### 使用LRUCache

* 创建一个LRUCache
`func NewLRUCache(size int, age int, expires time.Duration) *LRUCache`
size: Cache中能存放的最大记录量
age: 每条记录最多可被使用的次数，如果为0表示不限制使用次数
expires: 每条记录最长的有效期，如果为0表示不限制有效期

* 写记录入LRUCache
`func (lc *LRUCache) Set(key, val interface{}) {`

* 从LRUCache中读记录
`func (lc *LRUCache) Get(key interface{}) (interface{}, bool) {`

* 清空LRUCache
`func (lc *LRUCache) Purge() {`

* 统计LRUCache的计数
`func (lc *LRUCache) Count() (cntItems, cntQuery, cntHit uint64) {`

* 示例
```golang
// 设置一个最大容量为2000条记录，每条最多可以使用10次，且仅生效一分钟的LRUcache
mycache := excache.NewLRUCache(2000, 10, 60*time.Second)
// 将记录写入缓存
mycache.Set("someKey", "someValue")
// 从缓存中获取值
if v, ok := mycache.Get("someKey"); ok {
    print(v.(string))
}
```

## 测试
如何执行自动化测试

## 如何贡献
贡献patch流程、质量要求

## 讨论
百度Hi讨论群：XXXX

## 链接
[百度golang代码库组织和引用指南](http://wiki.baidu.com/pages/viewpage.action?pageId=515622823)

