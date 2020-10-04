package engin

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
)

var (
	blockCacheFile string
	blockAddressList map[string]bool
	blockLock sync.RWMutex
)

func init()  {
	blockCacheFile = "cache.db"
	blockAddressList = make(map[string]bool,1024)
	loadBlackCacheFile()
}

func parseBlack(line string) (string,bool) {
	hosts := strings.Split(line,"\t")
	if len(hosts) != 2 {
		return "", false
	}
	if 0 == strings.Compare(hosts[1],"true") {
		return hosts[0], true
	}
	return hosts[0], false
}

func loadBlackCacheFile()  {
	body, err := ioutil.ReadFile(blockCacheFile)
	if err != nil {
		return
	}
	lines := strings.Split(string(body),"\n")
	for _,v := range lines {
		address,blackFlag := parseBlack(v)
		if address == "" {
			continue
		}
		blockAddressList[address] = blackFlag
	}
}

func saveBlackCacheFile(address string, black bool)  {
	blockLock.Lock()
	defer blockLock.Unlock()

	file,err := os.OpenFile(blockCacheFile,os.O_APPEND|os.O_WRONLY,0644)
	if err != nil {
		file, err = os.Create(blockCacheFile)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
	defer file.Close()
	if black {
		fmt.Fprintf(file,"%s\t%s\n",address,"true")
	}else {
		fmt.Fprintf(file,"%s\t%s\n",address,"false")
	}
}

func IsSecondProxy(address string) bool {
	blockLock.RLock()
	black,flag := blockAddressList[address]
	blockLock.RUnlock()
	if flag == true {
		return black
	}

	black = !IsConnect(address, 5)
	blockLock.Lock()
	blockAddressList[address] = black
	blockLock.Unlock()

	saveBlackCacheFile(address,black)
	return black
}


