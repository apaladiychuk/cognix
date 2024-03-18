package main

import "cognix.ch/api/v2/core/repository"

type Conductor struct {
	connectorRepo repository.ConnectorRepository
}

func NewConductor(connectorRepo repository.ConnectorRepository) *Conductor {
	return &Conductor{
		connectorRepo: connectorRepo,
	}
}

func (c *Conductor) Start() {

}
