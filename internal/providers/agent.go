package providers

// https://hub.docker.com/r/datadog/agent/tags
var (
  // AgentVersions are a slices of the latest major versions of the Datadog
  // Agent. This should be used to set the Options field when creating a 
  // Question with QType, Picker.
  AgentVersions = []string{"latest", "7.54.0", "7.53.0", "7.52.0"}
)

