package db

import (
	"context"
	"fmt"
	"log"

	"github.com/seminhnva/e-com/internal/store"
)

var usernames = []string{
	"alex_chen",
	"maria_garcia",
	"james_wilson",
	"sarah_kim",
	"michael_brown",
	"emma_davis",
	"ryan_nguyen",
	"olivia_martin",
	"daniel_lee",
	"sophie_taylor",
}

var blogTags = []string{
	"Startups", "Remote Work", "AI Tools", "Side Hustle",
	"Cybersecurity", "Cloud Computing", "DevOps", "Open Source",
	"Machine Learning", "Blockchain", "UI/UX Design", "Productivity Hacks",
	"Entrepreneurship", "Career Growth", "Freelancing", "Digital Marketing",
	"Tech Trends", "Data Science", "Mobile Development", "Web Performance",
}

var comments = []string{
	"This is actually super useful 🔥",
	"I didn’t know this before lol",
	"Can someone explain this in simple terms?",
	"Finally something that makes sense",
	"100% agree with this",
	"This saved me so much time",
	"Why is this so underrated?",
	"I’ve been doing this wrong all along 😭",
	"Clean explanation, thanks!",
	"Does this work in production though?",
	"Bookmarking this for later",
	"Tried it and it works!",
	"This should be standard practice",
	"Not sure I fully get it yet",
	"This is next level stuff",
	"Anyone else confused by this?",
	"Short and clear, love it",
	"This broke my brain a bit lol",
	"Super helpful, appreciate it",
	"We need more posts like this",
}

func Seed(store store.Storage) error {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Printf("failed to create user %s: %v", user.Username, err)
		}
	}

	posts := generatePosts(1000, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Printf("failed to create post for user %d: %v", post.UserID, err)
		}
	}

	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Printf("failed to create comment for user %d on post %d: %v", comment.UserID, comment.PostID, err)
		}
	}
	log.Println("seeding completed")
	return nil
}

func generateUsers(n int) []*store.User {
	users := make([]*store.User, n)
	for i := 0; i < n; i++ {
		name := usernames[i%len(usernames)]
		users[i] = &store.User{
			Username: fmt.Sprintf("%s_%d", name, i+1),
			Email:    fmt.Sprintf("%s_%d@example.com", name, i+1),
			Password: "testing",
		}
	}
	return users
}
func generatePosts(n int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, n)
	for i := 0; i < n; i++ {
		user := users[i%len(users)]
		posts[i] = &store.Post{
			Content: fmt.Sprintf("This is the content of post %d from user %s", i+1, user.Username),
			Title:   fmt.Sprintf("Post Title %d ", i+1),
			UserID:  user.ID,
			Tags:    []string{blogTags[i%len(blogTags)], blogTags[(i+1)%len(blogTags)]},
			Version: 0,
		}
	}
	return posts
}

func generateComments(n int, users []*store.User, posts []*store.Post) []*store.Comment {
	commentsList := make([]*store.Comment, n)
	for i := 0; i < n; i++ {
		user := users[i%len(users)]
		post := posts[i%len(posts)]
		commentsList[i] = &store.Comment{
			Content: comments[i%len(comments)],
			UserID:  user.ID,
			PostID:  post.ID,
		}
	}
	return commentsList
}
