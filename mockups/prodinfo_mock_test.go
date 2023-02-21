package mock_mtls_proto

import (
	"context"
	"fmt"
	"testing"
	"time"

	pb "github.com/blablatov/stream-mtls-grpc/mtls-proto"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
)

// rpcMsg implements the gomock.Matcher interface
type rpcMsg struct {
	msg proto.Message
}

func (r *rpcMsg) Matches(msg interface{}) bool {
	m, ok := msg.(proto.Message)
	if !ok {
		return false
	}
	return proto.Equal(m, r.msg)
}

func (r *rpcMsg) String() string {
	return fmt.Sprintf("is %s", r.msg)
}

func TestAddProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	moclProdInfoClient := NewMockProductInfoClient(ctrl)

	name := "Sumsung S9999"
	description := "Samsung Galaxy S9999 is the latest smart phone, launched in February 2039"
	price := float32(7777.0)
	req := &pb.Product{Name: name, Description: description, Price: price}

	moclProdInfoClient.EXPECT().AddProduct(gomock.Any(), &rpcMsg{msg: req}).
		Return(&pb.ProductID{Value: "Product:" + name}, nil)
	testAddProduct(t, moclProdInfoClient)
}

func testAddProduct(t *testing.T, client pb.ProductInfoClient) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	name := "Sumsung S9999"
	description := "Samsung Galaxy S9999 is the latest smart phone, launched in February 2039"
	price := float32(7777.0)

	r, err := client.AddProduct(ctx, &pb.Product{Name: name, Description: description, Price: price})

	if err != nil || r.GetValue() == "Sumsung S9999" {
		t.Errorf("mocking failed")
	}
	t.Log("Reply : ", r.GetValue())
}
