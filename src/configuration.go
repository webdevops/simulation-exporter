package main

import (
	"github.com/docker/go-units"
	"github.com/prometheus/client_golang/prometheus"
	"math/rand"
	"regexp"
	"strconv"
)

var (
	regexpValueRange = regexp.MustCompile(`(?P<from>\d+[a-zA-Z]*)-(?P<to>\d+[a-zA-Z]*)`)
)

type Configuration struct {
	Version string
	Metrics map[string]ConfigurationMetric
}

type ConfigurationMetric struct {
	Name  string
	Help  string
	Type  string

	Labels []string

	Items []ConfigurationMetricItem

	prometheus struct{
		gauge *prometheus.GaugeVec
		summary *prometheus.SummaryVec
		histogram *prometheus.HistogramVec
	}
}

type ConfigurationMetricItem struct {
	Value string
	value *float64
	rangeFrom *float64
	rangeTo *float64
	Labels map[string]string
}

func (m *ConfigurationMetric) Init() {

	for index := range m.Items {
		item := &m.Items[index]

		// parse valuerange
		if item.Value != "" {
			item.parseValue()
		}
	}
}


func (m *ConfigurationMetricItem) parseValue() {
	match := regexpValueRange.FindStringSubmatch(m.Value)

	if len(match) >= 1 {
		result := make(map[string]string)
		for i, name := range regexpValueRange.SubexpNames() {
			if i != 0 && name != "" {
				result[name] = match[i]
			}
		}

		rangeFrom, err := m.parseFloatFromString(result["from"])
		if err != nil {
			panic(err)
		}

		rangeTo, err := m.parseFloatFromString(result["to"])
		if err != nil {
			panic(err)
		}

		m.rangeFrom = &rangeFrom
		m.rangeTo = &rangeTo
	} else {
		value, err := m.parseFloatFromString(m.Value)
		if err != nil {
			panic(err)
		}

		valueFloat := float64(value)
		m.value = &valueFloat
	}
}

func (m *ConfigurationMetricItem) parseValueRange() {



}

func (m *ConfigurationMetricItem) parseFloatFromString(value string) (ret float64, error error) {
	ret, err := strconv.ParseFloat(value, 64)
	if err != nil {
		tmp, err := units.FromHumanSize(value)
		if err != nil {
			error = err
			return
		}

		ret = float64(tmp)
	}

	return
}

func (m *ConfigurationMetricItem) GenerateValue() (value float64) {
	if m.value != nil {
		value = *m.value
	} else if m.rangeFrom != nil && m.rangeTo != nil {
		value = (rand.Float64() * (*m.rangeTo - *m.rangeFrom)) + *m.rangeFrom
	}

	return
}
