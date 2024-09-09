package project

import (
	http "github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
)

func init() {
	forward_ProjectService_List_0 = http.UnaryForwarder
}
