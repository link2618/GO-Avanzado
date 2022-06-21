package main

// Definir una familia de algoritmos en structs separados
// Clase base que las usa y pueda ir intercambiando

/*
Strategy
Es un patrón de diseño de comportamiento, consiste en definir una familia de funciones similares en una clase base, en el caso de Go seria un struct. Parte desde los principios SOLID.
*/

type PasswordProtector struct {
	user          string
	passwordName  string
	hashAlgorithm HashAlgorithm
}

type HashAlgorithm interface {
	Hash(p *PasswordProtector)
}

// Constructor
func NewPasswordProtector(user, passwordName string, hash HashAlgorithm) *PasswordProtector {
	return &PasswordProtector{
		user:          user,
		passwordName:  passwordName,
		hashAlgorithm: hash,
	}
}

// Poder intercalar el algoritmo de hash
func (p *PasswordProtector) SetHash(hash HashAlgorithm) {
	p.hashAlgorithm = hash
}

// Ejecutar el algoritmo de hash
func (p *PasswordProtector) Hash() {
	p.hashAlgorithm.Hash(p)
}

type SHA struct{}

func (SHA) Hash(p *PasswordProtector) {
	println("Using SHA hash for", p.user, p.passwordName)
}

type MD5 struct{}

func (MD5) Hash(p *PasswordProtector) {
	println("Using MD5 hash for", p.user, p.passwordName)
}

func main() {
	protector := NewPasswordProtector("link", "password", SHA{})
	protector.Hash()

	protector.SetHash(MD5{})
	protector.Hash()
}
