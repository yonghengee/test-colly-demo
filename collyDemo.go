package main

import (
    "container/list"
    "fmt"
    "github.com/gocolly/colly"
    "github.com/gocolly/colly/extensions"
    "github.com/gocolly/colly/proxy"
    "net"
    "net/http"
    "strconv"
    "strings"
    "test-colly/colly_database"
    "test-colly/entity"
    "test-colly/ip_proxy"
    "test-colly/utils"
    "time"
)

var c *colly.Collector

var ipProxy string

func main() {

    //该time为获取 ip代理网站的次数 通常也是一页一页的获取的
    timer := 2
    for index := 1; index<= timer;index++ {
       ip_proxy.InitIPProxy("https://www.kuaidaili.com/free/inha/" + strconv.Itoa(index) + "/")
       //获取到的IPproxy
    }
    links := list.New()
    getLimitPage(links,20)

    //_, u := SpliceTwoKey(`<a style="" target="_blank" class="p_author_face " href="/home/main?un=89b2aJ9H49&ie=utf-8&id=tb.1.7664345e.T9abDfXOeSCsgUlb1mfXsg&fr=pb&ie=utf-8"><img username="89b2aJ9H49" class="" src="//himg.bdimg.com/sys/portrait/item/tb.1.7664345e.T9abDfXOeSCsgUlb1mfXsg"/></a>`, `<img username="|" class=""`, `src="|"/></a>`)
    //for i := 0; i < len(u); i++ {
    //   err:=colly_database.InsertUserInfoSql(&u[i])
    //   if err !=nil {
    //      fmt.Println(err)
    //   }
    //}
}

/**
  获取第一页的链接 和 其他页的地址
*/
func getLimitPage(links *list.List, pageNum int) (visitLinks *list.List) {

    //手动模拟获取分页地址
    for index := 3; index < pageNum; index++ {
        html := "https://tieba.baidu.com/f?kw=%E6%AF%8D%E5%A9%B4&ie=utf-8&pn=" + strconv.Itoa(index*50)
        links.PushBack(html)
    }

    detailRelateLinks := list.New()
    for i := links.Front(); i != nil; i = i.Next() {

        //爬取所有详情页面link
        htmlLink := i.Value.(string)
        htmlResponse := getHtmlResponse(htmlLink)
        splice := utils.Splice(htmlResponse, `<div class="threadlist_lz clearfix">`, `</a>`)
        //insert
        detailRelateLinks.PushBackList(&splice)
    }

    detailRealLinkList := list.New()
    for i := detailRelateLinks.Front(); i != nil; i = i.Next() {
        fmt.Println(i.Value.(string))
        relateLinks := utils.Splice(i.Value.(string), `<a rel="noreferrer"`, `target="_blank" class="j_th_tit "`)
        for i := relateLinks.Front(); i != nil; i = i.Next() {
            details := utils.Splice(i.Value.(string), `href="`, `" title="`)
            detailRealLinkList.PushBack("https://tieba.baidu.com" + details.Front().Value.(string))
        }
    }

    infoData := list.New()
    for i := detailRealLinkList.Front(); i != nil; i = i.Next() {
        page := 5

        //模拟分页
        for j := 0; j < page; j++ {
            html := getHtmlResponse(i.Value.(string) + "?/pn=" + strconv.Itoa(j))

            splice := utils.Splice(html, `<a style="" target="_blank" class="p_author_face "`, `</a>`)
            infoData.PushBackList(&splice)
        }
    }

    for i := infoData.Front(); i != nil; i = i.Next() {
        _, userArray := SpliceTwoKey(i.Value.(string), `<img username="|" class=""`, `src="|"/>`)
        for i := 0; i < len(userArray); i++ {
            //插入数据库
            colly_database.InsertUserInfoSql(&userArray[i])
        }

    }

    return visitLinks
}

/**
  详情中的分页
*/
func getLimitDetailPage(links *list.List, pageNum int) {

    html := ""
    for index := 0; index < pageNum; index++ {
        html = "tieba.baidu.com/f?kw=%E6%AF%8D%E5%A9%B4&ie=utf-8&pn=" + strconv.Itoa(index*50)
        links.PushBack(html)
    }

}

