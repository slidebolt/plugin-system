package main

import "github.com/slidebolt/sdk-types"

func systemDomains() []types.DomainDescriptor {
	return []types.DomainDescriptor{
		{
			Domain:   "sensor.time",
			Commands: []types.ActionDescriptor{},
			Events: []types.ActionDescriptor{
				{
					Action: "tick",
					Fields: []types.FieldDescriptor{
						{Name: "value", Type: "string", Required: true},
						{Name: "ts", Type: "string", Required: true},
					},
				},
			},
		},
		{
			Domain:   "sensor.date",
			Commands: []types.ActionDescriptor{},
			Events: []types.ActionDescriptor{
				{
					Action: "tick",
					Fields: []types.FieldDescriptor{
						{Name: "value", Type: "string", Required: true},
						{Name: "ts", Type: "string", Required: true},
					},
				},
			},
		},
		{
			Domain:   "sensor.cpu",
			Commands: []types.ActionDescriptor{},
			Events: []types.ActionDescriptor{
				{
					Action: "tick",
					Fields: []types.FieldDescriptor{
						{Name: "percent", Type: "number", Required: true},
						{Name: "ts", Type: "string", Required: true},
					},
				},
			},
		},
	}
}
