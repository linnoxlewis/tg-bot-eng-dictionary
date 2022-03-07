package main

import (
	"fmt"
	"linnoxlewis/tg-bot-eng-dictionary/internal/api"
)

func main() {
	a := api.NewSkyEngApi()
	res, err := a.GetMeaning("girl")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)

}