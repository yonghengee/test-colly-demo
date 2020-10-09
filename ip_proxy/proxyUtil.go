package ip_proxy

import (
    "container/list"
    "fmt"
    "github.com/gocolly/colly"
    "github.com/gocolly/colly/extensions"
    "github.com/gocolly/colly/proxy"
    "net"
    "net/http"
    "test-colly/utils"
    "time"
)

var (
	IpProxyChannel chan string
	IpProxyArray []string
    ProxyIPList list.List
)

func init() {
    IpProxyChannel =make(chan string,100)
    IpProxyArray  = make([]string ,0,100)
}

func InitIPProxy(html string) []string{
    s := getHtml(html)
    splice := utils.Splice(s, `<td data-title="IP">`, `</td>`)

    ipList :=conventSlice(splice)
    splice2 := utils.Splice(s, `<td data-title="PORT">`, `</td>`)
    portList :=conventSlice(splice2)

    totalProxy := make([]string,0,100)
    for index:=0;index<splice.Len();index++{
        totalProxy = append(totalProxy,ipList[index]+":"+portList[index])
    }

    //获取到的IPproxy
    for index:=0;index<len(totalProxy);index++{
       result := CheckExist(totalProxy[index],`https://tieba.baidu.com/p/6916745651`)
       if !result {
           totalProxy = append(totalProxy[:index], totalProxy[index+1:]...)
           index--
           continue
       }else{
           IpProxyChannel <- totalProxy[index]
           IpProxyArray = append(IpProxyArray, totalProxy[index])
       }
    }

    for index:=0;index<len(totalProxy);index++{
        fmt.Println("final",totalProxy[index])
    }

    return totalProxy
}

func GetProxyIp() string {
    ipProxy := <-IpProxyChannel
    return ipProxy
}


//
//func init()  {
//    ProxyIPList = *GetIPProxy("https://www.kuaidaili.com/free/inha/")
//}

func conventSlice(list  list.List)([]string)  {
    result := make([]string,0,100)
    for i:=list.Front();i !=nil;i=i.Next() {
        if i.Value.(string) == "" {
            continue
        }
        result = append(result,i.Value.(string))
    }
    return result
}


func getHtml(html string)(string)  {

    result:=""
    c := colly.NewCollector(
        func(collector *colly.Collector) {
            extensions.RandomUserAgent(collector)
        })

    c.WithTransport(&http.Transport{
        Proxy: http.ProxyFromEnvironment,
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
    //c.SetProxy("http://119.57.156.90:53281")
    //init
    ////设置速率
    c.OnRequest(func(r *colly.Request) {
        fmt.Println(r.URL,"->>begin")
    })

    c.OnResponse(func(r *colly.Response) {
        result = string(r.Body)
        fmt.Println(r.StatusCode,"->>end")
    })
    c.OnHTML("", func(e *colly.HTMLElement) {
        //e.ForEach(".j_th_tit ", func(i int, element *colly.HTMLElement) {
        //    fmt.Println(element.Text)
        //})
        //fmt.Println(e.Text)
    })
    //入口
    c.Visit(html)

    return result
}

func CheckExist(proxyIp,html string)(bool)  {
    result:=false
    c := colly.NewCollector(
        func(collector *colly.Collector) {
            extensions.RandomUserAgent(collector)
        })

    c.WithTransport(&http.Transport{
        Proxy: http.ProxyFromEnvironment,
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
    if p, err := proxy.RoundRobinProxySwitcher(proxyIp); err == nil {
        c.SetProxyFunc(p)
    }
    //init
    ////设置速率
    c.OnRequest(func(r *colly.Request) {
        fmt.Println(r.URL,"->>begin")
    })

    c.OnResponse(func(r *colly.Response) {
        if r.StatusCode ==200 {
            result = true
        }
        fmt.Println(r.StatusCode,result)
    })
    //入口
    c.Visit(html)
    return result
}
