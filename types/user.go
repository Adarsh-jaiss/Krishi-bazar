package types

import "time"

type CompleteSignupRequest struct {
	User             User   `json:"user"`
	VerificationCode string `json:"verification_code"`
}

type Approve struct {
	UserID     string `json:"user_id" db:"user_id"`
	IsVerified bool   `json:"is_verified" db:"is_verified"`
}

type LoginRequest struct {
	PhoneNumber      string `json:"phone_number"`
	AadharNumber     string `json:"aadhar_number"`
	VerificationCode string `json:"verification_code"`
	IsVerified       string `json:"is_verified"`
}

// TODO:Add User's photo -> DONE
type User struct {
	ID           string    `json:"id" db:"id"`
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	AadharNumber string    `json:"aadhar_number" db:"aadhar_number"`
	Email        string    `json:"email" db:"email"`
	PhoneNumber  string    `json:"phone_number" db:"phone_number"`
	IsFarmer     bool      `json:"is_farmer" db:"is_farmer"`
	Address      string    `json:"address" db:"address"`
	City         string    `json:"city" db:"city"`
	State        string    `json:"state" db:"state"`
	PinCode      string    `json:"pin_code" db:"pin_code"`
	FarmSize     string    `json:"farm_size" db:"farm_size"`
	Password     string    `json:"password" db:"password"`
	IsVerified   bool      `json:"is_verified_by_admin" db:"is_verified_by_admin"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
	LastLoginAt  time.Time `json:"last_login_at" db:"last_login_at"`
	Image        string    `json:"img" db:"img"`
	UserType     string    `json:"user_type" db:"user_type"`
}

type Farmer struct {
	ID        string `json:"id" db:"id"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Image     string `json:"img" db:"img"`
}
