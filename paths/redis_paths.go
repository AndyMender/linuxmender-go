package paths

// ProjectRedisPath is the "root" of all Redis paths
var ProjectRedisPath = "linuxmender"

// SessionsRedisPath is the sub-path for user session keys
var SessionsRedisPath = ProjectRedisPath + ":sessions"

// PostsRedisPath is the sub-path for blog entry data
var PostsRedisPath = ProjectRedisPath + ":posts"

// VisitorsRedisPath is the sub-path for the visitor/hit counter
var VisitorsRedisPath = ProjectRedisPath + ":visitors"
