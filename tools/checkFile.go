package tools

import (
	"fmt"
	"time"
	"errors"
	"strings"
	"strconv"
	"os/exec"

	"xsbPro/log"
)

var (
	CurrentLine int64
)

func InitFileChecker(filePath string) {
	var err error
	CurrentLine, err = getFileLines(filePath)
	if err != nil {
		panic(err)
	}
	log.InfoF("当前nginx日志起始行数:%d", CurrentLine)
}

func CheckFileOnce(filePath string) {
	var infos = []*NginxLogInfo{}
	fileLineNow, err := getFileLines(filePath)
	if err != nil {
		log.SysF("fileTimer err: %s", err.Error())
		return
	}

	if fileLineNow < CurrentLine {
		fmt.Printf("Log restart(line %d >> line %d) at %s\n", CurrentLine, fileLineNow, time.Now().Format("2006-01-02 15:04:05"))
		ret, err := getLinesBeforeLine(filePath, fileLineNow)
		if err != nil {
			log.SysF("fileTimer err: %s", err.Error())
			return
		}
		infos = AnalyNginxLines(ret)

	} else if fileLineNow > CurrentLine {
		ret, err := getLinesBehindLine(filePath, CurrentLine + 1)
		if err != nil {
			log.SysF("fileTimer err: %s", err.Error())
			return
		}
		infos = AnalyNginxLines(ret)
	} else {
		return
	}
	SaveResultsToCache(infos)
	CurrentLine = fileLineNow
}

func getFileLines(file string) (int64, error) {
	out, err := exec.Command("/usr/bin/wc", file, "-l").Output()
	if err != nil {
		return 0, err
	}

	splitList := strings.Split(string(out), " ")
	if len(splitList) < 1 {
		fmt.Println("wc返回字符串拆分结果:", splitList)
		return 0, errors.New("wc非正常返回")
	}

	linesCount, err := strconv.ParseInt(splitList[0], 10, 64)
	if err != nil {
		return 0, err
	}
	return linesCount, nil
}

func getLinesBehindLine(file string, line int64) (string, error) {
	out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("sed -n -e '%d,$ p' %s", line, file)).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func getLinesBeforeLine(file string, line int64) (string, error) {
	out, err := exec.Command("/bin/sh", "-c", fmt.Sprintf("sed -n -e '1,%dp' %s", line, file)).Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}