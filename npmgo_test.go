package npmgo

import (
//	"encoding/json"
	"fmt"
	"net/url"
//	"strings"
	"time"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/streadway/amqp"
)

func FindTokenById(sl []Objects, id string) (t Objects) {
	for _, i := range sl {
		if id == i.Id {
			t = i
			break
		}
	}
	return t
}

func openConnection(vhost string) *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/" + url.QueryEscape(vhost))
	Ω(err).Should(BeNil())

	if err != nil {
		panic("failed to connect")
	}

	return conn
}

func ensureNonZeroMessageRate(ch *amqp.Channel) {
	for i := 0; i < 2000; i++ {
		q, _ := ch.QueueDeclare(
			"",    // name
			false, // durable
			false, // auto delete
			true,  // exclusive
			false,
			nil)
		ch.Publish("", q.Name, false, false, amqp.Publishing{Body: []byte("")})
	}
}

// Wait for the list of connections to reach the expected length
func listConnectionsUntil(c *Client, i int) {
	xs, _ := c.ListConnections()
	// Avoid infinity loops by breaking it after 30s
	breakLoop := 0
	for i != len(xs) {
		if breakLoop == 300 {
			fmt.Printf("Stopping listConnectionsUntil loop: expected %v obtained %v", i, len(xs))
			break
		}
		breakLoop += 1
		// Wait between calls
		time.Sleep(100 * time.Millisecond)
		xs, _ = c.ListConnections()
	}
}

func awaitEventPropagation() {
	time.Sleep(1150 * time.Millisecond)
}

type portTestStruct struct {
	Port Port `json:"port"`
}

var _ = Describe("Npmgo with username/password", func() {
	var (
		rmqc *Client
	)

	BeforeEach(func() {
		username := os.Getenv("NPMJS_USERNAME")
		password := os.Getenv("NPMJS_PASSWORD")
		rmqc, _ = NewClient("http://registry.npmjs.org", username, password)
	})

	Context("GET /whoami", func() {
		It("returns decoded response", func() {
			res, err := rmqc.Whoami()

			Ω(err).Should(BeNil())

			Ω(res.Username).ShouldNot(BeNil())
			Ω(res.Username).Should(Equal("chrismatteson"))
		})
	})

	Context("GET /npm/v1/tokens", func() {
		It("returns decoded response", func() {
			token := os.Getenv("NPMJS_EXISTINGTOKENID")
			xs, err := rmqc.ListTokens()
			Ω(err).Should(BeNil())

			t := FindTokenById(xs, token)
			Ω(t.Id).Should(BeEquivalentTo(token))
			Ω(t.Token).ShouldNot(BeNil())
			Ω(t.Created).ShouldNot(BeNil())
		})
	})

	Context("GET /tokens for {id} when id exists", func() {
		token := os.Getenv("NPMJS_EXISTINGTOKENID")
		It("returns decoded response", func() {
			t, err := rmqc.GetToken(token)
			Ω(err).Should(BeNil())

			Ω(t.Id).Should(BeEquivalentTo(token))
			Ω(t.Token).ShouldNot(BeNil())
			Ω(t.Created).ShouldNot(BeNil())
		})
	})

	Context("POST /tokens", func() {
		It("creates a token", func() {
                	password := os.Getenv("NPMJS_PASSWORD")
			info := TokenSettings{Password: password, Readonly: false}
			//var resp Objects
			resp, err := rmqc.CreateToken(info)
			Ω(err).Should(BeNil())
			//Ω(resp.Status).Should(HavePrefix("20"))
			
			// give internal events a moment to be
			// handled
			awaitEventPropagation()

			t, err := rmqc.GetToken(resp.Id)
			Ω(err).Should(BeNil())
			Ω(t.Id).ShouldNot(BeNil())

			rmqc.DeleteToken(resp.Id)
		})
	})
})

var _ = Describe("Npmgo with token", func() {
	var (
		rmqc *Client
	)

	BeforeEach(func() {
		token := os.Getenv("NPMJS_EXISTINGTOKEN")
		rmqc, _ = NewTokenClient("http://registry.npmjs.org", token)
	})

	Context("GET /whoami", func() {
		It("returns decoded response", func() {
			res, err := rmqc.Whoami()

			Ω(err).Should(BeNil())

			Ω(res.Username).ShouldNot(BeNil())
			Ω(res.Username).Should(Equal("chrismatteson"))
		})
	})

	Context("GET /npm/v1/tokens", func() {
		It("returns decoded response", func() {
			token := os.Getenv("NPMJS_EXISTINGTOKENID")
			xs, err := rmqc.ListTokens()
			Ω(err).Should(BeNil())

			t := FindTokenById(xs, token)
			Ω(t.Id).Should(BeEquivalentTo(token))
			Ω(t.Token).ShouldNot(BeNil())
			Ω(t.Created).ShouldNot(BeNil())
		})
	})

	Context("GET /tokens for {id} when id exists", func() {
		token := os.Getenv("NPMJS_EXISTINGTOKENID")
		It("returns decoded response", func() {
			t, err := rmqc.GetToken(token)
			Ω(err).Should(BeNil())

			Ω(t.Id).Should(BeEquivalentTo(token))
			Ω(t.Token).ShouldNot(BeNil())
			Ω(t.Created).ShouldNot(BeNil())
		})
	})

	Context("POST /tokens", func() {
		It("creates a token", func() {
                	password := os.Getenv("NPMJS_PASSWORD")
			info := TokenSettings{Password: password, Readonly: false}
			//var resp Objects
			resp, err := rmqc.CreateToken(info)
			Ω(err).Should(BeNil())
			//Ω(resp.Status).Should(HavePrefix("20"))
			
			// give internal events a moment to be
			// handled
			awaitEventPropagation()

			t, err := rmqc.GetToken(resp.Id)
			Ω(err).Should(BeNil())
			Ω(t.Id).ShouldNot(BeNil())

			rmqc.DeleteToken(resp.Id)
		})
	})
})
