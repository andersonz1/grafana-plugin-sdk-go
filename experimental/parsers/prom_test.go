package parsers

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadPromFrames(t *testing.T) {
	files := []string{
		"simple-labels",
		"simple-matrix",
		"simple-vector",
		"simple-streams",
	}

	for _, name := range files {
		t.Run(name, func(t *testing.T) {
			f, err := os.Open(path.Join("testdata", name+".json"))
			require.NoError(t, err)

			iter := jsoniter.Parse(jsoniter.ConfigDefault, f, 1024)
			rsp := ReadPrometheusResult(iter)

			out, err := jsoniter.MarshalIndent(rsp, "", "  ")
			require.NoError(t, err)

			save := false
			fpath := path.Join("testdata", name+"-frame.json")
			current, err := ioutil.ReadFile(fpath)
			if err == nil {
				same := assert.JSONEq(t, string(out), string(current))
				if !same {
					save = true
				}
			} else {
				assert.Fail(t, "missing file: "+fpath)
				save = true
			}

			if save {
				err = os.WriteFile(fpath, out, 0600)
				require.NoError(t, err)
			}

			fpath = path.Join("testdata", name+"-golden.txt")
			err = experimental.CheckGoldenDataResponse(fpath, rsp, true)
			assert.NoError(t, err)
		})
	}

	t.Fail()
}
