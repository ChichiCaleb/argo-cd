package notification

import (
	"context"

	"github.com/argoproj/notifications-engine/pkg/api"
	apierr "k8s.io/apimachinery/pkg/api/errors"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient/notification"
)

// Server provides an Application service and implements NotificationServiceServer
type Server struct {
	apiFactory api.Factory
	notification.UnimplementedNotificationServiceServer
}

// NewServer returns a new instance of the Application service
func NewServer(apiFactory api.Factory) notification.NotificationServiceServer {
	return &Server{apiFactory: apiFactory}
}

// ListTriggers returns a list of notification triggers
func (s *Server) ListTriggers(ctx context.Context, q *notification.TriggersListRequest) (*notification.TriggerList, error) {
	api, err := s.apiFactory.GetAPI()
	if err != nil {
		if apierr.IsNotFound(err) {
			return &notification.TriggerList{}, nil
		}
		return nil, err
	}
	triggers := []*notification.Trigger{}
	for trigger := range api.GetConfig().Triggers {
		triggers = append(triggers, &notification.Trigger{Name: trigger})
	}
	return &notification.TriggerList{Items: triggers}, nil
}

// ListServices returns a list of notification services
func (s *Server) ListServices(ctx context.Context, q *notification.ServicesListRequest) (*notification.ServiceList, error) {
	api, err := s.apiFactory.GetAPI()
	if err != nil {
		if apierr.IsNotFound(err) {
			return &notification.ServiceList{}, nil
		}
		return nil, err
	}
	services := []*notification.Service{}
	for svc := range api.GetConfig().Services {
		services = append(services, &notification.Service{Name: svc})
	}
	return &notification.ServiceList{Items: services}, nil
}

// ListTemplates returns a list of notification templates
func (s *Server) ListTemplates(ctx context.Context, q *notification.TemplatesListRequest) (*notification.TemplateList, error) {
	api, err := s.apiFactory.GetAPI()
	if err != nil {
		if apierr.IsNotFound(err) {
			return &notification.TemplateList{}, nil
		}
		return nil, err
	}
	templates := []*notification.Template{}
	for tmpl := range api.GetConfig().Templates {
		templates = append(templates, &notification.Template{Name: tmpl})
	}
	return &notification.TemplateList{Items: templates}, nil
}
