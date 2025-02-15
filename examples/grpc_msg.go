package examples

import (
	"github.com/yedf/dtm/dtmcli"
	dtmgrpc "github.com/yedf/dtm/dtmgrpc"
)

func init() {
	addSample("grpc_msg", func() string {
		req := dtmcli.MustMarshal(&TransReq{Amount: 30})
		gid := dtmgrpc.MustGenGid(DtmGrpcServer)
		msg := dtmgrpc.NewMsgGrpc(DtmGrpcServer, gid).
			Add(BusiGrpc+"/examples.Busi/TransOut", req).
			Add(BusiGrpc+"/examples.Busi/TransIn", req)
		err := msg.Submit()
		dtmcli.FatalIfError(err)
		return msg.Gid
	})
}
