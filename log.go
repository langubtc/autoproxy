package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

type LogFile struct {
	prefix  string /* 日志前缀信息 */
	dir     string /* 日志所在的目录 */
	name    string
	maxsize int64  /* 文件上限 */
	maxnum  int    /* 文件数量上限 */
	file    *os.File /* 当前正在写入的文件句柄 */
	cursize int64    /* 当前文件大小 */
	
	cache chan []byte
}

func NewLogFile(dir string, size int, num int) (*LogFile, error) {
	logfile := &LogFile{dir: dir, maxsize: int64(size), maxnum: num}
	fileinfo, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !fileinfo.IsDir() {
		return nil, errors.New(dir + "is not dir.")
	}
	logfile.cache = make(chan []byte, 1024)
	go logfile.logDumpTask()
	return logfile, nil
}

func timeStampGet() string {
	tm := time.Now()
	return fmt.Sprintf("%4d%02d%02d%02d%02d%02d",
		tm.Year(), tm.Month(), tm.Day(),
		tm.Hour(), tm.Minute(), tm.Second())
}

func getfiletimestamp(files []string) []time.Time {
	timestamp := make([]time.Time,0)	
	for _,v := range files {
		fileinfo, err := os.Stat(v) 
		if err != nil {
			timestamp = append(timestamp,time.Now())
			continue
		}
		timestamp = append(timestamp,fileinfo.ModTime())
	}
	return timestamp
}

func getoldtimestamp(timestamp []time.Time) int {
	oldIdx  := 0
	oldTime := timestamp[0]
	for i:=1; i<len(timestamp); i++ {
		if oldTime.After(timestamp[i]) {
			oldTime = timestamp[i]
			oldIdx = i
		}
	}
	return oldIdx
}

func cleanfile(dir string,num int) {
	filelist := make([]string,0)
	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasSuffix(file.Name(),".zip") {
			filename := fmt.Sprintf("%s%c%s",dir,os.PathSeparator,file.Name())
			filelist = append(filelist,filename)
		}
	}
	for {
		if len(filelist) <= num {
			break
		}
		timestamps := getfiletimestamp(filelist)
		idx := getoldtimestamp(timestamps)
		os.Remove(filelist[idx])
		filelist = append(filelist[0:idx],filelist[idx+1:]...)
	}
}

func (lf *LogFile) setPerfix(prefix string)  {
	lf.prefix = prefix
}

func (lf *LogFile) logDumpTask()  {
	for {
		body, b := <-lf.cache
		if b == false {
			return
		}
		lf.WriteToFile(body)
	}
}

func (lf *LogFile) packfile() error {
	filesrc := fmt.Sprintf("%s/%s.log", lf.dir, lf.name)
	filedest := fmt.Sprintf("%s/%s.zip", lf.dir, lf.name)
	err := Zip(filesrc, filedest)
	if err == nil {
		os.Remove(filesrc)
	}
	return err
}

func (lf *LogFile) openfile() error {
	lf.name = timeStampGet()
	path := fmt.Sprintf("%s/%s.log", lf.dir, lf.name)

	file, err := os.OpenFile(path, os.O_WRONLY, 0600)
	if err != nil {
		file, err = os.Create(path)
		if err != nil {
			return err
		}
		lf.cursize = 0
	} else {
		fileinfo, err := file.Stat()
		if err == nil {
			file.Seek(fileinfo.Size(), 0)
			lf.cursize = fileinfo.Size()
		} else {
			lf.cursize = 0
		}
	}
	lf.file = file
	return nil
}

func (lf *LogFile) Write(p []byte) (n int, err error) {
	buffer := make([]byte,len(p))
	copy(buffer,p)
	lf.cache <- buffer
	return len(p),nil
}

func (lf *LogFile) WriteToFile(p []byte) {
	for {
		if lf.file == nil {
			err := lf.openfile()
			if err != nil {
				os.Stderr.Write([]byte(err.Error()))
				return
			}
		} else {
			if lf.cursize > lf.maxsize {
				lf.file.Close()
				lf.file = nil
				err := lf.packfile()
				if err != nil {
					os.Stderr.Write([]byte(err.Error()))
					return
				}
				cleanfile(lf.dir, lf.maxnum)
			} else {
				os.Stderr.Write(p)
				cnt, _ := lf.file.Write(p)
				if cnt > 0 {
					lf.cursize += int64(cnt)
				}
				return
			}
		}
	}
}


type LOG_LEVEL int

const (
	INFO LOG_LEVEL = iota
	WARNING
	ERROR
	EXCEPT
)

func loglevel(level LOG_LEVEL) string {
	if level == INFO {
		return "INFO"
	} else if level == WARNING {
		return "WARNING"
	} else if level == ERROR {
		return "ERROR"
	} else {
		return "EXCEPT"
	}
}

var gFileLog *LogFile

func LogInit(cfg LogConfig) error {
	file, err := NewLogFile(cfg.Path, cfg.FileSize, cfg.FileNum)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	log.SetOutput(file)
	log.SetFlags(log.Lmicroseconds | log.LstdFlags)
	gFileLog = file
	return nil
}

func LogPrefix(s string)  {
	gFileLog.setPerfix(s)
}

func Infof(format string, v ...interface{})  {
	printf(INFO,format,v...)
}

func Info(v ...interface{})  {
	print(INFO,v...)
}

func Warnf(format string, v ...interface{})  {
	printf(WARNING,format,v...)
}

func Warn(v ...interface{})  {
	print(WARNING,v...)
}

func Errorf(format string, v ...interface{})  {
	printf(ERROR,format,v...)
}

func Error(v ...interface{})  {
	print(ERROR,v...)
}

func Fatalf(format string, v ...interface{})  {
	printf(EXCEPT,format,v...)
	os.Exit(-1)
}

func Fatal(v ...interface{})  {
	print(EXCEPT,v...)
	os.Exit(-1)
}

func perfix(level LOG_LEVEL) string {
	var prefix string
	if gFileLog.prefix != "" {
		prefix = fmt.Sprintf("[%s]",gFileLog.prefix)
	}
	return fmt.Sprintf("%s[%s]", prefix, loglevel(level))
}

func print(level LOG_LEVEL, v ...interface{}) {
	output := perfix(level)
	output += fmt.Sprint(v...)
	if level >= ERROR {
		output += fmt.Sprint("\n"+string(debug.Stack()))
	}
	log.Println(output)
}

func printf(level LOG_LEVEL, format string, v ...interface{}) {
	output := perfix(level)
	output += fmt.Sprintf(format, v...)
	if level >= ERROR {
		output += fmt.Sprint("\n"+string(debug.Stack()))
	}
	log.Printf(output)
}
