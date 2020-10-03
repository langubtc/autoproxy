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
			Text: LangValue("runningstatus"),
			Font: defaultFont,
		},
		Label{
			Text: LangValue("normal"),
			Font: defaultFont,
			TextColor: walk.RGB(250,100,100),
		},
		Label{
			Text: LangValue("requestcount"),
			Font: defaultFont,
		},
		Label {
			Text: "10294",
			Font: defaultFont,
		},
		Label{
			Text: LangValue("realtimeflow"),
			Font: defaultFont,
		},
		Label {
			Text: "0kb/s",
			Font: defaultFont,
		},
		Label{
			Text: LangValue("totalflow"),
			Font: defaultFont,
		},
		Label {
			Text: "1GB",
			Font: defaultFont,
		},
	}
}