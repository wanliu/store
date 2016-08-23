package query

type Map map[string]interface{}

type Query interface {
	Node
}

type Node interface {
	Evaluate([]byte) bool
	Elements() []Node
	Index() Index
}

type node struct {
}

type context struct {
	cache map[uint64]interface{}
	c     *bolt.Cursor
}

const PAGE_SIZE = 1024

// query op
// equal
// range
// like
// prefix
// suffix
// contains
// any

// logical concat
// Query{
// 	"Title": Like("Hysios"),
// 	"Email": Suffix("@gamil.com"),
// 	"Age": Range(18, 20),
// 	"Tag": Contains("Babel"),
// }

// Query{
// 	"Title": Like("Hysios"),
// 	"Email": Suffix("@gamil.com"),
// 	"$OR": Query(
// 		"Age": Range(18, 20)
// 	)
// }

// Title like "hysios" and Email like "*@gamil.com" or Age between 18 and 20

// q.New()

func (s *store) each(c *bolt.Cursor, q Query, handle IterFunc) (segments []uint64, err error) {
	var (
		ctx = newContext(c)
	)

	for i, node := range q.Elements() {
		first := func() { return i == 0 }

		switch node.Op {
		case NodeAnd:
			if first() {
				segments = ctx.each_all(node)
			} else {
				segments = ctx.each_intersect(node, segments)
			}
		case NodeOr:
			if first() {
				segments = ctx.each_all(node)
			} else {
				segments = ctx.each_merge(node, segments)
			}
		case NodeNot:
			if first() {
				segments = ctx.each_all(node)
			} else {
				segments = ctx.each_difference(node, segments)
			}
		// case NodeGroup:
		// 	segments = s.each(c, node, handle)
		default:
			segments = ctx.each_all(node)
		}
	}

	return segments, nil
}
func (ctx *context) each_all(node *Node) []uint64 {
	var (
		segments = make([]uint64, 0, PAGE_SIZE)
		idx      = node.Index()
		prefix   = idx.Classic()
	)

	for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
		if ok := node.Evaluate(v); ok {
			segments = append(segments, uint64(v))
		}
	}

	return segments
}

func (ctx *context) each_intersect(node *Node, segments []uint64) []uint64 {
	// var segments = make([]uint64, 0, PAGE_SIZE)
	idx := node.Index()
	prefix := idx.Classic()

	for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
		var i int
		if node.Evaluate(k) {
			for i, id = range segments {
				if id == uint64(v) {
					break
				}
			}
			if i == len(segments) {
				segments = append(segments, uint64(v))
			}
		} else {
			for i, id = range segments {
				if id == uint64(v) {
					segments = append(segments[:i], segments[i+1:]...)
					break
				}
			}
		}
	}

	return segments
}

func (ctx *context) each_difference(node *Node, segments []uint64) []uint64 {
	idx := node.Index()
	prefix := idx.Classic()

	for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
		var i int
		if !node.Evaluate(k) {
			for i, id = range segments {
				if id == uint64(v) {
					break
				}
			}
			if i == len(segments) {
				segments = append(segments, uint64(v))
			}
		} else {
			for i, id = range segments {
				if id == uint64(v) {
					segments = append(segments[:i], segments[i+1:]...)
					break
				}
			}
		}
	}

	return segments
}

func (ctx *context) each_merge(node *Node, segments []uint64) []uint64 {
	idx := node.Index()
	prefix := idx.Classic()

	for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
		var i int
		if node.Evaluate(k) {
			for i, id = range segments {
				if id == uint64(v) {
					break
				}
			}
			if i == len(segments) {
				segments = append(segments, uint64(v))
			}
		}
	}

	return segments
}

// func (ctx *context) each_cache(node *Node, segments []uint64) {
// 	var new_segments = make([]uint64, 0, PAGE_SIZE)
// 	for _, id := range segments {
// 		cache := ctx.cache[id]
// 		if cache == nil {
// 			return errors.New("missing cache")

// 		}
// 		v := cache[node.Name]

// 		if node.Evaluate(v) {
// 			new_segments = append(new_segments, id)
// 		}
// 	}

// 	return new_segments
// }

func (ctx *context) cache(id uint64, key string, val interface{}) {
	if ctx.cache[id] == nil {
		ctx.cache[id] = make(map[string]interface{})
	}

	ctx.cache[id][key] = val
}

type Node struct {
	Parent *Node
	Next   *Node
}

type EqNode struct {
	Node
	Value reflect.Value
}

type AndNode struct {
	Node
	Elements []*Node
}

type OrNode struct {
	Node
	Elements []*Node
}

func New(q Map) *Node {

}

func Or(q Map) *Node {

}

func And(q Map) *Node {

}
