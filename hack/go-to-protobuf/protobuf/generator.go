package protobuf

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"k8s.io/gengo/v2"
	"k8s.io/gengo/v2/generator"
	"k8s.io/gengo/v2/namer"
	"k8s.io/gengo/v2/types"
)

// genProtoIDL produces a .proto IDL.
type genProtoIDL struct {
	generator.GoGenerator
	localPackage   types.Name
	localGoPackage types.Name
	imports        namer.ImportTracker

	generateAll    bool
	omitFieldTypes map[types.Name]struct{}
}

func (g *genProtoIDL) PackageVars(c *generator.Context) []string {
	return []string{
		fmt.Sprintf("option go_package = %q;", g.localGoPackage.Package),
	}
}

func (g *genProtoIDL) Filename() string { return g.OutputFilename + ".proto" }

func (g *genProtoIDL) FileType() string { return "protoidl" }

func (g *genProtoIDL) Namers(c *generator.Context) namer.NameSystems {
	return namer.NameSystems{
		"local": localNamer{g.localPackage},
	}
}

// Filter ignores types that are identified as not exportable.
func (g *genProtoIDL) Filter(c *generator.Context, t *types.Type) bool {
	tagVals := gengo.ExtractCommentTags("+", t.CommentLines)["protobuf"]
	if tagVals != nil {
		if tagVals[0] == "false" {
			return false
		}
		if tagVals[0] == "true" {
			return true
		}
		log.Fatalf(`Comment tag "protobuf" must be true or false, found: %q`, tagVals[0])
	}
	if !g.generateAll {
		return false
	}
	seen := map[*types.Type]bool{}
	return isProtoable(seen, t)
}

func isProtoable(seen map[*types.Type]bool, t *types.Type) bool {
	if seen[t] {
		return true
	}
	seen[t] = true
	switch t.Kind {
	case types.Builtin:
		return true
	case types.Alias:
		return isProtoable(seen, t.Underlying)
	case types.Slice, types.Pointer:
		return isProtoable(seen, t.Elem)
	case types.Map:
		return isProtoable(seen, t.Key) && isProtoable(seen, t.Elem)
	case types.Struct:
		if len(t.Members) == 0 {
			return true
		}
		for _, m := range t.Members {
			if isProtoable(seen, m.Type) {
				return true
			}
		}
		return false
	case types.Func, types.Chan, types.Interface:
		return false
	default:
		log.Printf("WARNING: type %q is not portable: %s", t.Kind, t.Name)
		return false
	}
}

func (g *genProtoIDL) Imports(c *generator.Context) (imports []string) {
	lines := []string{}

	for _, line := range g.imports.ImportLines() {
		lines = append(lines, line)
	}

	return lines
}

// GenerateType makes the body of a file implementing a set for type t.
func (g *genProtoIDL) GenerateType(c *generator.Context, t *types.Type, w io.Writer) error {
	sw := generator.NewSnippetWriter(w, c, "$", "$")
	b := bodyGen{
		locator: &protobufLocator{
			namer:    c.Namers["proto"].(ProtobufFromGoNamer),
			tracker:  g.imports,
			universe: c.Universe,

			localGoPackage: g.localGoPackage.Package,
		},
		localPackage: g.localPackage,

		omitFieldTypes: g.omitFieldTypes,

		t: t,
	}
	switch t.Kind {
	case types.Alias:
		return b.doAlias(sw)
	case types.Struct:
		return b.doStruct(sw)
	default:
		return b.unknown(sw)
	}
}

type ProtobufFromGoNamer interface {
	GoNameToProtoName(name types.Name) types.Name
}

type protobufLocator struct {
	namer    ProtobufFromGoNamer
	tracker  namer.ImportTracker
	universe types.Universe

	localGoPackage string
}

func (p protobufLocator) CastTypeName(name types.Name) string {
	if name.Package == p.localGoPackage {
		return name.Name
	}
	return name.String()
}

func (p protobufLocator) GoTypeForName(name types.Name) *types.Type {
	if len(name.Package) == 0 {
		name.Package = p.localGoPackage
	}
	return p.universe.Type(name)
}

func (p protobufLocator) ProtoTypeFor(t *types.Type) (*types.Type, error) {
	if t.Kind == types.Protobuf || t.Kind == types.Map {
		p.tracker.AddType(t)
		return t, nil
	}
	if t, ok := isFundamentalProtoType(t); ok {
		p.tracker.AddType(t)
		return t, nil
	}
	if t.Kind == types.Struct {
		t := &types.Type{
			Name: p.namer.GoNameToProtoName(t.Name),
			Kind: types.Protobuf,

			CommentLines: t.CommentLines,
		}
		p.tracker.AddType(t)
		return t, nil
	}
	return nil, errUnrecognizedType
}

