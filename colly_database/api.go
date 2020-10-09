package colly_database

import (
    "test-colly/entity"
)

func InsertUserInfoSql(user *entity.UserInfo) error {
    stmt, err := dbConn.Prepare(
        "INSERT INTO `tests`.`ai_virtual_user`( `nick_name`, `head_pic`) VALUES (?, ?);")
    defer stmt.Close()
    if err != nil {
        return err
    }
    _, err = stmt.Exec(user.HeadPic, user.NickName)
    if err != nil {
        return err
    }
    return nil

}

//func BatchInsertUserInfoSql(user *list.List) error {
//
//
//    key:="(?,?)"
//    values :=""
//
//    for user := user.Front(); user != nil; user = user.Next() {
//
//    }
//
//    stmt, err := dbConn.Prepare(
//       "INSERT INTO `xinlingshou`.`templ_user`( `head_pic`, `nick_name`) VALUES (?, ?)")
//    defer stmt.Close()
//    if err != nil {
//       return err
//    }
//    _, err = stmt.exec(user.HeadPic, user.NickName)
//    if err != nil {
//       return err
//    }
//    return nil
//
//}