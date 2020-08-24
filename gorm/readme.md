## gorm [docs](https://gorm.io/docs/index.html)

go首选的ORM框架

- 连接池(database/sql的pool)
- 级联
- 钩子
- 预加载
- 事务
- 复合主键
- 自动聚合
- Logger
- 可基于回调扩展自己的插件
- 测试覆盖

### 高级用法 -- 给CRUD操作写回调函数

本身gorm已经写了很多CRUD对应回调，[参考这里](https://gorm.io/docs/write_plugins.html#Pre-Defined-Callbacks)

```go
func updateCreated(scope *Scope) {
  if scope.HasColumn("Created") {
    scope.SetColumn("Created", NowFunc())
  }
}

db.Callback().Create().Register("update_created_at", updateCreated)
// register a callback for Create process
```

Delete a callback from callbacks
```go
db.Callback().Create().Remove("gorm:create")
// delete callback `gorm:create` from Create callbacks
```

Replace a callback having same name with new one
```go
db.Callback().Create().Replace("gorm:create", newCreateFunction)
// replace callback `gorm:create` with new function `newCreateFunction` for Create process
```

Register callbacks with orders
```go
// Before, After参数意思是在这个参数代表回调函数之前/之后注册新的回调函数(参数对应回调不存在也不会报错)
db.Callback().Create().Before("gorm:create").Register("update_created_at", updateCreated)
db.Callback().Create().After("gorm:create").Register("update_created_at", updateCreated)
db.Callback().Query().After("gorm:query").Register("my_plugin:after_query", afterQuery)
db.Callback().Delete().After("gorm:delete").Register("my_plugin:after_delete", afterDelete)
db.Callback().Update().Before("gorm:update").Register("my_plugin:before_update", beforeUpdate)
db.Callback().Create().Before("gorm:create").After("gorm:before_create").Register("my_plugin:before_create", beforeCreate)
```

