package core_test

import (
	"testing"

	"github.com/shellvon/go-sender/core"
)

func TestProviderMeta_GetStrategy(t *testing.T) {
	c := &core.ProviderMeta{}
	if c.GetStrategy() != core.StrategyRoundRobin {
		t.Errorf("default strategy should be round_robin")
	}
	c.Strategy = core.StrategyRandom
	if c.GetStrategy() != core.StrategyRandom {
		t.Errorf("should get set strategy")
	}
}

func TestProviderMeta_IsDisabled(t *testing.T) {
	c := &core.ProviderMeta{}
	if c.IsDisabled() {
		t.Error("default should not be disabled")
	}
	c.Disabled = true
	if !c.IsDisabled() {
		t.Error("should be disabled")
	}
}

// Deprecated Account tests removed due to elimination of core.Account.
