package tools

import (
	"sync"
	"xsbPro/log"
	"strings"
)

var (
	api_file string = "./api.list"
	Cache    *SafeMap
)

func InitCache() {
	var err error
	Cache = NewSafeMap()
	err = Cache.InitWithFile(api_file)
	if err != nil {
		panic(err)
	}
}

type NginxStatistics struct {
	URL     string
	Success int
	Fail    int
	Delay   int
}

func newNginxStatistics(URL string, success, fail, delay int) *NginxStatistics {
	return &NginxStatistics{
		URL:     URL,
		Success: success,
		Fail:    fail,
		Delay:   delay,
	}
}

type SafeMap struct {
	lock *sync.RWMutex
	bm   map[string]*NginxStatistics
}

// NewBeeMap return new safemap
func NewSafeMap() *SafeMap {
	return &SafeMap{
		lock: new(sync.RWMutex),
		bm:   make(map[string]*NginxStatistics),
	}
}


//增加一个计数
func (m *SafeMap) AddCounter(URL string, success, fail, delay int) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if value, ok := m.bm[URL]; !ok {
		m.bm[URL] = newNginxStatistics(URL, success, fail, delay)
	} else {
		value.Success += 1
	}
}


//清零所有api计数
func (m *SafeMap) Reset() {
	m.lock.Lock()
	defer m.lock.Unlock()
	for k, v := range m.bm {
		m.bm[k] = &NginxStatistics{
			URL:     v.URL,
			Success: 0,
			Fail:    0,
			Delay:   0,
		}
	}
}

//利用文件内容初始化
func (m *SafeMap) InitWithFile(file_path string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	content, err := getLinesBehindLine(file_path, 1)
	if err != nil {
		return err
	}
	content = strings.Replace(content, "\n", "", -1)
	api_num := 0
	for _, single_api := range (strings.Split(content, ";")) {
		if len(single_api) > 1 {
			m.bm[single_api] = newNginxStatistics(single_api, 0, 0, 0)
			api_num += 1
		}
	}
	log.InfoF("共计读取到历史api个数(%d)个", api_num)
	return nil
}


// Items returns all items in safemap.
func (m *SafeMap) Items() map[string]*NginxStatistics {
	m.lock.RLock()
	defer m.lock.RUnlock()
	r := make(map[string]*NginxStatistics)
	for k, v := range m.bm {
		r[k] = v
	}
	return r
}

func SaveResultsToCache(infos []*NginxLogInfo) {
	for _, info := range (infos) {
		if len(info.URL) <= 0 {
			log.SysF("结构含有不正确的URL信息: %v", *info)
			continue
		}
		success := 0
		fail := 0
		delay := 0
		if info.Code == "200" || info.Code == "301" || info.Code == "302" || info.Code == "303" || info.Code == "304" {
			success = 1
		} else {
			fail = 1
		}
		//大于300ms的访问
		if info.ReqTime >= "0.300" {
			delay = 1
		}
		Cache.AddCounter(info.URL, success, fail, delay)
	}
}

func SaveApiHistory() {
	content := ""
	for _, item := range (Cache.Items()) {
		content = content + ";" + item.URL
	}
	writeLineToFile(api_file, content)
}