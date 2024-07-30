package providers

type QuestionType int

const (
	Picker QuestionType = iota
	Input
)

// A Question holds all the details relating to a question that a provide 
// needs a order to properly generate a working project.
type Question struct {
	// QType represents the type of question.
	QType         QuestionType

	// Prompt is for the question that will be presented to the user.
	Prompt        string

	// DefaultAnswer is the defualt answer for the question. It is only used for
	// the Input Question Type since it's the only question type that accepts an
	// empty answer from the user.
	DefaultAnswer string

	// Options is where all of the available options for the question are kept.
	// This is Only used for the QuestionType Picker
	Options []string

	// Answer is where the answer for the question will be placed.
	Answer string
}

// A RunnableQuestion is a question that can be ran to prompt the user for an
// answer.
//
// Run should be have a *Question receiver and should prompt the user with the 
// question, capture the response and load it into the *Question.Answer field
// or return an error if an error occured.
type RunnableQuestion interface {
	Run() error
}

