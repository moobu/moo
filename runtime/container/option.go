package container

import "context"

type CreateOptions struct {
	Context context.Context
	Name    string
}

type CreateOption func(*CreateOptions)

func CreateContext(c context.Context) CreateOption {
	return func(o *CreateOptions) {
		o.Context = c
	}
}

func Name(name string) CreateOption {
	return func(o *CreateOptions) {
		o.Name = name
	}
}

type DeleteOptions struct {
	Context context.Context
}

type DeleteOption func(*DeleteOptions)

func DeleteContext(c context.Context) DeleteOption {
	return func(o *DeleteOptions) {
		o.Context = c
	}
}

type InspectOptions struct {
	Context context.Context
}

type InspectOption func(*InspectOptions)

func InspectContext(c context.Context) InspectOption {
	return func(o *InspectOptions) {
		o.Context = c
	}
}

type ListOptions struct {
	Context context.Context
	Labels  []string
}

type ListOption func(*ListOptions)

func ListContext(c context.Context) ListOption {
	return func(o *ListOptions) {
		o.Context = c
	}
}

func Labels(labels ...string) ListOption {
	return func(o *ListOptions) {
		o.Labels = labels
	}
}

type WaitOptions struct {
	Context context.Context
}

type WaitOption func(*WaitOptions)

func WaitContext(c context.Context) WaitOption {
	return func(o *WaitOptions) {
		o.Context = c
	}
}

type StartOptions struct {
	Context context.Context
}

type StartOption func(*StartOptions)

func StartContext(c context.Context) StartOption {
	return func(o *StartOptions) {
		o.Context = c
	}
}

type StopOptions struct {
	Context context.Context
}

type StopOption func(*StopOptions)

func StopContext(c context.Context) StopOption {
	return func(o *StopOptions) {
		o.Context = c
	}
}
