package container

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	api "github.com/fsouza/go-dockerclient"
)

type Connector struct {
	client     *api.Client
	containers map[string]*Container
	mu         sync.Mutex
}

func NewDocker() *Connector {
	c, err := api.NewClientFromEnv()
	if err != nil {
		panic(err)
	}
	return &Connector{
		client:     c,
		containers: map[string]*Container{},
	}
}

func (c *Connector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for id, cont := range c.containers {
		cont.Stop()
		delete(c.containers, id)
	}
}

func (c *Connector) All() []*Container {
	c.mu.Lock()
	defer c.mu.Unlock()

	lca, err := c.client.ListContainers(api.ListContainersOptions{
		All:     true,
		Filters: map[string][]string{"status": []string{"running"}},
	})
	if err != nil {
		// FIXME return err to frontend (cannot connect to Docker endpoint)
		fmt.Println(err)
		return c.clonedContainers()
	}

	seen := map[string]struct{}{}
	for _, container := range lca {
		seen[container.ID] = struct{}{}

		// Append new containers
		if _, exist := c.containers[container.ID]; !exist {
			ci, err := c.client.InspectContainer(container.ID)
			if err != nil {
				if _, gone := err.(*api.NoSuchContainer); gone {
					delete(c.containers, container.ID)
				} else {
					fmt.Println(err)
				}
			}

			nc := &Container{
				conn:      c,
				ID:        ci.ID,
				Name:      strings.Replace(ci.Name, "/", "", 1),
				StartedAt: ci.State.StartedAt,
				Status:    ci.State.Status,
			}
			nc.Start()

			c.containers[container.ID] = nc
		}
	}

	// Delete old containers
	for id, cont := range c.containers {
		if _, ok := seen[id]; !ok {
			cont.Stop()
			delete(c.containers, id)
		}
	}

	return c.clonedContainers()
}

func (c *Connector) clonedContainers() []*Container {
	cs := []*Container{}

	for _, v := range c.containers {
		cs = append(cs, v.Clone())
	}

	sort.Slice(cs, func(i, j int) bool {
		return cs[i].StartedAt.Before(cs[j].StartedAt)
	})

	return cs
}
