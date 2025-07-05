package core_test

import (
	"context"
	"testing"

	"github.com/shellvon/go-sender/core"
)

// mockSelectable 用于测试的mock实现.
type mockSelectable struct {
	name    string
	weight  int
	enabled bool
	subType string
}

func (m *mockSelectable) GetName() string { return m.name }
func (m *mockSelectable) GetWeight() int  { return m.weight }
func (m *mockSelectable) IsEnabled() bool { return m.enabled }
func (m *mockSelectable) GetType() string { return m.subType }

func TestBaseConfig_Add(t *testing.T) {
	config := &core.BaseConfig[*mockSelectable]{}

	// 测试添加第一个item
	item1 := &mockSelectable{name: "test1", weight: 1, enabled: true}
	err := config.Add(item1)
	if err != nil {
		t.Errorf("Add() error = %v, want nil", err)
	}
	if len(config.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(config.Items))
	}

	// 测试添加重复名称的item
	item2 := &mockSelectable{name: "test1", weight: 2, enabled: true}
	err = config.Add(item2)
	if err == nil {
		t.Error("Expected error for duplicate name, got nil")
	}

	// 测试添加不同名称的item
	item3 := &mockSelectable{name: "test2", weight: 3, enabled: true}
	err = config.Add(item3)
	if err != nil {
		t.Errorf("Add() error = %v, want nil", err)
	}
	if len(config.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(config.Items))
	}
}

func TestBaseConfig_Delete(t *testing.T) {
	config := &core.BaseConfig[*mockSelectable]{
		Items: []*mockSelectable{
			{name: "test1", weight: 1, enabled: true},
			{name: "test2", weight: 2, enabled: true},
			{name: "test3", weight: 3, enabled: true},
		},
	}

	// 测试删除存在的item
	config.Delete("test2")
	if len(config.Items) != 2 {
		t.Errorf("Expected 2 items after delete, got %d", len(config.Items))
	}
	if config.Items[0].name != "test1" || config.Items[1].name != "test3" {
		t.Error("Delete() did not preserve order correctly")
	}

	// 测试删除不存在的item（应该无操作）
	config.Delete("nonexistent")
	if len(config.Items) != 2 {
		t.Errorf("Expected 2 items after deleting nonexistent, got %d", len(config.Items))
	}
}

func TestBaseConfig_Update(t *testing.T) {
	config := &core.BaseConfig[*mockSelectable]{
		Items: []*mockSelectable{
			{name: "test1", weight: 1, enabled: true},
			{name: "test2", weight: 2, enabled: true},
		},
	}

	// 测试更新存在的item
	updatedItem := &mockSelectable{name: "test1", weight: 10, enabled: false}
	err := config.Update(updatedItem)
	if err != nil {
		t.Errorf("Update() error = %v, want nil", err)
	}

	// 验证更新结果
	found := false
	for _, item := range config.Items {
		if item.name == "test1" && item.weight == 10 && !item.enabled {
			found = true
			break
		}
	}
	if !found {
		t.Error("Update() did not update item correctly")
	}

	// 测试更新不存在的item
	nonexistentItem := &mockSelectable{name: "nonexistent", weight: 1, enabled: true}
	err = config.Update(nonexistentItem)
	if err == nil {
		t.Error("Expected error for updating nonexistent item, got nil")
	}
}

func TestBaseConfig_Select(t *testing.T) {
	config := &core.BaseConfig[*mockSelectable]{
		Items: []*mockSelectable{
			{name: "test1", weight: 1, enabled: true},
			{name: "test2", weight: 2, enabled: true},
			{name: "test3", weight: 3, enabled: false},
		},
	}

	// 测试选择单个enabled item
	item, err := config.Select(context.Background(), nil)
	if err != nil {
		t.Errorf("Select() error = %v, want nil", err)
	}
	if item == nil {
		t.Error("Select() returned nil item")
	}

	// 测试选择单个disabled item
	config.Items = []*mockSelectable{
		{name: "test1", weight: 1, enabled: false},
	}
	_, err = config.Select(context.Background(), nil)
	if err == nil {
		t.Error("Expected error for selecting disabled item, got nil")
	}

	// 测试没有enabled items的情况
	config.Items = []*mockSelectable{
		{name: "test1", weight: 1, enabled: false},
		{name: "test2", weight: 2, enabled: false},
	}
	_, err = config.Select(context.Background(), nil)
	if err == nil {
		t.Error("Expected error for no enabled items, got nil")
	}

	// 测试使用filter函数
	config.Items = []*mockSelectable{
		{name: "test1", weight: 1, enabled: true, subType: "type1"},
		{name: "test2", weight: 2, enabled: true, subType: "type2"},
	}
	filter := func(item *mockSelectable) bool {
		return item.subType == "type1"
	}
	item, err = config.Select(context.Background(), filter)
	if err != nil {
		t.Errorf("Select() with filter error = %v, want nil", err)
	}
	if item.name != "test1" {
		t.Errorf("Expected item 'test1', got %s", item.name)
	}
}

func TestBaseConfig_Select_WithContext(t *testing.T) {
	config := &core.BaseConfig[*mockSelectable]{
		Items: []*mockSelectable{
			{name: "test1", weight: 1, enabled: true},
			{name: "test2", weight: 2, enabled: true},
		},
	}

	// 测试通过context指定item name
	ctx := core.WithCtxItemName(context.Background(), "test2")
	item, err := config.Select(ctx, nil)
	if err != nil {
		t.Errorf("Select() with context item name error = %v, want nil", err)
	}
	if item.name != "test2" {
		t.Errorf("Expected item 'test2', got %s", item.name)
	}

	// 测试通过context指定不存在的item name
	ctx = core.WithCtxItemName(context.Background(), "nonexistent")
	_, err = config.Select(ctx, nil)
	if err == nil {
		t.Error("Expected error for nonexistent item name, got nil")
	}
}

func TestBaseConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *core.BaseConfig[*mockSelectable]
		wantErr bool
	}{
		{
			name: "valid config",
			config: &core.BaseConfig[*mockSelectable]{
				Items: []*mockSelectable{
					{name: "test1", weight: 1, enabled: true},
					{name: "test2", weight: 2, enabled: true},
				},
			},
			wantErr: false,
		},
		{
			name: "disabled config",
			config: &core.BaseConfig[*mockSelectable]{
				ProviderMeta: core.ProviderMeta{Disabled: true},
				Items: []*mockSelectable{
					{name: "test1", weight: 1, enabled: true},
				},
			},
			wantErr: true,
		},
		{
			name: "no items",
			config: &core.BaseConfig[*mockSelectable]{
				Items: []*mockSelectable{},
			},
			wantErr: true,
		},
		{
			name: "empty item name",
			config: &core.BaseConfig[*mockSelectable]{
				Items: []*mockSelectable{
					{name: "", weight: 1, enabled: true},
				},
			},
			wantErr: true,
		},
		{
			name: "duplicate item names",
			config: &core.BaseConfig[*mockSelectable]{
				Items: []*mockSelectable{
					{name: "test1", weight: 1, enabled: true},
					{name: "test1", weight: 2, enabled: true},
				},
			},
			wantErr: true,
		},
		{
			name: "all items disabled",
			config: &core.BaseConfig[*mockSelectable]{
				Items: []*mockSelectable{
					{name: "test1", weight: 1, enabled: false},
					{name: "test2", weight: 2, enabled: false},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
