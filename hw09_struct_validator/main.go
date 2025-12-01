package main

type TestStruct struct {
	i int `validate:"len:32"`
	s string `validate:"len:32|;"`
	sS []string `validate:"len:32"`
}

func main() {
	a := TestStruct{
		i: 1,
		s: "hello",
		sS: []string{},
	}
	err := Validate(a)
	if err != nil {
		println(err.Error())
	}
}
