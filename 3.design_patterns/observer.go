package main

import "fmt"

/*
Observer
Es un patrón de diseño de comportamiento, permite que un conjunto de objetos se suscriban a otro objeto para tener notificaciones acerca de la ocurrencia de un evento.
*/

// Objects can suscribed to an event

type Topic interface {
	register(observer Observer) // Añadir observadores al objeto
	broadcast()                 // Notificar a todos los observadores
}

type Observer interface {
	getId() string      // Get the id of the observer
	updateValue(string) // Update the value of the observer, trigger the event
}

// Item -> No disponible
// Cuando tenga disponiblidad, avise a los observadores
type Item struct {
	observers []Observer // Lista de observadores
	name      string     // Nombre del item
	available bool       // Disponibilidad
}

// NewItem -> Crear un nuevo item
func NewItem(name string) *Item {
	return &Item{
		observers: make([]Observer, 0),
		name:      name,
		available: false,
	}
}

// Mandar el evento
func (i *Item) UpdateAvailable() {
	fmt.Println("Item", i.name, "is available")
	i.available = true
	i.broadcast()
}

// Mandar el evento a todos los observadores
func (i *Item) broadcast() {
	for _, o := range i.observers {
		o.updateValue(i.name)
	}
}

// Registrar un observador
func (i *Item) register(observer Observer) {
	i.observers = append(i.observers, observer)
}

type EmailClient struct {
	id string
}

func (e *EmailClient) getId() string {
	return e.id
}

func (e *EmailClient) updateValue(name string) {
	fmt.Println("Email to", e.id, "with item", name)
}

func main() {
	item := NewItem("TV")
	email1 := &EmailClient{id: "23ab"}
	item.register(email1)

	email2 := &EmailClient{id: "34dc"}
	item.register(email2)

	item.UpdateAvailable()
}
