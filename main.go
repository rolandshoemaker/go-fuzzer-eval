package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

func run(goBin string, targetName string, targetDir string, runtime time.Duration, workers int, num int) (string, error) {
	tmpDir, err := os.MkdirTemp(os.TempDir(), targetName)
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)
	cmd := exec.Command(goBin, "test", ".", "-run=XXXXXXXX", "-fuzz="+targetName, "-timeout=0", "-fuzztime="+runtime.String(), "-parallel="+fmt.Sprint(workers), "-test.fuzzcachedir="+tmpDir, "-v") // -keepfuzzing also, but not added yet?
	cmd.Env = append(os.Environ(), "GODEBUG=fuzzdebug=1")
	cmd.Dir = targetDir
	buf := bytes.NewBuffer(nil)
	cmd.Stderr = buf
	cmd.Stdout = buf
	err = cmd.Run()
	if err != nil {
		fmt.Println(buf.String())
		return "", err
	}
	return buf.String(), nil
}

type fuzzTarget struct {
	Dir  string `yaml:"dir"`
	Name string `yaml:"name"`
}

type config struct {
	WorkersPerRun int           `yaml:"workersPerRun"`
	RuntimePerRun time.Duration `yaml:"runtimePerRun"`
	Runs          int           `yaml:"runs"`
	// ResultsDir    string        `yaml:"resultsDir"`
	// GoBin         string        `yaml:"goBin"`
	Targets []fuzzTarget `yaml:"targets"`
}

func main() {
	experimentPath := flag.String("experiment", "", "Path to YAML file describing experiment to run")
	resultPath := flag.String("result", "", "Path to write results log file to")
	goPath := flag.String("go", "", "Path to Go binary")
	flag.Parse()

	b, err := os.ReadFile(*experimentPath)
	if err != nil {
		log.Fatalf("failed to read config from %q: %s", *experimentPath, err)
	}
	var c config
	if err := yaml.Unmarshal(b, &c); err != nil {
		log.Fatalf("failed to parse config from %q: %s", *experimentPath, err)
	}

	// TODO: validate flags / use a configuration file instead

	maxProcs := runtime.NumCPU()
	parallelRuns := 1
	if c.WorkersPerRun < maxProcs {
		parallelRuns = maxProcs / c.WorkersPerRun
	}

	totalRuns := c.Runs * len(c.Targets)
	work := make(chan func(), totalRuns)
	expectedRuntime := c.RuntimePerRun * time.Duration(totalRuns/parallelRuns)
	if totalRuns%parallelRuns > 0 {
		expectedRuntime += c.RuntimePerRun
	}
	log.Printf("expected runtime %s\n", expectedRuntime)

	results, err := os.OpenFile(*resultPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open results file %q: %s\n", *resultPath, err)
	}
	defer results.Close()
	var resultMu sync.Mutex

	for _, target := range c.Targets {
		for i := 0; i < c.Runs; i++ {
			runNum := i + 1
			t := target
			work <- func() {
				log.Printf("%s#%s run %d/%d started\n", t.Dir, t.Name, runNum, c.Runs)
				s := time.Now()
				l, err := run(*goPath, t.Name, t.Dir, c.RuntimePerRun, c.WorkersPerRun, runNum)
				if err != nil {
					panic(err)
				}
				resultMu.Lock()
				_, err = results.Write([]byte(l + "\n"))
				if err != nil {
					log.Fatalf("failed to write results to %q: %s\n", *resultPath, err)
				}
				resultMu.Unlock()
				if took := time.Since(s); took < c.RuntimePerRun {
					log.Printf("dubious runtime, took: %s, expected: %s\n", took, c.RuntimePerRun)
				}
				log.Printf("%s#%s run %d/%d finished\n", t.Dir, t.Name, runNum, c.Runs)
			}
		}
	}
	close(work)

	wg := new(sync.WaitGroup)
	for i := 0; i < parallelRuns; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for w := range work {
				w()
			}
		}()
	}
	wg.Wait()

	log.Println("done!")
}
