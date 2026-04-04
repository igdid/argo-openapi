package cmd

import (
	"os"
	"text/template"
)

type Project struct {
	// v2
	PkgName   string
	Copyright string
	Legal     License
	AppName   string
	rootDoc   *openapi3.T
}

var project Project

// Create OpenAPI protocol files
func (p *Project) createOpenAPIFiles(fn string) {
	p.rootDoc = load(opts.src)
	cfg := codegen.Configuration{
		PackageName: "proto",
		Generate: codegen.GenerateOptions{
			Models: true,
			Client: true,
			//GorillaServer:  true,
			Strict: true,
		},
		OutputOptions: codegen.OutputOptions{
			SkipPrune: true,
		},
	}
	code, err := codegen.Generate(p.rootDoc, cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(opts.target, fn), []byte(code), 0644); err != nil {
		log.Fatal(err)
	}
}

func showParams(method, name string, operation *openapi3.Operation) {
	if operation == nil {
		return
	}
	l := log.WithFields(logrus.Fields{
		"operationID": operation.OperationID,
	})
	if len(operation.Parameters) == 0 {
		l.Info(strings.ToUpper(method), " ", name)
	}
	for _, p := range operation.Parameters {
		l.WithFields(logrus.Fields{
			"param": p.Value.Name,
			"in":    p.Value.In,
		}).Info(strings.ToUpper(method), " ", name)
	}
}

func (p *Project) validateOpenAPI(cmd *cobra.Command, args []string) {
	p.rootDoc = load(opts.src)
	ctx := context.Background()

	if err := p.rootDoc.Validate(ctx); err != nil {
		log.Fatalf("invalid spec: %v", err)
	}
	log.Info("OpenAPI spec is valid. The following methods available:")
	for name, path := range p.rootDoc.Paths.Map() {
		showParams("get", name, path.Get)
		showParams("post", name, path.Post)
		showParams("put", name, path.Put)
		showParams("delete", name, path.Delete)
		showParams("options", name, path.Options)
		showParams("patch", name, path.Patch)
		showParams("trace", name, path.Trace)
	}
}

func isURL(s string) bool {
	u, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}

	return u.Scheme != "" && u.Host != ""
}

func load(src string) *openapi3.T {
	var doc *openapi3.T
	var err error
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	if isURL(src) {
		url, _ := url.ParseRequestURI(src)
		doc, err = loader.LoadFromURI(url)
	} else {
		doc, err = loader.LoadFromFile(src)
	}
	if err != nil {
		log.Fatalf("load spec: %v - %s", err, src)
	}
	if doc.Paths == nil {
		log.Fatal("No paths are inside the file. Nothing to do")
	}
	return doc
}
// Create common files
func generateSensorProject(targetDir string) {
	if !isURL(opts.src) {
		opts.target = filepath.Dir(opts.src)
	} else if opts.target == "" {
		log.Fatal("A target must be specified")
	}
	rootDir := filepath.Join(opts.target, "sensor")
	// Create dir hierarchy
	opts.target = filepath.Join(rootDir, "proto")
	err := os.MkdirAll(opts.target, 0754)
	if err != nil {
		log.Fatal(err)
	}
	cmdDir := filepath.Join(rootDir, "cmd")
	err = os.MkdirAll(cmdDir, 0754)
	if err != nil {
		log.Fatal(err)
	}
	// Create main file
	mainFile, err := os.Create(fmt.Sprintf("%s/main.go", rootDir))
	if err != nil {
		return err
	}
	defer mainFile.Close()

	// Render main template
	mainTemplate := template.Must(template.New("main").Parse(utils.Templates.ReadFile("templates/main.go.tpl")))
	err = mainTemplate.Execute(mainFile, p)
	if err != nil {
		return err
	}
	// Create main file
	rootCmdFile, err := os.Create(fmt.Sprintf("%s/root.go", cmdDir))
	if err != nil {
		return err
	}
	defer rootCmdDir.Close()

	// Render main template
	rootTemplate := template.Must(template.New("root").Parse(utils.Templates.ReadFile("templates/cmd/root.go.tpl")))
	err = rootTemplate.Execute(rootCmdFile, p)
	if err != nil {
		return err
	}

}
