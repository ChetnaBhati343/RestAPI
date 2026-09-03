package main

import "fmt"

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

func main() {
	students := []Student{
		{ID: 1, Name: "Chetna Bhati", Age: 23, Email: "chetnabhati343@gmail.com"},
		{ID: 2, Name: "test01", Age: 82, Email: "test01@gmail.com"},
		{ID: 3, Name: "test02", Age: 41, Email: "test00@gmail.com"},
	}
	fmt.Println(students)
}
