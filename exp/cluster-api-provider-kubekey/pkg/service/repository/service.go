/*
 Copyright 2022 The KubeSphere Authors.

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

package repository

import (
	"github.com/kubesphere/kubekey/exp/cluster-api-provider-kubekey/pkg/clients/ssh"
	"github.com/kubesphere/kubekey/exp/cluster-api-provider-kubekey/pkg/scope"
	"github.com/kubesphere/kubekey/exp/cluster-api-provider-kubekey/pkg/service/operation"
	"github.com/kubesphere/kubekey/exp/cluster-api-provider-kubekey/pkg/service/operation/file"
	"github.com/kubesphere/kubekey/exp/cluster-api-provider-kubekey/pkg/service/operation/repository"
	"github.com/kubesphere/kubekey/exp/cluster-api-provider-kubekey/pkg/util/osrelease"
)

// Service holds a collection of interfaces.
// The interfaces are broken down like this to group functions together.
type Service struct {
	sshClient     ssh.Interface
	scope         scope.KKInstanceScope
	instanceScope *scope.InstanceScope

	os        *osrelease.Data
	mountPath string

	repositoryFactory func(sshClient ssh.Interface, os *osrelease.Data) operation.Repository
	isoFactory        func(sshClient ssh.Interface, arch, isoName string) (operation.Binary, error)
}

// NewService returns a new service given the remote instance kubekey build-in repository client.
func NewService(sshClient ssh.Interface, scope scope.KKInstanceScope, instanceScope *scope.InstanceScope) *Service {
	return &Service{
		sshClient:     sshClient,
		scope:         scope,
		instanceScope: instanceScope,
	}
}

func (s *Service) getRepositoryService(os *osrelease.Data) operation.Repository {
	if s.repositoryFactory != nil {
		return s.repositoryFactory(s.sshClient, os)
	}
	return repository.NewService(s.sshClient, os)
}

func (s *Service) getISOService(sshClient ssh.Interface, os *osrelease.Data, arch string, isoName string) (operation.Binary, error) {
	if s.isoFactory != nil {
		return s.isoFactory(sshClient, arch, isoName)
	}
	return file.NewISO(sshClient, s.scope.RootFs(), os, arch, isoName)
}
