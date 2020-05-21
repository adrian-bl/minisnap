package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/adrian-bl/minisnap/lib/snapobj"

	"gopkg.in/yaml.v2"
)

// VolumePolicy describes the per volume policy we return.
type VolPolicy map[string]*VolPolicyEntry

type VolPolicyEntry struct {
	Schedule map[snapobj.Type]int
	Options  VolOptions
}

type VolOptions struct {
	Recursive bool
}

// yamlConfig is used to unmarshal the user config.
type yamlConf struct {
	Targets map[string]struct {
		Schedule map[string]int
		Options  VolOptions
	}
}

// parseConfig converts the YAML encoded config at path and returns a volume policy.
func parseConfig(path string) (VolPolicy, error) {
	fh, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	pl, err := ioutil.ReadAll(fh)
	if err != nil {
		return nil, err
	}

	c := yamlConf{}
	if err := yaml.Unmarshal(pl, &c); err != nil {
		return nil, err
	}

	vp := make(VolPolicy)
	for k, tg := range c.Targets {
		if _, ok := vp[k]; ok {
			return nil, fmt.Errorf("Volume '%s' defined multiple times", k)
		}

		vp[k] = &VolPolicyEntry{
			Schedule: make(map[snapobj.Type]int),
			Options:  tg.Options,
		}

		for t, v := range tg.Schedule {
			st, err := snapobj.ToType(t)
			if err != nil {
				return nil, err
			}
			if _, ok := vp[k].Schedule[st]; ok {
				return nil, fmt.Errorf("Volume '%s' defines target '%s' multiple times", k, st)
			}
			vp[k].Schedule[st] = v
		}
	}
	return vp, nil
}
