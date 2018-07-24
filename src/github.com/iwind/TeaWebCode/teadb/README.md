# 存储
## 配置
~~~
$ = {
    "id": "000001" // 自增ID
}
~~~

## 数据
格式：[DATA TYPE].DATA.[DATA ID]=[DATA]
例子：
~~~
LOG.DATA.1 = {
    "meta": {
        "id": "自动生成的ID",
        "modifiedAt": "最后修改时间，单位为纳秒",
        "indexes": {
            "name": [ "719109643.123" ],
            "books": [ "2742599054.123", "2830114965.123", "3224789851.123" ],
            "orders": [ "2707236321.123", "2212294583.123" ]
        }
    },
    "value": {    // 数据
        "name": "Liu",
        "age": 20,
        "books": [
           {
             "name": GoLang"
           },
           {
            "name": "Python"
           },
           {
            "name": "Java"
           }
        ]
    }        
}

LOG.DATA.2 = {
    ...
}
~~~

## 索引
格式：[DATA TYPE].INDEX.[[KEY]].[FIELD CRC32 VALUE].[DATA ID]=[Field VALUE]
~~~
LOG.INDEX.[name].719109643.1 = "liu"
LOG.INDEX.[name].2835264487.2 = "lu"
LOG.INDEX.[name].634732029.3 = "ping"
LOG.INDEX.[books.name].142774187.1 = "GoLang"
...
~~~

## API
### 打开数据库
~~~
var db = teadb.Open("my.db")
defer db.Close()
~~~

### 写入内容
~~~
db.Put(dataType string, data map[string]interface{})
~~~

### 设置部分字段内容
~~~
db.Update(dataType string, id int64, field string, value interface{])
db.Append(dataType string, id int64, field string, value string) // 待实现
db.Increase(dataType string, id int64, field string, delta int64) // 待实现
~~~

### 删除内容
~~~
db.Delete(dataType, id int)
~~~

### 读取内容
~~~
var one = db.NewQuery("log").Attr("name", "Liu").FindOne()
...
~~~
支持以下方法：
* `Attr(field string, value interface{})` - 指定查询条件
* `In(field string, value []interface{})` - 指定查询条件
* `Offset(offset int64)` - 限定开始读取的位置
* `Limit(size int64)` - 限定读取的条数
* `Find()` - 查找单条记录
* `FindAll()` - 查找所有记录
* `FindField(field string)` - 查询某个字段值
* `FindFields(fields ... string)` - 查询某些字段值 // 待实现
* `Delete()` - 删除匹配的记录
* `Replace(value map[string]interface{})` - 更新全部数据 // 待实现
* `Update(value map[string]interface{})` - 更新部分字段数据 // 待实现
* `Result(fields ... string)` - 指定返回的字段 // 待实现
* `Count()` - 计算记录数量


