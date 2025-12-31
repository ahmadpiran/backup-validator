package validator

type Validator interface {
	validate(backupPath string) error
}
