package kafka

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestConsumer_ReadMessage(t *testing.T) {
	ctx := context.Background()
	// 确保本地有服务, topic with 15 partition
	t.Log("start read message")
	topic := "local_test_topic"
	brokers := []string{"localhost:32400", "localhost:32401", "localhost:32402"}
	c, err := NewConsumer(brokers, topic)
	require.NoError(t, err)

	p, err := NewProducer(brokers, topic)
	require.NoError(t, err)
	for i := 0; i < 1000; i++ {
		require.NoError(t, p.SendMessage(ctx, strconv.Itoa(i)))
	}
	t.Log("send message completed")
	//
	wg := sync.WaitGroup{}
	m := sync.Map{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for {
				msg := c.ReadMessage(time.Second * 10)
				if msg == "" {
					break
				}
				_, ok := m.Load(msg)
				// notice: 这里是判断 有没有出现消息重复消费的问题，但是注意确保topic一开始没有数据
				assert.False(t, ok)
				m.Store(msg, msg)
				fmt.Println("read msg: ", msg)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
