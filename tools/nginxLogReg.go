package tools

import (
	"regexp"
	"fmt"
	"errors"
	"time"
	"strings"

	"xsbPro/log"
)

type NginxLogInfo struct {
	IP        string
	TimeStart int64
	Method    string
	URL       string
	Code      string
	BodySize  string
	Refer     string
	UA        string
	Real_IP   string
	ReqTime   string
}

type myRegexp struct {
	*regexp.Regexp
}

func (r *myRegexp)FindStringSubmatchMap(s string) map[string]string {
	captures := make(map[string]string)

	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures
	}

	for i, name := range r.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		//去掉所有双引号
		captures[name] = strings.Replace(match[i], `"`, ``, -1)
	}
	return captures
}

const (
	IP = `?P<IP>[\d.]*`
	TimeStart = `?P<TimeStart>\[[^\[\]]*\]`
	Request = `?P<Request>\"[^\"]*\"`
	Code = `?P<Code>\d+`
	BodySize = `?P<BodySize>\d+`
	Refer = `?P<Refer>\"[^\"]*\"`
	UA = `?P<UA>\"[^\"]*\"`
	Real_IP = `?P<Real_IP>([\d.]*|-)`
	ReqTime = `?P<ReqTime>\"[\d.]*\"`
)

var myExp = myRegexp{regexp.MustCompile(fmt.Sprintf(`(%s)\ -\ -\ (%s)\ (%s)\ (%s)\ (%s)\ (%s)\ (%s)\ (%s)\ (%s)`,
	IP, TimeStart, Request, Code, BodySize, Refer, UA, Real_IP, ReqTime))}

func AnalyNginxLines(lines string) []*NginxLogInfo {
	var infos []*NginxLogInfo
	for _, line := range (strings.Split(lines, "\n")) {
		if len(line) > 0 {
			result, err := analyANginxLine(line)
			if err != nil {
				log.SysF("analyANginxLine error:%s", err.Error())
				log.SysF("error line:%s", line)
				continue
			}
			if strings.Contains(result.URL, "/api/v") {
				api_name, err := getApiNameInfo(result.URL)
				if err != nil {
					log.SysF("analyANginxLine error:%s", err.Error())
					log.SysF("error line:%s", line)
					continue
				}
				result.URL = api_name
				infos = append(infos, result)
			}
		}
	}
	return infos
}

// 多于4个分隔标志  `/api/v1/bookShelf/authedBookShelfToCompany?bookShelf=new_books&authorized=1&_=1489136033875`
// 4个分隔标志  `/api/v1/banner/banners`
func getApiNameInfo(line string) (string, error) {
	lenth := len(line)
	var markList = []int{}
	j := 0

	for i := 0; i < lenth; i ++ {
		if string(line[i]) == "/" || string(line[i]) == "?" {
			markList = append(markList, i)
			j ++
		}
		if j >= 5 {
			break
		}
	}
	if j < 4 {
		return "", errors.New("未能正确找到api名称")
	} else if j == 4 {
		return line[markList[2] + 1:], nil
	}
	return line[markList[2] + 1:markList[4]], nil
}

func analyANginxLine(line string) (*NginxLogInfo, error) {
	mmap := myExp.FindStringSubmatchMap(line)
	if len(mmap) <= 0 {
		return nil, errors.New("未能完全匹配")
	}

	reqTime, err := convertTime(mmap["TimeStart"])
	if err != nil {
		return nil, err
	}

	method, url, err := analyRequestPart(mmap["Request"])
	if err != nil {
		return nil, err
	}

	return &NginxLogInfo{
		IP:        mmap["IP"],
		TimeStart: reqTime,
		Method:    method,
		URL:       url,
		Code:      mmap["Code"],
		BodySize:  mmap["BodySize"],
		Refer:     mmap["Refer"],
		UA:        mmap["UA"],
		Real_IP:   mmap["Real_IP"],
		ReqTime:   mmap["ReqTime"],
	}, nil
}

func convertTime(raw_string string) (int64, error) {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return 0, err
	}
	theTime, err := time.ParseInLocation("[02/Jan/2006:15:04:05 -0700]", raw_string, loc)
	if err != nil {
		return 0, err
	}
	return theTime.Unix(), nil
}

func analyRequestPart(raw_string string) (string, string, error) {
	alist := strings.Split(raw_string, " ")
	if len(alist) != 3 {
		return "", "", errors.New("analyRequestPart失败，不能取得正确信息")
	}
	return alist[0], alist[1], nil
}
