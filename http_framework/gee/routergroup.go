package gee

type RouterGroup struct {
	prefix     string  // 组前缀
	engine     *Engine // 公用一个Engine
	parent     *RouterGroup
	middleware []HandlerFunc // 中间件
}

func (g *RouterGroup) Append(prefix string) *RouterGroup {
	engine := g.engine
	routerGroup := &RouterGroup{
		prefix: g.prefix + prefix,
		engine: engine,
		parent: g,
	}
	g.engine.groups = append(g.engine.groups, routerGroup)

	return routerGroup
}

func (g *RouterGroup) Use(middleware ...HandlerFunc) {
	g.middleware = append(g.middleware, middleware...)
}

func (g *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := g.prefix + comp
	g.engine.router.addRoute(method, pattern, handler)
}
