/*
 * @Author: chenYY
 * @Date: 2023-05-06 09:33:53
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2023-05-09 15:01:37
 * @FilePath: /go/gormWithGin/main.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Class struct {
	gorm.Model
	Name     string
	Students []Student
}
type Student struct {
	gorm.Model
	Name     string
	IDcard   IDcard
	ClassID  uint
	Teachers []*Teacher `gorm:"many2many:students_teachers;"`
}
type Teacher struct {
	gorm.Model
	Name     string
	Students []*Student `gorm:"many2many:students_teachers;"`
}
type IDcard struct {
	gorm.Model
	StudentID uint
	CardNum   int
}
type MyClaims struct {
	Username string `json:username`
	jwt.StandardClaims
}

func main() {
	fmt.Println("my way to go")
	//  DB.AutoMigrate(&Teacher{}, &Class{}, &Student{}, &IDcard{})
	DB = DB.Session(&gorm.Session{
		Logger: mysqlLogger,
	})

	//  i := IDcard{
	//  	CardNum: 12343242,
	//  }
	//  s := Student{
	//  	Name:   "华晨宇",
	//  	IDcard: i,
	//  }
	//  t := Teacher{
	//  	Name: "王子豪",
	//  }
	//  c := Class{
	//  	Name:     "二年三班",
	//  	Students: []Student{s},
	//  }
	//  DB.Create(&c)
	//  DB.Create(&t)
	// defer DB.Close()

	//添加关系
	//  t := Teacher{
	//  	ID: 1,
	//  }
	// var tercher Teacher
	// var student Student
	// student.ID = 1
	// tercher.ID = 1
	// DB.Preload("Teachers").Find(&student)
	// DB.Model(&student).Association("Teachers").Append(&tercher)
	// fmt.Println(student)

	//查询指定记录,打印指针
	// var student Student
	// DB.Preload("Teachers").First(&student, 1)
	// fmt.Println(student.Teachers)
	// for _, teacher := range student.Teachers {
	// 	fmt.Println(*teacher)
	// }
	//初始加密密钥
	mySigningKey := []byte("getpostdelect")
	//使用jwt校验用户信息
	c := MyClaims{
		Username: "chen",
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 60, //当前一分钟生效
			ExpiresAt: time.Now().Unix() + 5,  //有效时间2小时
			Issuer:    "chenGo",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {
		panic(err)
	}
	fmt.Println(ss)
	//解析token
	// time.Sleep(6 * time.Second)
	t, err := jwt.ParseWithClaims(ss, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Println(t.Claims.(*MyClaims).Username)
	return

	//使用gin搭配gorm食用
	router := gin.Default()
	router.POST("/creatStudent", func(c *gin.Context) {
		var student Student
		if err := c.ShouldBindJSON(&student); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
				"msg":   "错误请求",
			})
			return
		} else {
			DB.Create(&student)
			c.JSON(http.StatusOK, gin.H{
				"msg": "creat success",
			})
		}

	})
	router.GET("/student", func(c *gin.Context) {
		id := c.Query("ID")
		// lastname := c.Query("lastname") // 是 c.Request.URL.Query().Get("lastname") 的简写
		var student Student
		DB.Preload("Teachers").First(&student, "id = ?", id)
		var teachers []string
		for _, teacher := range student.Teachers {
			teachers = append(teachers, teacher.Name)
		}
		c.JSON(http.StatusOK, gin.H{
			"msg":         "success",
			"studentName": student.Name,
			"studentID":   student.ID,
			"creatTime":   student.CreatedAt,
			"teachers":    teachers,
		})
	})
	router.GET("/class", func(c *gin.Context) {
		id := c.Query("ID")
		// lastname := c.Query("lastname") // 是 c.Request.URL.Query().Get("lastname") 的简写
		var class Class
		DB.Preload("Students").Preload("Students.IDcard").Preload("Students.Teachers").First(&class, "id = ?", id)
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
			"c":   class,
		})
	})
	router.Run(":80")
}
