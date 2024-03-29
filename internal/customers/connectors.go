package customers

type Named interface {
	Name() string
}

type Registry[T Named] struct {
	connectors map[string]T
}

var globalRegistry = Registry[GDS]{
	connectors: map[string]GDS{},
}

func GlobalRegistry() *Registry[GDS] {
	return &globalRegistry
}

func (registry *Registry[T]) Add(cs ...T) {
	for _, c := range cs {
		name := c.Name()
		if _, ok := registry.connectors[name]; ok {
			continue
		}

		registry.connectors[name] = c
	}
}

func (registry *Registry[T]) Get(name string) T {
	return registry.connectors[name]
}
