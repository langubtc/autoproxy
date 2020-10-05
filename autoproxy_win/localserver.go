package main

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/easymesh/autoproxy/autoproxy_win/engin"
	"net/http"
)

var access engin.Access

func AuthSwitch(auth *engin.AuthInfo) bool {
	if auth == nil {
		return false
	}
	return AuthCheck(auth.User, auth.Token)
}

var LocalForward engin.Forward

func ForwardFunc(address string, r *http.Request) engin.Forward {
	return LocalForward
}

func ServerStart() error {
	var err error

	if access != nil {
		logs.Error("server has beed start")
		return fmt.Errorf("server has been start")
	}

	address := fmt.Sprintf("%s:%d",
		IfaceOptions()[LocalIfaceOptionsIdx()],
		PortOptionGet())

	logs.Info("server start %s", address)

	access, err = engin.NewHttpsAccess(address, 60, AuthSwitchGet())
	if err != nil {
		return err
	}

	if AuthSwitchGet() {
		access.AuthHandlerSet(AuthSwitch)
	}
	access.ForwardHandlerSet(ForwardFunc)

	LocalForward, _ = engin.NewDefault(60)

	logs.Info("server start %s success", address)
	return nil
}

func ServerShutdown() error {
	if access == nil {
		return fmt.Errorf("server has been stop")
	}
	err := access.Shutdown()
	if err != nil {
		logs.Error("shutdown fail, %s", err.Error())
		return err
	}
	access = nil

	LocalForward.Close()
	LocalForward = nil
	return nil
}
