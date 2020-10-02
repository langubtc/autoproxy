package main

import (
	"github.com/lxn/walk"
	"os"
	"log"
)

func IconLoadFromBox(filename string) *walk.Icon {
	body, err := BoxFile().Bytes(filename)
	if err != nil {
		log.Fatalln(err.Error())
		return nil
	}
	dir := DEFAULT_HOME + "\\icon\\"
	_, err = os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 644)
		if err != nil {
			log.Fatalln(err.Error())
			return nil
		}
	}
	filepath := dir + filename
	err = SaveToFile(filepath, body)
	if err != nil {
		log.Fatalln(err.Error())
		return nil
	}
	icon, err := walk.NewIconFromFileWithSize(filepath, walk.Size{
		Width: 128, Height: 128,
	})
	if err != nil {
		log.Fatalln(err.Error())
		return nil
	}
	return icon
}

var ICON_Main            *walk.Icon
var ICON_Network_Disable *walk.Icon
var ICON_Network_Flow *walk.Icon
var ICON_Network_Full *walk.Icon
var ICON_Network_High *walk.Icon
var ICON_Network_MID  *walk.Icon
var ICON_Network_LOW  *walk.Icon

func IconInit() error {
	ICON_Main = IconLoadFromBox("main.ico")
	ICON_Network_Disable = IconLoadFromBox("network_disable.ico")
	ICON_Network_Flow = IconLoadFromBox("network_flow.ico")
	ICON_Network_Full = IconLoadFromBox("network_full.ico")
	ICON_Network_High = IconLoadFromBox("network_high.ico")
	ICON_Network_MID  = IconLoadFromBox("network_medium.ico")
	ICON_Network_LOW  = IconLoadFromBox("network_low.ico")
	return nil
}
