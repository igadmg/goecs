package cmd

import (
	"github.com/igadmg/gogen/core"
	"gopkg.in/yaml.v3"
)

const (
	Tag_Archetype = "archetype"
	Tag_Feature   = "feature"
	Tag_Component = "component"
	Tag_Query     = "query"
	Tag_System    = "system"

	Tag_Reference = "reference" // fields marked as reference are not calling Prepare, Defer, Store/Restore methods but are saved
	Tag_Transient = "transient" // fields marked as transient are not Store()'d or Restore()'d nor saved to file

	Tag_Cached = "cached" // queries marked as cached get *Cached structs generated
)

const (
	Tag_Abstract = "abstract" // field is abstract - it does not have it's own storage but can be overrided by subarchetypes
	Tag_Virtual  = "virtual"  // field is virtual - it does have it's own storage and can be overrided by subarchetypes
)

const (
	Tag_Fn_RefCall = "fn_ref_call"
)

type Tag core.Tag

func (t Tag) GetEcsTag() EcsType {
	if _, ok := t.Data[Tag_Archetype]; ok {
		return EcsEntity
	}
	if _, ok := t.Data[Tag_Feature]; ok {
		return EcsFeature
	}
	if _, ok := t.Data[Tag_Component]; ok {
		return EcsComponent
	}
	if _, ok := t.Data[Tag_Query]; ok {
		return EcsQuery
	}
	if _, ok := t.Data[Tag_System]; ok {
		return EcsSystem
	}

	return EcsTypeInvalid
}

func (t Tag) GetEcs() (core.Tag, bool) {
	gt := func(tag string, v any) (core.Tag, bool) {
		switch vt := v.(type) {
		case yaml.Node:
			return core.Tag(t).GetObject(tag)
		case map[string]any:
			vt["."] = tag
			return core.Tag{Data: vt}, true
		}
		return core.Tag{Data: core.TagData{".": tag}}, true
	}

	if v, ok := t.Data[Tag_Archetype]; ok {
		return gt(Tag_Archetype, v)
	}
	if v, ok := t.Data[Tag_Feature]; ok {
		return gt(Tag_Feature, v)
	}
	if v, ok := t.Data[Tag_Component]; ok {
		return gt(Tag_Component, v)
	}
	if v, ok := t.Data[Tag_Query]; ok {
		return gt(Tag_Query, v)
	}
	if v, ok := t.Data[Tag_System]; ok {
		return gt(Tag_System, v)
	}

	return core.Tag{}, false
}
