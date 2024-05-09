package testutils

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type (
	Dependencies[T any] interface {
		Build(t *testing.T) (T, error)
	}

	Test[Patient any, Data any, Deps Dependencies[Patient]] struct {
		Desc     string
		EditData func(*Data)
		EditFlow func(Data, Deps, *HookSet)
		TestFunc func(Patient, Data)
	}
)

type (
	HookSet struct {
		hooks []hook
	}

	hook struct {
		name     string
		call     *gomock.Call
		returns  []any
		disabled bool
	}
)

func (test Test[Patient, Data, Deps]) Run(
	t *testing.T,
	initDeps func(t *testing.T) Deps,
	initData func(t *testing.T) Data,
	initHookSet func(deps Deps, data Data) HookSet,
) {
	t.Helper()

	if test.TestFunc == nil {
		t.Skip("no test function provided")
	}

	deps := initDeps(t)
	data := initData(t)
	if test.EditData != nil {
		test.EditData(&data)
	}
	hooks := initHookSet(deps, data)
	if test.EditFlow != nil {
		test.EditFlow(data, deps, &hooks)
	}
	hooks.build()

	instance, err := deps.Build(t)
	require.NoErrorf(t, err, "building instance")
	require.NotNilf(t, instance, "instance is nil")

	test.TestFunc(instance, data)
}

func (hookSet *HookSet) build() {
	for _, hook := range hookSet.hooks {
		times := 1
		if hook.disabled {
			times = 0
		}
		hook.call.Return(hook.returns...).Times(times)
	}
}

// ReturnLast sets the return for a given hook and make the call to it the last one in the chain
func (hookSet *HookSet) ReturnLast(name string, results ...any) {
	found := false
	for i := range hookSet.hooks {
		if found {
			hookSet.hooks[i].disabled = true
		} else if hookSet.hooks[i].name == name {
			hookSet.hooks[i].returns = results
			found = true
		}
	}
}

func (hookSet *HookSet) DisableAll() {
	for i := range hookSet.hooks {
		hookSet.hooks[i].disabled = true
	}
}

func (hookSet *HookSet) DisableAllAfter(name string) {
	found := false
	for i := range hookSet.hooks {
		if found {
			hookSet.hooks[i].disabled = true
		} else if hookSet.hooks[i].name == name {
			found = true
		}
	}
}

func (hookSet *HookSet) Returns(name string, results ...any) {
	for i := range hookSet.hooks {
		if hookSet.hooks[i].name == name {
			hookSet.hooks[i].returns = results
		}
	}
}

func (hookSet *HookSet) Set(name string, call *gomock.Call, results ...any) {
	for i := range hookSet.hooks {
		if hookSet.hooks[i].name == name {
			hookSet.hooks[i].call.Times(0)
			hookSet.hooks[i].call = call
			hookSet.hooks[i].returns = results
		}
	}
}

func (hookSet *HookSet) Add(name string, call *gomock.Call, results ...any) {
	hookSet.hooks = append(hookSet.hooks,
		hook{
			name:    name,
			call:    call,
			returns: results,
		},
	)
}
