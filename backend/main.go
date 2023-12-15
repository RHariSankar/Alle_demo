package main

import (
	"alle/client/azure"
	"fmt"
)

func main() {

	fmt.Print("Hello world")

	azureCLU := azure.AzureCLU{}
	x, y, z := azureCLU.GetIntentAndEntity("Get me pictures of dog from my saved data")
	fmt.Printf("%t %s %s", x, y, z)

}
