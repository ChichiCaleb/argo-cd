package e2e

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	argov1alpha1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	appsetfixture "github.com/argoproj/argo-cd/v2/test/e2e/fixture/applicationsets"
	"github.com/argoproj/argo-cd/v2/test/e2e/fixture/applicationsets/utils"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application"
)

func TestListMatrixGenerator(t *testing.T) {
	// Define fields to ignore for protobuf types
	// To avoid copying impl.MessageState sync.Mutex
	opts := cmp.Options{
		cmpopts.IgnoreFields(argov1alpha1.Application{}, "state", "sizeCache", "unknownFields"),
	}

	generateExpectedApp := func(cluster, name string) argov1alpha1.Application {
		return argov1alpha1.Application{
			TypeMeta: metav1.TypeMeta{
				Kind:       application.ApplicationKind,
				APIVersion: "argoproj.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:       fmt.Sprintf("%s-%s", cluster, name),
				Namespace:  utils.TestNamespace(),
				Finalizers: []string{"resources-finalizer.argocd.argoproj.io"},
			},
			Spec: argov1alpha1.ApplicationSpec{
				Project: "default",
				Source: &argov1alpha1.ApplicationSource{
					RepoURL:        "https://github.com/argoproj/argocd-example-apps.git",
					TargetRevision: "HEAD",
					Path:           name,
				},
				Destination: argov1alpha1.ApplicationDestination{
					Server:    "https://kubernetes.default.svc",
					Namespace: name,
				},
			},
		}
	}

	expectedApps := []argov1alpha1.Application{
		generateExpectedApp("cluster1", "kustomize-guestbook"),
		generateExpectedApp("cluster1", "helm-guestbook"),
		generateExpectedApp("cluster1", "ksonnet-guestbook"),

		generateExpectedApp("cluster2", "kustomize-guestbook"),
		generateExpectedApp("cluster2", "helm-guestbook"),
		generateExpectedApp("cluster2", "ksonnet-guestbook"),
	}

	var expectedAppsNewNamespace []argov1alpha1.Application
	var expectedAppsNewMetadata []argov1alpha1.Application

	appsetfixture.Given(t).
		// Create a ClusterGenerator-based ApplicationSet
		When().
		Create(v1alpha1.ApplicationSet{
			ObjectMeta: metav1.ObjectMeta{
				Name: "matrix-generator",
			},
			Spec: v1alpha1.ApplicationSetSpec{
				Template: v1alpha1.ApplicationSetTemplate{
					ApplicationSetTemplateMeta: v1alpha1.ApplicationSetTemplateMeta{Name: "{{values.name}}-{{path.basename}}"},
					Spec: argov1alpha1.ApplicationSpec{
						Project: "default",
						Source: &argov1alpha1.ApplicationSource{
							RepoURL:        "https://github.com/argoproj/argocd-example-apps.git",
							TargetRevision: "HEAD",
							Path:           "{{path}}",
						},
						Destination: argov1alpha1.ApplicationDestination{
							Server:    "https://kubernetes.default.svc",
							Namespace: "{{path.basename}}",
						},
					},
				},
				Generators: []v1alpha1.ApplicationSetGenerator{
					{
						Matrix: &v1alpha1.MatrixGenerator{
							Generators: []v1alpha1.ApplicationSetNestedGenerator{
								{
									List: &v1alpha1.ListGenerator{
										Elements: []apiextensionsv1.JSON{
											{Raw: []byte(`{"cluster": "my-cluster","url": "https://kubernetes.default.svc", "values": {"name": "cluster1"}}`)},
											{Raw: []byte(`{"cluster": "my-cluster","url": "https://kubernetes.default.svc", "values": {"name": "cluster2"}}`)},
										},
									},
								},
								{
									Git: &v1alpha1.GitGenerator{
										RepoURL: "https://github.com/argoproj/argocd-example-apps.git",
										Directories: []v1alpha1.GitDirectoryGeneratorItem{
											{
												Path: "*guestbook*",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}).Then().Expect(appsetfixture.ApplicationsExist(expectedApps, opts)).

		// Update the ApplicationSet template namespace, and verify it updates the Applications
		When().
		And(func() {
			for _, expectedApp := range expectedApps {
				newExpectedApp := expectedApp.DeepCopy()
				newExpectedApp.Spec.Destination.Namespace = "guestbook2"
				expectedAppsNewNamespace = append(expectedAppsNewNamespace, *newExpectedApp)
			}
		}).
		Update(func(appset *v1alpha1.ApplicationSet) {
			appset.Spec.Template.Spec.Destination.Namespace = "guestbook2"
		}).Then().Expect(appsetfixture.ApplicationsExist(expectedAppsNewNamespace, opts)).

		// Update the metadata fields in the appset template, and make sure it propagates to the apps
		When().
		And(func() {
			for _, expectedApp := range expectedAppsNewNamespace {
				expectedAppNewMetadata := expectedApp.DeepCopy()
				expectedAppNewMetadata.ObjectMeta.Annotations = map[string]string{"annotation-key": "annotation-value"}
				expectedAppNewMetadata.ObjectMeta.Labels = map[string]string{"label-key": "label-value"}
				expectedAppsNewMetadata = append(expectedAppsNewMetadata, *expectedAppNewMetadata)
			}
		}).
		Update(func(appset *v1alpha1.ApplicationSet) {
			appset.Spec.Template.Annotations = map[string]string{"annotation-key": "annotation-value"}
			appset.Spec.Template.Labels = map[string]string{"label-key": "label-value"}
		}).Then().Expect(appsetfixture.ApplicationsExist(expectedAppsNewMetadata, opts)).

		// Delete the ApplicationSet, and verify it deletes the Applications
		When().
		Delete().Then().Expect(appsetfixture.ApplicationsDoNotExist(expectedAppsNewNamespace, opts))
}

func TestClusterMatrixGenerator(t *testing.T) {
	// Define fields to ignore for protobuf types
	// To avoid copying impl.MessageState sync.Mutex
	opts := cmp.Options{
		cmpopts.IgnoreFields(argov1alpha1.Application{}, "state", "sizeCache", "unknownFields"),
	}

	generateExpectedApp := func(cluster, name string) argov1alpha1.Application {
		return argov1alpha1.Application{
			TypeMeta: metav1.TypeMeta{
				Kind:       application.ApplicationKind,
				APIVersion: "argoproj.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:       fmt.Sprintf("%s-%s", cluster, name),
				Namespace:  utils.TestNamespace(),
				Finalizers: []string{"resources-finalizer.argocd.argoproj.io"},
			},
			Spec: argov1alpha1.ApplicationSpec{
				Project: "default",
				Source: &argov1alpha1.ApplicationSource{
					RepoURL:        "https://github.com/argoproj/argocd-example-apps.git",
					TargetRevision: "HEAD",
					Path:           name,
				},
				Destination: argov1alpha1.ApplicationDestination{
					Name:      cluster,
					Namespace: name,
				},
			},
		}
	}

	expectedApps := []argov1alpha1.Application{
		generateExpectedApp("cluster1", "kustomize-guestbook"),
		generateExpectedApp("cluster1", "helm-guestbook"),
		generateExpectedApp("cluster1", "ksonnet-guestbook"),

		generateExpectedApp("cluster2", "kustomize-guestbook"),
		generateExpectedApp("cluster2", "helm-guestbook"),
		generateExpectedApp("cluster2", "ksonnet-guestbook"),
	}

	var expectedAppsNewNamespace []argov1alpha1.Application
	var expectedAppsNewMetadata []argov1alpha1.Application

	appsetfixture.Given(t).
		// Create a ClusterGenerator-based ApplicationSet
		When().
		CreateClusterSecret("my-secret", "cluster1", "https://kubernetes.default.svc").
		CreateClusterSecret("my-secret2", "cluster2", "https://kubernetes.default.svc").
		Create(v1alpha1.ApplicationSet{
			ObjectMeta: metav1.ObjectMeta{
				Name: "matrix-generator",
			},
			Spec: v1alpha1.ApplicationSetSpec{
				Template: v1alpha1.ApplicationSetTemplate{
					ApplicationSetTemplateMeta: v1alpha1.ApplicationSetTemplateMeta{Name: "{{name}}-{{path.basename}}"},
					Spec: argov1alpha1.ApplicationSpec{
						Project: "default",
						Source: &argov1alpha1.ApplicationSource{
							RepoURL:        "https://github.com/argoproj/argocd-example-apps.git",
							TargetRevision: "HEAD",
							Path:           "{{path}}",
						},
						Destination: argov1alpha1.ApplicationDestination{
							Name:      "{{name}}",
							Namespace: "{{path.basename}}",
						},
					},
				},
				Generators: []v1alpha1.ApplicationSetGenerator{
					{
						Matrix: &v1alpha1.MatrixGenerator{
							Generators: []v1alpha1.ApplicationSetNestedGenerator{
								{
									Clusters: &v1alpha1.ClusterGenerator{
										Selector: metav1.LabelSelector{
											MatchLabels: map[string]string{
												"argocd.argoproj.io/secret-type": "cluster",
											},
										},
									},
								},
								{
									Git: &v1alpha1.GitGenerator{
										RepoURL: "https://github.com/argoproj/argocd-example-apps.git",
										Directories: []v1alpha1.GitDirectoryGeneratorItem{
											{
												Path: "*guestbook*",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}).Then().Expect(appsetfixture.ApplicationsExist(expectedApps, opts)).

		// Update the ApplicationSet template namespace, and verify it updates the Applications
		When().
		And(func() {
			for _, expectedApp := range expectedApps {
				newExpectedApp := expectedApp.DeepCopy()
				newExpectedApp.Spec.Destination.Namespace = "guestbook2"
				expectedAppsNewNamespace = append(expectedAppsNewNamespace, *newExpectedApp)
			}
		}).
		Update(func(appset *v1alpha1.ApplicationSet) {
			appset.Spec.Template.Spec.Destination.Namespace = "guestbook2"
		}).Then().Expect(appsetfixture.ApplicationsExist(expectedAppsNewNamespace, opts)).

		// Update the metadata fields in the appset template, and make sure it propagates to the apps
		When().
		And(func() {
			for _, expectedApp := range expectedAppsNewNamespace {
				expectedAppNewMetadata := expectedApp.DeepCopy()
				expectedAppNewMetadata.ObjectMeta.Annotations = map[string]string{"annotation-key": "annotation-value"}
				expectedAppNewMetadata.ObjectMeta.Labels = map[string]string{"label-key": "label-value"}
				expectedAppsNewMetadata = append(expectedAppsNewMetadata, *expectedAppNewMetadata)
			}
		}).
		Update(func(appset *v1alpha1.ApplicationSet) {
			appset.Spec.Template.Annotations = map[string]string{"annotation-key": "annotation-value"}
			appset.Spec.Template.Labels = map[string]string{"label-key": "label-value"}
		}).Then().Expect(appsetfixture.ApplicationsExist(expectedAppsNewMetadata, opts)).

		// Delete the ApplicationSet, and verify it deletes the Applications
		When().
		Delete().Then().Expect(appsetfixture.ApplicationsDoNotExist(expectedAppsNewNamespace, opts))
}

func TestMatrixTerminalMatrixGeneratorSelector(t *testing.T) {
	// Define fields to ignore for protobuf types
	// To avoid copying impl.MessageState sync.Mutex
	opts := cmp.Options{
		cmpopts.IgnoreFields(argov1alpha1.Application{}, "state", "sizeCache", "unknownFields"),
	}

	generateExpectedApp := func(cluster, name string) argov1alpha1.Application {
		return argov1alpha1.Application{
			TypeMeta: metav1.TypeMeta{
				Kind:       application.ApplicationKind,
				APIVersion: "argoproj.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:       fmt.Sprintf("%s-%s", cluster, name),
				Namespace:  utils.TestNamespace(),
				Finalizers: []string{"resources-finalizer.argocd.argoproj.io"},
			},
			Spec: argov1alpha1.ApplicationSpec{
				Project: "default",
				Source: &argov1alpha1.ApplicationSource{
					RepoURL:        "https://github.com/argoproj/argocd-example-apps.git",
					TargetRevision: "HEAD",
					Path:           name,
				},
				Destination: argov1alpha1.ApplicationDestination{
					Server:    "https://kubernetes.default.svc",
					Namespace: name,
				},
			},
		}
	}

	expectedApps1 := []argov1alpha1.Application{
		generateExpectedApp("cluster1", "kustomize-guestbook"),
		generateExpectedApp("cluster1", "helm-guestbook"),
		generateExpectedApp("cluster1", "ksonnet-guestbook"),
	}
	expectedApps2 := []argov1alpha1.Application{
		generateExpectedApp("cluster2", "kustomize-guestbook"),
		generateExpectedApp("cluster2", "helm-guestbook"),
		generateExpectedApp("cluster2", "ksonnet-guestbook"),
	}

	appsetfixture.Given(t).
		// Create ApplicationSet with LabelSelector on an ApplicationSetTerminalGenerator
		When().
		Create(v1alpha1.ApplicationSet{
			ObjectMeta: metav1.ObjectMeta{
				Name: "matrix-generator-nested-matrix",
			},
			Spec: v1alpha1.ApplicationSetSpec{
				ApplyNestedSelectors: true,
				Template: v1alpha1.ApplicationSetTemplate{
					ApplicationSetTemplateMeta: v1alpha1.ApplicationSetTemplateMeta{Name: "{{values.name}}-{{path.basename}}"},
					Spec: argov1alpha1.ApplicationSpec{
						Project: "default",
						Source: &argov1alpha1.ApplicationSource{
							RepoURL:        "https://github.com/argoproj/argocd-example-apps.git",
							TargetRevision: "HEAD",
							Path:           "{{path}}",
						},
						Destination: argov1alpha1.ApplicationDestination{
							Server:    "https://kubernetes.default.svc",
							Namespace: "{{path.basename}}",
						},
					},
				},
				Generators: []v1alpha1.ApplicationSetGenerator{
					{
						Matrix: &v1alpha1.MatrixGenerator{
							Generators: []v1alpha1.ApplicationSetNestedGenerator{
								{
									Matrix: toAPIExtensionsJSON(t, &v1alpha1.NestedMatrixGenerator{
										Generators: []v1alpha1.ApplicationSetTerminalGenerator{
											{
												List: &v1alpha1.ListGenerator{
													Elements: []apiextensionsv1.JSON{
														{Raw: []byte(`{"cluster": "my-cluster","url": "https://kubernetes.default.svc", "values": {"name": "cluster1"}}`)},
														{Raw: []byte(`{"cluster": "my-cluster","url": "https://kubernetes.default.svc", "values": {"name": "cluster2"}}`)},
													},
												},
												Selector: &metav1.LabelSelector{
													MatchLabels: map[string]string{
														"values.name": "cluster1",
													},
												},
											},
											{
												Git: &v1alpha1.GitGenerator{
													RepoURL: "https://github.com/argoproj/argocd-example-apps.git",
													Directories: []v1alpha1.GitDirectoryGeneratorItem{
														{
															Path: "*guestbook*",
														},
													},
												},
											},
										},
									}),
								},
								{
									List: &v1alpha1.ListGenerator{
										Elements: []apiextensionsv1.JSON{
											{Raw: []byte(`{}`)},
										},
									},
								},
							},
						},
					},
				},
			},
		}).Then().Expect(appsetfixture.ApplicationsExist(expectedApps1, opts)).Expect(appsetfixture.ApplicationsDoNotExist(expectedApps2, opts)).

		// Update the ApplicationSetTerminalGenerator LabelSelector, and verify the Applications are deleted and created
		When().
		Update(func(appset *v1alpha1.ApplicationSet) {
			appset.Spec.Generators[0].Matrix.Generators[0].Matrix = toAPIExtensionsJSON(t, &v1alpha1.NestedMatrixGenerator{
				Generators: []v1alpha1.ApplicationSetTerminalGenerator{
					{
						List: &v1alpha1.ListGenerator{
							Elements: []apiextensionsv1.JSON{
								{Raw: []byte(`{"cluster": "my-cluster","url": "https://kubernetes.default.svc", "values": {"name": "cluster1"}}`)},
								{Raw: []byte(`{"cluster": "my-cluster","url": "https://kubernetes.default.svc", "values": {"name": "cluster2"}}`)},
							},
						},
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"values.name": "cluster2",
							},
						},
					},
					{
						Git: &v1alpha1.GitGenerator{
							RepoURL: "https://github.com/argoproj/argocd-example-apps.git",
							Directories: []v1alpha1.GitDirectoryGeneratorItem{
								{
									Path: "*guestbook*",
								},
							},
						},
					},
				},
			})
		}).Then().Expect(appsetfixture.ApplicationsExist(expectedApps2, opts)).Expect(appsetfixture.ApplicationsDoNotExist(expectedApps1, opts)).

		// Set ApplyNestedSelector to false and verify all Applications are created
		When().
		Update(func(appset *v1alpha1.ApplicationSet) {
			appset.Spec.ApplyNestedSelectors = false
		}).Then().Expect(appsetfixture.ApplicationsExist(expectedApps1, opts)).Expect(appsetfixture.ApplicationsExist(expectedApps2, opts)).

		// Delete the ApplicationSet, and verify it deletes the Applications
		When().
		Delete().Then().Expect(appsetfixture.ApplicationsDoNotExist(expectedApps1, opts)).Expect(appsetfixture.ApplicationsDoNotExist(expectedApps2, opts))
}

func TestMatrixTerminalMergeGeneratorSelector(t *testing.T) {
	// Define fields to ignore for protobuf types
	// To avoid copying impl.MessageState sync.Mutex
	opts := cmp.Options{
		cmpopts.IgnoreFields(argov1alpha1.Application{}, "state", "sizeCache", "unknownFields"),
	}

	generateExpectedApp := func(name, nameSuffix string) argov1alpha1.Application {
		return argov1alpha1.Application{
			TypeMeta: metav1.TypeMeta{
				Kind:       application.ApplicationKind,
				APIVersion: "argoproj.io/v1alpha1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:       fmt.Sprintf("%s-%s", name, nameSuffix),
				Namespace:  utils.TestNamespace(),
				Finalizers: []string{"resources-finalizer.argocd.argoproj.io"},
			},
			Spec: argov1alpha1.ApplicationSpec{
				Project: "default",
				Source: &argov1alpha1.ApplicationSource{
					RepoURL:        "https://github.com/argoproj/argocd-example-apps.git",
					TargetRevision: "HEAD",
					Path:           name,
				},
				Destination: argov1alpha1.ApplicationDestination{
					Server:    "https://kubernetes.default.svc",
					Namespace: name,
				},
			},
		}
	}

	expectedApps1 := []argov1alpha1.Application{
		generateExpectedApp("kustomize-guestbook", "1"),
	}
	expectedApps2 := []argov1alpha1.Application{
		generateExpectedApp("helm-guestbook", "2"),
	}

	appsetfixture.Given(t).
		// Create ApplicationSet with LabelSelector on an ApplicationSetTerminalGenerator
		When().
		Create(v1alpha1.ApplicationSet{
			ObjectMeta: metav1.ObjectMeta{
				Name: "matrix-generator-nested-merge",
			},
			Spec: v1alpha1.ApplicationSetSpec{
				ApplyNestedSelectors: true,
				Template: v1alpha1.ApplicationSetTemplate{
					ApplicationSetTemplateMeta: v1alpha1.ApplicationSetTemplateMeta{Name: "{{path.basename}}-{{name-suffix}}"},
					Spec: argov1alpha1.ApplicationSpec{
						Project: "default",
						Source: &argov1alpha1.ApplicationSource{
							RepoURL:        "https://github.com/argoproj/argocd-example-apps.git",
							TargetRevision: "HEAD",
							Path:           "{{path}}",
						},
						Destination: argov1alpha1.ApplicationDestination{
							Server:    "https://kubernetes.default.svc",
							Namespace: "{{path.basename}}",
						},
					},
				},
				Generators: []v1alpha1.ApplicationSetGenerator{
					{
						Matrix: &v1alpha1.MatrixGenerator{
							Generators: []v1alpha1.ApplicationSetNestedGenerator{
								{
									Merge: toAPIExtensionsJSON(t, &v1alpha1.NestedMergeGenerator{
										MergeKeys: []string{"path.basename"},
										Generators: []v1alpha1.ApplicationSetTerminalGenerator{
											{
												Git: &v1alpha1.GitGenerator{
													RepoURL: "https://github.com/argoproj/argocd-example-apps.git",
													Directories: []v1alpha1.GitDirectoryGeneratorItem{
														{
															Path: "*guestbook*",
														},
													},
												},
												Selector: &metav1.LabelSelector{
													MatchLabels: map[string]string{
														"path.basename": "kustomize-guestbook",
													},
												},
											},
											{
												List: &v1alpha1.ListGenerator{
													Elements: []apiextensionsv1.JSON{
														{Raw: []byte(`{"path.basename": "kustomize-guestbook", "name-suffix": "1"}`)},
														{Raw: []byte(`{"path.basename": "helm-guestbook", "name-suffix": "2"}`)},
													},
												},
											},
										},
									}),
								},
								{
									List: &v1alpha1.ListGenerator{
										Elements: []apiextensionsv1.JSON{
											{Raw: []byte(`{}`)},
										},
									},
								},
							},
						},
					},
				},
			},
		}).Then().Expect(appsetfixture.ApplicationsExist(expectedApps1, opts)).Expect(appsetfixture.ApplicationsDoNotExist(expectedApps2, opts)).

		// Update the ApplicationSetTerminalGenerator LabelSelector, and verify the Applications are deleted and created
		When().
		Update(func(appset *v1alpha1.ApplicationSet) {
			appset.Spec.Generators[0].Matrix.Generators[0].Merge = toAPIExtensionsJSON(t, &v1alpha1.NestedMergeGenerator{
				MergeKeys: []string{"path.basename"},
				Generators: []v1alpha1.ApplicationSetTerminalGenerator{
					{
						Git: &v1alpha1.GitGenerator{
							RepoURL: "https://github.com/argoproj/argocd-example-apps.git",
							Directories: []v1alpha1.GitDirectoryGeneratorItem{
								{
									Path: "*guestbook*",
								},
							},
						},
						Selector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"path.basename": "helm-guestbook",
							},
						},
					},
					{
						List: &v1alpha1.ListGenerator{
							Elements: []apiextensionsv1.JSON{
								{Raw: []byte(`{"path.basename": "kustomize-guestbook", "name-suffix": "1"}`)},
								{Raw: []byte(`{"path.basename": "helm-guestbook", "name-suffix": "2"}`)},
							},
						},
					},
				},
			})
		}).Then().Expect(appsetfixture.ApplicationsExist(expectedApps2, opts)).Expect(appsetfixture.ApplicationsDoNotExist(expectedApps1, opts)).

		// Set ApplyNestedSelector to false and verify all Applications are created
		When().
		Update(func(appset *v1alpha1.ApplicationSet) {
			appset.Spec.ApplyNestedSelectors = false
		}).Then().Expect(appsetfixture.ApplicationsExist(expectedApps1, opts)).Expect(appsetfixture.ApplicationsExist(expectedApps2, opts)).

		// Delete the ApplicationSet, and verify it deletes the Applications
		When().
		Delete().Then().Expect(appsetfixture.ApplicationsDoNotExist(expectedApps1, opts)).Expect(appsetfixture.ApplicationsDoNotExist(expectedApps2, opts))
}
