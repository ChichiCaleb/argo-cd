package protobuf

import (
	"k8s.io/gengo/v2/namer"
	"k8s.io/gengo/v2/types"
)

type ImportTracker struct {
	namer.DefaultImportTracker
}

func NewImportTracker(local types.Name, typesToAdd ...*types.Type) *ImportTracker {
	tracker := namer.NewDefaultImportTracker(local)
	tracker.IsInvalidType = func(t *types.Type) bool { return t.Kind != types.Protobuf }
	tracker.LocalName = func(name types.Name) string { return name.Package }
	tracker.PrintImport = func(path, name string) string { return path }

	tracker.AddTypes(typesToAdd...)
	return &ImportTracker{
		DefaultImportTracker: tracker,
	}
}
