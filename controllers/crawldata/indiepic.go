package crawldata

import(
	"sitepointgoapp/controllers/crawldata"
	"fmt"
	"src/github.com/go-martini/martini"
	"net/http"
	s "strings"
	"src/github.com/martini-contrib/encoder"
)

type Results struct {
	Err   int                  // 错误码
	Msg   string               // 错误信息
	Datas crawldata.ImageDatas // 数据，无数据时为nil
}

func main() {
	crawldata.Crawl()

	m := martini.New()
	route := martini.NewRouter()

	var (
		results Results
		err     error
	)

	m.Use(func(c martini.Context, w http.ResponseWriter, r *http.Request) {
		// 将encoder.JsonEncoder{}按照encoder.Encoder接口（注意大小写）类型注入到内部
		c.MapTo(encoder.JsonEncoder{}, (*encoder.Encoder)(nil))
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	})

	route.Get("/", func(enc encoder.Encoder) (int, []byte) {
		result := Results{10001, "Not Found Data", nil}
		return http.StatusOK, encoder.Must(enc.Encode(result))
	})

	route.Get("/api", func(enc encoder.Encoder) (int, []byte) {
		results.Datas, err = crawldata.GetAllImages()
		if err != nil {
			fmt.Println(s.Join([]string{"获取数据失败", err.Error()}, "-->"))
			result := Results{10001, "Data Error", nil}
			return http.StatusOK, encoder.Must(enc.Encode(result))
		} else {
			results.Err = 10001
			results.Msg = "获取数据成功"
			return http.StatusOK, encoder.Must(enc.Encode(results))
		}
	})

	m.Action(route.Handle)
	m.Run()

}