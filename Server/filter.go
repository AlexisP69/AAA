package forum

import "database/sql"

func FilterByCategory(db *sql.DB, posts []Posts) []PostWithComments {
	var postSlice []PostWithComments
	for _, post := range posts {
		var t PostWithComments
		t.Post = post
		t.EveryComments = SelectAllComments(db, post.Id)
		postSlice = append(postSlice, t)
	}
	return postSlice
}
