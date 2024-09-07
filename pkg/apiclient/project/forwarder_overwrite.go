package project

import (
	
	http "github.com/argoproj/argo-cd/v2/pkg/apiclient"
)

func init() {
	forward_ProjectService_List_0 = http.UnaryForwarder
}
