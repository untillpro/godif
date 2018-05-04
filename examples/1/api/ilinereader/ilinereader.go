package ilinereader

import "context"

// NextLine returns next line to be processed, null if EOF reached
var NextLine func(ctx context.Context) string
