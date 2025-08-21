package app

import "time"

var EMAILS []string = []string{
	"test1@example.com",
	"test2@example.com",
	"test3@example.com",
	"test4@example.com",
	"test5@example.com",
	"test6@example.com",
	"test7@example.com",
	"test8@example.com",
	"test9@example.com",
	"test10@example.com",
	"test11@example.com",
	"test12@example.com",
	"test13@example.com",
	"test14@example.com",
	"test15@example.com",
	"test16@example.com",
	"test17@example.com",
	"test18@example.com",
	"test19@example.com",
	"test20@example.com",
	"test21@example.com",
	"test22@example.com",
	"test23@example.com",
	"test24@example.com",
	"test25@example.com",
	"test26@example.com",
	"test27@example.com",
	"test28@example.com",
	"test29@example.com",
	"test30@example.com",
}

var TOPICS []string = []string{
	"go",
	"java",
	"python",
	"c",
	"c++",
	"zig",
	"rust",
	"ocaml",
	"javascript",
	"ruby",
	"kotlin",
	"lua",
	"shell",
}

func getNumberOfEmailsAndTopics() (emails int, topics int) {
	emails = len(EMAILS)
	topics = len(TOPICS)
	return emails, topics
}

// getTimeDiff return difference in milli seconds between t2 and t1 (t2 - t1)
func getTimeDiff(t1, t2 time.Time) int64 {
	return int64(t2.Sub(t1).Milliseconds())
}
