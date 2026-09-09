package runner

import (
	"github.com/buildbeaver/buildbeaver/common/models"
	"github.com/buildbeaver/buildbeaver/server/api/rest/documents"
)

// collectReferencedSecretNames returns the plaintext names of every secret that job's
// configuration could reference: its own environment and Docker registry auth, and those of each
// of its services. This lets the runner fetch only the secrets a job actually needs, rather than
// every secret in the repo. Always includes models.RepoSSHKeySecretName, since that internal
// secret is needed for git checkout regardless of what the job's own config references.
func collectReferencedSecretNames(job *documents.Job) []string {
	names := map[string]bool{
		models.RepoSSHKeySecretName: true,
	}
	addEnvVarSecretNames(names, job.Environment)
	addDockerConfigSecretNames(names, job.DockerConfig)
	for _, service := range job.Services {
		addEnvVarSecretNames(names, service.Environment)
		addDockerConfigSecretNames(names, service.DockerConfig)
	}

	result := make([]string, 0, len(names))
	for name := range names {
		result = append(result, name)
	}
	return result
}

func addEnvVarSecretNames(names map[string]bool, envVars []*documents.EnvVar) {
	for _, envVar := range envVars {
		if envVar.ValueFromSecret != "" {
			names[envVar.ValueFromSecret] = true
		}
	}
}

func addDockerConfigSecretNames(names map[string]bool, configOrNil *documents.DockerConfig) {
	if configOrNil == nil {
		return
	}
	if configOrNil.BasicAuth != nil {
		addSecretStringName(names, configOrNil.BasicAuth.Username)
		addSecretStringName(names, configOrNil.BasicAuth.Password)
	}
	if configOrNil.AWSAuth != nil {
		addSecretStringName(names, configOrNil.AWSAuth.AWSAccessKeyID)
		addSecretStringName(names, configOrNil.AWSAuth.AWSSecretAccessKey)
	}
}

func addSecretStringName(names map[string]bool, secretStringOrNil *documents.SecretString) {
	if secretStringOrNil != nil && secretStringOrNil.FromSecret != "" {
		names[secretStringOrNil.FromSecret] = true
	}
}
