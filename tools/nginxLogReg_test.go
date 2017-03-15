package tools

import (
	"testing"
	"fmt"
)

func TestAnalyLineWithRealIP(t *testing.T) {
	r, err := analyANginxLine(`211.137.119.241 - - [12/Mar/2017:18:24:37 +0800] "GET /api/v1/group/info/1460000879_1460018631308714750 HTTP/1.1" 301 190 "-" "Dalvik/1.6.0 (Linux; U; Android 4.4.4; SM-W2015 Build/KTU84P)" 10.103.87.188 "0.000"`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(*r)
}

func TestAnalyLineWithoutRealIP(t *testing.T) {
	r, err := analyANginxLine(`117.35.118.74 - - [13/Mar/2017:10:16:32 +0800] "GET /api/v1/taskrecord/taskRecordsForUser/1460024917208992512?expired=false HTTP/1.1" 200 4044 "-" "Dalvik/2.1.0 (Linux; U; Android 5.1.1; OPPO A33 Build/LMY47V)" - "0.002"`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(*r)
}

func Test2(t *testing.T) {
	r, err := analyANginxLine(`117.136.50.214 - - [15/Mar/2017:11:12:27 +0800] "GET /api/v1/category/categoryReadingMaterials/rmcbs/1459838108 HTTP/1.1" 200 10026 "-" "Dalvik/2.1.0 (Linux; U; Android 5.1; m2 note Build/LMY47D)" - "0.001"`)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(*r)
	fmt.Println(getApiNameInfo(r.URL))
}