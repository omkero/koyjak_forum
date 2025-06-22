package internal

import (
	"context"
	"fmt"
	"koyjak/config"
	"koyjak/internal/functions"
	"time"
)

type ForumType struct {
	ForumID          int        `json:"forum_id"`
	ForumCategory    string     `json:"forum_category"`
	ForumDescription string     `json:"forum_description"`
	ForumTitle       string     `json:"forum_title"`
	ForumPathTitle   string     `json:"forum_path_title"`
	ThreadsCount     int        `json:"threads_count"`
	PostsCount       int        `json:"posts_count"`
	CreatedAt        time.Time  `json:"created_at"`
	LatestThread     ThreadType `json:"thread_type"`
}

type ForumsResult struct {
	ForumType []ForumType
	Err       error
}

type ForumTitleExistResult struct {
	IsExist bool
	Err     error
}

type ForumData struct {
	Forum map[string][]ForumType `json:"forums"`
}

var ForumList = []ForumType{
	{
		ForumID:          1,
		ForumTitle:       "Frontend",
		ForumPathTitle:   "Frontend",
		ForumCategory:    "Web Development",
		ForumDescription: "Tihs is Web dev",
		ThreadsCount:     0,
		PostsCount:       0,
		LatestThread:     ThreadType{},
	},
	{
		ForumID:          2,
		ForumTitle:       "Back end",
		ForumPathTitle:   "Back-end",
		ForumCategory:    "Web Development",
		ForumDescription: "Tihs is Back Web dev",
		ThreadsCount:     0,
		PostsCount:       0,
		LatestThread:     ThreadType{},
	},
	{
		ForumID:          2,
		ForumTitle:       "Low Level ESP32",
		ForumPathTitle:   "Low-Level-ESP32",
		ForumCategory:    "Hardware",
		ForumDescription: "Tihs is Back Web dev",
		ThreadsCount:     0,
		PostsCount:       0,
		LatestThread:     ThreadType{},
	},
}

func (Th *App) get_forums() (map[string][]ForumType, error) {
	if config.Pool == nil {
		functions.Failed_db_connection()
	}

	// make sure to select * fromus and it will return array of forums loop

	var result []ForumType

	var sql_query string = `
    SELECT 
        forums.*,
        threads.*,
        users.user_id,
	    users.username,
	    users.email_address,
        users.created_at
        FROM forums
    INNER JOIN threads ON forums.forum_title = threads.forum_title
    INNER JOIN users ON threads.user_id = users.user_id;

	`
	row, err := config.Pool.Query(context.Background(), sql_query)
	if err != nil {
		fmt.Println(err)
		return map[string][]ForumType{}, fmt.Errorf("ops something went wrong.")
	}

	for row.Next() {
		var temp ForumType

		// if one missing or not scanned its going to be bad
		err = row.Scan(
			&temp.ForumID, &temp.ForumTitle, &temp.ThreadsCount, &temp.PostsCount, &temp.CreatedAt,
			&temp.LatestThread.ThreadID, &temp.LatestThread.UserID, &temp.LatestThread.ThreadTitle, &temp.LatestThread.ThreadContent, &temp.LatestThread.CreatedAt, &temp.LatestThread.SafeUrl, &temp.LatestThread.ThreadToken,
			&temp.LatestThread.ForumTitle, &temp.LatestThread.Member.UserID, &temp.LatestThread.Member.UserName, &temp.LatestThread.Member.EmailAddress,
			&temp.LatestThread.Member.CreatedAt,
		)
		if err != nil {
			// scanning issue handle here
			fmt.Println(err)
		}

		result = append(result, temp)
	}

	// group all forums
	var grouped = make(map[string][]ForumType)

	for _, forum := range ForumList {
		temp := ForumType{
			ForumID:          forum.ForumID,
			ForumTitle:       forum.ForumTitle,
			ForumPathTitle:   forum.ForumPathTitle,
			ForumCategory:    forum.ForumCategory,
			ForumDescription: forum.ForumDescription,
		}

		for _, item := range result {
			if item.ForumTitle == forum.ForumTitle {
				// convert timestamptz into readable stirng
				var lastTime = functions.TimeAgo(item.LatestThread.CreatedAt)

				temp.ThreadsCount = item.ThreadsCount // set threads_count
				temp.PostsCount = item.PostsCount

				// assign LatestThread
				temp.LatestThread = ThreadType{
					ThreadID:       item.LatestThread.ThreadID,
					ThreadTitle:    item.LatestThread.ThreadTitle,
					CreatedAtSince: lastTime,
					Member:         item.LatestThread.Member,
					SafeUrl:        item.LatestThread.SafeUrl,
				}
			}

		}

		// finally append the forum into the map based on ForumCategory
		grouped[forum.ForumCategory] = append(grouped[forum.ForumCategory], temp)
	}

	return grouped, nil
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

func (Th *App) is_forum_exist(forum_title string) bool {
	var count int

	for _, item := range ForumList {
		if item.ForumTitle == forum_title {
			count += 1
		}
	}

	return count >= 1
}

/*
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
		data = append(data, tempForum)
	}

	return data, nil
}
*/
