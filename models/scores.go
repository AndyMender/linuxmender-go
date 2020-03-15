package models

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"

	"linuxmender/paths"
	"linuxmender/utilities"
)

var sessionTTL, _ = time.ParseDuration("24h")

// ScoreManager is responsible for setting/getting metrics in Redis
type ScoreManager struct {
	Conn *redis.Client
}

// SetSession sets the user session key in Redis
func (mgr *ScoreManager) SetSession(sessionID string) {
	mgr.Conn.Set(
		fmt.Sprintf("%v:%v", paths.SessionsRedisPath, sessionID), 
		time.Now().Format(time.RFC3339),
		sessionTTL, // 24 hours
	)
}

// GetSession gets the user session key from Redis
func (mgr *ScoreManager) GetSession(sessionID string) (string, bool) {
	sessionTime, err := mgr.Conn.Get(
		fmt.Sprintf("%v:%v", paths.SessionsRedisPath, sessionID),
	).Result()
	if err == redis.Nil {
		return "", false
	} else if err != nil {
		log.Println(err)
		return "", false
	}

	return sessionTime, true
}

// WasLikedBy checks whether a blog entry was liked by a specific user
func (mgr *ScoreManager) WasLikedBy(sessionID string, entryID int) bool {
	key, err := mgr.Conn.Exists(
		fmt.Sprintf(
			"%v:%v:likedby:%v", 
			paths.PostsRedisPath, 
			entryID, 
			sessionID,
		),
	).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		log.Println(err)
		return false
	}

	// 1 - exists, 0 - doesn't exist
	return key == 1
}

// SetLikedBy attaches a user session ID to a blog entry key
func (mgr *ScoreManager) SetLikedBy(sessionID string, entryID int) {
	mgr.Conn.Set(
		fmt.Sprintf(
			"%v:%v:likedby:%v", 
			paths.PostsRedisPath, 
			entryID,
			sessionID,
		), 
		time.Now().Format(time.RFC3339),
		sessionTTL, // 24 hours
	)
}

// LikeEntry bumps the "likes" score for a blog entry key
func (mgr *ScoreManager) LikeEntry(sessionID string, entryID int) {
	// Get + set "likedby" entry + session combination session to avoid superfluous counts
	if ok := mgr.WasLikedBy(sessionID, entryID); ok {
		return
	}
	mgr.SetLikedBy(sessionID, entryID)

	mgr.Conn.Incr(
		fmt.Sprintf("%v:%v:likes", paths.PostsRedisPath, entryID),
	)
}

// GetLikes returns the "likes" score for a blog entry key
func (mgr *ScoreManager) GetLikes(entryID int) int {
	// Get raw score
	scoreRaw, err := mgr.Conn.Get(
		fmt.Sprintf("%v:%v:likes", paths.PostsRedisPath, entryID),
	).Result()
	if err == redis.Nil {
		return 0
	} else if err != nil {
		log.Println(err)
		return 0
	}

	// Convert score to a number
	scoreNum, err := strconv.Atoi(scoreRaw)
	if err != nil {
		log.Println(err)
		return 0
	}

	return scoreNum
}

// AddVisit increments the global visitor/hit counter and the daily counter
func (mgr *ScoreManager) AddVisit() {
	// Increment daily counter
	mgr.Conn.Incr(
		fmt.Sprintf("%v:%v", paths.VisitorsRedisPath, time.Now().Format(utilities.TimeFormat)),
	)

	// Increment global counter
	mgr.Conn.Incr(
		fmt.Sprintf("%v:all", paths.VisitorsRedisPath),
	)
}

// GetVisits gets the value of the global visitor/hit counter
func (mgr *ScoreManager) GetVisits() int {
	// Get raw counter
	counterRaw, err := mgr.Conn.Get(
		fmt.Sprintf("%v:all", paths.VisitorsRedisPath),
	).Result()
	if err == redis.Nil {
		return 0
	} else if err != nil {
		log.Println(err)
		return 0
	}

	// Convert counter to a number
	counterNum, err := strconv.Atoi(counterRaw)
	if err != nil {
		log.Println(err)
		return 0
	}
	
	return counterNum
}
