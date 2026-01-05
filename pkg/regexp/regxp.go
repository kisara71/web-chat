package regexp

import regexp "github.com/dlclark/regexp2"

type Handler struct {
	validDatePhone *regexp.Regexp
	validDateEmail *regexp.Regexp
}

func NewHandler() *Handler {
	return &Handler{
		validDatePhone: regexp.MustCompile(`^1[3-9]\d{9}$`, regexp.Compiled),
		validDateEmail: regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, regexp.Compiled),
	}
}
func (h *Handler) ValidatePhone(phone string) (bool, error) {
	return h.validDatePhone.MatchString(phone)
}
func (h *Handler) ValidateEmail(email string) (bool, error) {
	return h.validDateEmail.MatchString(email)
}
