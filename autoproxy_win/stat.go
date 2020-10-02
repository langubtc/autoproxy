package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var defaultFont = Font{
	PointSize: 16,
	Bold: true,
}

func StatWidget() []Widget {
	return []Widget{
		Label{
			Text: "运行状态",
			Font: defaultFont,
		},
		Label{
			Text: "正常",
			Font: defaultFont,
			TextColor: walk.RGB(250,100,100),
		},
		Label{
			Text: "请求次数",
			Font: defaultFont,
		},
		Label {
			Text: "10294",
			Font: defaultFont,
		},
		Label{
			Text: "实时流量",
			Font: defaultFont,
		},
		Label {
			Text: "0kb/s",
			Font: defaultFont,
		},
		Label{
			Text: "总体流量",
			Font: defaultFont,
		},
		Label {
			Text: "1GB",
			Font: defaultFont,
		},
	}
}