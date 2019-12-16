package main

import (
	"fmt"

	"github.com/g-harel/coco/collectors"
	"github.com/g-harel/coco/collectors/github"
	"github.com/g-harel/coco/collectors/npm"
	"github.com/g-harel/coco/internal/exec"
	"github.com/g-harel/coco/internal/log"
)

func main() {
	c := []collectors.Collector{
		&github.Collector{},
		&npm.Collector{},
	}

	exec.ParallelN(len(c), func(n int) {
		c[n].Collect(func(err error) {
			log.Error("%v\n", err)
		})
	})

	for i := 0; i < len(c); i++ {
		fmt.Print(c[i].Format())
	}
}
