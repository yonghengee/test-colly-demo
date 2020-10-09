package utils

import (
    "container/list"
    "strings"
)

/**
  提取中间有效值
*/
func Splice(html, head, tail string) (list.List) {
    result := list.New()
    count := strings.Count(html, head)

    htmlTempl := html
    headLen := len(head)
    tailLen := len(tail)

    for count > 0 {
        //fmt.Println("打印遍历 ： ", count)
        count--
        element := ""
        headIndex := strings.Index(htmlTempl, head)
        if headIndex == -1 {
            result.PushBack(element)
            continue
        }
        //截断头部之前的长度
        htmlTempl = htmlTempl[headIndex:]
        tailIndex := strings.Index(htmlTempl, tail)
        if tailIndex == -1 {
            result.PushBack(element)
            continue
        }
        element = htmlTempl[headLen: tailIndex]
        result.PushBack(element)
        //删除之前的
        htmlTempl = htmlTempl[tailIndex+tailLen:]

    }
    return *result

}

