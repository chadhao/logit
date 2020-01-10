package utils

// Roles 用户角色组
type Roles []int

// RolesAssert 类型转换为Roles
func RolesAssert(in interface{}) Roles {
	r, ok := in.([]int)
	if !ok {
		r = []int{}
	}
	return r
}

// IsGuest 用户是否访客
func (r Roles) IsGuest() bool {
	return len(r) == 0
}

// Is 用户是否是指定角色
func (r Roles) Is(role int) bool {
	for _, v := range r {
		if v == role {
			return true
		}
	}
	return false
}

// Are 用户是否是指定角色的其中一个
func (r Roles) Are(rs []int) bool {
	for _, v := range r {
		for _, s := range rs {
			if v == s {
				return true
			}
		}
	}
	return false
}
