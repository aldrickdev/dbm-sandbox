package providers

import (
	"bytes"
	"dbm-sandbox/internal/utils/helpers"
	"embed"
	"fmt"
	"os"
	"strings"
	"text/template"
)

var (
	// templateFS is where the template files used to create to the project are
	// embeded.
	//go:embed embed/docker/*
	templateFS embed.FS

	// The indexes of where the provider questions can be found, these are set
	// when questions are generated.
	ProjectNameIndex  uint8
	AgentVersionIndex uint8
	DBMSIndex         uint8
	DBMSVersionIndex  uint8
)

// DockerProvider implements the Provider Interface and holds all the required
// information needed to create a Docker Project.
type DockerProvider struct {
	// QuestionAnswers is a slice of *Question for this provider.
	QuestionAnswers []*Question
	// supportedDBMS is a slice of names for all of the DBMS's that this provider
	// supports.
	supportedDBMS []string
	// questionFuncs is a slice of functions that returns a *Question.
	questionFuncs []func() *Question

	// templatePath is the location in the templateFS where the templates for
	// this provider are located.
	templatePath string
	// templateData is the data required to fill the project template files for
	// this provider.
	templateData dockerTemplateData
	// templateFS is where the template files for this provider are located.
	templateFS embed.FS
}

// dockerTemplateData is used to contain the data for the DockerProvider
// template files. Most of these will be filled in using the field
// Question.Answer.
type dockerTemplateData struct {
	Agent agentTemplateData
	DB    dbTemplateData
}

// agentTemplateData is used to contain the agent data for the
// dockerTemplateData.Agent.
type agentTemplateData struct {
	// Version is used to contain the version of the agent that the user would
	// like to deploy.
	Version string
	// DDAPIKey is used to contain the Datadog API Key of the user.
	DDAPIKey string
	// ProjectName is used to contain the name of the directory for the project.
	ProjectName string
}

// dbTemplateData is used to contain all of the data for the
// dockerTemplateData.DB.
type dbTemplateData struct {
	// DBMS is used to contain the Database Management System that will be used
	// for the project.
	DBMS string
	// Version is used to contain the Version of the DBMS to use for the project.
	Version string
}

// GetDockerProvider will initiallize a new DockerProvider instance and return
// a pointer to it.
func GetDockerProvider() *DockerProvider {
	dp := new(DockerProvider)
	supportedDBMS := dp.getSupportedDMBS()
	supportedDBMSNames := []string{}
	for _, dbms := range supportedDBMS {
		supportedDBMSNames = append(supportedDBMSNames, dbms.Name)
	}
	dp.supportedDBMS = supportedDBMSNames
	dp.generateProviderQuestions()
	dp.templatePath = "embed/docker/"
	dp.templateFS = templateFS

	return dp
}

// getSupportedDMBS returns a slice of all of the DBMS's that this provider
// supports.
func (d *DockerProvider) getSupportedDMBS() []DBMS {
	return []DBMS{
		PostgresDBMS(),
		MySQLDBMS(),
		SQLServerDBMS(),
	}
}

// generateProviderQuestions generates all of the questions that are required
// to fill the providers template data.
func (d *DockerProvider) generateProviderQuestions() {
	projectName := func() *Question {
		question := &Question{
			QType:         Input,
			Prompt:        "What is your project name?",
			DefaultAnswer: "dbm-sandbox",
		}

		d.QuestionAnswers = append(d.QuestionAnswers, question)

		return question
	}

	agentVersion := func() *Question {
		question := &Question{
			QType:   Picker,
			Prompt:  "What version of the agent would you like to use?",
			Options: AgentVersions,
		}
		d.QuestionAnswers = append(d.QuestionAnswers, question)

		return question
	}
	dbmsPicker := func() *Question {
		question := &Question{
			QType:   Picker,
			Prompt:  "What DBMS would you like to use?",
			Options: d.supportedDBMS,
		}
		d.QuestionAnswers = append(d.QuestionAnswers, question)

		return question
	}
	dbmsVersionInput := func() *Question {
		selectedDBMS := d.QuestionAnswers[DBMSIndex].Answer
		dbmsInfo := GetDBMS(selectedDBMS)

		question := &Question{
			QType:   Picker,
			Prompt:  "What version of the DBM would you like to use?",
			Options: dbmsInfo.versions,
		}
		d.QuestionAnswers = append(d.QuestionAnswers, question)

		return question
	}

	d.addQuestion(projectName, &ProjectNameIndex)
	d.addQuestion(agentVersion, &AgentVersionIndex)
	d.addQuestion(dbmsPicker, &DBMSIndex)
	d.addQuestion(dbmsVersionInput, &DBMSVersionIndex)
}

// addQuestion will set the proper index for the location of the question so
// that they can be accessed correctly in fillTemplateData.
func (d *DockerProvider) addQuestion(q func() *Question, indexVariable *uint8) {
	*indexVariable = uint8(len(d.questionFuncs))
	d.questionFuncs = append(d.questionFuncs, q)
}

// GetProviderQuestions provides a slice of functions that return a *Question
// for the provider.
func (d *DockerProvider) GetProviderQuestions() []func() *Question {
	return d.questionFuncs
}

// fillTemplateData will fill the DockerProvider.templateData with the answers
// from the DockerProvider.QuestionAnswers.
func (d *DockerProvider) fillTemplateData(ddapikey string) {
	d.templateData = dockerTemplateData{
		Agent: agentTemplateData{
			Version:     d.QuestionAnswers[AgentVersionIndex].Answer,
			DDAPIKey:    ddapikey,
			ProjectName: d.QuestionAnswers[ProjectNameIndex].Answer,
		},
		DB: dbTemplateData{
			DBMS:    d.QuestionAnswers[DBMSIndex].Answer,
			Version: d.QuestionAnswers[DBMSVersionIndex].Answer,
		},
	}

}

// GenerateProject will generate the project directory on the users machine
// using the templates and template data.
func (d *DockerProvider) GenerateProject(ddapikey string) error {
	d.fillTemplateData(ddapikey)

	if err := helpers.CheckDirectory(d.templateData.Agent.ProjectName); err != nil {
		return err
	}

	if err := helpers.CreateDirectory(d.templateData.Agent.ProjectName); err != nil {
		return err
	}

	var content bytes.Buffer

	providerDirectory := d.templatePath
	composeTemplatePath := d.templatePath + "docker-compose.tmpl"

	dbms := providerDirectory + strings.ToLower(d.templateData.DB.DBMS)
	if err := helpers.CopyDirectoryFS(d.templateFS, dbms, d.templateData.Agent.ProjectName); err != nil {
		return err
	}

	temp := template.Must(template.New("docker-compose.tmpl").ParseFS(d.templateFS, composeTemplatePath))

	if err := temp.Execute(&content, d.templateData); err != nil {
		return fmt.Errorf("Failed to execute the template: %q, error: %q", "docker-compose.yaml", err)
	}

	fullProviderLocation := d.templateData.Agent.ProjectName + "/"

	// Creates template
	newFile := fullProviderLocation + "docker-compose.yaml"
	f, err := os.Create(newFile)
	if err != nil {
		return fmt.Errorf("Failed to create file: %q, error: %q", "docker-compose.yaml", err)
	}
	defer f.Close()

	if _, err := f.WriteString(content.String()); err != nil {
		return fmt.Errorf("Failed to write to file: %q, error: %q", "docker-compose.yaml", err)
	}

	return nil
}
