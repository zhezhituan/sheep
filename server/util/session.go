package util

import (
	"github.com/astaxie/beego/session"
)

var GlobalSessions *session.Manager

func init() {
	sessionConfig := &session.ManagerConfig{
		CookieName:      "gosessionid",
		EnableSetCookie: true,
		Gclifetime:      3600,
		Maxlifetime:     3600,
		Secure:          false,
		CookieLifeTime:  0,
		ProviderConfig:  "./tmp",
	}
	GlobalSessions, _ = session.NewManager("memory", sessionConfig)
	go GlobalSessions.GC()
}

// globalSessions 有多个函数如下所示：

// SessionStart 根据当前请求返回 session 对象
// SessionDestroy 销毁当前 session 对象
// SessionRegenerateId 重新生成 sessionID
// GetActiveSession 获取当前活跃的 session 用户
// SetHashFunc 设置 sessionID 生成的函数
// SetSecure 设置是否开启 cookie 的 Secure 设置
// 返回的 session 对象是一个 Interface，包含下面的方法

// Set(key, value interface{}) error
// Get(key interface{}) interface{}
// Delete(key interface{}) error
// SessionID() string
// SessionRelease()
// Flush() error
