package main

import (
	"fmt"
	"os"
)

func main() {
  args := os.Args[1:]
  
  if len(args) > 0 { //En caso de que haya mas de 1 palabra
    var name = "" //Variable en la que se guardaran los argumentos
    for _, element := range args {//Se recorre cada elemento 
      name += element + " " // Y se agrega a la variable
    }
    fmt.Println("Welcome to the jungle " + name)
  } else { //si no hay elementos en args entonces no se a puesto ningun nombre
    fmt.Println("No name found")
  }
}

//Bibliography 
// https://gobyexample.com/command-line-arguments
// https://www.digitalocean.com/community/tutorials/how-to-use-variadic-functions-in-go
