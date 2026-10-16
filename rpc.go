package main

import (
	"fmt"
	"time"

	"github.com/hugolgst/rich-go/client"
)

const appID = "1560387470300946443"
func RPCLogin() bool {
	return client.Login(appID) == nil
}

const logoImageName string = "rosewatericon"
const programName string = "Rosewater"

type RPCQuery struct {
	Action	     string
	Where	     string
}

func RPCActivity(query RPCQuery) {

	start := time.Now()
	var err = client.SetActivity(client.Activity{
		Details: query.Action+programName,
		State: query.Where,
		LargeText: programName,
		LargeImage: logoImageName,
		Timestamps: &client.Timestamps{Start: &start},
	})

	if err != nil {
		fmt.Println(red+underline+"INFO"+reset+red+" : "+reset+yellow+"Could not access RPC"+reset)
		fmt.Println(err)
	}
}
