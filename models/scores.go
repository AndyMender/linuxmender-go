package models

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"

	"linuxmender/paths"
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
	_, err := mgr.Conn.Exists(
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

	return true
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
	scoreNum, _ := strconv.Atoi(scoreRaw)

	return scoreNum
}
