package core

type Instructor interface {
	List()
	Add()
	Remove()
	Rename()
	Edit()
	Copy()
	Reorganize()
	Execute(string)
	Help()
}
