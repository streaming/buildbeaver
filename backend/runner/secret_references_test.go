package runner

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/buildbeaver/buildbeaver/common/models"
	"github.com/buildbeaver/buildbeaver/server/api/rest/documents"
)

func sortedNames(names []string) []string {
	sort.Strings(names)
	return names
}

func TestCollectReferencedSecretNames_AlwaysIncludesRepoSSHKey(t *testing.T) {
	job := &documents.Job{}
	names := collectReferencedSecretNames(job)
	require.Equal(t, []string{models.RepoSSHKeySecretName}, names)
}

func TestCollectReferencedSecretNames_JobEnvironment(t *testing.T) {
	job := &documents.Job{
		Environment: []*documents.EnvVar{
			{Name: "PLAIN", Value: "not-a-secret"},
			{Name: "FROM_SECRET", ValueFromSecret: "job-secret"},
		},
	}
	names := collectReferencedSecretNames(job)
	require.Equal(t, sortedNames([]string{models.RepoSSHKeySecretName, "job-secret"}), sortedNames(names))
}

func TestCollectReferencedSecretNames_JobDockerAuth(t *testing.T) {
	job := &documents.Job{
		DockerConfig: &documents.DockerConfig{
			BasicAuth: &documents.DockerBasicAuth{
				Username: &documents.SecretString{FromSecret: "docker-username"},
				Password: &documents.SecretString{FromSecret: "docker-password"},
			},
			AWSAuth: &documents.DockerAWSAuth{
				AWSAccessKeyID:     &documents.SecretString{FromSecret: "aws-access-key"},
				AWSSecretAccessKey: &documents.SecretString{FromSecret: "aws-secret-key"},
			},
		},
	}
	names := collectReferencedSecretNames(job)
	require.Equal(t,
		sortedNames([]string{models.RepoSSHKeySecretName, "docker-username", "docker-password", "aws-access-key", "aws-secret-key"}),
		sortedNames(names))
}

func TestCollectReferencedSecretNames_ServiceEnvironmentAndDockerAuth(t *testing.T) {
	job := &documents.Job{
		Services: []*documents.Service{
			{
				Name: "db",
				Environment: []*documents.EnvVar{
					{Name: "DB_PASSWORD", ValueFromSecret: "service-secret"},
				},
				DockerConfig: &documents.DockerConfig{
					BasicAuth: &documents.DockerBasicAuth{
						Password: &documents.SecretString{FromSecret: "service-docker-password"},
					},
				},
			},
		},
	}
	names := collectReferencedSecretNames(job)
	require.Equal(t,
		sortedNames([]string{models.RepoSSHKeySecretName, "service-secret", "service-docker-password"}),
		sortedNames(names))
}

func TestCollectReferencedSecretNames_ExplicitValuesAreNotSecretReferences(t *testing.T) {
	job := &documents.Job{
		DockerConfig: &documents.DockerConfig{
			BasicAuth: &documents.DockerBasicAuth{
				Username: &documents.SecretString{Value: "plain-username"},
			},
		},
	}
	names := collectReferencedSecretNames(job)
	require.Equal(t, []string{models.RepoSSHKeySecretName}, names)
}
