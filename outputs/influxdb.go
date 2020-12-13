package outputs

import (
	"log"
	"strings"

	"github.com/falcosecurity/falcosidekick/types"
)

type influxdbPayload string

func newInfluxdbPayload(falcopayload types.FalcoPayload, config *types.Configuration) influxdbPayload {
	s := "events,rule=" + strings.Replace(falcopayload.Rule, " ", "_", -1) + ",priority=" + strings.Replace(falcopayload.Priority, " ", "_", -1)

	for i, j := range falcopayload.OutputFields {
		switch v := j.(type) {
		case string:
			s += "," + i + "=" + strings.Replace(v, " ", "_", -1)
		default:
			continue
		}
	}

	s += " value=\"" + falcopayload.Output + "\""

	return influxdbPayload(s)
}

// InfluxdbPost posts event to InfluxDB
func (c *Client) InfluxdbPost(falcopayload types.FalcoPayload) {
	c.Stats.Influxdb.Add(Total, 1)

	err := c.Post(newInfluxdbPayload(falcopayload, c.Config))
	if err != nil {
		go c.CountMetric(Outputs, 1, []string{"output:influxdb", "status:error"})
		c.Stats.Influxdb.Add(Error, 1)
		c.PromStats.Outputs.With(map[string]string{"destination": "influxdb", "status": Error}).Inc()
		log.Printf("[ERROR] : InfluxDB - %v\n", err)
		return
	}

	// Setting the success status
	go c.CountMetric(Outputs, 1, []string{"output:influxdb", "status:ok"})
	c.Stats.Influxdb.Add(OK, 1)
	c.PromStats.Outputs.With(map[string]string{"destination": "influxdb", "status": OK}).Inc()
	log.Printf("[INFO] : InfluxDB - Publish OK\n")
}
