package kafka

import "context"

// CommitFunc - commits kafka message.
type CommitFunc func(context.Context) error
