package rabbit

//import (
//	"github.com/streadway/amqp"
//	"testing"
//	"time"
//)
//
//var (
//	s            Session
//	exchangeName = "session-exchange"
//	exchangeKind = "topic"
//	queueName    = "session-queue"
//	routingKey   = "#"
//	message      = []byte("hello")
//)
//
//const URL = "amqp://guest:guest@localhost:5672/"
//
//func init() {
//	url := "amqp://guest:guest@localhost:5672/"
//	s = New(url)
//
//	for {
//		if r := IsClosed(); r == false {
//			break
//		}
//	}
//}
//
//func BenchmarkSession_Publish(b *testing.B) {
//	b.Run("benchmark publish", func(b *testing.B) {
//		b.ResetTimer()
//		for i := 0; i < b.N; i++ {
//			if err := Publish(exchangeName, routingKey, message); err != nil {
//				b.Errorf("publish failed, %s\n", err)
//			}
//		}
//	})
//
//	b.Run("benchmark parallel 50 publish", func(b *testing.B) {
//		b.ResetTimer()
//		b.SetParallelism(100)
//		b.RunParallel(func(pb *testing.PB) {
//			for pb.Next() {
//				if err := Publish(exchangeName, exchangeKind, message); err != nil {
//					b.Errorf("publish failed, %s\n", err)
//				}
//			}
//		})
//	})
//}
//
//var consumeMsgBody string
//
//func TestSession_Consume(t *testing.T) {
//	t.Run("normal consume", func(t *testing.T) {
//		go Consume(queueName, consumeHandle)
//
//		time.Sleep(2 * time.Second)
//		if consumeMsgBody != "hello" {
//			t.Errorf("expect `hello`, wot `%s`\n", consumeMsgBody)
//		}
//	})
//}
//
//func consumeHandle(d amqp.Delivery) {
//	defer d.Ack(false)
//	consumeMsgBody = string(d.Body)
//}
