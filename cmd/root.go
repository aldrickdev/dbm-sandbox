package cmd

import (
	"fmt"
	"os"

	"github.com/aldrickdev/dbm-sandbox/internal/providers"
	"github.com/aldrickdev/dbm-sandbox/internal/styles"
	"github.com/aldrickdev/dbm-sandbox/internal/utils/components/picker"
	"github.com/aldrickdev/dbm-sandbox/internal/utils/components/textInput"

	"github.com/spf13/cobra"
)

const (
	DATADOG_API_KEY_ENV = "DD_API_KEY"
)

var rootCmd = &cobra.Command{
	Use:   "dbm-sandbox",
	Short: "Create a DBM sandbox",
	Long:  `A tool for automating the creation of a DBM sandbox`,
	Run:   run,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	// Displays the initial Banner
	title := `       ____                                        ____              
  ____/ / /_  ____ ___       _________ _____  ____/ / /_  ____  _  __
 / __  / __ \/ __ '__ \     / ___/ __ '/ __ \/ __  / __ \/ __ \| |/_/
/ /_/ / /_/ / / / / / /    (__  ) /_/ / / / / /_/ / /_/ / /_/ />  <  
\__,_/_.___/_/ /_/ /_/    /____/\__,_/_/ /_/\__,_/_.___/\____/_/|_|
`
	fmt.Print(styles.ProjectTitle.Render(title))

	initialText := "Welcome to the dbm sandboxing tool, where the goal is to help you create DBM sandboxes."
	fmt.Print(styles.Question.Render(initialText))

	// Gets the Datadog API Key
	ddapikey, ok := os.LookupEnv(DATADOG_API_KEY_ENV)
	if !ok {
		errorMsg := fmt.Sprintf("Failed to find your %q, please make sure to have the environment variable set", DATADOG_API_KEY_ENV)
		fmt.Println(styles.Error.Render(errorMsg))
		return
	}

	// Gets a list of the available providers
	var selectedProvider string
	pickProvider := picker.NewPicker(
		providers.GetAvailableProviders(),
		providers.GetProviderDescriptions(),
		"What provider would you like to use?",
		&selectedProvider,
	)

	// Prompt the user to select a provider
	if err := pickProvider.Run(); err != nil {
		errorMsg := fmt.Sprintf("Error Running Program: %q", err)
		fmt.Println(styles.Error.Render(errorMsg))
		return
	}

	if selectedProvider == "" {
		return
	}

	// Get the Provider instance for the provider the user selected
	provider := providers.GetProvider(selectedProvider)
	if provider == nil {
		errorMsg := fmt.Sprintf("Provider %q not implemented", selectedProvider)
		fmt.Println(styles.Error.Render(errorMsg))
		return
	}

	// Get the providers questions
	questions := provider.GetProviderQuestions()

	// Loop over each of the provider questions
	for _, questionFunc := range questions {
		question := questionFunc()

		var runner providers.RunnableQuestion

		switch question.QType {
		case providers.Picker:
			runner = picker.NewPicker(question.Options, nil, question.Prompt, &question.Answer)

		case providers.Input:
			runner = textInput.NewTextInput(question.Prompt, question.DefaultAnswer, &question.Answer)
		}

		if err := runner.Run(); err != nil {
			errorMsg := fmt.Sprintf("Error Running Program: %q", err)
			fmt.Println(styles.Error.Render(errorMsg))
			return
		}
		if question.Answer == "" {
			return
		}
	}

	// Have the provider generate the project directory
	if err := provider.GenerateProject(ddapikey); err != nil {
		errorMsg := fmt.Sprintf("Error generating project: %q", err)
		fmt.Println(styles.Error.Render(errorMsg))
		return
	}

	fmt.Println(styles.Success.Render("Your project has been created"))
}
