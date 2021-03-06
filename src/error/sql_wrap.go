package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

// 定义一个全局对象db
var db *sql.DB

// 定义一个初始化数据库的函数
func initDB() error {
	// DSN:Data Source Name
	dsn := "user:password@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	return db.Ping()
}

// 数据结构
type User struct {
	Id   int
	Age  int
	Name string
}

// ErrUserNotFound 未找到用户的错误
var ErrUserNotFound = errors.New("未找到用户")

// getById 根据主键查询单个用户信息
func getById(id int) (*User, error) {
	sqlStr := "select id, name, age from user where id=?"
	var user User
	// 包装错误信息
	if err := db.QueryRow(sqlStr, id).Scan(&user.Id, &user.Name, &user.Age); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrapf(ErrUserNotFound, "未找到用户,id为 %d\n", id)
		}
		return nil, errors.Wrap(err, fmt.Sprintf("查询单个用户失败, id为 %d\n", id))
	}
	return &user, nil
}

// getPage 根据主键获取后续用户
func getPage(id, limit int) (*[]User, error) {
	sqlStr := "select id, name, age from user where id > ? limit ?"
	rows, err := db.Query(sqlStr, id, limit)
	users := []User{}
	// 包装错误信息
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &users, errors.Wrap(ErrUserNotFound, fmt.Sprintf("查询全部个用户失败, id为 %d, 参数为 %d", id, limit))
		}
		return nil, errors.Wrap(err, fmt.Sprintf("查询全部个用户失败, id为 %d, 参数为 %d", id, limit))
	}

	// 关闭rows释放持有的数据库链接
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return &users, nil
}

func main() {
	defer func() {
		db.Close()
	}()
	if err := initDB(); err != nil {
		log.Fatalf("连接数据库失败,错误信息为:%+v\n", err)
	}

	id := 1
	user, err := getById(id)
	if err == nil {
		log.Printf("id为%d的用户为 %v\n", id, user)
	} else if errors.Is(err, ErrUserNotFound) {
		log.Printf("id为%d的用户不存在\n", id)
	} else {
		log.Printf("获取用户信息失败,错误信息为:%+v\n", err)
	}

	size := 10
	users, err := getPage(id, size)
	if err != nil {
		log.Printf("获取用户信息失败,错误信息为:%+v\n", err)
	} else {
		log.Printf("id为%d后%d个的用户为 %v\n", id, size, users)
	}
}
