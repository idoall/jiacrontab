package kproc

import (
	"context"
	"fmt"
	"os/exec"
)

func CommandContext(ctx context.Context, name string, arg ...string) *KCmd {
	cmd := exec.CommandContext(ctx, name, arg...)
	return &KCmd{
		Cmd:  cmd,
		ctx:  ctx,
		done: make(chan struct{}),
	}
}

func (k *KCmd) KillAll() {
	select {
	case k.done <- struct{}{}:
	default:
	}
	if k.Process == nil {
		return
	}
	c := exec.Command("taskkill", "/t", "/f", "/pid", fmt.Sprint(k.Process.Pid))
	c.Stdout = k.Cmd.Stdout
	c.Stderr = k.Cmd.Stderr
}

func (k *KCmd) Wait() error {
	defer k.KillAll()
	go func() {
		select {
		case <-k.ctx.Done():
			k.KillAll()
		case <-k.done:
		}
	}()
	return k.Cmd.Wait()
}
