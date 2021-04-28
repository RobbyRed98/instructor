package core

type Instructor interface {
	List() error
	Add() error
	Remove() error
	Rename() error
	Edit() error
	Copy() error
	Reorganize() error
	Execute(string) error
	Help()
}
