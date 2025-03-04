package cmd

import (
	"github.com/igadmg/gogen/core"
	"gopkg.in/yaml.v3"
)

const (
	Tag_Archetype = "ecsa"
	Tag_Feature   = "ecsf"
	Tag_Component = "ecsc"
	Tag_Query     = "ecsq"
	Tag_System    = "ecss"
	Tag_Mixin     = "ecsm"

	Tag_Reference = "reference" // fields marked as reference are not calling Prepare, Defer, Store/Restore methods but are saved
	Tag_Transient = "transient" // fields marked as transient are not Store()'d or Restore()'d nor saved to file
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
	_, ok := t.Data[Tag_Archetype]
	if ok {
		return EcsEntity
	}
	_, ok = t.Data[Tag_Feature] // feature declined // why? rebirth for DrawCalEntity reuse
	if ok {
		return EcsFeature
	}
	_, ok = t.Data[Tag_Component]
	if ok {
		return EcsComponent
	}
	_, ok = t.Data[Tag_Query]
	if ok {
		return EcsQuery
	}

	return EcsTypeInvalid
}

func (t Tag) GetEcs() (core.Tag, bool) {
	v, ok := t.Data[Tag_Archetype]
	if ok {
		switch vt := v.(type) {
		case yaml.Node:
			return core.Tag(t).GetObject(Tag_Archetype)
		case map[string]any:
			vt["."] = Tag_Archetype
			return core.Tag{Data: vt}, true
		}
		return core.Tag{Data: core.TagData{".": Tag_Archetype}}, true
	}
	v, ok = t.Data[Tag_Component]
	if ok {
		switch vt := v.(type) {
		case yaml.Node:
			return core.Tag(t).GetObject(Tag_Component)
		case map[string]any:
			vt["."] = Tag_Component
			return core.Tag{Data: vt}, true
		}
		return core.Tag{Data: core.TagData{".": Tag_Component}}, true
	}

	return core.Tag{}, false
}
