package main

import "fmt"

type Human struct {
	Name    string
	Surname string
	Age     int
}

type Action struct {
	Human
	ActionType string
	Status     bool
}

func (h *Human) Info() {
	fmt.Printf("Имя: %s, Фамилия: %s, Возраст: %d лет\n", h.Name, h.Surname, h.Age)
}

func (a *Action) DoStuff() {
	fmt.Printf("%s занимется следующим: %s, статус - %t", a.Name, a.ActionType, a.Status)
}

func main() {
	person := Human{
		Name:    "Макар",
		Surname: "Соловьев",
		Age:     26,
	}

	person.Info()

	job := Action{
		Human:      person,
		ActionType: "Программирование",
		Status:     true,
	}

	job.DoStuff()
}
