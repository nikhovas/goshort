package test

import (
	"github.com/spf13/viper"
	. "gopkg.in/check.v1"
	"goshort"
	"goshort/utils"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) SetUpSuite(_ *C) {
	utils.SetupViper("")
	viper.SetDefault("token", "demo")
	goshort.AppObject = goshort.App{}
	goshort.AppObject.Initialize()
}
