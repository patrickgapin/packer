// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See the LICENSE file in builder/azure for license information.

package arm

import (
	"fmt"
	"testing"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/builder/azure/common/constants"
)

func TestStepValidateTemplateShouldFailIfValidateFails(t *testing.T) {
	var testSubject = &StepValidateTemplate{
		validate: func(string, string) error { return fmt.Errorf("!! Unit Test FAIL !!") },
		say:      func(message string) {},
		error:    func(e error) {},
	}

	stateBag := createTestStateBagStepValidateTemplate()

	var result = testSubject.Run(stateBag)
	if result != multistep.ActionHalt {
		t.Fatalf("Expected the step to return 'ActionHalt', but got '%d'.", result)
	}

	if _, ok := stateBag.GetOk(constants.Error); ok == false {
		t.Fatalf("Expected the step to set stateBag['%s'], but it was not.", constants.Error)
	}
}

func TestStepValidateTemplateShouldPassIfValidatePasses(t *testing.T) {
	var testSubject = &StepValidateTemplate{
		validate: func(string, string) error { return nil },
		say:      func(message string) {},
		error:    func(e error) {},
	}

	stateBag := createTestStateBagStepValidateTemplate()

	var result = testSubject.Run(stateBag)
	if result != multistep.ActionContinue {
		t.Fatalf("Expected the step to return 'ActionContinue', but got '%d'.", result)
	}

	if _, ok := stateBag.GetOk(constants.Error); ok == true {
		t.Fatalf("Expected the step to not set stateBag['%s'], but it was.", constants.Error)
	}
}

func TestStepValidateTemplateShouldTakeStepArgumentsFromStateBag(t *testing.T) {
	var actualResourceGroupName string
	var actualDeploymentName string

	var testSubject = &StepValidateTemplate{
		validate: func(resourceGroupName string, deploymentName string) error {
			actualResourceGroupName = resourceGroupName
			actualDeploymentName = deploymentName

			return nil
		},
		say:   func(message string) {},
		error: func(e error) {},
	}

	stateBag := createTestStateBagStepValidateTemplate()
	var result = testSubject.Run(stateBag)

	if result != multistep.ActionContinue {
		t.Fatalf("Expected the step to return 'ActionContinue', but got '%d'.", result)
	}

	var expectedDeploymentName = stateBag.Get(constants.ArmDeploymentName).(string)
	var expectedResourceGroupName = stateBag.Get(constants.ArmResourceGroupName).(string)

	if actualDeploymentName != expectedDeploymentName {
		t.Fatalf("Expected the step to source 'constants.ArmDeploymentName' from the state bag, but it did not.")
	}

	if actualResourceGroupName != expectedResourceGroupName {
		t.Fatalf("Expected the step to source 'constants.ArmResourceGroupName' from the state bag, but it did not.")
	}
}

func createTestStateBagStepValidateTemplate() multistep.StateBag {
	stateBag := new(multistep.BasicStateBag)

	stateBag.Put(constants.ArmDeploymentName, "Unit Test: DeploymentName")
	stateBag.Put(constants.ArmResourceGroupName, "Unit Test: ResourceGroupName")

	return stateBag
}
