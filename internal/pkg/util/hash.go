package util

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// ใช้ DefaultCost (10) ซึ่งเหมาะสมกับ Server ส่วนใหญ่ในปัจจุบัน
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	// ถ้าคืนค่า nil แสดงว่ารหัสถูกต้อง
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
