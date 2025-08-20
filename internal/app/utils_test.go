package app

import "testing"

func Test_app_utils_getNumberOfEmailsAndTopics(t *testing.T) {
	EMAILS = []string{
		"test1@example.com",
		"test2@example.com",
		"test3@example.com",
		"test4@example.com",
	}
	TOPICS = []string{
		"go",
		"java",
		"python",
	}
	emails, topics := getNumberOfEmailsAndTopics()
	if emails != 4 {
		t.Errorf("Expected 4 emails, got %d", emails)
	}
	if topics != 3 {
		t.Errorf("Expected 3 topics, got %d", topics)
	}
}
