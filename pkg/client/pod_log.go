/*
Copyright 2023 The Kubernetes Authors.

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

package client

import (
	"bufio"

	"github.com/dims/hydrophone/pkg/service"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// List pod resource with the given namespace
func getPodLogs(clientset *kubernetes.Clientset, stream streamLogs) {
	podLogOpts := v1.PodLogOptions{
		Container: "conformance-container",
		Follow:    true,
	}

	req := clientset.CoreV1().Pods(service.Namespace).GetLogs(service.PodName, &podLogOpts)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		stream.errCh <- err
	}
	defer podLogs.Close()

	reader := bufio.NewScanner(podLogs)

	for reader.Scan() {
		line := reader.Text()
		stream.logCh <- line + "\n"
	}
	stream.doneCh <- true
}
