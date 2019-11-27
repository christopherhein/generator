/*
Copyright © 2019 AWS Controller authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

import (
	"github.com/spf13/afero"

	kbinput "sigs.k8s.io/kubebuilder/pkg/scaffold/input"

	"go.awsctrl.io/generator/pkg/controller"
	"go.awsctrl.io/generator/pkg/controllermanager"
	"go.awsctrl.io/generator/pkg/group"
	"go.awsctrl.io/generator/pkg/stackobject"
	"go.awsctrl.io/generator/pkg/types"

	"go.awsctrl.io/generator/pkg/input"
	"go.awsctrl.io/generator/pkg/resource"
	"go.awsctrl.io/generator/pkg/scaffold"
)

type API struct {
	fs afero.Fs

	// Project loads the project file for adding resources
	project *input.ProjectFile

	// options contains CLI params
	options input.Options
}

// New will generate an API builder
func New(fs afero.Fs, options input.Options) *API {
	return &API{
		fs:      fs,
		options: options,
	}
}

// Build will generate all the
func (a *API) Build(r *resource.Resource, rs []resource.Resource) (err error) {
	var in *input.Input
	if in, err = a.setDefaults(); err != nil {
		return err
	}

	files := []input.File{
		&types.Types{Resource: r, Input: *in, Resources: rs},
		&group.Group{Resource: r, Input: *in, Resources: rs},
		&stackobject.StackObject{Resource: r, Input: *in, Resources: rs},
		&controller.Controller{Resource: r, Input: *in, Resources: rs},
		&controllermanager.ControllerManager{Resource: r, Input: *in, Resources: rs},
	}

	s := scaffold.New(a.fs, r)

	if err := s.Execute(files...); err != nil {
		return err
	}

	return nil
}

func (a *API) setDefaults() (i *input.Input, err error) {
	i = &input.Input{Input: kbinput.Input{
		Domain: "awsctrl.io",
	}}

	var boilerplate string
	if boilerplate, err = a.getBoilerplate(a.options); err != nil {
		return i, err
	}

	i.SetBoilerplate(boilerplate)

	return i, nil
}

func (a *API) getBoilerplate(e input.Options) (string, error) {
	afs := afero.Afero{
		Fs: a.fs,
	}

	b, err := afs.ReadFile(e.BoilerplatePath) // nolint: gosec
	return string(b), err
}
