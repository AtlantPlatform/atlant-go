package state

type storeOptions struct {
	SyncWrites bool
}

type storeOpt func(o *storeOptions)

func defaultStoreOptions() *storeOptions {
	return &storeOptions{
		SyncWrites: true,
	}
}

func NoSyncOption() storeOpt {
	return func(o *storeOptions) {
		o.SyncWrites = false
	}
}
