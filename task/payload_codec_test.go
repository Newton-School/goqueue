package task

import "testing"

func TestPayloadCodecInterfaceAcceptsTestCodec(t *testing.T) {
	var _ PayloadCodec = testPayloadCodec{}
}

type testPayloadCodec struct{}

func (testPayloadCodec) EncodePayload(TaskPayload) ([]byte, error) {
	return []byte("{}"), nil
}

func (testPayloadCodec) DecodePayload([]byte) (TaskPayload, error) {
	return NewTaskPayload(nil, nil), nil
}
