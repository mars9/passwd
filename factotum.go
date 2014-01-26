package passwd

import (
	"bytes"
	"fmt"

	"code.google.com/p/goplan9/plan9"
	"code.google.com/p/goplan9/plan9/client"
)

type rpc struct {
	fid *client.Fid
}

func newRpc(name string) (*rpc, error) {
	fsys, err := client.MountService("factotum")
	if err != nil {
		return nil, err
	}
	fid, err := fsys.Open(name, plan9.ORDWR)
	if err != nil {
		return nil, err
	}
	return &rpc{fid: fid}, nil
}

func (c *rpc) Close() error { return c.fid.Close() }

type factotum struct{}

func (f *factotum) Get(service, username string) (string, error) {
	params := fmt.Sprintf("dom=%s proto=pass role=client", service)

	ctl, err := newRpc("rpc")
	if err != nil {
		return "", err
	}
	defer ctl.Close()

	_, err = ctl.fid.Write([]byte("start " + params))
	if err != nil {
		return "", err
	}
	buf := make([]byte, 4096)
	n, err := ctl.fid.Read(buf)
	if err != nil {
		return "", err
	}
	if !bytes.HasPrefix(buf, []byte("ok")) {
		return "", fmt.Errorf("start failed: %s", buf[:n])
	}

	_, err = ctl.fid.Write([]byte("read"))
	if err != nil {
		return "", err
	}
	n, err = ctl.fid.Read(buf)
	if err != nil {
		return "", err
	}
	if !bytes.HasPrefix(buf, []byte("ok")) {
		return "", fmt.Errorf("read failed: %s", buf[:n])
	}

	elems := bytes.Split(buf[:n], []byte(" "))
	if len(elems) != 3 {
		return "", fmt.Errorf("split response failed")
	}

	return string(elems[2]), nil
}

func (f *factotum) Set(service, username, password string) error {
	params := fmt.Sprintf("dom=%s proto=pass role=client", service)
	key := params + fmt.Sprintf(" user=%s !password=%s", username, password)

	ctl, err := newRpc("ctl")
	if err != nil {
		return err
	}
	defer ctl.Close()

	_, err = ctl.fid.Write([]byte("key " + key))
	return err
}

func (f *factotum) Delete(service, username string) error {
	params := fmt.Sprintf("dom=%s proto=pass role=client user=%s",
		service, username)

	ctl, err := newRpc("ctl")
	if err != nil {
		return err
	}
	defer ctl.Close()

	_, err = ctl.fid.Write([]byte("delkey " + params))
	return err
}
