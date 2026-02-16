package nav

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

// mockScreen is a test implementation of Screen.
type mockScreen struct {
	name        string
	initCalled  bool
	updateCount int
}

func (m *mockScreen) Init() tea.Cmd {
	m.initCalled = true
	return nil
}

func (m *mockScreen) Update(msg tea.Msg) (Screen, tea.Cmd) {
	m.updateCount++
	return m, nil
}

func (m *mockScreen) View() tea.View {
	return tea.NewView(m.name)
}

// mockLifecycleScreen is a test implementation of LifecycleScreen.
type mockLifecycleScreen struct {
	mockScreen
	appearedCalled    bool
	disappearedCalled bool
}

func (m *mockLifecycleScreen) Appeared() tea.Cmd {
	m.appearedCalled = true
	return nil
}

func (m *mockLifecycleScreen) Disappeared() {
	m.disappearedCalled = true
}

func TestNewStack(t *testing.T) {
	tests := []struct {
		name    string
		root    Screen
		wantErr bool
	}{
		{
			name:    "valid root screen",
			root:    &mockScreen{name: "root"},
			wantErr: false,
		},
		{
			name:    "nil root screen panics",
			root:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("NewStack() should panic with nil root")
					}
				}()
			}
			s := NewStack(tt.root)
			if !tt.wantErr {
				if s.Depth() != 1 {
					t.Errorf("NewStack() Depth = %d, want 1", s.Depth())
				}
			}
		})
	}
}

func TestStackInit(t *testing.T) {
	root := &mockScreen{name: "root"}
	s := NewStack(root)

	cmd := s.Init()
	if cmd != nil {
		t.Errorf("Init() should return nil cmd, got %v", cmd)
	}

	if !root.initCalled {
		t.Error("Init() should call root screen Init()")
	}
}

func TestStackPush(t *testing.T) {
	root := &mockScreen{name: "root"}
	s := NewStack(root)

	// Simulate push message
	pushMsg := PushMsg{Screen: &mockScreen{name: "child"}}
	model, _ := s.Update(pushMsg)
	s = model.(Stack)

	if s.Depth() != 2 {
		t.Errorf("After Push, Depth = %d, want 2", s.Depth())
	}
}

func TestStackPop(t *testing.T) {
	root := &mockScreen{name: "root"}
	child := &mockScreen{name: "child"}
	s := NewStack(root)

	// Push child
	pushMsg := PushMsg{Screen: child}
	model, _ := s.Update(pushMsg)
	s = model.(Stack)

	// Pop child
	popMsg := PopMsg{}
	model, _ = s.Update(popMsg)
	s = model.(Stack)

	if s.Depth() != 1 {
		t.Errorf("After Pop, Depth = %d, want 1", s.Depth())
	}
}

func TestStackPopRoot(t *testing.T) {
	root := &mockScreen{name: "root"}
	s := NewStack(root)

	// Try to pop root (should not happen)
	popMsg := PopMsg{}
	model, _ := s.Update(popMsg)
	s = model.(Stack)

	if s.Depth() != 1 {
		t.Errorf("Pop on root should not change Depth, got %d", s.Depth())
	}
}

func TestStackReplace(t *testing.T) {
	root := &mockScreen{name: "root"}
	replacement := &mockScreen{name: "replacement"}
	s := NewStack(root)

	// Replace root
	replaceMsg := ReplaceMsg{Screen: replacement}
	model, _ := s.Update(replaceMsg)
	s = model.(Stack)

	if s.Depth() != 1 {
		t.Errorf("After Replace, Depth = %d, want 1", s.Depth())
	}
}

func TestStackMessageForwarding(t *testing.T) {
	root := &mockScreen{name: "root"}
	s := NewStack(root)

	// Send a non-navigation message
	testMsg := struct{}{}
	model, _ := s.Update(testMsg)
	s = model.(Stack)

	if root.updateCount != 1 {
		t.Errorf("Message should be forwarded to screen, updateCount = %d", root.updateCount)
	}
}

func TestLifecycleAppeared(t *testing.T) {
	root := &mockLifecycleScreen{}
	child := &mockLifecycleScreen{}
	s := NewStack(root)

	// Appeared is called during navigation operations, not Init
	// Init just calls the screen's Init method
	s.Init()

	// Verify Init was called
	if !root.initCalled {
		t.Error("Init() should call root screen's Init()")
	}

	// Now push a screen - this should trigger Appeared on child and Disappeared on root
	pushMsg := PushMsg{Screen: child}
	model, _ := s.Update(pushMsg)
	s = model.(Stack)

	if !root.disappearedCalled {
		t.Error("Push should cause root to Disappear()")
	}

	if !child.appearedCalled {
		t.Error("Push should cause new screen to Appear()")
	}
}

func TestLifecyclePushDisappeared(t *testing.T) {
	root := &mockLifecycleScreen{mockScreen: mockScreen{name: "root"}}
	child := &mockLifecycleScreen{mockScreen: mockScreen{name: "child"}}
	s := NewStack(root)

	// Push child - root should disappear
	pushMsg := PushMsg{Screen: child}
	model, _ := s.Update(pushMsg)
	s = model.(Stack)

	if !root.disappearedCalled {
		t.Error("Push should cause old screen to Disappear()")
	}

	if !child.appearedCalled {
		t.Error("Push should cause new screen to Appear()")
	}
}

func TestStackView(t *testing.T) {
	root := &mockScreen{name: "root screen"}
	s := NewStack(root)

	view := s.View()
	// In v2, tea.View doesn't have a Render method - it's just a wrapper
	// We can check that it's not nil
	if view == (tea.View{}) {
		t.Error("View() should return non-zero tea.View")
	}
}

func TestPushCommand(t *testing.T) {
	screen := &mockScreen{name: "test"}
	cmd := Push(screen)

	if cmd == nil {
		t.Fatal("Push() should return non-nil command")
	}

	msg := cmd()
	if _, ok := msg.(PushMsg); !ok {
		t.Errorf("Push() command should return PushMsg, got %T", msg)
	}
}

func TestPopCommand(t *testing.T) {
	cmd := Pop()

	if cmd == nil {
		t.Fatal("Pop() should return non-nil command")
	}

	msg := cmd()
	if _, ok := msg.(PopMsg); !ok {
		t.Errorf("Pop() command should return PopMsg, got %T", msg)
	}
}

func TestReplaceCommand(t *testing.T) {
	screen := &mockScreen{name: "test"}
	cmd := Replace(screen)

	if cmd == nil {
		t.Fatal("Replace() should return non-nil command")
	}

	msg := cmd()
	if _, ok := msg.(ReplaceMsg); !ok {
		t.Errorf("Replace() command should return ReplaceMsg, got %T", msg)
	}
}
