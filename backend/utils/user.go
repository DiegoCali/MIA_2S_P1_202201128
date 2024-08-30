package utils

type User struct {
	Name string
	Id   string
}

func (user *User) Set(name string, id string) {
	user.Name = name
	user.Id = id
}

func (user *User) GetId() string {
	return user.Id
}

func (user *User) GetName() string {
	return user.Name
}

func (user *User) Print() {
	println("Name: ", user.Name)
	println("Id: ", user.Id)
}
