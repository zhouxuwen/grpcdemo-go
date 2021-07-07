package service

import (
	"context"
	"grpcdemo-go/protobuf"
)

var HelloService = helloService{}

type helloService struct {
	protobuf.HelloServiceServer
}

func (h *helloService) Hello(c context.Context, requestMessage *protobuf.Massage) (responseMessage *protobuf.Massage, err error) {
	msg := "hello " + requestMessage.Request
	responseMessage = &protobuf.Massage{Request: msg}
	return
}
