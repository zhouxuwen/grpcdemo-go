package interceptor

import "context"

var AuthCredential = authCredential{}

type authCredential struct{}

func (*authCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"app_id": "123456", "secret_key": "secret_key"}, nil
}

func (*authCredential) RequireTransportSecurity() bool {
	return true
}
