package structure

import (
	tfjson "github.com/hashicorp/terraform-json"
	alicloud "github.com/lonegunmanb/terraform-alicloud-schema/generated"
	aws "github.com/lonegunmanb/terraform-aws-schema/v5/generated"
	azurerm "github.com/lonegunmanb/terraform-azurerm-schema/v3/generated"
	google "github.com/lonegunmanb/terraform-google-schema/v4/generated"
)

var resources = func() map[string]*tfjson.Schema {
	return new(mapMerger[*tfjson.Schema]).
		Merge(aws.Resources,
			alicloud.Resources,
			azurerm.Resources,
			google.Resources).m
}()

type mapMerger[T any] struct {
	m map[string]T
}

func (merger *mapMerger[T]) Merge(inputs ...map[string]T) *mapMerger[T] {
	if merger.m == nil {
		merger.m = make(map[string]T)
	}
	for _, im := range inputs {
		for k, v := range im {
			merger.m[k] = v
		}
	}
	return merger
}
