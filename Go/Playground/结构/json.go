package main

import (
	"encoding/json"
	"fmt"
)

type Address struct {
	City string `json:"city"`
}

type User struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email,omitempty"`
	Password string   `json:"-"`
	Tags     []string `json:"tags"`
	Address  Address  `json:"address"`
	active   bool     // 未导出，被忽略
}

func main() {
	u := User{
		ID:       1,
		Name:     "AnAn",
		Password: "secret",
		Tags:     []string{"dev", "go"},
		Address:  Address{City: "ChengDu"},
	}

	b, err := json.MarshalIndent(u, "", "  ")    // func MarshalIndent(v any, prefix, indent string) ([]byte, error)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}