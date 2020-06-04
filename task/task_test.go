package task

import (
	"fmt"
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	Run("hello", func(task *Task) error {
		fmt.Println("hello async world")
		return nil
	})
	time.Sleep(1 * time.Second)
	fmt.Println("hello sync world")
}
