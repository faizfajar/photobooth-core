package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Define the file templates
const (
	repoTemplate = `package repository

import (
	"photobooth-core/internal/domain"
	"gorm.io/gorm"
)

type {{.Name}}Repository interface {
}

type {{.NameLower}}Repository struct {
	db *gorm.DB
}

func New{{.Name}}Repository(db *gorm.DB) {{.Name}}Repository {
	return &{{.NameLower}}Repository{db}
}
`
	usecaseTemplate = `package usecase

import (
	"photobooth-core/internal/{{.NameLower}}/repository"
)

type {{.Name}}Usecase interface {
}

type {{.NameLower}}Usecase struct {
	repo repository.{{.Name}}Repository
}

func New{{.Name}}Usecase(repo repository.{{.Name}}Repository) {{.Name}}Usecase {
	return &{{.NameLower}}Usecase{repo}
}
`
	handlerTemplate = `package handler

import (
	"photobooth-core/internal/{{.NameLower}}/usecase"
	"github.com/gin-gonic/gin"
)

type {{.Name}}Handler struct {
	u usecase.{{.Name}}Usecase
}

func New{{.Name}}Handler(u usecase.{{.Name}}Usecase) *{{.Name}}Handler {
	return &{{.Name}}Handler{u}
}
`
)

type Config struct {
	Name      string
	NameLower string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run scripts/gen.go <module_name>")
		return
	}

	name := strings.Title(os.Args[1])
	nameLower := strings.ToLower(os.Args[1])
	config := Config{Name: name, NameLower: nameLower}

	// 1. Create Folders
	folders := []string{
		filepath.Join("internal", nameLower, "repository"),
		filepath.Join("internal", nameLower, "usecase"),
		filepath.Join("internal", nameLower, "handler"),
	}

	for _, folder := range folders {
		os.MkdirAll(folder, 0755)
	}

	// 2. Create Files
	createFile(filepath.Join("internal", nameLower, "repository", nameLower+"_repository.go"), repoTemplate, config)
	createFile(filepath.Join("internal", nameLower, "usecase", nameLower+"_usecase.go"), usecaseTemplate, config)
	createFile(filepath.Join("internal", nameLower, "handler", nameLower+"_handler.go"), handlerTemplate, config)

	// 3. Create Domain Model (Empty)
	domainPath := filepath.Join("internal", "domain", nameLower+".go")
	if _, err := os.Stat(domainPath); os.IsNotExist(err) {
		os.WriteFile(domainPath, []byte(fmt.Sprintf("package domain\n\ntype %s struct {\n}\n", name)), 0644)
	}

	fmt.Printf("Successfully scaffolded module: %s\n", name)
}

func createFile(path, tpl string, config Config) {
	f, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	t := template.Must(template.New("tpl").Parse(tpl))
	t.Execute(f, config)
}
