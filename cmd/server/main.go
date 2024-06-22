// Сервер должен реализовывать следующую бизнес-логику:
// - регистрация, аутентификация и авторизация пользователей;
// - хранение приватных данных;
// - синхронизация данных между несколькими авторизованными клиентами одного владельца;
// - передача приватных данных владельцу по запросу.

package main

import (
	"fmt"
	"github.com/awnumar/memguard"
)

//func main() {
//	fmt.Print("Hello World")
//}

//todo хранить логины, пароли, бинарные данные и прочую приватную информацию.

// endpoints
// - login
// - auth
// - save all
// - get all
// keys exchange
func main() {
	// Safely terminate in case of an interrupt signal
	memguard.CatchInterrupt()

	// Purge the session when we return
	defer memguard.Purge()

	// Generate a key sealed inside an encrypted container
	key := memguard.NewEnclave([]byte("sekretkey"))

	b, err := key.Open()
	if err != nil {
		panic(err)
	}
	fmt.Println(b)

	//// Passing the key off to another function
	//key = invert(key)
	//
	//// Decrypt the result returned from invert
	//keyBuf, err := key.Open()
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, err)
	//	return
	//}
	//defer keyBuf.Destroy()
	//
	//// Um output it
	//fmt.Println(keyBuf.Bytes())
}

func invert(key *memguard.Enclave) *memguard.Enclave {
	// Decrypt the key into a local copy
	b, err := key.Open()
	if err != nil {
		memguard.SafePanic(err)
	}
	defer b.Destroy() // Destroy the copy when we return

	// Open returns the data in an immutable buffer, so make it mutable
	b.Melt()

	// Set every element to its complement
	for i := range b.Bytes() {
		b.Bytes()[i] = ^b.Bytes()[i]
	}

	// Return the new data in encrypted form
	return b.Seal() // <- sealing also destroys b
}
