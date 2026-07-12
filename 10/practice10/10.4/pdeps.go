package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Package struct {
	ImportPath   string   `json:"ImportPath"`
	Imports      []string `json:"Imports"`
	TestImports  []string `json:"TestImports"`
	XTestImports []string `json:"XTestImports"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <package> [<package>...]\n", os.Args[0])
		os.Exit(1)
	}

	targetPackages := os.Args[1:]

	allPackages := getAllPackages()
	if allPackages == nil {
		os.Exit(1)
	}

	importers := make(map[string][]string)
	for _, pkg := range allPackages {
		for _, imp := range pkg.Imports {
			importers[imp] = append(importers[imp], pkg.ImportPath)
		}

		for _, imp := range pkg.TestImports {
			importers[imp] = append(importers[imp], pkg.ImportPath)
		}
		for _, imp := range pkg.XTestImports {
			importers[imp] = append(importers[imp], pkg.ImportPath)
		}
	}

	dependents := make(map[string]bool)
	for _, target := range targetPackages {
		findDependents(target, importers, dependents)
	}

	for pkg := range dependents {
		fmt.Println(pkg)
	}
}

// getAllPackages возвращает все пакеты в рабочей области
func getAllPackages() []Package {
	cmd := exec.Command("go", "list", "-json", "./...")
	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "go list error: %v\n", err)
		return nil
	}

	var packages []Package
	decoder := json.NewDecoder(strings.NewReader(string(output)))
	for {
		var pkg Package
		if err := decoder.Decode(&pkg); err != nil {
			break
		}
		packages = append(packages, pkg)
	}
	return packages
}

// findDependents рекурсивно находит все пакеты, которые зависят от target
func findDependents(target string, importers map[string][]string, result map[string]bool) {
	for _, importer := range importers[target] {
		if !result[importer] {
			result[importer] = true
			findDependents(importer, importers, result)
		}
	}
}
