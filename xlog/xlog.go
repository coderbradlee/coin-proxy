package xlog

import (
	"fmt"
	"github.com/robfig/cron"
	"log"
	"os"
	"path/filepath"
	"time"
)

/*
	path/bin/xxx
	path/logs/xxx_20060102.log
	path/logs/xxx_sign_20060102.log
or
 	path/xxx
	path/logs/xxx_20060102.log
	path/logs/xxx_sign_20060102.log

time-hour-log
	path/logs/xxx_200601021504sign.log
*/
// logFilePrefix abspath  /.../path/logs/xxx

// 日志文件定时器
var logCron *cron.Cron

//
var loggers map[string]*LoggerMan

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func Create(_dir string) {

	exist, err := PathExists(_dir)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
		return
	}

	if !exist {
		fmt.Printf("no dir![%v]\n", _dir)
		// 创建文件夹
		err := os.Mkdir(_dir, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		} else {
			fmt.Printf("mkdir success!\n")
		}
	}
}
func XX() {

}
func init() {
	Create("./logs")
	// logFile, err := os.OpenFile(getLogFilePrefix()+"_"+time.Now().Format("20060102")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, os.ModePerm)
	logFile, err := os.OpenFile(getLogFilePrefix()+"_"+time.Now().Format("20060102")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println("[xlog] init OpenFile err", err)
		return
	}
	log.SetOutput(logFile)

	loggers = make(map[string]*LoggerMan)

	//启动定时任务
	logCron = cron.New()
	logCron.Start()

	//
	AddSign("", "0 0 0 * * *", "20060102")
	// AddSign("share", "0 0 0 * * *", "20060102")
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
}

func AddSign(sign string, cronStr string, format string) {
	loggerMan := &LoggerMan{Sign: sign, CronStr: cronStr, FileFormat: format}
	loggerMan.Run()
	loggers[sign] = loggerMan
	err := logCron.AddJob(cronStr, loggerMan)
	if err != nil {
		log.Println("[xlog] AddSign AddJob err", err)
	}
}

func Println(sign string, args ...interface{}) {
	if _, ok := loggers[sign]; !ok {
		if loggers[""].Logger == nil {
			loggers[""].Run()
			return
		}

		loggers[""].Logger.Printf("["+sign+"] %v", args)
		return
	}

	if loggers[sign].Logger == nil {
		loggers[sign].Run()
		return
	}
	loggers[sign].Logger.Println(args)
}

func Printf(sign, format string, args ...interface{}) {
	// if loggers[""].Logger == nil {
	// 	loggers[""].Run()
	// 	log.Printf(format, args...)
	// 	return
	// }
	// loggers[""].Logger.Printf(format, args...)
	// if loggers[sign].Logger == nil {
	// 	loggers[sign].Run()
	// 	return
	// }
	// loggers[sign].Logger.Println(args)
	if _, ok := loggers[sign]; !ok {
		if loggers[""].Logger == nil {
			loggers[""].Run()
			return
		}
		loggers[""].Logger.Printf(format, args...)
		return
	}

	if loggers[sign].Logger == nil {
		loggers[sign].Run()
		return
	}
	loggers[sign].Logger.Printf(format, args...)
}

func Fatal(args ...interface{}) {
	if loggers[""].Logger == nil {
		loggers[""].Run()
		log.Fatal(args)
		return
	}
	loggers[""].Logger.Fatal(args...)
}

//
type LoggerMan struct {
	Sign       string      //日志标识
	CronStr    string      //定时
	FileFormat string      //日志文件名 时间 格式化
	Logger     *log.Logger //
	LogFile    *os.File    //日志文件句柄
}

func (p *LoggerMan) Run() {
	defer func() {
		if rev := recover(); rev != nil {
			log.Println("[xlog] Run recover", rev)
		}
	}()

	var err error

	//make new log file
	now := time.Now()
	var fileName string
	fileName = getLogFilePrefix() + "_" + now.Format(p.FileFormat) + p.Sign + ".log"
	newLogFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0644)
	if err != nil {
		log.Println("[xlog] Run os.OpenFile err", err)
		return
	}

	//
	oldLogFile := p.LogFile

	//
	p.LogFile = newLogFile
	if p.Logger != nil {
		p.Logger.SetOutput(p.LogFile)
	} else {
		p.Logger = log.New(p.LogFile, "", log.Ldate|log.Ltime|log.Llongfile)
	}

	//close old logfile
	if oldLogFile != nil {
		err = oldLogFile.Close()
		if err != nil {
			log.Println("[xlog] Run oldLogFile close", err)
		}
	}

	//默认log 使用每日log的file
	if p.Sign == "" {
		log.SetOutput(p.LogFile)
	}
}

// getLogFilePrefix
// path/bin/xxx  path/logs/xxx
func getLogFilePrefix() string {
	instanceAbsPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Println("[xlog] getLogFilePrefix err", err)
		return ""
	}
	instanceName := filepath.Base(instanceAbsPath)
	// topDir := instanceName
	topDir := filepath.Dir(instanceAbsPath)
	logsDir := filepath.Join(instanceAbsPath, "logs")

	if _, errStat := os.Stat(logsDir); os.IsNotExist(errStat) {
		err := os.Mkdir(logsDir, os.ModePerm)
		if err != nil {
			log.Println("[xlog] getLogFilePrefix Mkdir logs", err)
			return ""
		}
	}
	return filepath.Join(topDir, "logs", instanceName)
}

// getLogFilePrefix
// path/xxx path/logs/xxx
func getLogFilePrefix2() string {
	instanceAbsPath, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Println("[xlog] getLogFilePrefix2 err", err)
		return ""
	}
	instanceName := filepath.Base(instanceAbsPath)
	topDir := filepath.Dir(filepath.Dir(instanceAbsPath))

	return filepath.Join(topDir, "logs", instanceName)
}
