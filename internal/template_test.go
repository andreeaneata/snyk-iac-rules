package internal

import (
	"errors"
	"strings"
	"testing"

	"github.com/snyk/snyk-iac-rules/util"
	"github.com/stretchr/testify/assert"
)

func mockTemplateParams() *TemplateCommandParams {
	return &TemplateCommandParams{
		RuleID:       "Test Rule ID",
		RuleTitle:    "Test Rule Title",
		RuleSeverity: util.NewEnumFlag(LOW, []string{LOW, MEDIUM, HIGH, CRITICAL}),
	}
}

var directories = []struct {
	workingDirectory string
	name             string
}{
	{
		workingDirectory: "test",
		name:             "rules",
	},
	{
		workingDirectory: "test/rules",
		name:             "Test Rule ID",
	},
	{
		workingDirectory: "test/rules/Test Rule ID",
		name:             "fixtures",
	},
	{
		workingDirectory: "test",
		name:             "lib",
	},
	{
		workingDirectory: "test/lib",
		name:             "testing",
	},
}

var files = []struct {
	workingDirectory string
	name             string
	template         string
}{
	{
		workingDirectory: "test/rules/Test Rule ID",
		name:             "main.rego",
		template:         "templates/main.tpl.rego",
	},
	{
		workingDirectory: "test/rules/Test Rule ID",
		name:             "main_test.rego",
		template:         "templates/main_test.tpl.rego",
	},
	{
		workingDirectory: "test/rules/Test Rule ID/fixtures",
		name:             "allowed.json",
		template:         "templates/fixtures/allowed.json",
	},
	{
		workingDirectory: "test/rules/Test Rule ID/fixtures",
		name:             "denied1.yaml",
		template:         "templates/fixtures/denied1.yaml",
	},
	{
		workingDirectory: "test/rules/Test Rule ID/fixtures",
		name:             "denied2.tf",
		template:         "templates/fixtures/denied2.tf",
	},
	{
		workingDirectory: "test/rules/Test Rule ID/fixtures",
		name:             "denied.json.tfplan",
		template:         "templates/fixtures/denied.json.tfplan",
	},
	{
		workingDirectory: "test/lib",
		name:             "main.rego",
		template:         "templates/lib/main.tpl.rego",
	},
	{
		workingDirectory: "test/lib/testing",
		name:             "main.rego",
		template:         "templates/lib/testing/main.tpl.rego",
	},
	{
		workingDirectory: "test/lib/testing",
		name:             "tfplan.rego",
		template:         "templates/lib/testing/tfplan.tpl.rego",
	},
}

func TestTemplateInEmptyDirectory(t *testing.T) {
	directoriesIndex := 0
	oldCreateDirectory := createDirectory
	defer func() {
		createDirectory = oldCreateDirectory
	}()
	createDirectory = func(workingDirectory string, name string, strict bool) (string, error) {
		if directoriesIndex >= len(directories) {
			return "", errors.New("Tried to create more directories than expected")
		}
		assert.Equal(t, directories[directoriesIndex].workingDirectory, workingDirectory)
		assert.Equal(t, directories[directoriesIndex].name, name)
		directoriesIndex++
		return workingDirectory + "/" + name, nil
	}

	filesIndex := 0
	oldTemplateFile := templateFile
	defer func() {
		templateFile = oldTemplateFile
	}()
	templateFile = func(workingDirectory string, name string, template string, templating util.Templating) error {
		if filesIndex >= len(files) {
			return errors.New("Tried to create more files than expected")
		}
		assert.Equal(t, files[filesIndex].workingDirectory, workingDirectory)
		assert.Equal(t, files[filesIndex].name, name)
		assert.Equal(t, files[filesIndex].template, template)
		assert.Equal(t, "Test Rule ID", templating.RuleID)
		assert.Equal(t, "Test Rule Title", templating.RuleTitle)
		assert.Equal(t, LOW, templating.RuleSeverity)
		filesIndex++
		return nil
	}

	oldtTPlanExists := tfPlanExists
	defer func() {
		tfPlanExists = oldtTPlanExists
	}()
	tfPlanExists = func(string) bool {
		return false
	}

	templateParams := mockTemplateParams()
	err := RunTemplate([]string{"test"}, templateParams)
	assert.Nil(t, err)
	assert.Equal(t, len(directories), directoriesIndex)
	assert.Equal(t, len(files), filesIndex)
}

