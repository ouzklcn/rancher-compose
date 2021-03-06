package rancher

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/docker/libcompose/config"
	"github.com/docker/libcompose/project"
	"github.com/ouzklcn/rancher-compose/preprocess"
	"io/ioutil"
)

func NewProject(context *Context) (*project.Project, error) {
	context.ServiceFactory = &RancherServiceFactory{
		Context: context,
	}

	context.VolumesFactory = &RancherVolumesFactory{
		Context: context,
	}

	if context.Binding != nil {
		bindingBytes, err := json.Marshal(context.Binding)
		if err != nil {
			return nil, err
		}
		context.BindingsBytes = bindingBytes
	}

	if context.BindingsBytes == nil {
		if context.BindingsFile != "" {
			bindingsContent, err := ioutil.ReadFile(context.BindingsFile)
			if err != nil {
				return nil, err
			}
			context.BindingsBytes = bindingsContent
		}
	}

	preProcessServiceMap := preprocess.PreprocessServiceMap(context.BindingsBytes)
	p := project.NewProject(&context.Context, nil, &config.ParseOptions{
		Interpolate: true,
		Validate:    true,
		Preprocess:  preProcessServiceMap,
	})

	err := p.Parse()
	if err != nil {
		return nil, err
	}

	if err = context.open(); err != nil {
		logrus.Errorf("Failed to open project %s: %v", p.Name, err)
		return nil, err
	}

	p.Name = context.ProjectName

	context.SidekickInfo = NewSidekickInfo(p)

	return p, err
}

func DeleteProject(context *Context) (error) {
	context.ServiceFactory = &RancherServiceFactory{
		Context: context,
	}

	context.VolumesFactory = &RancherVolumesFactory{
		Context: context,
	}

	if context.Binding != nil {
		bindingBytes, err := json.Marshal(context.Binding)
		if err != nil {
			return err
		}
		context.BindingsBytes = bindingBytes
	}

	if context.BindingsBytes == nil {
		if context.BindingsFile != "" {
			bindingsContent, err := ioutil.ReadFile(context.BindingsFile)
			if err != nil {
				return err
			}
			context.BindingsBytes = bindingsContent
		}
	}

	preProcessServiceMap := preprocess.PreprocessServiceMap(context.BindingsBytes)
	p := project.NewProject(&context.Context, nil, &config.ParseOptions{
		Interpolate: true,
		Validate:    true,
		Preprocess:  preProcessServiceMap,
	})


	if err := p.Parse(); err != nil {
		return  err
	}

	if err := context.delete(); err != nil {
		logrus.Errorf("Failed to open project %s: %v", p.Name, err)
		return  err
	}

	return nil
}