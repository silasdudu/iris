package view

import (
	"github.com/flosch/pongo2"
	"github.com/kataras/go-template/django"
)

// DjangoAdaptor is the  adaptor for the Django engine.
// Read more about the Django Go Template at:
// https://github.com/flosch/pongo2
// and https://github.com/kataras/go-template/tree/master/django
type DjangoAdaptor struct {
	*Adaptor
	engine *django.Engine
}

// Django returns a new kataras/go-template/django template engine
// with the same features as all iris' view engines have:
// Binary assets load (templates inside your executable with .go extension)
// Layout, Funcs, {{ url }} {{ urlpath}} for reverse routing and much more.
//
// Read more: https://github.com/flosch/pongo2
func Django(directory string, extension string) *DjangoAdaptor {
	e := django.New()
	return &DjangoAdaptor{
		Adaptor: NewAdaptor(directory, extension, e),
		engine:  e,
	}
}

// Filters for pongo2, map[name of the filter] the filter function . The filters are auto register.
func (d *DjangoAdaptor) Filters(filtersMap map[string]pongo2.FilterFunction) *DjangoAdaptor {

	if len(filtersMap) == 0 {
		return d
	}
	// configuration maps are never nil, because
	// they are initialized at each of the engine's New func
	// so we're just passing them inside it.
	for k, v := range filtersMap {
		d.engine.Config.Filters[k] = v
	}

	return d
}

// Globals share context fields between templates. https://github.com/flosch/pongo2/issues/35
func (d *DjangoAdaptor) Globals(globalsMap map[string]interface{}) *DjangoAdaptor {
	if len(globalsMap) == 0 {
		return d
	}

	for k, v := range globalsMap {
		d.engine.Config.Globals[k] = v
	}

	return d
}

// DebugTemplates enables template debugging.
// The verbose error messages will appear in browser instead of quiet passes with error code.
func (d *DjangoAdaptor) DebugTemplates(debug bool) *DjangoAdaptor {
	d.engine.Config.DebugTemplates = debug
	return d
}
