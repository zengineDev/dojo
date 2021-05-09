package cmd

import (
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func init() {
	rootCmd.AddCommand(kustomizeCmd)
}

type kustomizeConfig struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Images     []struct {
		Name    string `yaml:"name"`
		NewName string `yaml:"newName"`
		NewTag  string `yaml:"newTag"`
	} `yaml:"images"`
	Resources       []string      `yaml:"resources"`
	PatchesJson6902 []interface{} `yaml:"patchesJson6902"`
}

var kustomizeCmd = &cobra.Command{
	Use:   "kustomize",
	Short: "",
	Long:  `All software has versions. This is Hugo's`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		stage := args[0]
		tag := args[1]

		switch stage {
		case "staging":
			path := "./deploy/k8s/staging/kustomization.yaml"
			config := readKustomization(path)
			for _, i := range config.Images {
				i.NewTag = tag
			}
			writeKostumization(config, path)
		case "production":
			path := "./deploy/k8s/production/kustomization.yaml"
			config := readKustomization(path)
			for _, i := range config.Images {
				i.NewTag = tag
			}
			writeKostumization(config, path)
		}
	},
}

func readKustomization(path string) kustomizeConfig {
	var config kustomizeConfig
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return config
}

func writeKostumization(conf kustomizeConfig, path string) {
	d, err := yaml.Marshal(&conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = ioutil.WriteFile(path, d, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
