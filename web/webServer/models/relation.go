package models

import (
	"fmt"
)

// 点赞信息
type LikeInfo struct {
	FvID        int `json:"fvId"`
	UserAct     int `json:"userAct"`
	FavorNoteID int `json:"favorNoteId"`
}

// 根据笔记编号获取作者账号
func NoteToUser(noteId int) int {
	var userId int
	sltsql := "SELECT creatorAccount FROM noteInfo WHERE noteId=?"
	err := db.QueryRow(sltsql, noteId).Scan(&userId)
	if err != nil {
		fmt.Printf("scan failed, err:%v\n", err)
		return -1
	} else {
		return userId
	}
}

// 插入点赞信息
func NewLike(nl LikeInfo, noteId int) bool {
	sqlstr := `INSERT INTO favorTable (userAct, favorNoteId) VALUES (?,?)`
	ret, err := db.Exec(sqlstr, nl.UserAct, noteId)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return false
	}
	// 新插入数据的id
	theID, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return false
	}
	fmt.Printf("点赞成功！评论在数据库行数：%d\n", theID)
	return true
}

// 删除点赞信息
func DeleteLike(nl LikeInfo, noteId int) bool {
	// userId := NoteToUser(noteId)
	sqlstr := "DELETE from favorTable WHERE userAct=? AND favorNoteId=?"
	ret, err := db.Exec(sqlstr, nl.UserAct, noteId)
	if err != nil {
		fmt.Printf("delete failed, err:%v\n", err)
		return false
	}
	// 操作影响的行数
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return false
	}
	fmt.Printf("点赞信息 delete success, affected rows:%d\n", n)
	return true
}
