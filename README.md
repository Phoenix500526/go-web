# go-web
## 描述

go-web 是一个学习型项目，实现了一个简易的 web 框架。go-web 在 go 语言强大的 http package 的支持下，借鉴了 gin 和 gee 的设计和思想

## 功能
* 利用 Trie 实现了动态路由功能，能够根据某一条路由规则匹配某一类路由，而非某一个路由，主要包含了:
    * 标签匹配, 如 /:lang/doc 可以匹配到 /go/doc, /C++/doc 以及 /java/doc 等路由。
    * 模糊查找，目前仅支持 * 通配符，如 /assets/\*filepath 将会匹配到 /assets 下的所有路由
* 提供组(Group)定义以实现路由的分组控制功能,支持分组嵌套
* 支持中间件功能，目前 go-web 中使用的通用的中间件有两个，分别是 Logger 和 Recovery。前者用于对请求进行一些相应的簿记工作，而后者则用于错误恢复
* 实现了静态资源服务，支持 HTML 模板渲染
* 使用 go test 进行单元测试

## 目录结构
```bash
$ tree
.
├── go.mod
├── LICENSE
├── main.go
├── README.md
├── static
│   ├── css
│   │   └── demo.css
│   └── file1.txt
├── templates
│   ├── arr.tmpl
│   ├── css.tmpl
│   └── custom_func.tmpl
└── webpack
    ├── context.go
    ├── go.mod
    ├── logger.go
    ├── recovery.go
    ├── router.go
    ├── router_test.go
    ├── trie.go
    ├── webpack.go
    └── webpack_test.go

4 directories, 18 files
```
其中，static 与 templates 目录存放静态资源文件及相应的 html 模板。webpack 为 go-web 框架的主要组成，其中以 webpack.go 为主体， 为用户提供框架的入口。main.go 文件展示了 go-web 的一些基本用法，相当于功能测试