package render

import "github.com/firmeve/firmeve/kernel/contract"

type Json struct {

}

func (j *Json) Name() string {
	panic("implement me")
}

func (j *Json) Render(protocol contract.Protocol, v interface{}) error {
	_, err := protocol.Write(v)
	return err
}

