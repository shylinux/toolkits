package conf

import (
	"context"
)

func (conf *Conf) Context() context.Context {
	return conf.ctx
}

func (conf *Conf) Daemon(name string, cb func(context.Context)) {
	conf.add()
	go func() {
		defer conf.done()
		cb(context.WithValue(conf.ctx, "daemon", name))
	}()
}
func (conf *Conf) add() {
	for c := conf; c != nil; c = c.sup {
		c.wg.Add(1)
	}
}
func (conf *Conf) done() {
	for c := conf; c != nil; c = c.sup {
		c.wg.Done()
	}
}
func (conf *Conf) Wait()  { conf.wg.Wait() }
func (conf *Conf) Close() { conf.cancel() }
