package main

import "log"

func main()  {
	err := FileInit()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = BoxInit()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = AuthInit()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = IconInit()
	if err != nil {
		log.Fatal(err.Error())
	}
	mainWindows()
}


