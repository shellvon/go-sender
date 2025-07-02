package core_test

import (
	"testing"

	"github.com/shellvon/go-sender/core"
)

func TestBaseConfig_GetStrategy(t *testing.T) {
	c := &core.BaseConfig{}
	if c.GetStrategy() != core.StrategyRoundRobin {
		t.Errorf("default strategy should be round_robin")
	}
	c.Strategy = core.StrategyRandom
	if c.GetStrategy() != core.StrategyRandom {
		t.Errorf("should get set strategy")
	}
}

func TestBaseConfig_IsDisabled(t *testing.T) {
	c := &core.BaseConfig{}
	if c.IsDisabled() {
		t.Error("default should not be disabled")
	}
	c.Disabled = true
	if !c.IsDisabled() {
		t.Error("should be disabled")
	}
}

func TestAccount_Methods(t *testing.T) {
	a := &core.Account{Name: "n", Weight: 0, Disabled: false, Type: "t"}
	if !a.IsEnabled() {
		t.Error("should be enabled")
	}
	a.Disabled = true
	if a.IsEnabled() {
		t.Error("should not be enabled")
	}
	if a.GetName() != "n" {
		t.Error("GetName failed")
	}
	if a.GetWeight() != 1 {
		t.Error("default weight should be 1")
	}
	a.Weight = 5
	if a.GetWeight() != 5 {
		t.Error("GetWeight failed")
	}
	if a.GetType() != "t" {
		t.Error("GetType failed")
	}
}
