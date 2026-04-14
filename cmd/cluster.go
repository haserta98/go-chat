package cmd

import (
	"log"
	"os"
	"time"

	"github.com/haserta98/go-rest/internal"
)

type Cluster struct {
	redis *internal.RedisClient
}

func NewCluster(redis *internal.RedisClient) *Cluster {
	return &Cluster{
		redis: redis,
	}
}

func (c *Cluster) SendHeartbeat() {
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for range ticker.C {

			nodeID := os.Getenv("NODE_ID")
			if nodeID == "" {
				log.Fatal("NODE_ID environment variable is not set")
			}
			activeKey := "active_nodes" + nodeID
			if err := c.redis.Set(activeKey, "alive", 10*time.Second); err != nil {
				log.Printf("Error setting heartbeat in Redis: %v", err)
			}
		}
	}()
}

func (c *Cluster) IsTargetNodeAlive(nodeID string) (bool, error) {
	activeKey := "active_nodes" + nodeID
	val, err := c.redis.Get(activeKey)
	if err != nil {
		return false, err
	}
	return val == "alive", nil
}
