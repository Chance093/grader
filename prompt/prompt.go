package prompt

import "github.com/manifoldco/promptui"

func List(label string, list []string) (string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: list,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

func Input(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}
