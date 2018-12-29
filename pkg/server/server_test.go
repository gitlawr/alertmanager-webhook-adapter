package server

import (
	"encoding/json"
	"github.com/prometheus/alertmanager/template"
	"testing"
)

var sampleInput = `{
  "receiver": "c-lhr4b:cag-mpsdv",
  "status": "firing",
  "alerts": [
    {
      "status": "firing",
      "labels": {
        "alert_name": "tttest",
        "alert_type": "nodeMemory",
        "cluster_name": "nn",
        "group_id": "c-lhr4b:cag-mpsdv",
        "mem_threshold": "1",
        "node_name": "192.168.99.210",
        "rule_id": "c-lhr4b:cag-mpsdv_car-w5r56",
        "severity": "warning",
        "total_mem": "8072684Ki",
        "used_mem": "862Mi"
      },
      "annotations": {},
      "startsAt": "2018-12-25T08:47:48.214259395Z",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": ""
    }
  ],
  "groupLabels": {},
  "commonLabels": {
    "alert_name": "tttest",
    "alert_type": "nodeMemory",
    "cluster_name": "nn",
    "group_id": "c-lhr4b:cag-mpsdv",
    "mem_threshold": "1",
    "node_name": "192.168.99.210",
    "rule_id": "c-lhr4b:cag-mpsdv_car-w5r56",
    "severity": "warning",
    "total_mem": "8072684Ki",
    "used_mem": "862Mi"
  },
  "commonAnnotations": {},
  "externalURL": "http://alertmanager-cluster-alerting-0:9093",
  "version": "4",
  "groupKey": "{}/{group_id=\"c-lhr4b:cag-mpsdv\"}/{rule_id=\"c-lhr4b:cag-mpsdv_car-w5r56\"}:{}"
}`

var sampleTemplate = `
{
    "@type": "MessageCard",
    "@context": "http://schema.org/extensions",
    "themeColor": "8C1A1A",
    "summary": "Server High Memory usage",
    "title": "Prometheus Alert (firing)",
    "sections": [
        {
			{{$size := len (index .Alerts 0).Labels}}
			{{$cur := 0}}
            "facts": [
				{{range $key, $value := (index .Alerts 0).Labels }}
				{{$cur = add $cur 1}}
				{
					"name":"{{$key}}",
					"value":"{{$value}}"
				}{{if lt $cur $size}},{{end}}
				{{end}}
            ],
            "markdown": false
        }
    ]
}`

func TestRender(t *testing.T) {
	data := template.Data{}
	json.Unmarshal([]byte(sampleInput), &data)
	b, err := renderTemplate(sampleTemplate, data)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(b))
}
