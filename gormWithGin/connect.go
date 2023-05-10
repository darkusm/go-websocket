/*
 * @Author: chenYY
 * @Date: 2023-04-14 15:36:09
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2023-05-09 09:34:20
 * @FilePath: /go/gormWithGin/connect.go
 * @Description:
 *
 * Copyright (c) 2023 by ${git_name_email}, All Rights Reserved.
 */
 package main

 import (
	 "fmt"
 
	 "gorm.io/driver/mysql"
	 "gorm.io/gorm"
	 "gorm.io/gorm/logger"
	 "gorm.io/gorm/schema"
 )
 
 var DB *gorm.DB
 var mysqlLogger logger.Interface
 
 func init() {
	 username := "root"   //账号
	 password := "739090" //密码
	 host := "127.0.0.1"  //数据库地址，可以是Ip或者域名
	 port := 3306         //数据库端口
	 Dbname := "gorm"     //数据库名
	 timeout := "10s"     //连接超时，10秒
 
	 // 要显示的日志等级
	 mysqlLogger = logger.Default.LogMode(logger.Info)
	 // root:root@tcp(127.0.0.1:3306)/gorm?
	 dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	 //连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
 
	 db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		 // SkipDefaultTransaction: true,
		 NamingStrategy: schema.NamingStrategy{
			 // TablePrefix:   "f_", //表名前缀	自动添加前缀
			 SingularTable: true, //单数表名	单数表名
			 NoLowerCase:   true, //大小写转换 不允许转换
		 },
		 // Logger: mysqlLogger,
	 })
	 if err != nil {
		 panic("连接数据库失败, error=" + err.Error())
	 }
	 // 连接成功
	 fmt.Println(db)
	 DB = db
 }
 