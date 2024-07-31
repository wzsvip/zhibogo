package handler

import(
	"Golang/liveurls"
	"Golang/list"
	"Golang/utils"
	"fmt"
	"net/http"
	"strings"
	"encoding/json"
	"strconv"
	"log"
	"os"
)

func duanyan(adurl string, realurl any) string {
	var liveurl string
	if str, ok := realurl.(string); ok {
		liveurl = str
	} else {
		liveurl = adurl
	}
	return liveurl
}

// vercel 平台会将请求传递给该函数，这个函数名随意，但函数参数必须按照该规则。
func Handler(w http.ResponseWriter, r * http.Request) {
    path: = r.URL.Path
    switch path {
        // 虎牙一起看
        case "/huyayqk.m3u":
            yaobj: = & list.HuyaYqk {}
            res,
            _: = yaobj.HuYaYqk("https://live.cdn.huya.com/liveHttpUI/getLiveList?iGid=2135")
            var result list.YaResponse
            json.Unmarshal(res, & result)
            pageCount: = result.ITotalPage
            pageSize: = result.IPageSize
            w.Header().Set("Content-Type", "application/octet-stream")
            w.Header().Set("Content-Disposition", "attachment; filename=huyayqk.m3u")
            utils.GetTestVideoUrl(w)

            for i: = 1;
            i <= pageCount;
            i++{
                apiRes, _: = yaobj.HuYaYqk(fmt.Sprintf("https://live.cdn.huya.com/liveHttpUI/getLiveList?iGid=2135&iPageNo=%d&iPageSize=%d", i, pageSize))
                var res list.YaResponse
                json.Unmarshal(apiRes, & res)
                data: = res.VList
                for _,
                value: = range data {
                    fmt.Fprintf(w, "#EXTINF:-1 tvg-logo=\"%s\" group-title=\"%s\", %s\n", value.SAvatar180, value.SGameFullName, value.SNick)
                    fmt.Fprintf(w, "%s/huya/%v\n", utils.GetLivePrefix(r), value.LProfileRoom)
                }
            }
            // 斗鱼一起看
        case "/douyuyqk.m3u":
            yuobj: = & list.DouYuYqk {}
            resAPI,
            _: = yuobj.Douyuyqk("https://www.douyu.com/gapi/rkc/directory/mixList/2_208/list")

            var result list.DouYuResponse
            json.Unmarshal(resAPI, & result)
            pageCount: = result.Data.Pgcnt

            w.Header().Set("Content-Type", "application/octet-stream")
            w.Header().Set("Content-Disposition", "attachment; filename=douyuyqk.m3u")
            utils.GetTestVideoUrl(w)

            for i: = 1;
            i <= pageCount;
            i++{
                apiRes, _: = yuobj.Douyuyqk("https://www.douyu.com/gapi/rkc/directory/mixList/2_208/" + strconv.Itoa(i))

                var res list.DouYuResponse
                json.Unmarshal(apiRes, & res)
                data: = res.Data.Rl

                for _,
                value: = range data {
                    fmt.Fprintf(w, "#EXTINF:-1 tvg-logo=\"https://apic.douyucdn.cn/upload/%s_big.jpg\" group-title=\"%s\", %s\n", value.Av, value.C2name, value.Nn)
                    fmt.Fprintf(w, "%s/douyu/%v\n", utils.GetLivePrefix(r), value.Rid)
                }
            }
            // YY轮播
        case "/yylunbo.m3u":
            yylistobj: = & list.Yylist {}
            w.Header().Set("Content-Type", "application/octet-stream")
            w.Header().Set("Content-Disposition", "attachment; filename=yylunbo.m3u")
            utils.GetTestVideoUrl(w)

            i: = 1
            for {
                apiRes: = yylistobj.Yylb(fmt.Sprintf("https://rubiks-idx.yy.com/nav/other/pnk1/448772?channel=appstore&compAppid=yymip&exposured=80&hdid=8dce117c5c963bf9e7063e7cc4382178498f8765&hostVersion=8.25.0&individualSwitch=1&ispType=2&netType=2&openCardLive=1&osVersion=16.5&page=%d&stype=2&supportSwan=0&uid=1834958700&unionVersion=0&y0=8b799811753625ef70dbc1cc001e3a1f861c7f0261d4f7712efa5ea232f4bd3ce0ab999309cac0d7869449a56b44c774&y1=8b799811753625ef70dbc1cc001e3a1f861c7f0261d4f7712efa5ea232f4bd3ce0ab999309cac0d7869449a56b44c774&y11=9c03c7008d1fdae4873436607388718b&y12=9d8393ec004d98b7e20f0c347c3a8c24&yv=1&yyVersion=8.25.0", i))
                var res list.ApiResponse
                json.Unmarshal([] byte(apiRes), & res)
                for _, value: = range res.Data.Data {
                        fmt.Fprintf(w, "#EXTINF:-1 tvg-logo=\"%s\" group-title=\"%s\", %s\n", value.Avatar, value.Biz, value.Desc)
                        fmt.Fprintf(w, "%s/yy/%v\n", utils.GetLivePrefix(r), value.Sid)
                    }
                if res.Data.IsLastPage == 1 {
                    break
                }
                i++
            }
            // 其他链接
        default:
            adurl: = "http://159.75.85.63:5680/d/ad/roomad/playlist.m3u8"
            params: = strings.Split(path, "/")

            // log.Println("request url: ", path)

            if len(params) >= 3 {
                // 解析成功
                // 平台
                platform: = params[1]
                // 房间号
                rid: = params[2]
                // fmt.Fprintf(w, "parsed platform=%s, room=%s", platform, rid)
                switch platform {
                    case "itv":
			if enableTV {
				itvobj := &liveurls.Itv{}
				cdn := c.Query("cdn")
				if ts == "" {
					itvobj.HandleMainRequest(c, cdn, rid)
				} else {
					itvobj.HandleTsRequest(c, ts)
				}
			} else {
				c.String(http.StatusForbidden, "公共服务不提供TV直播")
			}
		case "ysptp":
			if enableTV {
				ysptpobj := &liveurls.Ysptp{}
				if ts == "" {
					ysptpobj.HandleMainRequest(c, rid)
				} else {
					ysptpobj.HandleTsRequest(c, ts, c.Query("wsTime"))
				}
			} else {
				c.String(http.StatusForbidden, "公共服务不提供TV直播")
			}
		case "douyin":
			douyinobj := &liveurls.Douyin{}
			douyinobj.Rid = rid
			douyinobj.Stream = c.DefaultQuery("stream", "flv")
			c.Redirect(http.StatusMovedPermanently, duanyan(adurl, douyinobj.GetDouYinUrl()))
		case "douyu":
			douyuobj := &liveurls.Douyu{}
			douyuobj.Rid = rid
			douyuobj.Stream_type = c.DefaultQuery("stream", "flv")
			c.Redirect(http.StatusMovedPermanently, duanyan(adurl, douyuobj.GetRealUrl()))
		case "huya":
			huyaobj := &liveurls.Huya{}
			huyaobj.Rid = rid
			huyaobj.Cdn = c.DefaultQuery("cdn", "hwcdn")
			huyaobj.Media = c.DefaultQuery("media", "flv")
			huyaobj.Type = c.DefaultQuery("type", "nodisplay")
			if huyaobj.Type == "display" {
				c.JSON(200, huyaobj.GetLiveUrl())
			} else {
				c.Redirect(http.StatusMovedPermanently, duanyan(adurl, huyaobj.GetLiveUrl()))
			}
		case "bilibili":
			biliobj := &liveurls.BiliBili{}
			biliobj.Rid = rid
			biliobj.Platform = c.DefaultQuery("platform", "web")
			biliobj.Quality = c.DefaultQuery("quality", "10000")
			biliobj.Line = c.DefaultQuery("line", "first")
			c.Redirect(http.StatusMovedPermanently, duanyan(adurl, biliobj.GetPlayUrl()))
		case "youtube":
			ytbObj := &liveurls.Youtube{}
			ytbObj.Rid = rid
			ytbObj.Quality = c.DefaultQuery("quality", "1080")
			c.Redirect(http.StatusMovedPermanently, duanyan(adurl, ytbObj.GetLiveUrl()))
		case "yy":
			yyObj := &liveurls.Yy{}
			yyObj.Rid = rid
			yyObj.Quality = c.DefaultQuery("quality", "4")
			c.Redirect(http.StatusMovedPermanently, duanyan(adurl, yyObj.GetLiveUrl()))
                    default:
                        fmt.Fprintf(w, "Unknown platform=%s, room=%s", platform, rid)
                }
            } else {
                log.Println("Invalid path:", path)
                //http.Error(w, "welcome", http.StatusInternalServerError)
                fmt.Fprintf(w, "<a href='/huyayqk.m3u'>虎牙一起看</a><br>")
                fmt.Fprintf(w, "<a href='/douyuyqk.m3u'>斗鱼一起看</a><br>")
                fmt.Fprintf(w, "<a href='/yylunbo.m3u'>YY轮播</a><br>")
                fmt.Fprintf(w, os.Getenv("LIVE_URL"))
            }
            // log.Println("Invalid path:", path)
            // fmt.Fprintf(w, "<h1>链接错误!</h1>")
    }
    return
}
