//go:build js && wasm

package main

import (
	_ "embed"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/glasslabs/client-go"
	"github.com/pawal/go-hass"
)

var (
	//go:embed assets/style.css
	css []byte

	//go:embed assets/index.html
	html []byte
)

// Config is the module configuration.
type Config struct {
	URL       string `yaml:"url"`
	Token     string `yaml:"token"`
	SensorIDs struct {
		Load          string `yaml:"load"`
		PV            string `yaml:"pv"`
		Battery       string `yaml:"battery"`
		BatterySoC    string `yaml:"batterySoC"`
		Grid          string `yaml:"grid"`
		GridFrequency string `yaml:"gridFrequency"`
	} `yaml:"sensorIds"`
	Battery struct {
		Warning int `yaml:"warning"`
		Low     int `yaml:"low"`
	} `yaml:"battery"`
	MaxWatts int `yaml:"maxWatts"`
}

// NewConfig creates a default configuration for the module.
func NewConfig() *Config {
	return &Config{}
}

func main() {
	log := client.NewLogger()
	mod, err := client.NewModule()
	if err != nil {
		log.Error("Could not create module", "error", err.Error())
		return
	}

	cfg := NewConfig()
	if err = mod.ParseConfig(&cfg); err != nil {
		log.Error("Could not parse config", "error", err.Error())
		return
	}

	log.Info("Loading Module", "module", mod.Name())

	m := &Module{
		mod: mod,
		cfg: cfg,
		log: log,
	}

	if err = m.setup(); err != nil {
		log.Error("Could not setup module", "error", err.Error())
		return
	}

	first := true
	for {
		if !first {
			time.Sleep(10 * time.Second)
		}
		first = false

		if err = m.syncStates(); err != nil {
			log.Error("Could not sync states", "error", err.Error())
			continue
		}

		if err = m.listenStates(); err != nil {
			log.Error("Could not listen to states", "error", err.Error())
			continue
		}
	}
}

// Module runs the module.
type Module struct {
	mod *client.Module
	cfg *Config

	ha *hass.Access

	log *client.Logger
}

func (m *Module) setup() error {
	if err := m.mod.LoadCSS(string(css)); err != nil {
		return fmt.Errorf("loading css: %w", err)
	}
	m.mod.Element().SetInnerHTML(string(html))

	ha := hass.NewAccess(m.cfg.URL, "")
	ha.SetBearerToken(m.cfg.Token)
	if err := ha.CheckAPI(); err != nil {
		return fmt.Errorf("could not connect to home assistant: %w", err)
	}
	m.ha = ha

	return nil
}

func (m *Module) syncStates() error {
	states, err := m.ha.FilterStates("sensor")
	if err != nil {
		return fmt.Errorf("getting states: %w", err)
	}

	for _, state := range states {
		m.updateState(state.EntityID, state.State)
	}
	return nil
}

func (m *Module) listenStates() error {
	l, err := m.ha.ListenEvents()
	if err != nil {
		return fmt.Errorf("calling listen: %w", err)
	}
	defer func() { _ = l.Close() }()

	for {
		event, err := l.NextStateChanged()
		if err != nil {
			return fmt.Errorf("listening for event: %w", err)
		}

		if event.EventType != "state_changed" {
			continue
		}
		if strings.TrimSuffix(strings.SplitAfter(event.Data.EntityID, ".")[0], ".") != "sensor" {
			continue
		}

		m.updateState(event.Data.EntityID, event.Data.NewState.State)
	}
}

const percentageVar = "--percentage: "

func (m *Module) updateState(id, state string) {
	switch id {
	case m.cfg.SensorIDs.Load:
		w, err := strconv.ParseInt(state, 10, 64)
		if err != nil {
			return
		}
		kw := int(w / 1000)
		cw := int((w % 1000) / 10)
		per := float64(w) / float64(m.cfg.MaxWatts) * 100
		perStr := strconv.FormatFloat(per, 'f', 2, 64)

		if elem := m.mod.Element().QuerySelector("#load"); elem != nil {
			elem.SetAttribute("style", percentageVar+perStr)
		}
		if elem := m.mod.Element().QuerySelector("#loadText .super"); elem != nil {
			elem.SetTextContent(strconv.Itoa(kw))
		}
		if elem := m.mod.Element().QuerySelector("#loadText .sub"); elem != nil {
			elem.SetTextContent("." + strconv.Itoa(cw))
		}
	case m.cfg.SensorIDs.PV:
		w, err := strconv.ParseInt(state, 10, 64)
		if err != nil {
			return
		}
		per := float64(w) / float64(m.cfg.MaxWatts) * 100
		perStr := strconv.FormatFloat(per, 'f', 2, 64)

		if elem := m.mod.Element().QuerySelector("#pv"); elem != nil {
			elem.SetAttribute("style", percentageVar+perStr)
		}
	case m.cfg.SensorIDs.Battery:
		w, err := strconv.ParseInt(state, 10, 64)
		if err != nil {
			return
		}
		per := float64(w) / float64(m.cfg.MaxWatts) * 100
		perStr := strconv.FormatFloat(per, 'f', 2, 64)

		if elem := m.mod.Element().QuerySelector("#battery"); elem != nil {
			elem.SetAttribute("style", percentageVar+perStr)
		}
	case m.cfg.SensorIDs.BatterySoC:
		per, err := strconv.ParseInt(state, 10, 64)
		if err != nil {
			return
		}
		perStr := strconv.Itoa(int(per))

		var class string
		if per <= int64(m.cfg.Battery.Low) {
			class = "low"
		} else if per <= int64(m.cfg.Battery.Warning) {
			class = "warning"
		}

		if elem := m.mod.Element().QuerySelector("#batterySoC"); elem != nil {
			elem.SetAttribute("style", percentageVar+perStr)
			elem.Class().Remove("low")
			elem.Class().Remove("warning")
			if class != "" {
				elem.Class().Add(class)
			}
		}
		if elem := m.mod.Element().QuerySelector("#batterySoC .value"); elem != nil {
			elem.SetTextContent(perStr + "%")
		}
	case m.cfg.SensorIDs.Grid:
		w, err := strconv.ParseInt(state, 10, 64)
		if err != nil {
			return
		}
		per := float64(w) / float64(m.cfg.MaxWatts) * 100
		perStr := strconv.FormatFloat(per, 'f', 2, 64)

		if elem := m.mod.Element().QuerySelector("#grid"); elem != nil {
			elem.SetAttribute("style", percentageVar+perStr)
		}
	case m.cfg.SensorIDs.GridFrequency:
		hz, err := strconv.ParseFloat(state, 64)
		if err != nil {
			return
		}

		class := "off"
		if hz < 10 {
			class = "on"
		}

		if elem := m.mod.Element().QuerySelector("#icons #grid-disconnect"); elem != nil {
			elem.Class().Remove("on")
			elem.Class().Remove("off")
			elem.Class().Add(class)
		}
	}
}
