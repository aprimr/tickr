package hash

import "golang.org/x/crypto/bcrypt"

// HashPassword converts a plain text password to hashed string
func HashPasword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err
}

// ComparePasswordHash compares the plain password and hashed password
func ComparePasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}
