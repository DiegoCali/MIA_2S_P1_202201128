package utils

var (
	GlobalMounts map[string]string = make(map[string]string)
	ActualUser   User              = User{
		Name: "",
		Id:   "",
	}
)
