package domain

type Account struct {
	Email string
	Role  string
}

func (a Account) DocID() DocumentID {
	return HashInputs(a.Email)
}
