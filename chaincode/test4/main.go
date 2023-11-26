package main

import (
	"fmt"
	"os"
	"time"

	"fabric-go-sdk/sdkInit"
)

const (
	cc_name    = "simplecc"
	cc_version = "1.0.0"
)

var App sdkInit.Application

func main() {
	// init orgs information

	orgs := []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "/home/chain2/go/src/github.com/hyperledger/fabric-samples/fabric-docker-multiple/channel-artifacts/Org1MSPanchors.tx",
		},
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org2",
			OrgMspId:      "Org2MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "/home/chain2/go/src/github.com/hyperledger/fabric-samples/fabric-docker-multiple/channel-artifacts/Org2MSPanchors.tx",
		},
	}

	// init sdk env info
	info := sdkInit.SdkEnvInfo{
		ChannelID: "businesschannel",
		// ChannelID:        "fabric-channel",
		ChannelConfig: "/home/chain2/go/src/github.com/hyperledger/fabric-samples/fabric-docker-multiple/channel-artifacts/businesschannel.tx",
		// ChannelConfig:    "/home/chain2/go/src/github.com/hyperledger/fabric-samples/fabric-docker-multiple/channel-artifacts/orderer.genesis.block",
		Orgs:             orgs,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer0.example.com",
		ChaincodeID:      cc_name,
		ChaincodePath:    "/home/go/src/github.com/fabric-go-sdk/chaincode/test3/",
		ChaincodeVersion: cc_version,
	}
	// sdk setup
	sdk, err := sdkInit.Setup("config.yaml", &info)
	if err != nil {
		fmt.Println(">> SDK setup error:", err)
		os.Exit(-1)
	}

	// create channel and join
	if err := sdkInit.CreateAndJoinChannel(&info); err != nil {
		fmt.Println(">> Create channel and join error:", err)
		os.Exit(-1)
	}

	// create chaincode lifecycle
	if err := sdkInit.CreateCCLifecycle(&info, 1, false, sdk); err != nil {
		fmt.Println(">> create chaincode lifecycle error: %v", err)
		os.Exit(-1)
	}

	// invoke chaincode set status
	fmt.Println(">> 通过链码外部服务设置链码状态......")

	if err := info.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[0], sdk); err != nil {

		fmt.Println("InitService successful")
		os.Exit(-1)
	}

	App = sdkInit.Application{
		SdkEnvInfo: &info,
	}
	fmt.Println(">> 设置链码状态完成")

	defer info.EvClient.Unregister(sdkInit.BlockListener(info.EvClient))
	defer info.EvClient.Unregister(sdkInit.ChainCodeEventListener(info.EvClient, info.ChaincodeID))
	// g := map[string]map[string]int{
	// 	"a": {"b": 20, "c": 80},
	// 	"b": {"a": 20, "c": 20},
	// 	"c": {"a": 80, "b": 20},
	// }
	// graph,_:=json.Marshal(g)
	// graphstr:=string(graph)
	// a := []string{"path", "a", "c",graphstr}
	// ret, err := App.Path(a)
	// if err != nil {
	// 	fmt.Println("调用path失败",err)

	// }
	// fmt.Println("<--- 添加信息　--->：", ret)

	// // a := []string{"set", "ID2", "456"}
	// // ret, err := App.Set(a)
	// // if err != nil {
	// // 	fmt.Println(err)
	// // }
	// // fmt.Println("<--- 添加信息　--->：", ret)
	// // // a = []string{"get", "ID2"}
	// // // response, err := App.Get(a)
	// // // if err != nil {
	// // // 	fmt.Println(err)
	// // // }
	// // // fmt.Println("<--- 查询信息　--->：", response)

	// // a = []string{"set", "ID3", "7899"}
	// // ret, err = App.Set(a)
	// // if err != nil {
	// // 	fmt.Println(err)
	// // }
	// // fmt.Println("<--- 添加信息　--->：", ret)

	// a = []string{"get", "b"}
	// response, err := App.Get(a)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("<--- 查询信息　--->：", response)

	a := []string{"get", "graph"}
	response, err := App.Get(a)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("<--- 查询信息　--->：", response)

	// a = []string{"get", "ID3"}
	// response, err = App.Get(a)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("<--- 查询信息　--->：", response)

	time.Sleep(time.Second * 10)

}
