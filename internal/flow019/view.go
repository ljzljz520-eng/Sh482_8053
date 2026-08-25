package flow019

import "enterpriselead/internal/domain"

type View struct {
	Record  domain.Record
	Summary string
}

func (b *Board) View() (View, error) {
	record, summary, err := b.Current()
	if err != nil {
		return View{}, err
	}
	return View{Record: record, Summary: summary}, nil
}

func (b *Board) SelectByID(id string) (View, error) {
	for index, record := range b.records {
		if record.ID == id {
			_, summary, err := b.Switch(index)
			if err != nil {
				return View{}, err
			}
			return View{Record: record, Summary: summary}, nil
		}
	}
	return View{}, domainNotFound(id)
}

type notFound string

func (e notFound) Error() string { return "record not found: " + string(e) }

func domainNotFound(id string) error { return notFound(id) }
