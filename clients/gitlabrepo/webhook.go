// Copyright 2022 OpenSSF Scorecard Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gitlabrepo

import (
	"fmt"
	"sync"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/ossf/scorecard/v5/clients"
)

type webhookHandler struct {
	glClient *gitlab.Client
	once     *sync.Once
	errSetup error
	repourl  *Repo
	webhooks []clients.Webhook
}

func (handler *webhookHandler) init(repourl *Repo) {
	handler.repourl = repourl
	handler.errSetup = nil
	handler.once = new(sync.Once)
}

func (handler *webhookHandler) setup() error {
	handler.once.Do(func() {
		projectHooks, _, err := handler.glClient.Projects.ListProjectHooks(
			handler.repourl.projectID, &gitlab.ListProjectHooksOptions{})
		if err != nil {
			handler.errSetup = fmt.Errorf("request for project hooks failed with %w", err)
			return
		}

		// TODO: make sure that enablesslverification is similarly equivalent to auth secret.
		for _, hook := range projectHooks {
			handler.webhooks = append(handler.webhooks,
				clients.Webhook{
					Path:           hook.URL,
					ID:             int64(hook.ID),
					UsesAuthSecret: hook.EnableSSLVerification,
				})
		}
	})

	return handler.errSetup
}

func (handler *webhookHandler) listWebhooks() ([]clients.Webhook, error) {
	if err := handler.setup(); err != nil {
		return nil, fmt.Errorf("error during webhookHandler.setup: %w", err)
	}

	return handler.webhooks, nil
}
