package controllers

type IController interface {
	Prepare() error
	Run() error
}