type bodyGen struct {
	locator        *protobufLocator
	localPackage   types.Name
	omitFieldTypes map[types.Name]struct{}

	t *types.Type
}

func (b bodyGen) unknown(sw *generator.SnippetWriter) error {
	return fmt.Errorf("not sure how to generate: %#v", b.t)
}

func (b bodyGen) doAlias(sw *generator.SnippetWriter) error {
	return nil
}

func (b bodyGen) doStruct(sw *generator.SnippetWriter) error {
	if len(b.t.Name.Name) == 0 {
		return nil
	}
	if namer.IsPrivateGoName(b.t.Name.Name) {
		return nil
	}

	var alias *types.Type
	var fields []protoField
	options := []string{}
	allOptions := gengo.ExtractCommentTags("+", b.t.CommentLines)

	for k, v := range allOptions {
		switch {
		case strings.HasPrefix(k, "protobuf.options."):
			key := strings.TrimPrefix(k, "protobuf.options.")
			switch key {
			case "marshal":

			default:
				options = append(options, fmt.Sprintf("%s = %s", key, v[0]))
			}
		case k == "protobuf.as":
			fields = nil
			if alias = b.locator.GoTypeForName(types.Name{Name: v[0]}); alias == nil {
				return fmt.Errorf("type %v references alias %q which does not exist", b.t, v[0])
			}
		case k == "protobuf.embed":
			fields = []protoField{
				{
					Tag:  1,
					Name: v[0],
					Type: &types.Type{
						Name: types.Name{
							Name:    v[0],
							Package: b.localPackage.Package,
							Path:    b.localPackage.Path,
						},
					},
				},
			}
		}
	}

	if alias == nil {
		alias = b.t
	}

	if fields == nil {
		memberFields, err := membersToFields(b.locator, alias, b.localPackage, b.omitFieldTypes)
		if err != nil {
			return fmt.Errorf("type %v cannot be converted to protobuf: %v", b.t, err)
		}
		fields = memberFields
	}

	out := sw.Out()
	genComment(out, b.t.CommentLines, "")
	sw.Do(`message $.Name.Name$ {`, b.t)

	if len(options) > 0 {
		sort.Strings(options)
		for _, s := range options {
			fmt.Fprintf(out, "  option %s;\n", s)
		}
		fmt.Fprintln(out)
	}

	for i, field := range fields {
		genComment(out, field.CommentLines, "  ")
		fmt.Fprintf(out, "  ")
		switch {
		case field.Map:
			// Map fields are handled directly.
		case field.Repeated:
			fmt.Fprintf(out, "repeated ")
		}
		sw.Do(`$.Type|local$ $.Name$ = $.Tag$`, field)
		if len(field.Extras) > 0 {
			extras := []string{}
			for k, v := range field.Extras {
				extras = append(extras, fmt.Sprintf("%s = %s", k, v))
			}
			sort.Strings(extras)
			fmt.Fprintf(out, " [%s]", strings.Join(extras, ", "))
		}
		fmt.Fprintln(out, ";")
		if i == len(fields)-1 {
			fmt.Fprintln(out)
		}
	}
	fmt.Fprintln(out, "}")
	return nil
}

type protoField struct {
	LocalPackage types.Name
	Tag          int
	Name         string
	Type         *types.Type
	Extras       map[string]string
	CommentLines []string
	Repeated     bool
	Map          bool
}

var errUnrecognizedType = fmt.Errorf("did not recognize the provided type")

func isFundamentalProtoType(t *types.Type) (*types.Type, bool) {
	switch t.Kind {
	case types.Slice:
		if t.Elem.Name.Name == "byte" && len(t.Elem.Name.Package) == 0 {
			// Go slice of bytes maps to `bytes` in Proto3
			return &types.Type{Name: types.Name{Name: "bytes"}, Kind: types.Protobuf}, true
		}
	case types.Builtin:
		switch t.Name.Name {
		case "string":
			return &types.Type{Name: types.Name{Name: "string"}, Kind: types.Protobuf}, true
		case "uint32":
			return &types.Type{Name: types.Name{Name: "uint32"}, Kind: types.Protobuf}, true
		case "int32":
			return &types.Type{Name: types.Name{Name: "int32"}, Kind: types.Protobuf}, true
		case "uint64":
			return &types.Type{Name: types.Name{Name: "uint64"}, Kind: types.Protobuf}, true
		case "int64":
			return &types.Type{Name: types.Name{Name: "int64"}, Kind: types.Protobuf}, true
		case "bool":
			return &types.Type{Name: types.Name{Name: "bool"}, Kind: types.Protobuf}, true
		case "int":
			// Go int maps to Proto3 int64
			return &types.Type{Name: types.Name{Name: "int64"}, Kind: types.Protobuf}, true
		case "uint":
			// Go uint maps to Proto3 uint64
			return &types.Type{Name: types.Name{Name: "uint64"}, Kind: types.Protobuf}, true
		case "float64":
			return &types.Type{Name: types.Name{Name: "double"}, Kind: types.Protobuf}, true
		case "float32":
			return &types.Type{Name: types.Name{Name: "float"}, Kind: types.Protobuf}, true
		case "uintptr":
			// Go uintptr maps to Proto3 uint64
			return &types.Type{Name: types.Name{Name: "uint64"}, Kind: types.Protobuf}, true
		}
	}
	return nil, false
}

