package models

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"

	"linuxmender/paths"
)

// ScoreManager is responsible for setting/getting metrics in Redis
type ScoreManager struct {
	Conn *redis.Client
}

// SetSession sets the user session key in Redis
func (mgr *ScoreManager) SetSession(uuidString string) {
	mgr.Conn.Set(
		fmt.Sprintf("%v:%v", paths.SessionsRedisPath, uuidString), 
		time.Now().Format(time.RFC3339),
		24 * 60 * 60, // 24 hours
	)
}

// GetSession gets the user session key from Redis
func (mgr *ScoreManager) GetSession(uuidString string) (string, bool) {
	sessionTime, err := mgr.Conn.Get(
		fmt.Sprintf("%v:%v", paths.SessionsRedisPath, uuidString),
	).Result()
	if err == redis.Nil {
		return "", false
	} else if err != nil {
		log.Println(err)
		return "", false
	}

	return sessionTime, true
}

// LikeEntry bumps the "like" score for a blog entry key
func (mgr *ScoreManager) LikeEntry(sessionID string, entryID int) {
	// Get + set user session to avoid superfluous counts
	if _, ok := mgr.GetSession(sessionID); ok {
		return
	}
	mgr.SetSession(sessionID)

	mgr.Conn.Incr(
		fmt.Sprintf("%v:%v:likes", paths.PostsRedisPath, entryID),
	)
}
