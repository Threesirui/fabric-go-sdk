package main

import (
	"fmt"
	"os"
	// "time"
	"net"
	"fabric-go-sdk/sdkInit"
	"encoding/json"
)

const (
	cc_name    = "simplecc"
	cc_version = "1.0.0"
)

var Path []string


type TransInfo struct {
	Src_ip string  `json:"Src_ip"`
	Dst_ip string  `json:"Dst_ip"`
}

var App sdkInit.Application

func handleConnection(conn net.Conn,tranmap interface{})  {
	fmt.Println("进入了连接")
	defer conn.Close()
	flag := 0
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
		return
	}
	message := buffer[:n]
	var transinfo TransInfo
	err = json.Unmarshal([]byte(message),&transinfo)
	tranmap1 := tranmap.(map[TransInfo]int)
	value ,ok := tranmap1[transinfo]
	for k,v := range tranmap1{
		fmt.Println("map信息",k.Src_ip,k.Dst_ip,v)
	}
	if ok {
		flag=value
	}else{
		tranmap1[transinfo]=0
	}
	if	err!=nil{
		println("Errorjj",err)
		// return
	}

	// message := buffer[:n]
	// transinfo := &TransInfo{}
	// err = json.Unmarshal(message,transinfo)
	// if err != nil{
	// 	fmt.Println("transinfo unmarshal failed")
	// 	return
	// }
	// fmt.Println("Received message from Python:", message)

	// You can process the received message and prepare the response here.
	// response := "[10.10.4.1, 10.10.5.1]"

	// _, err = conn.Write([]byte(response))
	// if err != nil {
	// 	fmt.Println("Error writing to connection:", err)
	// 	return
	// }
	// message := string(buffer[:n])
	// transinfo := &TransInfo{}
	// transinfo.Src_ip= message[2:10]
	// transinfo.Dst_ip= message [11:19]
	// fmt.Println("Received message from Python:", message)
	fmt.Println("Python_transinfo.SRC",transinfo.Src_ip)
	fmt.Println("Python_transinfo.Dst",transinfo.Dst_ip)
	var src string
	var dst string
	if transinfo.Src_ip != transinfo.Dst_ip{
		fmt.Println("message:",transinfo.Src_ip[0:5])
		
		if transinfo.Src_ip[0:7] == "10.10.1" {
			src = "chain1"	
		}else if transinfo.Src_ip[0:7] == "10.10.2"{
			src = "chain2"
		}else if transinfo.Src_ip[0:7] == "10.10.3"{
			src = "chain3"
		}else{
			fmt.Printf("错误的源地址")
			
		}
		fmt.Println("transinfo:",transinfo.Dst_ip[0:5])
		if transinfo.Dst_ip[0:7] == "10.10.1" {
			dst = "chain1"	
		}else if transinfo.Dst_ip[0:7] == "10.10.2"{
			dst = "chain2"
		}else if transinfo.Dst_ip[0:7] == "10.10.3"{
			dst = "chain3"
		}else{
			fmt.Printf("错误的目的地址")
			
		}
		if flag ==0{
			a := []string{"route",src,dst}
			response, err := App.Route(a)
			if err!=nil{
				fmt.Println("计算路由错误")
			}
			tranmap1[transinfo] = 1
			for  i := 2; i<len(response)-5 ; i=i+9{
				if response[i:i+6] == "chain1"{
					Path = append(Path,"s1")
				}else if response[i:i+6] =="chain2" {
					Path = append(Path,"s2")
				}else if response[i:i+6] =="chain3" {
					Path = append(Path,"s3")
				}else {
					fmt.Println("出错了")
				}
			}
			// Path = append(Path,"s1")
			// Path = append(Path,"s2")
			fmt.Println("智能合约计算-----path-------",Path)
			path,err := json.Marshal(Path)
			if err != nil {
				fmt.Println("path marshal fail")
			}
			handleresponse(conn,path)
			Path = nil
		}else {
			if transinfo.Src_ip[0:7] == "10.10.1" {
				src = "s1"	
			}else if transinfo.Src_ip[0:7] == "10.10.2"{
				src = "s2"
			}else if transinfo.Src_ip[0:7] == "10.10.3"{
				src = "s3"
			}else{
				fmt.Printf("错误的源地址")
				
			}
			fmt.Println("message:",transinfo.Dst_ip[0:5])
			if transinfo.Dst_ip[0:7] == "10.10.1" {
				dst = "s1"	
			}else if transinfo.Dst_ip[0:7] == "10.10.2"{
				dst = "s2"
			}else if transinfo.Dst_ip[0:7] == "10.10.3"{
				dst = "s3"
			}else{
				fmt.Printf("错误的目的地址")
				
			}
			Path = append(Path,src)

			Path = append(Path,dst)
			fmt.Println("无智能合约-----path-------",Path)
			path,err := json.Marshal(Path)
			if err != nil {
				fmt.Println("path marshal fail")
			}
			handleresponse(conn,path)
			Path = nil
		}
	}


}
	


