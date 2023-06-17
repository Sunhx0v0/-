package mysqldb

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Password string `json:"password"` // 密码
	UserName string `json:"userName"` // 用户名
}

// 网站用户登录
func Login(ctx *gin.Context) {
	// 传入并获取前端数据
	var form LoginRequest
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if form.UserName == "" {
		ctx.JSON(400, gin.H{"message": "请输入用户名"})
		return
	}
	if form.Password == "" {
		ctx.JSON(400, gin.H{"message": "请输入密码"})
		return
	}
	// 在数据库中查找用户
	username := form.UserName
	sqlstr := "select password from easy where username = ?"
	rows, err := db.Query(sqlstr, username)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		ctx.JSON(400, gin.H{"message": "没有该用户"})
		return
	}
	var password string
	for rows.Next() {
		err := rows.Scan(&password)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			ctx.JSON(400, gin.H{"message": "用户解析失败"})
			return
		}
	}

	// 关闭rows释放持有的数据库链接
	defer rows.Close()
}

// // 定义结构
// type user struct {
// 	Username string `json:"username"`
// 	Password string `json:"password"`
// }

// // 插入数据
// func insert(ctx *gin.Context) {
// 	// 连接mongodb服务
// 	url := "mongodb://127.0.0.1"
// 	// 设置数据库一致性模式
// 	// 连接数据库操作，该操作赋值给session
// 	// err值必写，用于错误处理
// 	session, err := mgo.Dial(url)
// 	// 后边程序执行的err与go程序比对，若有错误则返回错误内容
// 	if err != nil {
// 		panic(err)
// 	} else {
// 		// 若没有错误，则在页面返回字符串，显示插入成功
// 		ctx.String(http.StatusOK, "插入成功")
// 	}
// 	// defer用法大家自行百度，我解释不清
// 	defer session.Close()

// 	// 设置数据库一致性模式，就当作打开数据库
// 	session.SetMode(mgo.Monotonic, true)
// 	// 找到某数据库下的某数据表
// 	c := session.DB("db_go").C("user")

// 	// 将insert状态传值给err
// 	err = c.Insert(&user{"admin", "123456"}, &user{"Johnson", "Johnson"})
// }

// // 查询数据
// func find(ctx *gin.Context) {
// 	url := "mongodb://127.0.0.1"
// 	session, err := mgo.Dial(url)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer session.Close()

// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB("db_go").C("user")

// 	// 定义查询结构，根据上方定义的结构查询usr
// 	usr := user{}
// 	// 查找数据
// 	err = c.Find(bson.M{"username": "Johnson"}).One(&usr)
// 	ctx.JSON(http.StatusOK, usr)
// }

// // 查询全部数据
// func findAll(ctx *gin.Context) {
// 	url := "mongodb://127.0.0.1"
// 	session, err := mgo.Dial(url)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer session.Close()

// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB("db_go").C("user")

// 	// 查找全部
// 	usrs := make([]user, 10)
// 	// 查找全部
// 	err = c.Find(nil).All(&usrs)
// 	ctx.JSON(http.StatusOK, usrs)
// }

// // 删除数据
// func delete(ctx *gin.Context) {
// 	url := "mongodb://127.0.0.1"
// 	session, err := mgo.Dial(url)
// 	if err != nil {
// 		panic(err)
// 	} else {
// 		ctx.String(http.StatusOK, "删除成功")
// 	}
// 	defer session.Close()

// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB("db_go").C("user")

// 	err = c.Remove(bson.M{"username": "Johnson"})

// }

// // 修改数据
// func update(ctx *gin.Context) {
// 	url := "mongodb://127.0.0.1"
// 	session, err := mgo.Dial(url)
// 	if err != nil {
// 		panic(err)
// 	} else {
// 		ctx.String(http.StatusOK, "修改成功")
// 	}
// 	defer session.Close()

// 	session.SetMode(mgo.Monotonic, true)
// 	c := session.DB("db_go").C("user")

// 	err = c.Update(bson.M{"username": "Johnson"}, bson.M{"$set": bson.M{"password": "123456"}})

// }

// // 处理跨域请求,将跨域请求函数作为中间件处理
// func Cors() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		method := c.Request.Method

// 		c.Header("Access-Control-Allow-Origin", "*")                                                                                         // 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
// 		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")                              //header的类型
// 		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")                                                          //允许请求方法
// 		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type") //返回数据格式
// 		c.Header("Access-Control-Allow-Credentials", "true")                                                                                 //设置为true，允许ajax异步请求带cookie信息

// 		//放行所有OPTIONS方法
// 		if method == "OPTIONS" {
// 			c.AbortWithStatus(http.StatusNoContent)
// 		}
// 		// 处理请求
// 		c.Next()
// 	}
// }

// // 配置路由
// func Main(e *gin.Engine) {
// 	// 全局使用中间件
// 	e.Use(Cors())
// 	// 定义路由，调用接口函数
// 	e.GET("/mongo/insert", insert)
// 	// 定义json函数接口
// 	e.GET("/mongo/find", find)
// 	// 定义json结构体函数接口
// 	e.POST("/mongo/find_all", findAll)
// 	// 定义json结构体函数接口
// 	e.GET("/mongo/delete", delete)
// 	// 定义json结构体函数接口
// 	e.GET("/mongo/update", update)
// 	// 将路由信息return回调
// 	// return r
// }
