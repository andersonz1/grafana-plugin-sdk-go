package numeric

import "github.com/andersonz1/grafana-plugin-sdk-go/data"

func emptyFrameWithTypeMD(t data.FrameType) *data.Frame {
	return data.NewFrame("").SetMeta(&data.FrameMeta{Type: t})
}