func handleresponse(conn net.Conn, response []byte){
	_, err := conn.Write([]byte(response))
	fmt.Println("发消息了")
	if err != nil {
		fmt.Println("Error writing to connection:", err)
		return
	}
}

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
    
	count := 1;
	if count == 1{


		// create channel and join
		if err := sdkInit.CreateAndJoinChannel(&info); err != nil {
			fmt.Println(">> Create channel and join error:", err)
			os.Exit(-1)
		}

		// create chaincode lifecycle
		if err := sdkInit.CreateCCLifecycle(&info, 1, false, sdk); err != nil {
			fmt.Println(">> create chaincode lifecycle error: ", err)
			os.Exit(-1)
		}
		count++
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
	a := []string{"route", "chain1", "chain2"}
	ret, err := App.Route(a)
	if err != nil {
		fmt.Println("调用route失败",err)

	}
	
	fmt.Println("<--- 添加信息　--->：", ret)
	// var res []string
	// err = json.Unmarshal([]byte(ret),res)
	// if err != nil {
	// 	fmt.Println("res unmarshal fail")
	// }
	// fmt.Println("<--- res　--->：", res)

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
	//tong xing mo kuai

	TranMap := make(map[TransInfo]int,10)

	host := "localhost"
	port := "7736"
	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		fmt.Println("Error creating listener:", err)
		return
	}
	// defer listener.Close()

	fmt.Println("Go server listening on " + host + ":" + port)

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error accepting connection:", err)
			return
		}
		fmt.Println("重新连接")
		go handleConnection(conn,TranMap)
		fmt.Println("退出了连接")

		// var src string
		// var dst string
		// if message.Src_ip != message.Dst_ip{
		// 	fmt.Println("message:",message.Src_ip[0:5])
			
			// if message.Src_ip[0:7] == "10.10.1" {
			// 	src = "chain1"	
			// }else if message.Src_ip[0:7] == "10.10.2"{
			// 	src = "chain2"
			// }else if message.Src_ip[0:7] == "10.10.3"{
			// 	src = "chain3"
			// }else{
			// 	fmt.Printf("错误的源地址")
				
			// }
			// fmt.Println("message:",message.Dst_ip[0:5])
			// if message.Dst_ip[0:7] == "10.10.1" {
			// 	dst = "chain1"	
			// }else if message.Dst_ip[0:7] == "10.10.2"{
			// 	dst = "chain2"
			// }else if message.Dst_ip[0:7] == "10.10.3"{
			// 	dst = "chain3"
			// }else{
			// 	fmt.Printf("错误的目的地址")
				
			// }

			// a := []string{"route",src,dst}
			// response, err := App.Route(a)
			// if err!=nil{
			// 	fmt.Println("计算路由错误")
			// }
			// for  i := 2; i<len(response)-5 ; i=i+9{
			// 	if response[i:i+6] == "chain1"{
			// 		Path = append(Path,"s1")
			// 	}else if response[i:i+6] =="chain2" {
			// 		Path = append(Path,"s2")
			// 	}else if response[i:i+6] =="chain3" {
			// 		Path = append(Path,"s3")
			// 	}else {
			// 		fmt.Println("出错了")
			// 	}
			// }
		// 	Path = append(Path,"s1")
		// 	Path = append(Path,"s2")
		// 	fmt.Println("-----path-------",Path)
		// 	path,err := json.Marshal(Path)
		// 	if err != nil {
		// 		fmt.Println("path marshal fail")
		// 	}
		// 	handleresponse(conn,path)
		// 	Path = nil

		// }
		// conn.Close()
		
		
		// a := []string{"route", "chain1","chain2"}
		// _, err = App.Route(a)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// b := "[10.10.4.1, 10.10.5.1]"
		// fmt.Println("<--- 计算得到的最短路径　--->：", b)
		// time.Sleep(time.Second * 10)
	

		
	}


	// a := []string{"get", "graph"}
	// response, err := App.Get(a)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("<--- 查询信息　--->：", response)

	// a = []string{"get", "ID3"}
	// response, err = App.Get(a)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("<--- 查询信息　--->：", response)

	// time.Sleep(time.Second * 10)

}
