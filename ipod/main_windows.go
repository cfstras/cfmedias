// +build windows

package ipod

import (
	"github.com/cfstras/cfmedias/errrs"
)

func (p *IPod) Sync(mountpoint string) error {
	return errrs.New("iPod sync not supported on Windows. I'm sorry.")
}
