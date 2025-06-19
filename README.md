@@ 1,410
# 基于Gin、Gorm、Vue 实现的在线练习系统
> 演示地址：[http://getcharzp.gitee.io/gin-gorm-oj](http://getcharzp.gitee.io/gin-gorm-oj)
>
> 后台语言：Golang、框架：Gin、GORM
>
> 前台框架：Vue、ElementUI
>

## 参考链接
> GO下载网址： [https://golang.google.cn/dl/](https://golang.google.cn/dl/)
>
> 参考文档 Module：[https://www.kancloud.cn/aceld/golang/1958311](https://www.kancloud.cn/aceld/golang/1958311)
>
> GORM中文官网：[https://gorm.io/zh_CN/docs/](https://gorm.io/zh_CN/docs/)
>
> GIN中文官网：[https://gorm.io/zh_CN/docs/](https://gorm.io/zh_CN/docs/)
>
> INI中文官网 : [https://ini.unknwon.io/docs/intro](https://ini.unknwon.io/docs/intro)
>
> GIN : [https://github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
>
> GORM : [https://github.com/go-gorm/gorm](https://github.com/go-gorm/gorm)
>
> JWT : [https://github.com/golang-jwt/jwt](https://github.com/golang-jwt/jwt)
>
> EMAIL : [https://github.com/jordan-wright/email](https://github.com/jordan-wright/email)
>
> UUID : [https://github.com/satori/go.uuid](https://github.com/satori/go.uuid)
>
> SWAG : [https://github.com/swaggo/swag](https://github.com/swaggo/swag)
>
> GIN-SWAGGWE :  [https://github.com/swaggo/gin-swagger](https://github.com/swaggo/gin-swagger)
>



## 安装指南
### Gin框架使用
使用 Go 的模块支持 ，  
当你在代码中添加 import 时，go [build|run|test] 会自动获取必要的依赖项：

```go
import "github.com/gin-gonic/gin"
```

```go
go get -u github.com/gin-gonic/gin
```

**运行 Gin**

```go
package main

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()
  r.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })
  r.Run() // 在 0.0.0.0:8080 监听（Windows系统使用 "localhost:8080"）
}
```

**运行方式**

1. 使用 `go run` 命令执行：

```bash
go run example.go
```

### gorm框架使用
#### 安装gorm
```go
go get -u gorm.io/gorm
```

#### 安装mysql驱动
```go
go get -u gorm.io/driver/mysql
```

```go
package main

import (
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
)

type Product struct {
    gorm.Model
    Code  string
    Price uint
}

func main() {
    db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }

    // 迁移 schema
    db.AutoMigrate(&Product{})

    // Create
    db.Create(&Product{Code: "D42", Price: 100})

    // Read
    var product Product
    db.First(&product, 1) // 根据整型主键查找
    db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录

    // Update - 将 product 的 price 更新为 200
    db.Model(&product).Update("Price", 200)
    // Update - 更新多个字段
    db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
    db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

    // Delete - 删除 product
    db.Delete(&product, 1)
}
```

### JWt
1. 可以使用下面的命令将 jwt-go 添加为 Go 程序中的依赖项。

```go
go get -u github.com/golang-jwt/jwt/v5
```

2. 将其导入到您的代码中：

```go
import "github.com/golang-jwt/jwt/v5"
```

[<font style="color:rgb(9, 105, 218);">解析和验证令牌的简单示例</font>](https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-Parse-Hmac)

```go
// sample token string taken from the New example
tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

// Parse takes the token string and a function for looking up the key. The latter is especially
// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
// head of the token to identify which key to use, but the parsed token (head and claims) is provided
// to the callback, providing flexibility.
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
	return hmacSampleSecret, nil
}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
if err != nil {
	log.Fatal(err)
}

if claims, ok := token.Claims.(jwt.MapClaims); ok {
	fmt.Println(claims["foo"], claims["nbf"])
} else {
	fmt.Println(err)
}

```

[<font style="color:rgb(9, 105, 218);">构建和签署令牌的简单示例</font>](https://pkg.go.dev/github.com/golang-jwt/jwt/v5#example-New-Hmac)

```go
// Create a new token object, specifying signing method and the claims
// you would like it to contain.
token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	"foo": "bar",
	"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
})

// Sign and get the complete encoded token as a string using the secret
tokenString, err := token.SignedString(hmacSampleSecret)

fmt.Println(tokenString, err)

```

### emil
```go
get github.com/jordan-wright/email
```

```go
e := email.NewEmail()
e.From = "Jordan Wright <test@gmail.com>"
e.To = []string{"test@example.com"}
e.Bcc = []string{"test_bcc@example.com"}
e.Cc = []string{"test_cc@example.com"}
e.Subject = "Awesome Subject"
e.Text = []byte("Text Body is, of course, supported!")
e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "test@gmail.com", "password123", "smtp.gmail.com"))
```

### uuid
```go
go get github.com/satori/go.uuid
```

```go
package main

import (
    "fmt"
    "github.com/satori/go.uuid"
)

func main() {
    // Creating UUID Version 4
    // panic on error
    u1 := uuid.Must(uuid.NewV4())
    fmt.Printf("UUIDv4: %s\n", u1)

    // or error handling
    u2, err := uuid.NewV4()
    if err != nil {
        fmt.Printf("Something went wrong: %s", err)
        return
    }
    fmt.Printf("UUIDv4: %s\n", u2)

    // Parsing UUID from string input
    u2, err := uuid.FromString("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
    if err != nil {
        fmt.Printf("Something went wrong: %s", err)
        return
    }
    fmt.Printf("Successfully parsed: %s", u2)
}
```

### swag
1. 向 API 源代码添加注释，[请参阅 声明性注释格式](https://github.com/swaggo/swag/blob/master/README.md#declarative-comments-format)。

```go
// GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":"200","msg","","data":""}"
// @Router /problem-list [get]
```

2. 使用以下方法下载 [Swag](https://github.com/swaggo/swag) for Go：

```go
go get -u github.com/swaggo/swag/cmd/swag
```



从 Go 1.17 开始，使用 go get 安装可执行文件已被弃用。可以使用 `go install` 代替：

```go
go install github.com/swaggo/swag/cmd/swag@latest
```



在你的 Go 项目根路径（例如 `~/root/go-project-name`）运行 [Swag](https://github.com/swaggo/swag)， [Swag](https://github.com/swaggo/swag) 将在 `<font style="color:rgb(240, 246, 252);background-color:rgba(101, 108, 118, 0.2);">~</font>/root/go-project-name/docs` 解析评论并生成所需的文件（`docs` 文件夹和 `docs/doc.go`）。

```go
swag init
```



使用以下方法下载 [gin-swagger](https://github.com/swaggo/gin-swagger)：

```go
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
```



在代码中导入以下内容：

```go
import "github.com/swaggo/gin-swagger" // gin-swagger middleware
import "github.com/swaggo/files" // swagger embed files
```



### <font style="color:rgb(55,163,240)">  规范示例：</font>
<font style="color:rgb(55,163,240);">现在假设您已经实现了一个简单的 API，如下所示：</font>

```plain
// A get function which returns a hello world string by json
func Helloworld(g *gin.Context)  {
    g.JSON(http.StatusOK,"helloworld")
}
```

1. <font style="color:rgb(55,163,240);">使用 gin-swagger 规则为 api 和 main 函数添加注释，如下所示：</font>

```plain
// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Helloworld(g *gin.Context)  {
    g.JSON(http.StatusOK,"helloworld")
}
```

1. <font style="color:rgb(55,163,240);">使用 </font>`swag init`<font style="color:rgb(55,163,240);"> 命令生成一个文档，生成的文档会存储在 </font>`docs`<font style="color:rgb(31, 35, 40);"> 目录下。</font>
2. <font style="color:rgb(55,163,240);">按如下方式导入文档：我假设您的项目名为 </font>`github.com/go-project-name/docs`

```plain
import (
   docs "github.com/go-project-name/docs"
)
```

1. <font style="color:rgb(55,163,240);">构建您的应用程序，然后转到 </font>[<font style="color:rgb(9, 105, 218);">http://localhost:8080/swagger/index.html</font>](http://localhost:8080/swagger/index.html)<font style="color:rgb(31, 35, 40);"> ，查看您的 Swagger UI。</font>
2. <font style="color:rgb(55,163,240)">完整的代码和文件夹相对值如下：</font>

```plain
package main

import (
   "github.com/gin-gonic/gin"
   docs "github.com/go-project-name/docs"
   swaggerfiles "github.com/swaggo/files"
   ginSwagger "github.com/swaggo/gin-swagger"
   "net/http"
)
// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /example/helloworld [get]
func Helloworld(g *gin.Context)  {
   g.JSON(http.StatusOK,"helloworld")
}

func main()  {
   r := gin.Default()
   docs.SwaggerInfo.BasePath = "/api/v1"
   v1 := r.Group("/api/v1")
   {
      eg := v1.Group("/example")
      {
         eg.GET("/helloworld",Helloworld)
      }
   }
   r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
   r.Run(":8080")

}
```

```plain
.
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
└── main.go
```

## 系统模块
- [x] 用户模块
    - [x] 密码登录
    - [x] 邮箱注册
    - [x] 用户详情
- [x] 题目管理模块
    - [x] 题目列表、题目详情
    - [x] 题目创建、题目修改
- [x] 分类管理模块
    - [x] 分类列表
    - [x] 分类创建、分类修改、分类删除
- [x] 判题模块
    - [x] 提交记录列表
    - [x] 代码的提交及判断
- [x] 排名模块
    - [x] 排名的列表情况
- [x] 竞赛模块
    - [x] 竞赛列表
    - [x] 竞赛管理
    - [x] 竞赛报名