func memberTypeToProtobufField(locator *protobufLocator, field *protoField, t *types.Type) error {
	var err error

	switch t.Kind {
	case types.Protobuf:
		field.Type, err = locator.ProtoTypeFor(t)
	case types.Builtin:
		field.Type, err = locator.ProtoTypeFor(t)
	case types.Map:
		valueField := &protoField{}
		if err := memberTypeToProtobufField(locator, valueField, t.Elem); err != nil {
			return err
		}
		keyField := &protoField{}
		if err := memberTypeToProtobufField(locator, keyField, t.Key); err != nil {
			return err
		}
		// Map handling for Proto3
		field.Type = &types.Type{
			Kind: types.Protobuf,
			Key:  keyField.Type,
			Elem: valueField.Type,
		}

		field.Map = true
	case types.Pointer:
		if err := memberTypeToProtobufField(locator, field, t.Elem); err != nil {
			return err
		}

	case types.Alias:

		if err := memberTypeToProtobufField(locator, field, t.Underlying); err != nil {
			log.Printf("failed to alias: %s %s: err %v", t.Name, t.Underlying.Name, err)
			return err
		}
		// Handling for aliases in Proto3
		if !field.Repeated {
			if field.Extras == nil {
				field.Extras = make(map[string]string)
			}
		}

	case types.Slice:
		if t.Elem.Name.Name == "byte" && len(t.Elem.Name.Package) == 0 {
			// Go slice of bytes maps to `bytes` in Proto3
			field.Type = &types.Type{Name: types.Name{Name: "bytes"}, Kind: types.Protobuf}
			return nil
		}
		if err := memberTypeToProtobufField(locator, field, t.Elem); err != nil {
			return err
		}
		field.Repeated = true
	case types.Struct:
		if len(t.Name.Name) == 0 {
			return errUnrecognizedType
		}
		// Proto3 does not use the `Struct` type directly in most cases
		field.Type, err = locator.ProtoTypeFor(t)

	default:
		return errUnrecognizedType
	}
	return err
}

// protobufTagToField extracts information from an existing protobuf tag
func protobufTagToField(tag string, field *protoField, m types.Member, t *types.Type, localPackage types.Name) error {
	if len(tag) == 0 || tag == "-" {
		return nil
	}

	parts := strings.Split(tag, ",")
	if len(parts) < 3 {
		return fmt.Errorf("member %q of %q malformed 'protobuf' tag, not enough segments", m.Name, t.Name)
	}
	protoTag, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("member %q of %q malformed 'protobuf' tag, field ID is %q which is not an integer: %w", m.Name, t.Name, parts[1], err)
	}
	field.Tag = protoTag

	// In general there is doesn't make sense to parse the protobuf tags to get the type,
	// as all auto-generated once will have wire type "bytes", "varint" or "fixed64".
	// However, sometimes we explicitly set them to have a custom serialization, e.g.:
	//   type Time struct {
	//     time.Time `protobuf:"Timestamp,1,req,name=time"`
	//   }
	// to force the generator to use a given type (that we manually wrote serialization &
	// deserialization methods for).
	switch parts[0] {
	case "varint", "fixed32", "fixed64", "bytes", "group":
	default:
		var name types.Name
		if last := strings.LastIndex(parts[0], "."); last != -1 {
			prefix := parts[0][:last]
			name = types.Name{
				Name:    parts[0][last+1:],
				Package: prefix,
				Path:    strings.ReplaceAll(prefix, ".", "/"),
			}
		} else {
			name = types.Name{
				Name:    parts[0],
				Package: localPackage.Package,
				Path:    localPackage.Path,
			}
		}
		field.Type = &types.Type{
			Name: name,
			Kind: types.Protobuf,
		}
	}

	protoExtra := make(map[string]string)
	for i, extra := range parts[3:] {
		parts := strings.SplitN(extra, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("member %q of %q malformed 'protobuf' tag, tag %d should be key=value, got %q", m.Name, t.Name, i+4, extra)
		}
		switch parts[0] {
		case "name":
			protoExtra[parts[0]] = parts[1]
		}
	}

	field.Extras = protoExtra
	if name, ok := protoExtra["name"]; ok {
		field.Name = name
		delete(protoExtra, "name")
	}

	return nil
}

