package mock_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/vbogretsov/go-mail"
	"github.com/vbogretsov/go-mail/mock"
)

func TestSendConcurrent(t *testing.T) {
	sender := mock.New()

	to := []mailcd.Address{{Email: "to1@mail.com"}, {Email: "to2@mail.com"}}
	cc := []mailcd.Address{{Email: "cc1@mail.com"}, {Email: "cc2@mail.com"}}

	args := map[string]interface{}{"test": "test"}

	exp := mailcd.Request{
		TemplateLang: "en",
		TemplateName: "test",
		TemplateArgs: args,
		To:           to,
		Cc:           cc,
	}

	err := sender.Send(exp)

	checkInbox := func(addrs []mailcd.Address) {
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
