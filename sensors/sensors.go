package sensors

import (
	"bufio"
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

const (
	takeChip = iota
	takeAdapter
	takeLabel
	takeFeature
)

// Sensor sensor
type Sensor struct {
	Chip     string
	Adapter  string
	Features map[string]Feature // label => feature
}

// Feature feature values, example:
//
//	temp1_input: 77.000
//	temp1_max: 100.000
//	temp1_crit: 100.000
//	temp1_crit_alarm: 0.000
type Feature map[string]float64

func Get() ([]Sensor, error) {
	var buf bytes.Buffer
	cmd := exec.Command("sensors", "-u")
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	var ret []Sensor
	s := bufio.NewScanner(bytes.NewReader(buf.Bytes()))
	var sr Sensor
	var label string
	var f Feature
	next := func() {
		sr.Features[label] = f
		ret = append(ret, sr)
		sr.Chip = ""
		sr.Adapter = ""
		sr.Features = make(map[string]Feature)
		label = ""
	}
	sr.Features = make(map[string]Feature)
	take := takeChip
	for s.Scan() {
		if len(s.Text()) == 0 {
			take = takeChip
			next()
			continue
		}
		switch take {
		case takeChip:
			sr.Chip = s.Text()
			take = takeAdapter
		case takeAdapter:
			if strings.HasPrefix(s.Text(), "Adapter: ") {
				sr.Adapter = s.Text()[9:]
				take = takeLabel
			} else {
				label = s.Text()[:len(s.Text())-1]
				take = takeFeature
				f = make(Feature)
			}
		case takeLabel:
			label = s.Text()[:len(s.Text())-1]
			take = takeFeature
			f = make(Feature)
		case takeFeature:
			if !strings.HasPrefix(s.Text(), "  ") {
				sr.Features[label] = f
				label = s.Text()[:len(s.Text())-1]
				f = make(Feature)
				continue
			}
			tmp := strings.SplitN(s.Text()[2:], ":", 2)
			if len(tmp) != 2 {
				continue
			}
			n, _ := strconv.ParseFloat(tmp[1][1:], 64)
			f[tmp[0]] = n
		}
	}
	return ret, nil
}