func TestTemplateInDirectoryWithLibWithTfPlan(t *testing.T) {
	directoriesIndex := 0
	oldCreateDirectory := createDirectory
	defer func() {
		createDirectory = oldCreateDirectory
	}()
	createDirectory = func(workingDirectory string, name string, strict bool) (string, error) {
		if directoriesIndex >= len(directories) {
			return "", errors.New("Tried to create more directories than expected")
		}
		if name == "lib" || strings.Contains(workingDirectory, "lib") {
			return "", errors.New("Directory already exists at location")
		}
		assert.Equal(t, directories[directoriesIndex].workingDirectory, workingDirectory)
		assert.Equal(t, directories[directoriesIndex].name, name)
		directoriesIndex++
		return workingDirectory + "/" + name, nil
	}

	filesIndex := 0
	oldTemplateFile := templateFile
	defer func() {
		templateFile = oldTemplateFile
	}()
	templateFile = func(workingDirectory string, name string, template string, templating util.Templating) error {
		if filesIndex >= len(files) {
			return errors.New("Tried to create more files than expected")
		}
		assert.Equal(t, files[filesIndex].workingDirectory, workingDirectory)
		assert.Equal(t, files[filesIndex].name, name)
		assert.Equal(t, files[filesIndex].template, template)
		assert.Equal(t, "Test Rule ID", templating.RuleID)
		assert.Equal(t, "Test Rule Title", templating.RuleTitle)
		assert.Equal(t, LOW, templating.RuleSeverity)
		filesIndex++
		return nil
	}

	oldtTPlanExists := tfPlanExists
	defer func() {
		tfPlanExists = oldtTPlanExists
	}()
	tfPlanExists = func(string) bool {
		return true
	}

	templateParams := mockTemplateParams()
	err := RunTemplate([]string{"test"}, templateParams)
	assert.Nil(t, err)
	assert.Equal(t, len(directories)-2, directoriesIndex)
	assert.Equal(t, len(files)-3, filesIndex)
}

func TestTemplateInDirectoryWithLibWithoutTfPlan(t *testing.T) {
	directoriesIndex := 0
	oldCreateDirectory := createDirectory
	defer func() {
		createDirectory = oldCreateDirectory
	}()
	createDirectory = func(workingDirectory string, name string, strict bool) (string, error) {
		if directoriesIndex >= len(directories) {
			return "", errors.New("Tried to create more directories than expected")
		}
		if name == "lib" || strings.Contains(workingDirectory, "lib") {
			return "", errors.New("Directory already exists at location")
		}
		assert.Equal(t, directories[directoriesIndex].workingDirectory, workingDirectory)
		assert.Equal(t, directories[directoriesIndex].name, name)
		directoriesIndex++
		return workingDirectory + "/" + name, nil
	}

	filesIndex := 0
	oldTemplateFile := templateFile
	defer func() {
		templateFile = oldTemplateFile
	}()
	templateFile = func(workingDirectory string, name string, template string, templating util.Templating) error {
		if filesIndex >= len(files) {
			return errors.New("Tried to create more files than expected")
		}

		// if creating the tfplan testing file then change its order
		var oldFilesIndex int
		if strings.Contains(name, "tfplan.rego") {
			oldFilesIndex = filesIndex
			filesIndex = len(files) - 1
		}

		assert.Equal(t, files[filesIndex].workingDirectory, workingDirectory)
		assert.Equal(t, files[filesIndex].name, name)
		assert.Equal(t, files[filesIndex].template, template)
		assert.Equal(t, "Test Rule ID", templating.RuleID)
		assert.Equal(t, "Test Rule Title", templating.RuleTitle)
		assert.Equal(t, LOW, templating.RuleSeverity)

		// if creating the tfplan testing file then change its order
		if strings.Contains(workingDirectory, "/lib/testing") {
			filesIndex = oldFilesIndex
		}
		filesIndex++

		return nil
	}

	oldtTPlanExists := tfPlanExists
	defer func() {
		tfPlanExists = oldtTPlanExists
	}()
	tfPlanExists = func(string) bool {
		return false
	}

	templateParams := mockTemplateParams()
	err := RunTemplate([]string{"test"}, templateParams)
	assert.Nil(t, err)
	assert.Equal(t, len(directories)-2, directoriesIndex)
	assert.Equal(t, len(files)-2, filesIndex)
}

func TestTemplateWithExistingRule(t *testing.T) {
	directoriesIndex := 0
	oldCreateDirectory := createDirectory
	defer func() {
		createDirectory = oldCreateDirectory
	}()
	createDirectory = func(workingDirectory string, name string, strict bool) (string, error) {
		if directoriesIndex >= len(directories) {
			return "", errors.New("Tried to create more directories than expected")
		}
		if name == "Test Rule ID" {
			return "", errors.New("Directory already exists at location")
		}
		assert.Equal(t, directories[directoriesIndex].workingDirectory, workingDirectory)
		assert.Equal(t, directories[directoriesIndex].name, name)
		directoriesIndex++
		return workingDirectory + "/" + name, nil
	}

	templateFile = func(workingDirectory string, name string, template string, templating util.Templating) error {
		return errors.New("This function should not be called")
	}

	templateParams := mockTemplateParams()
	err := RunTemplate([]string{"test"}, templateParams)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Rule with the provided name already exists")
	assert.Equal(t, 1, directoriesIndex)
}
