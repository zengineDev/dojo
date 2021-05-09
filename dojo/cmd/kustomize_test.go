package cmd

import (
	"fmt"
	"math/rand"
	"testing"
)

func Test_readKustomization(t *testing.T) {
	config := readKustomization("./../../_test/kustomization.yaml")
	if len(config.Images) == 0 {
		t.Error("The file must have a image")
		t.Fail()
	}
}

func Test_writeKostumization(t *testing.T) {
	newVersion := fmt.Sprintf("v0.1.%v", rand.Intn(100))
	configFirst := readKustomization("./../../_test/kustomization.yaml")

	configFirst.Images[0].NewTag = newVersion
	writeKostumization(configFirst, "./../../_test/kustomization.yaml")

	configSecond := readKustomization("./../../_test/kustomization.yaml")
	if configSecond.Images[0].NewTag != newVersion {
		t.Error("The Tag for the images has to be changed")
		t.Fail()
	}
}
