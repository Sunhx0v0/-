package models

import (
	"fmt"
)

// 评论信息
type Comment struct {
	CommentatorID   int64    `json:"commentatorId"`   // 评论者账号
	CommentatorName string   `json:"commentatorName"` // 评论者
	CommentID       int64    `json:"commentId"`       // 评论编号
	CommentTime     string   `json:"commentTime"`     // 评论时间
	Content         string   `json:"content"`         // 评论内容
	Portrait        string   `json:"portrait"`        // 用户头像
	State           int      `json:"state"`           //是否已读
	AtName          []string `json:"atName"`          //@人的名字
	AtLocation      []int    `json:"atLocation"`      //@的位置
}

// 获取评论
func GetCommentInfo(Id, option int) (comments []Comment, totalState int, ok bool) {
	//option字段用来判断是笔记获取评论还是消息获取评论
	ok = true
	var sqlstr string
	//笔记中获取评论
	notesql := `SELECT c.commentId, c.commentatorId, c.content, c.commentTime, c.state, u.userName, u.portrait
	FROM commentInfo c, userInfo u
	WHERE c.noteID = ? AND c.commentatorId = u.userAccount`
	//消息列表中获取评论
	messagesql := `SELECT c.commentId, c.commentatorId, c.content, c.commentTime, c.state, u.userName, u.portrait
	FROM commentInfo c, userInfo u, noteInfo n
	WHERE u.userAccount=? AND u.userAccount=n.creatorAccount AND c.noteID = n.noteId`
	if option == 0 {
		sqlstr = notesql
	} else {
		sqlstr = messagesql
	}
	rows, err := db.Query(sqlstr, Id)
	if err != nil {
		fmt.Printf("评论query failed, err:%v\n", err)
		ok = false
		return
	}
	// 关闭rows释放持有的数据库链接
	defer rows.Close()

	//整体的状态，1表示已读
	totalState = 1
	// 循环读取结果集中的数据
	for rows.Next() {
		var cmt Comment
		err := rows.Scan(&cmt.CommentID, &cmt.CommentatorID, &cmt.Content, &cmt.CommentTime, &cmt.State, &cmt.CommentatorName, &cmt.Portrait)
		if err != nil {
			fmt.Printf("评论scan failed, err:%v\n", err)
			ok = false
			return
		}
		totalState = totalState * cmt.State
		comments = append(comments, cmt)
	}
	return comments, totalState, ok
}

// 插入评论信息
func NewComment(nc Comment, noteId int) (int, bool) {
	sqlstr := `INSERT INTO commentInfo (noteID, commentatorId, content, commentTime) VALUES (?,?,?,?)`
	ret, err := db.Exec(sqlstr, noteId, nc.CommentatorID, nc.Content, nc.CommentTime)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return -1, false
	}
	// 新插入数据的id
	theID, err := ret.LastInsertId()
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return -1, false
	}
	fmt.Printf("评论成功！评论在数据库行数：%d\n", theID)
	return int(theID), true
}

// 删除评论信息
func DeleteComment(commentId int) bool {
	sqlstr := "DELETE FROM commentInfo WHERE commentId=?"
	ret, err := db.Exec(sqlstr, commentId)
	if err != nil {
		fmt.Printf("评论删除失败, err:%v\n", err)
		return false
	}
	// 操作影响的行数
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("评论删除Get RowsAffected failed, err:%v\n", err)
		return false
	}
	fmt.Printf("评论信息 delete success, affected rows:%d\n", n)
	return true
}

// 把某条评论设为已读
func SetCommentState(commentId int) bool {
	sqlstr := "UPDATE commentInfo SET state=1 WHERE commentId=?"
	ret, err := db.Exec(sqlstr, commentId)
	if err != nil {
		fmt.Printf("评论状态update failed, err:%v\n", err)
		return false
	}
	// 操作影响的行数
	n, err := ret.RowsAffected()
	if err != nil {
		fmt.Printf("评论状态get RowsAffected failed, err:%v\n", err)
		return false
	}
	fmt.Printf("评论状态修改编号：%d\n", n)
	return true
}
