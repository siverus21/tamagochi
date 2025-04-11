package main

import "fmt"

type Age int

func (a Age) isAdult() bool{
	if a >= 18{
		return true
	} else {
		return false
	}
}

type User struct {
	name   string
	age    Age
	sex    string
	weight int
	heigth int
}

func (u *User) printUserInfo(){
	fmt.Println(u.name, u.age, u.sex, u.weight, u.heigth)
}

func NewUser(name string, age int, sex string, weight int, heigth int) User{
	return User{
		name: name,
		age: Age(age),
		sex: sex,
		weight: weight,
		heigth: heigth,
	}
}

func main() {

	user := NewUser("Vasya", 23, "Male", 75, 185)
	user1 := NewUser("Max", 21, "Male", 23, 150)
	
	fmt.Println(user.age.isAdult())

	user.printUserInfo()
	user1.printUserInfo()
}