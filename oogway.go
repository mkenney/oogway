package main

import (
	"errors"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// oogway masters your checks
type oogway struct {
	CheckInterval time.Duration

	ChecksDir       string
	ChecksExtension string
	checks          map[string]*check
	checksLock      *sync.Mutex
}

// loadChecks reads in all config files and reloads the configuration
func (o *oogway) loadChecks() error {
	checks := make(map[string]*instructions)
	checksInFiles := combineConfigFiles(o.ChecksDir, o.ChecksExtension)

	o.checksLock.Lock()
	defer o.checksLock.Unlock()

	yamlError := yaml.Unmarshal(checksInFiles, &checks)
	if yamlError != nil {
		log.Error("Unable to parse yaml")
		log.Error(yamlError.Error())
		return errors.New("Unable to parse yaml")
	}

	// delete checks that are no longer needed
	for preExistingCheck := range o.checks {
		del := true

		// go through the new checks, make sure it exists
		for newCheck := range checks {
			if newCheck == preExistingCheck {
				del = false
				break
			}
		}

		// do we need to delete it?
		if del {
			delete(o.checks, preExistingCheck)
			log.WithFields(log.Fields{
				"name": preExistingCheck,
			}).Error("Removing check")

		}
	}

	for name, li := range checks {
		// does this check exist?
		if _, exists := o.checks[name]; !exists {
			// ok. lets add it
			o.checks[name] = &check{
				Name:         name,
				Instructions: li,
				ExecLock:     &sync.Mutex{},
			}
			log.WithFields(log.Fields{
				"name": name,
			}).Info("A new check was found")
		} else {
			// kk, something was updated ...
			if o.checks[name].id() != li.id() {
				o.checks[name].Instructions = li
				// reset the logic ...
				o.checks[name].Attempts = 0
				log.WithFields(log.Fields{
					"name": name,
				}).Info("Check was modified")
			}
		}
	}
	return nil
}

// instruct is the central hub for dispatching checks
func (o *oogway) instruct() {
	for {
		o.loadChecks()
		for _, check := range o.checks {
			if check.LastChecked.Add(check.every()).Before(time.Now()) && !check.isMuted() {
				check.check()
			}
		}
		<-time.After(o.CheckInterval)
	}
}
