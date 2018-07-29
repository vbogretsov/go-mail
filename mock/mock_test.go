package mock_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/labstack/gommon/random"

	"github.com/stretchr/testify/require"

	"github.com/vbogretsov/go-mail"
	"github.com/vbogretsov/go-mail/mock"
)

func TestSendConcurrent(t *testing.T) {
	sender := mock.New()

	to := []mail.Address{{Email: "to1@mail.com"}, {Email: "to2@mail.com"}}
	cc := []mail.Address{{Email: "cc1@mail.com"}, {Email: "cc2@mail.com"}}

	args := map[string]interface{}{"test": "test"}

	exp := mail.Request{
		TemplateLang: "en",
		TemplateName: "test",
		TemplateArgs: args,
		To:           to,
		Cc:           cc,
	}

	err := sender.Send(exp)

	checkInbox := func(addrs []mail.Address) {
		for _, addr := range addrs {
			act, ok := sender.ReadMail(addr.Email)
			require.True(t, ok)
			require.Equal(t, exp, act)
		}
	}

	require.Nil(t, err, "send error: %v", err)
	checkInbox(to)
	checkInbox(cc)
}

func TestConcurrentAccess(t *testing.T) {
	const n = 10000

	wg := sync.WaitGroup{}
	wg.Add(n)

	sender := mock.New()

	for i := 0; i < n; i++ {
		go func() {
			t.Run("Concurrent", func(t *testing.T) {
				email := fmt.Sprintf(
					"%s@mail.com",
					random.String(10, random.Alphabetic))
				exp := mail.Request{
					TemplateLang: "en",
					TemplateName: "test",
					TemplateArgs: map[string]interface{}{"test": "test"},
					To:           []mail.Address{{Email: email}},
				}

				sender.Send(exp)

				act, ok := sender.ReadMail(email)
				require.True(t, ok)

				require.Equal(t, exp, act)
				wg.Done()
			})
		}()
	}

	wg.Wait()
}
