package jwt_test

import (
	"fmt"
	"testing"

	"github.com/douglaszuqueto/go-grpc-user/pkg/util/jwt"

	"github.com/stretchr/testify/require"
)

func TestJWT(t *testing.T) {
	t.Parallel()

	// Test cases
	// * secret empty
	// * newJwt not nil
	// * generate not error
	// * generate not empty
	// * token invalid signature

	secret := "secret"

	j := jwt.New(secret)
	require.NotNil(t, j)

	token, err := j.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	j = jwt.New("s")
	require.NotNil(t, j)

	err = j.Verify(token)
	require.Error(t, err)
}

func BenchmarkBinary(b *testing.B) {
	loop(b, func() {
		j := jwt.New("secret")
		err := j.Verify("eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE1ODgyMDcwODMsImV4cCI6MTYxOTc0MzA4MywiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoianJvY2tldEBleGFtcGxlLmNvbSJ9.TOWNAKZwSwgolXJMu2UdYZ69WkcrAkBJvpM1EZfVFiU")
		require.NoError(b, err)
	})
}

func loop(b *testing.B, cb func()) {
	b.Run(fmt.Sprintf("%d", 1), func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cb()
		}
	})
}