func membersToFields(locator *protobufLocator, t *types.Type, localPackage types.Name, omitFieldTypes map[types.Name]struct{}) ([]protoField, error) {
	fields := []protoField{}

	for _, m := range t.Members {
		if namer.IsPrivateGoName(m.Name) {
			// skip private fields
			continue
		}
		if _, ok := omitFieldTypes[types.Name{Name: m.Type.Name.Name, Package: m.Type.Name.Package}]; ok {
			continue
		}
		tags := reflect.StructTag(m.Tags)
		field := protoField{
			LocalPackage: localPackage,

			Tag:    -1,
			Extras: make(map[string]string),
		}

		protobufTag := tags.Get("protobuf")
		if protobufTag == "-" {
			continue
		}

		if err := protobufTagToField(protobufTag, &field, m, t, localPackage); err != nil {
			return nil, err
		}

		// extract information from JSON field tag
		if tag := tags.Get("json"); len(tag) > 0 {
			parts := strings.Split(tag, ",")
			if len(field.Name) == 0 && len(parts[0]) != 0 {
				field.Name = parts[0]
			}
			if field.Tag == -1 && field.Name == "-" {
				continue
			}
		}

		if field.Type == nil {
			if err := memberTypeToProtobufField(locator, &field, m.Type); err != nil {
				return nil, fmt.Errorf("unable to embed type %q as field %q in %q: %v", m.Type, field.Name, t.Name, err)
			}
		}
		if len(field.Name) == 0 {
			field.Name = namer.IL(m.Name)
		}

		field.CommentLines = m.CommentLines
		fields = append(fields, field)
	}

	// assign tags
	highest := 0
	byTag := make(map[int]*protoField)
	// fields are in Go struct order, which we preserve
	for i := range fields {
		field := &fields[i]
		tag := field.Tag
		if tag != -1 {
			if existing, ok := byTag[tag]; ok {
				return nil, fmt.Errorf("field %q and %q both have tag %d", field.Name, existing.Name, tag)
			}
			byTag[tag] = field
		}
		if tag > highest {
			highest = tag
		}
	}
	// starting from the highest observed tag, assign new field tags
	for i := range fields {
		field := &fields[i]
		if field.Tag != -1 {
			continue
		}
		highest++
		field.Tag = highest
		byTag[field.Tag] = field
	}
	return fields, nil
}

func genComment(out io.Writer, lines []string, indent string) {
	for {
		l := len(lines)
		if l == 0 || len(lines[l-1]) != 0 {
			break
		}
		lines = lines[:l-1]
	}
	for _, c := range lines {
		if len(c) == 0 {
			fmt.Fprintf(out, "%s//\n", indent) // avoid trailing whitespace
			continue
		}
		fmt.Fprintf(out, "%s// %s\n", indent, c)
	}
}

func formatProtoFile(source []byte) ([]byte, error) {
	// TODO; Is there any protobuf formatter?
	return source, nil
}

func assembleProtoFile(w io.Writer, f *generator.File) {
	// Write the header, if any
	w.Write(f.Header)

	// Update syntax to proto3
	fmt.Fprint(w, "syntax = \"proto3\";\n\n")

	// Write the package name if present
	if len(f.PackageName) > 0 {
		fmt.Fprintf(w, "package %s;\n\n", f.PackageName)
	}

	// Write imports if there are any
	if len(f.Imports) > 0 {
		imports := []string{}
		for i := range f.Imports {
			imports = append(imports, i)
		}
		sort.Strings(imports)
		for _, s := range imports {
			fmt.Fprintf(w, "import %q;\n", s)
		}
		fmt.Fprint(w, "\n")
	}

	// Write additional variables or settings if present
	if f.Vars.Len() > 0 {
		fmt.Fprintf(w, "%s\n", f.Vars.String())
	}

	// Write the main body of the proto file
	w.Write(f.Body.Bytes())
}

func NewProtoFile() *generator.DefaultFileType {
	return &generator.DefaultFileType{
		Format:   formatProtoFile,
		Assemble: assembleProtoFile,
	}
}
