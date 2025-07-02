package core_test

import (
	"testing"

	"github.com/shellvon/go-sender/core"
)

type fakeSelectable struct {
	name    string
	weight  int
	enabled bool
}

func (f *fakeSelectable) GetName() string { return f.name }
func (f *fakeSelectable) GetWeight() int  { return f.weight }
func (f *fakeSelectable) IsEnabled() bool { return f.enabled }
func (f *fakeSelectable) GetType() string { return "fake" }

func TestRoundRobinStrategy(t *testing.T) {
	s := core.NewRoundRobinStrategy()
	items := []core.Selectable{&fakeSelectable{"a", 1, true}, &fakeSelectable{"b", 1, true}}
	if s.Name() != core.StrategyRoundRobin {
		t.Error("Name should be round_robin")
	}
	first := s.Select(items)
	second := s.Select(items)
	if first == nil || second == nil || first == second {
		t.Error("RoundRobin should rotate")
	}
	if s.Select([]core.Selectable{}) != nil {
		t.Error("empty select should return nil")
	}
}

func TestRandomStrategy(t *testing.T) {
	s := core.NewRandomStrategy()
	items := []core.Selectable{&fakeSelectable{"a", 1, true}, &fakeSelectable{"b", 1, true}}
	if s.Name() != core.StrategyRandom {
		t.Error("Name should be random")
	}
	_ = s.Select(items)
	if s.Select([]core.Selectable{}) != nil {
		t.Error("empty select should return nil")
	}
}

func TestWeightedStrategy(t *testing.T) {
	s := core.NewWeightedStrategy()
	items := []core.Selectable{&fakeSelectable{"a", 2, true}, &fakeSelectable{"b", 1, true}}
	if s.Name() != core.StrategyWeighted {
		t.Error("Name should be weighted")
	}
	_ = s.Select(items)
	if s.Select([]core.Selectable{}) != nil {
		t.Error("empty select should return nil")
	}
}

func TestStrategyRegistry(t *testing.T) {
	r := core.NewStrategyRegistry()
	r.Register("custom", core.NewRoundRobinStrategy())
	if _, ok := r.Get("custom"); !ok {
		t.Error("custom strategy should be registered")
	}
	if r.GetDefault() == nil {
		t.Error("GetDefault should not be nil")
	}
}
