
## 说明
协议文件结构如下
```
pb/protocol
├── sub1
│   └── item.proto
└── sub2
    └── req.proto
```

要求生成代码到pb_gen_code/目录下，并保持一致的目录结构。请参考协议文件头部的 go_package 标识以及 protoc 命令。