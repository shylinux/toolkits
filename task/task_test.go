package task

import (
	"fmt"
	"testing"
	"time"

	"github.com/shylinux/toolkits/conf"
)

func TestTask(t *testing.T) {
	p := New(conf.New(nil))
	p.Run("hello", func(task *Task) error {
		fmt.Println("hello async world")
		return nil
	})

	time.Sleep(1 * time.Second)
	fmt.Println("hello sync world")
}