//func randomProxySwitcher(_ *http.Request) (*url.URL, error) {
//    return ip_proxy.IpProxyArray[rand.Intn(len(ip_proxy.IpProxyArray))], nil
//}

func getHtmlResponse(htmLink string) (resultHtml string) {

    //正式请求前先判断是否能正常请求
    for !ip_proxy.CheckExist(ip_proxy.CurrentIpProxy, htmLink) {
        ip_proxy.CurrentIpProxy = ip_proxy.GetProxyIp()
    }

    c = colly.NewCollector(
        func(collector *colly.Collector) {
            extensions.RandomUserAgent(collector)
        })

    c.WithTransport(&http.Transport{
        //Proxy: http.ProxyFromEnvironment,
        DialContext: (&net.Dialer{
            Timeout:   30 * time.Second,
            KeepAlive: 30 * time.Second,
            DualStack: true,
        }).DialContext,
        MaxIdleConns:          100,
        IdleConnTimeout:       90 * time.Second,
        TLSHandshakeTimeout:   10 * time.Second,
        ExpectContinueTimeout: 1 * time.Second,
    })

    if p, err := proxy.RoundRobinProxySwitcher( ip_proxy.IpProxyArray...); err == nil {
        c.SetProxyFunc(p)
    }

    //c.SetProxy(ip_proxy.CurrentIpProxy)

    //init
    ////设置速率
    c.OnRequest(func(r *colly.Request) {
        fmt.Println(r.URL, "->>begin")
    })

    c.OnResponse(func(r *colly.Response) {
        resultHtml = string(r.Body)
        fmt.Println(r.StatusCode, "->>end")
    })
    c.OnHTML("", func(e *colly.HTMLElement) {
        //e.ForEach(".j_th_tit ", func(i int, element *colly.HTMLElement) {
        //    fmt.Println(element.Text)
        //})
        //fmt.Println(e.Text)
    })
    //入口
    c.Visit(htmLink)

    return resultHtml
}

func getLowestPage(html string) {
    //init
    c := colly.NewCollector(
        colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36"), )

    ////设置速率
    c.OnRequest(func(r *colly.Request) {
        fmt.Println("begin")
    })

    c.OnResponse(func(r *colly.Response) {
        html = string(r.Body)
    })
    c.OnHTML("", func(e *colly.HTMLElement) {
        //e.ForEach(".j_th_tit ", func(i int, element *colly.HTMLElement) {
        //    fmt.Println(element.Text)
        //})
        //fmt.Println(e.Text)
    })

    //入口
    c.Visit("https://tieba.baidu.com/p/6817072978")

    //爬取最下层数据
    //printf(html,`<a style="" target="_blank" class="p_author_face " |<li class="d_nameplate">`,`<img username="^" class=""|src="//^"/></a>`)
}

func SpliceTwoKey(html, firstKey, secondKey string) (list.List, []entity.UserInfo) {
    result := list.New()
    userInfoArray := make([]entity.UserInfo, 0, 1000)

    firsts := strings.Split(firstKey, "|")
    seconds := strings.Split(secondKey, "|")
    firstsResult := utils.Splice(html, firsts[0], firsts[1])
    secondsResult := utils.Splice(html, seconds[0], seconds[1])

    length := 0
    if firstsResult.Len() > secondsResult.Len() {
        length = firstsResult.Len()
    } else {
        length = secondsResult.Len()
    }

    fir := firstsResult.Front()
    sec := secondsResult.Front()
    for i := 0; i < length; i++ {
        nick := sec.Value
        head := fir.Value
        userInfo := entity.UserInfo{HeadPic: head.(string), NickName: nick.(string)}
        userInfoArray = append(userInfoArray, userInfo)
        result.PushBack(userInfo)
        fir = fir.Next()
        sec = sec.Next()
    }

    return *result, userInfoArray
}
