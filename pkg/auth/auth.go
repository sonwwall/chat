package auth

import "golang.org/x/crypto/bcrypt"

//该文件对密码进行哈希加密处理

// HashPassword 加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) //14是盐值
	//14 代表的是 bcrypt 哈希算法的代价因子（cost factor）。
	//代价因子决定了哈希计算的复杂度和所需的时间，
	//从而影响密码哈希的安全性和性能。bcrypt 使用这个因子来增加计算哈希的时间，使得暴力破解更加困难
	return string(bytes), err
}

// CheckPasswordHash 检查密码
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
