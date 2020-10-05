package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"time"
)

var defaultFont = Font{
	PointSize: 14,
	Bold: true,
}

var TotalFlowSize uint64
var TotalReqCnt   uint64
var RealTimeFlow  uint64

var LastUpdate     time.Time

func StatUpdate(requst uint64, flowsize uint64)  {
	if requst == 0 || flowsize == 0 {
		return
	}
	now := time.Now()

	TotalReqCnt += requst
	TotalFlowSize += flowsize
	RealTimeFlow = flowsize / uint64(now.Sub(LastUpdate).Seconds())

	LastUpdate = now
	requestCount.SetText(requestShow())
	realtimeflow.SetText(realTimeShow())
	totalflow.SetText(totalFlowShow())
	NotifyUpdateFlow(realTimeShow())
}

func StatInit() error {
	TotalFlowSize = DataLongValueGet("statflowsize")
	TotalReqCnt = DataLongValueGet("stattotalreq")

	go func() {
		for  {
			flowsize := TotalFlowSize
			reqcnt := TotalReqCnt

			time.Sleep(time.Minute)
			if flowsize != TotalFlowSize {
				DataLongValueSet("statflowsize", TotalFlowSize)
			}
			if reqcnt != TotalReqCnt {
				DataLongValueSet("stattotalreq", TotalReqCnt)
			}
		}
	}()

	return nil
}

func totalFlowShow() string {
	return ByteView(int64(TotalFlowSize))
}

func requestShow() string {
	return fmt.Sprintf("%d", TotalReqCnt)
}

func realTimeShow() string {
	return fmt.Sprintf("%s/s",
		ByteViewLite(int64(RealTimeFlow * 8)))
}

func StatRunningStatus(idx int)  {
	var image *walk.Icon
	switch idx {
	case 0:
		image = ICON_Network_Disable
	case 1:
		image = ICON_Network_LOW
	case 2:
		image = ICON_Network_MID
	case 3:
		image = ICON_Network_High
	case 4:
		image = ICON_Network_Full
	default:
		image = ICON_Network_Disable
	}
	runningStatus.SetImage(image)
	NotifyUpdateIcon(image)
}

var runningStatus *walk.ImageView
var requestCount  *walk.Label
var realtimeflow  *walk.Label
var totalflow     *walk.Label

func StatWidget() []Widget {
	return []Widget{
		Label{
			Text: LangValue("runningstatus"),
			Font: defaultFont,
		},
		Label {
			MinSize: Size{Width: 10},
		},
		ImageView{
			AssignTo: &runningStatus,
			Image:    ICON_Network_Disable,
			MaxSize:  Size{16, 16},
		},
		Label{
			Text: LangValue("requestcount"),
			Font: defaultFont,
		},
		Label {
			MinSize: Size{Width: 10},
		},
		Label {
			AssignTo: &requestCount,
			Text: requestShow(),
			Font: defaultFont,
			MinSize: Size{Width: 100},
		},
		Label{
			Text: LangValue("realtimeflow"),
			Font: defaultFont,
		},
		Label {
			MinSize: Size{Width: 10},
		},
		Label {
			AssignTo: &realtimeflow,
			Text: realTimeShow(),
			Font: defaultFont,
			MinSize: Size{Width: 100},
		},
		Label{
			Text: LangValue("totalflow"),
			Font: defaultFont,
		},
		Label {
			MinSize: Size{Width: 10},
		},
		Label {
			AssignTo: &totalflow,
			Text: totalFlowShow(),
			Font: defaultFont,
			MinSize: Size{Width: 100},
		},
	}
}