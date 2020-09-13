package main

import (
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

func writeFull(conn net.Conn, buf []byte) error {
	totallen := len(buf)
	sendcnt := 0
	for {
		cnt, err := conn.Write(buf[sendcnt:])
		if err != nil {
			return err
		}
		if cnt+sendcnt >= totallen {
			return nil
		}
		sendcnt += cnt
	}
}

type WaitGroupTimeout struct {
	timeout time.Duration
	stat int64
	sync.WaitGroup
}

func NewWaitGroupTimeout(cnt int, timeout time.Duration, event func()) *WaitGroupTimeout {
	wt := new(WaitGroupTimeout)
	wt.timeout = timeout
	wt.Add(cnt)
	go func() {
		for {
			befor := wt.stat
			time.Sleep(wt.timeout)
			if wt.stat == befor {
				event()
				return
			}
		}
	}()
	return wt
}

func ConnectCopy(dst net.Conn, src net.Conn, wt *WaitGroupTimeout) {
	defer wt.Done()

	var cnt int
	var srcErr error
	var dstErr error

	body := make([]byte, 8192)
	for {
		cnt, srcErr = src.Read(body)
		if cnt > 0 {
			dstErr = writeFull(dst, body[:cnt])
			atomic.AddInt64(&wt.stat, int64(cnt))
		}

		if srcErr != nil || dstErr != nil {
			if srcErr != nil && srcErr != io.EOF {
				Warnf("connect %s copy error: %s",
					src.RemoteAddr(), srcErr.Error())
			}
			if dstErr != nil && dstErr != io.EOF {
				Warnf("connect %s copy error: %s",
					dst.RemoteAddr(), dstErr.Error())
			}
			return
		}
	}
}