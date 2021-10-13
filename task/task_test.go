package task

import (
	"log"
	"testing"

	"github.com/fghwett/heytap/config"
)

var conf *config.Conf

func init() {
	var err error
	conf, err = config.Init("../config.yml")
	if err != nil {
		log.Println("read config err: ", err)
		return
	}
}

func TestNew(t *testing.T) {
	task := New(conf.Config)
	err := task.signTask()
	if err != nil {
		t.Error(err)
	}
}
