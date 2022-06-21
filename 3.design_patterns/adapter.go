package main

import "fmt"

/*
Adapter
Es un patrón de diseño estructural, es un intermedio que permite adaptar el comportamiento de un struct a una interfaz, cuando existe una incompatibilidad sin la necesidad de tener que refactorizar código.
*/

type Payment interface {
	Pay()
}

type CashPayment struct{}

func (CashPayment) Pay() {
	fmt.Println("Payment using Cash")
}

func ProcessPayment(p Payment) {
	p.Pay()
}

type BankPayment struct{}

func (BankPayment) Pay(bankAccount int) {
	fmt.Printf("Paying using Bankaccount %d\n", bankAccount)
}

type BankPaymentAdapter struct {
	BankPayment *BankPayment
	bankAccount int
}

func (bpa *BankPaymentAdapter) Pay() {
	bpa.BankPayment.Pay(bpa.bankAccount)
}

func main() {
	cash := &CashPayment{}
	ProcessPayment(cash)
	// bank := &BankPayment{}
	// ProcessPayment(bank)
	bpa := &BankPaymentAdapter{
		bankAccount: 5,
		BankPayment: &BankPayment{},
	}
	ProcessPayment(bpa)
}
