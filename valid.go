package store

import valid "github.com/asaskevich/govalidator"

// func init() {
// 	valid.SetFieldsRequiredByDefault(true)
// }

func (s *store) Validate(obj interface{}) (bool, error) {
	return valid.ValidateStruct(obj)
}
