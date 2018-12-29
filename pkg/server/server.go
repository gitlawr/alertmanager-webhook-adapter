package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
	alerttmpl "github.com/prometheus/alertmanager/template"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Server struct {
	Debug      bool
	Port       int
	ConfigPath string
}

func (s *Server) Start() {
	logrus.Infof("Listening on port %d", s.Port)
	http.HandleFunc("/", s.WebhookAdapter)
	logrus.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", s.Port), nil))
}

func (s *Server) handle(rw http.ResponseWriter, r *http.Request) error {
	notifier := r.URL.Query().Get("notifier")
	config, err := s.getConfig(notifier)
	if err != nil {
		return err
	}
	data := alerttmpl.Data{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return fmt.Errorf("Failed decoding Prometheus alert message: %v", err)
	}
	return sendAlert(config, data)
}

func (s *Server) WebhookAdapter(rw http.ResponseWriter, r *http.Request) {
	if err := s.handle(rw, r); err != nil {
		s.handleError(rw, err)
	}
}

func (s *Server) handleError(rw http.ResponseWriter, err error) {
	http.Error(rw, err.Error(), http.StatusInternalServerError)
}

func (s *Server) getConfig(notifierName string) (*v3.WebhookTemplateConfig, error) {
	var configs map[string]*v3.WebhookTemplateConfig
	data, err := ioutil.ReadFile(s.ConfigPath)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, &configs); err != nil {
		return nil, err
	}
	if configs[notifierName] == nil {
		return nil, fmt.Errorf("config for notifier %q not found", notifierName)
	}
	return configs[notifierName], nil
}

func sendAlert(config *v3.WebhookTemplateConfig, data alerttmpl.Data) error {
	payload, err := renderTemplate(config.Template, data)
	if err != nil {
		return err
	}
	logrus.Info(string(payload))
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Post(config.URL, "application/json", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode > http.StatusBadRequest {
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("request to %s get response code %d", string(respBody), resp.StatusCode)
	}
	return nil
}

func renderTemplate(text string, data alerttmpl.Data) ([]byte, error) {
	tmpl, err := template.New("").Funcs(sprig.TxtFuncMap()).Option("missingkey=zero").Parse(text)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	return buf.Bytes(), err
}
