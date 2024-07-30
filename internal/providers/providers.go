package providers


const (
	DOCKER = "Docker"

	// In the future...
	RDS    = "RDS"
	AURORA = "Aurora"
	AZURE  = "Azure"
	GCP    = "GCP"
)

// A Provider handles gathering the information required for it to 
// then generate a project structure.
//
// GetProviderQuestions should provide the caller a slice of functions
// that returns a pointer to a Question when executed. This allows the client
// to interate over all of the questions that it needs to present to the user.
// 
// GenerateProject should create a project directory on the users machine that
// matches the required files that the provider needs to deploy the sandboxed 
// environment. The name of the project directory should match the string 
// passed to this function.
type Provider interface {
	GetProviderQuestions() []func() *Question
	GenerateProject(string) error
}

// Returns a list of the available providers that this tool currently supports.
// Note that the order here matters and should match with the provider 
// descriptions in the GetProviderDescriptions function.
func GetAvailableProviders() []string {
	return []string{
		DOCKER,
	}
}

// Returns the description for each of the supported providers.
// Note that the order here matters and should match with the providers returned
// by the GetAvailableProviders function.
func GetProviderDescriptions() []string {
	return []string{
		"Uses Docker locally to create your project",

		// In the future...
		// "Uses Amazon RDS in our Amazon Sandbox to create your project",
		// "Uses Amazon Aurora in our Amazon Sandbox to create your project",
		// "Uses our Microsoft Azure Sandbox to create your project",
		// "Uses our Google Cloud Sandbox to create your project",
	}
}

// Returns the concrete implementation of the Provider Interface for the 
// provider name passed in.
func GetProvider(providerName string) Provider {
	switch providerName {
	case DOCKER:
		return GetDockerProvider()

	default:
		return nil
	}
}
