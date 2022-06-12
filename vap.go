package vap

import (
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

// Vap is a structure format is designed for small arenas. It contains the actual arena information, along with the
// blocks inside the structure. It is very similar to the legacy schematic format.
type Vap struct {
	// games contain the games supported by this Vap.
	games []byte
	// name is the display name of the arena.
	name string
	// positions contains two positions, one for each player with the second position being the middle position.
	positions [3]mgl64.Vec3
	// bounds contains the bounds of the arena.
	bounds [3]uint32
	// palette is the block palette of the arena.
	palette []state
	// data contains the actual block data of the arena.
	data []byte
}

// state is a structure that contains block state information.
type state struct {
	// name is the name of the block.
	name string
	// properties is a map of properties of the block.
	properties map[string]interface{}
}

// currentVersion is the current version of Vap.
var currentVersion = uint16(0x01)

// New creates a new Vap and initialises it with air blocks. The Vap returned may be written to using Vap.Set and read
// using Vap.At, and the arena information may be loaded using Vap.Arena.
func New(name string, games []byte, positions [3]mgl64.Vec3, bounds [3]uint32) *Vap {
	return &Vap{
		games:     games,
		name:      name,
		positions: positions,
		bounds:    bounds,
		palette: []state{{
			name:       "minecraft:air",
			properties: map[string]interface{}{},
		}},
		data: make([]byte, bounds[0]*bounds[1]*bounds[2]),
	}
}

// Arena returns the arena information of the Vap.
func (v *Vap) Arena() (string, []byte, [3]mgl64.Vec3) {
	return v.name, v.games, v.positions
}

// Dimensions returns the bounds of the Vap.
func (v *Vap) Dimensions() [3]int {
	return [3]int{int(v.bounds[0]), int(v.bounds[1]), int(v.bounds[2])}
}

// Set sets the block at a specific position within the Vap to the world.Block passed. Set will panic if the x, y or z
// exceed the bounds of the structure.
func (v *Vap) Set(x, y, z int, b world.Block) {
	l, w := int(v.bounds[2]), int(v.bounds[0])
	offset := (y*l+z)*w + x

	name, properties := b.EncodeBlock()
	index, ok := v.lookup(name, properties)
	if !ok {
		index = byte(len(v.palette))
		v.palette = append(v.palette, state{
			name:       name,
			properties: properties,
		})
	}
	v.data[offset] = index
}

// At returns the block at the x, y and z passed in the structure.
func (v *Vap) At(x, y, z int, _ func(x int, y int, z int) world.Block) (world.Block, world.Liquid) {
	l, w := int(v.bounds[2]), int(v.bounds[0])
	offset := (y*l+z)*w + x

	index := v.data[offset]
	if index == 0 {
		// Vap structures ensure that the index zero is air, so we can cheat here a little.
		return nil, nil
	}

	st := v.palette[index]
	b, ok := world.BlockByName(st.name, st.properties)
	if !ok {
		return nil, nil
	}
	return b, nil
}

// lookup looks up the world.Block passed in the palette of the Vap. If not found, the second return value will be false.
func (v *Vap) lookup(name string, properties map[string]interface{}) (byte, bool) {
	for index, block := range v.palette {
		if block.name == name {
			allEqual := true
			for k, v := range block.properties {
				if bVal, _ := properties[k]; bVal != v {
					allEqual = false
					break
				}
			}
			if allEqual {
				return byte(index), true
			}
		}
	}
	return 0, false
}
