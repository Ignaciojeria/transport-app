package domain

type Account struct {
	Email string
}

func (a Account) DocID() DocumentID {
	return HashInputs(a.Email)
}
