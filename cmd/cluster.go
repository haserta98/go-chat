package cmd

import (
	"log"
	"time"

	"github.com/haserta98/go-rest/internal"
)

type Cluster struct {
	redis  *internal.RedisClient
	nodeID string
}

func NewCluster(redis *internal.RedisClient, nodeID string) *Cluster {
	return &Cluster{
		redis:  redis,
		nodeID: nodeID,
	}
}

func (c *Cluster) SendHeartbeat() {
	ticker := time.NewTicker(5 * time.Second)

	go func() {
		for range ticker.C {
			activeKey := "active_nodes:" + c.nodeID
			if err := c.redis.Set(activeKey, "alive", 10*time.Second); err != nil {
				log.Printf("Error setting heartbeat in Redis: %v", err)
			}
		}
	}()
}

func (c *Cluster) IsTargetNodeAlive(nodeID string) (bool, error) {
	activeKey := "active_nodes:" + nodeID
	val, err := c.redis.Get(activeKey)
	if err != nil {
		return false, err
	}
	return val == "alive", nil
}
