package biz

// 用户状态, 二进制开关
// 8 7 6 5 4 3 2  1
//             ^  ^
//             封 存
//             禁 在

const (
	indexExist = uint8(0)
	indexBan   = uint8(1)
)

type UserStatus uint8

func NewUserStatus(i uint8) *UserStatus {
	return (*UserStatus)(&i)
}

// IsExist 用户是否存在
// true: 存在, false: 不存在
func (s *UserStatus) IsExist() bool {
	return s.is(indexExist)
}

// SetExist 设置用户是否存在
func (s *UserStatus) SetExist(is bool) {
	s.set(indexExist, is)
}

// IsBan 用户是否封禁
// true: 封禁, false: 未封禁
func (s *UserStatus) IsBan() bool {
	return s.is(indexBan)
}

// SetBan 设置用户是否封禁
func (s *UserStatus) SetBan(is bool) {
	s.set(indexBan, is)
}

func (s *UserStatus) is(index uint8) bool {
	if index >= 8 {
		return false
	}

	return ((*s)>>index)&1 == 1
}
func (s *UserStatus) set(index uint8, is bool) {
	*s = (*s) | (1 << index)
	if !is {
		*s = (*s) ^ (1 << index)
	}
}
