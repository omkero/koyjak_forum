package internal

import (
	"context"
	"fmt"
	"koyjak/config"
	"koyjak/internal/functions"
	"time"
)

type ForumType struct {
	ForumID          int       `json:"forum_id"`
	ForumCategory    string    `json:"forum_category"`
	ForumTitle       string    `json:"forum_title"`
	ForumDescription string    `json:"forum_description"`
	ThreadsCount     int       `json:"threads_count"`
	PostsCount       int       `json:"posts_count"`
	CreatedAt        time.Time `json:"created_at"`
	SafeTitle string `json:"safe_title"`
}

type ForumData struct {
	ForumCategory string `json:"forum_category"`
	Forum []ForumType `json:"forums"`
}
type ForumsResult struct {
	ForumType []ForumType
	Err error
}

type ForumTitleExistResult struct {
	IsExist bool
	Err error
}

func (Th *App) create_forum() {

}

func (Th *App) get_forums() ([]ForumData, error) {
	if config.Pool == nil {
		functions.Failed_db_connection()
	}

	var data []ForumData
	var sql_query string = `
	SELECT forum_category, 
       ARRAY_AGG(
         JSON_BUILD_OBJECT(
           'forum_id', forum_id,
           'forum_title', forum_title,
           'forum_description', forum_description,
           'threads_count', threads_count,
           'posts_count', posts_count,
           'created_at', created_at,
		   'safe_title', safe_title
         )
       ) AS forums
    FROM forums
    GROUP BY forum_category;
	`
	row, err := config.Pool.Query(context.Background(), sql_query)
	if err != nil {
		return []ForumData{}, err
	} 

	for row.Next() {
		var tempForum ForumData
		err := row.Scan(&tempForum.ForumCategory, &tempForum.Forum)
		 if err != nil {
	    	return []ForumData{}, err
		}
		 data = append(data,tempForum )
	}

	return data, nil
}

func (Th *App) is_forum_title_exist(forum_title string) (bool, error) {
	if config.Pool == nil {
		functions.Failed_db_connection()
	}

	var count int
	var sql_query string = `SELECT COUNT(*) AS count FROM forums WHERE forum_title = $1`

	err := config.Pool.QueryRow(context.Background(), sql_query, forum_title).Scan(&count)
	if err != nil {
		fmt.Println(err)
		return false, fmt.Errorf("ops somethings wnet wrong")
	}

	if count < 1 {
		return false, fmt.Errorf("forum_title not exist")
	}

	return count >= 1, nil
}

func (Th *App) filter_forums(forums []ForumType) {
	// Map of category => list of forums
	groupedForums := make(map[string][]ForumType)

	for _, forum := range forums {
		groupedForums[forum.ForumCategory] = append(groupedForums[forum.ForumCategory], forum)
	}

	// Just to print the result
	for category, forumsInCategory := range groupedForums {
		fmt.Printf("Category: %s\n", category)
		for _, f := range forumsInCategory {
			fmt.Printf("  - %s (ID: %d)\n", f.ForumTitle, f.ForumID)
		}
	}
}
