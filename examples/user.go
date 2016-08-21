package examples

import (
	"time"

	"github.com/wanliu/store"
)

type User struct {
	ID                uint64
	Login             string `index:"unique"`
	Avatar            string `valid:"requri,optional"`
	Realname          string `valid:"-"`
	Email             string `index:"unique" valid:"email,optional"`
	Title             string `index:"index"`
	Mobile            string `index:"unique"`
	Phone             string `index:"unique"`
	RoleName          string
	RoleID            uint64
	OrgID             uint64
	password          string
	HashedPassword    [48]byte      // 加密密码
	EnablePassword    bool          // 是否设置密码开关
	PasswordAlive     time.Duration // 密码有效周期  0 表示永久
	PasswordExpiredAt time.Time     // 密码有效时间，过期会要更改
	store.Timestamp
	// Sex               SexType       // 性别
}
