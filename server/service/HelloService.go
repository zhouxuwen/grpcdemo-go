package service

import (
	"context"
	"grpcdemo-go/protobuf"
)

var HelloService = helloService{}

type helloService struct {
	protobuf.HelloServiceServer
}

func (h *helloService) Hello(ctx context.Context, requestMessage *protobuf.Massage) (responseMessage *protobuf.Massage, err error) {
	msg := "hello " + requestMessage.Massage + " !"
	responseMessage = &protobuf.Massage{Massage: msg}
	return
}
