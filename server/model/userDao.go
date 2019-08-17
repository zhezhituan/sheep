package model

import (
	"encoding/json"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

// UserDao 实例，全局唯一
var CurrentUserDao *UserDao

//UserDao 放redisPool
type UserDao struct {
	pool *redis.Pool
}

//InitUserDao 初始化一个 UserDao 结构体示例，
func InitUserDao(pool *redis.Pool) (currentUserDao *UserDao) {
	currentUserDao = &UserDao{pool: pool}
	return
}

func idIncr(conn redis.Conn) (id int, err error) {
	res, err := conn.Do("incr", "Uid")
	fmt.Println(res)
	id = int(res.(int64))
	if err != nil {
		fmt.Printf("uid incr error: %v\n", err)
		return
	}
	return
}

//GetUsrById 根据用户 id 获取用户信息
// 获取成功返回 user 信息，err nil
// 获取失败返回 err，user 为 nil
func (ud *UserDao) GetUsrByUid(Uid int) (user User, err error) {
	conn := ud.pool.Get()
	defer conn.Close()

	res, err := redis.String(conn.Do("hget", "users", Uid))
	if err != nil {
		err = ERROR_USER_DOES_NOT_EXIST
		return
	}
	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		fmt.Printf("Unmarshal user info error: %v\n", err)
		return
	}
	return
}

// GetUsrByUserName 根据用户 username 获取用户信息
// 获取成功返回 user 信息，err nil
// 获取失败返回 err，user 为 nil
func (ud *UserDao) GetUsrByName(Name string) (user User, err error) {
	conn := ud.pool.Get()
	defer conn.Close()

	res, err := redis.String(conn.Do("hget", "users", Name))
	if err != nil {
		err = ERROR_USER_DOES_NOT_EXIST
		return
	}
	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		fmt.Printf("Unmarshal user info error: %v\n", err)
		return
	}
	return
}

// Register 注册用户
// 用户名不能重复
func (ud *UserDao) Register(Name, Pw, Pw_sure string) (user User, err error) {
	// 判断密码是否正确
	if Pw != Pw_sure {
		err = ERROR_PASSWORD_DOES_NOT_MATCH
		return
	}

	// 保证用户名不重复
	user, err = ud.GetUsrByName(Name)
	if err == nil {
		fmt.Printf("User already exists!\n")
		err = ERROR_USER_ALREADY_EXISTS
		return
	}

	conn := ud.pool.Get()
	defer conn.Close()
	// id 自增 1，作为下个用户 id
	id, err := idIncr(conn)
	if err != nil {
		return
	}
	user = User{Uid: id, Name: Name, Pw: Pw}
	info, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("json marshal error: %v", err)
		return
	}
	_, err = conn.Do("hset", "users", Name, info)
	if err != nil {
		fmt.Printf("set user to reids error: %v", err)
		return
	}
	return
}

//Login 登陆
func (ud *UserDao) Login(userName, password string) (user User, err error) {
	user, err = ud.GetUsrByName(userName)
	if err != nil {
		fmt.Printf("get user by id error: %v\n", err)
		return
	}

	if user.Pw != password {
		err = ERROR_USER_PWD
		return
	}

	return
}
