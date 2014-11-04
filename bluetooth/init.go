package bluetooth

import "pkg.linuxdeepin.com/dde-daemon"

func init() {
	loader.Register(&loader.Module{
		Name:   "bluetooth",
		Start:  Start,
		Stop:   Stop,
		Enable: true,
	})
}
