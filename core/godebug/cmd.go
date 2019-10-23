package godebug

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"go/ast"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/jmigpin/editor/core/godebug/debug"
	"github.com/jmigpin/editor/util/goutil"
	"github.com/jmigpin/editor/util/osutil"
)

type Cmd struct {
	Client    *Client
	NoModules bool // not in go.mod modules

	Dir    string // "" will use current dir
	Stdout io.Writer
	Stderr io.Writer

	annset *AnnotatorSet

	tmpDir       string
	tmpBuiltFile string // file built and exec'd

	start struct {
		cancel    context.CancelFunc
		waitg     sync.WaitGroup
		serverErr error
	}

	flags struct {
		mode struct {
			run     bool
			test    bool
			build   bool
			connect bool
		}
		verbose   bool
		filename  string
		work      bool
		dirs      []string
		files     []string
		address   string   // build/connect
		env       []string // build
		otherArgs []string
		runArgs   []string
	}
}

func NewCmd() *Cmd {
	cmd := &Cmd{
		annset: NewAnnotatorSet(),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if os.Getenv("GO111MODULE") == "off" {
		cmd.NoModules = true
	}
	return cmd
}

//------------

func (cmd *Cmd) Printf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(cmd.Stdout, format, a...)
}

//------------

func (cmd *Cmd) Start(ctx context.Context, args []string) (done bool, _ error) {
	// parse arguments
	done, err := cmd.parseArgs(args)
	if done || err != nil {
		return done, err
	}

	// absolute dir
	if u, err := filepath.Abs(cmd.Dir); err == nil {
		cmd.Dir = u
	}

	if cmd.flags.verbose {
		cmd.Printf("nomodules=%v\n", cmd.NoModules)
	}

	// tmp dir for building
	if cmd.NoModules {
		d := "editor_godebug_gopath_work"
		tmpDir, err := ioutil.TempDir(os.TempDir(), d)
		if err != nil {
			return true, err
		}
		cmd.tmpDir = tmpDir
	} else {
		d := "editor_godebug_mod_work"

		// The fixed directory will improve the file sync performance since modules require the whole directory to be there (not like gopath)
		// TODO: will have problems running more then one debug session in different editor sessions
		//d += "_"+md5.Sum([]byte(cmd.Dir))

		//cmd.tmpDir = filepath.Join(os.TempDir(), d)

		tmpDir, err := ioutil.TempDir(os.TempDir(), d)
		if err != nil {
			return true, err
		}
		cmd.tmpDir = tmpDir
	}

	// print tmp dir if work flag is present
	if cmd.flags.work {
		fmt.Fprintf(cmd.Stdout, "work: %v\n", cmd.tmpDir)
	}

	m := &cmd.flags.mode

	if m.run || m.test || m.build {
		setupServerNetAddr(cmd.flags.address)
		err := cmd.initAndAnnotate(ctx)
		if err != nil {
			return true, err
		}
	}

	// just building: inform the address used in the binary
	if m.build {
		fmt.Fprintf(cmd.Stdout, "build: %v (builtin address: %v, %v)",
			cmd.tmpBuiltFile,
			debug.ServerNetwork,
			debug.ServerAddress,
		)
		return true, err
	}

	if m.run || m.test || m.connect {
		err = cmd.startServerClient(ctx)
		return false, err
	}

	return false, nil
}

//------------

func (cmd *Cmd) Wait() error {
	cmd.start.waitg.Wait()
	cmd.start.cancel() // ensure resources are cleared
	return cmd.start.serverErr
}

//------------

func (cmd *Cmd) initAndAnnotate(ctx context.Context) error {
	files := NewFiles(cmd.annset.FSet)
	files.Dir = cmd.Dir

	files.Add(cmd.flags.files...)
	files.Add(cmd.flags.dirs...)

	mainFilename := files.absFilename(cmd.flags.filename)

	err := files.Do(ctx, mainFilename, cmd.flags.mode.test, cmd.NoModules)
	if err != nil {
		return err
	}

	// pre-build without annotations for better errors (result is ignored)
	if err := cmd.preBuild(ctx, mainFilename, cmd.flags.mode.test); err != nil {
		return err
	}

	if cmd.flags.verbose {
		files.verbose(cmd)
	}

	// copy
	for filename := range files.copyFilenames {
		dst := cmd.tmpDirBasedFilename(filename)
		if err := mkdirAllCopyFileSync(filename, dst); err != nil {
			return err
		}
	}
	for filename := range files.modFilenames {
		dst := cmd.tmpDirBasedFilename(filename)
		if err := mkdirAllCopyFileSync(filename, dst); err != nil {
			return err
		}
	}

	// annotate
	for filename := range files.annFilenames {
		dst := cmd.tmpDirBasedFilename(filename)
		typ := files.annTypes[filename]
		astFile, err := files.fullAstFile(filename)
		if err != nil {
			return err
		}
		if err := cmd.annset.AnnotateAstFile(astFile, typ); err != nil {
			return err
		}
		if err := cmd.mkdirAllWriteAstFile(dst, astFile); err != nil {
			return err
		}
	}

	// write config file after annotations
	if err := cmd.writeGoDebugConfigFilesToTmpDir(); err != nil {
		return err
	}

	// create testmain file
	if cmd.flags.mode.test && !cmd.annset.InsertedExitIn.TestMain {
		if err := cmd.writeTestMainFilesToTmpDir(); err != nil {
			return err
		}
	}

	// main must have exit inserted
	if !cmd.flags.mode.test && !cmd.annset.InsertedExitIn.Main {
		return fmt.Errorf("have not inserted debug exit in main()")
	}

	if !cmd.NoModules {
		if err := SetupGoMods(ctx, cmd, files, mainFilename, cmd.flags.mode.test); err != nil {
			return err
		}
	}

	return cmd.doBuild(ctx, mainFilename, cmd.flags.mode.test)
}

func (cmd *Cmd) doBuild(ctx context.Context, mainFilename string, tests bool) error {
	filename := cmd.filenameForBuild(mainFilename, tests)
	filenameAtTmp := cmd.tmpDirBasedFilename(filename)

	// create parent dirs
	if err := os.MkdirAll(filepath.Dir(filenameAtTmp), 0755); err != nil {
		return err
	}

	// build
	filenameAtTmpOut, err := cmd.runBuildCmd(ctx, filenameAtTmp, tests)
	if err != nil {
		return err
	}

	// move filename to working dir
	filenameWork := filepath.Join(cmd.Dir, filepath.Base(filenameAtTmpOut))
	if err := os.Rename(filenameAtTmpOut, filenameWork); err != nil {
		return err
	}

	// keep moved filename that will run in working dir for later cleanup
	cmd.tmpBuiltFile = filenameWork

	return nil
}

func (cmd *Cmd) filenameForBuild(mainFilename string, tests bool) string {
	if tests {
		// final filename will include extension replacement with "_godebug"
		return filepath.Join(cmd.Dir, "pkgtest")
	}
	return mainFilename
}

// pre-build without annotations for better errors (result is ignored)
func (cmd *Cmd) preBuild(ctx context.Context, mainFilename string, tests bool) error {
	filename := cmd.filenameForBuild(mainFilename, tests)
	filenameOut, err := cmd.runBuildCmd(ctx, filename, tests)
	if err != nil {
		return err
	}
	os.Remove(filenameOut)
	return nil
}

//------------

func (cmd *Cmd) startServerClient(ctx context.Context) error {
	filenameWork := cmd.tmpBuiltFile

	// server/client context to cancel the other when one of them ends
	ctx2, cancel := context.WithCancel(ctx)
	cmd.start.cancel = cancel

	// arguments
	filenameWork2 := normalizeFilenameForExec(filenameWork)
	args := []string{filenameWork2}
	if cmd.flags.mode.test {
		args = append(args, cmd.flags.runArgs...)
	} else {
		args = append(args, cmd.flags.otherArgs...)
	}

	// start server
	var serverCmd *exec.Cmd
	if !cmd.flags.mode.connect {
		u, err := cmd.startCmd(ctx2, cmd.Dir, args, nil)
		if err != nil {
			// cmd.Wait() won't be called, need to clear resources
			cmd.start.cancel()
			return err
		}
		serverCmd = u

		// output cmd pid
		fmt.Fprintf(cmd.Stdout, "# pid %d\n", serverCmd.Process.Pid)
	}

	// setup address to connect to
	if cmd.flags.mode.connect && cmd.flags.address != "" {
		debug.ServerNetwork = "tcp"
		debug.ServerAddress = cmd.flags.address
	}
	// start client (blocking connect)
	client, err := NewClient(ctx2)
	if err != nil {
		// cmd.Wait() won't be called, need to clear resources
		cmd.start.cancel()
		return err
	}
	cmd.Client = client

	// from this point, cmd.Wait() clears resources from cmd.start.cancel

	// server done
	if serverCmd != nil {
		cmd.start.waitg.Add(1)
		go func() {
			defer cmd.start.waitg.Done()
			// wait for server to finish
			cmd.start.serverErr = serverCmd.Wait()
		}()
	}

	// client done
	cmd.start.waitg.Add(1)
	go func() {
		defer cmd.start.waitg.Done()
		cmd.Client.Wait() // wait for client to finish
	}()

	return nil
}

//------------

func (cmd *Cmd) RequestFileSetPositions() error {
	msg := &debug.ReqFilesDataMsg{}
	encoded, err := debug.EncodeMessage(msg)
	if err != nil {
		return err
	}
	_, err = cmd.Client.Conn.Write(encoded)
	return err
}

func (cmd *Cmd) RequestStart() error {
	msg := &debug.ReqStartMsg{}
	encoded, err := debug.EncodeMessage(msg)
	if err != nil {
		return err
	}
	_, err = cmd.Client.Conn.Write(encoded)
	return err
}

//------------

func (cmd *Cmd) tmpDirBasedFilename(filename string) string {
	// remove volume name
	v := filepath.VolumeName(filename)
	if len(v) > 0 {
		filename = filename[len(v):]
	}
	if cmd.NoModules {
		// trim filename when inside a src dir
		_, rest := goutil.ExtractSrcDir(filename)
		return filepath.Join(cmd.tmpDir, "src", rest)
	}
	return filepath.Join(cmd.tmpDir, filename)
}

//------------

func (cmd *Cmd) environ() []string {
	env := os.Environ()
	if cmd.NoModules {
		// gopath
		env = append(env, cmd.environGoPath())
	}
	// add cmd line env vars
	for _, s := range cmd.flags.env {
		env = append(env, s)
	}
	return env
}

func (cmd *Cmd) environGoPath() string {
	if !cmd.NoModules {
		panic("must be in nomodules mode")
	}

	goPath := []string{}
	// Add a fixed gopath directory in a temporary location for caching downloaded modules packages for godebug.
	// This solves downloading modules every time a godebug session is started since the default directory for downloading is the first directory defined in the GOPATH env.
	cacheTmpDir := filepath.Join(os.TempDir(), "editor_godebug_gopath_cache")
	goPath = append(goPath, cacheTmpDir)

	// add tmpdir to gopath to allow the compiler to give priority to the annotated files
	goPath = append(goPath, cmd.tmpDir)

	// add already defined gopath
	goPath = append(goPath, goutil.GoPath()...)
	// build gopath string
	return "GOPATH=" + strings.Join(goPath, string(os.PathListSeparator))
}

//------------

func (cmd *Cmd) Cleanup() {
	// cleanup unix socket in case of bad stop
	if debug.ServerNetwork == "unix" {
		_ = os.Remove(debug.ServerAddress)
	}

	if cmd.flags.work {
		// don't cleanup work dir
	} else {
		if cmd.tmpDir != "" {
			_ = os.RemoveAll(cmd.tmpDir)
		}
	}

	if cmd.tmpBuiltFile != "" && !cmd.flags.mode.build {
		_ = os.Remove(cmd.tmpBuiltFile)
	}
}

//------------

func (cmd *Cmd) runBuildCmd(ctx context.Context, filename string, tests bool) (string, error) {
	filenameOut := cmd.execName(filename)

	args := []string{}
	if tests {
		args = []string{
			osutil.GoExec(), "test",
			"-c", // compile binary but don't run
			// TODO: faster dummy pre-builts?
			// "-toolexec", "", // don't run asm?
			"-o", filenameOut,
		}
	} else {
		args = []string{
			osutil.GoExec(), "build",
			"-o", filenameOut,
			filename,
		}
	}

	// append otherargs in test mode
	if tests {
		args = append(args, cmd.flags.otherArgs...)
	}

	dir := filepath.Dir(filenameOut)
	if cmd.flags.verbose {
		cmd.Printf("runBuildCmd: dir=%v\n", dir)
	}
	err := cmd.runCmd(ctx, dir, args, cmd.environ())
	if err != nil {
		err = fmt.Errorf("runBuildCmd: %v", err)
	}
	return filenameOut, err
}

func (cmd *Cmd) execName(name string) string {
	return replaceExt(name, osutil.ExecName("_godebug"))
}

//------------

func (cmd *Cmd) runCmd(ctx context.Context, dir string, args, env []string) error {
	// ctx with early cancel for startcmd to clear inner goroutine resource
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()

	ecmd, err := cmd.startCmd(ctx2, dir, args, env)
	if err != nil {
		return err
	}
	return ecmd.Wait()
}

func (cmd *Cmd) startCmd(ctx context.Context, dir string, args, env []string) (*exec.Cmd, error) {
	cargs := osutil.ShellRunArgs(args...)
	ecmd := exec.CommandContext(ctx, cargs[0], cargs[1:]...)

	ecmd.Env = env
	ecmd.Dir = dir
	ecmd.Stdout = cmd.Stdout
	ecmd.Stderr = cmd.Stderr
	osutil.SetupExecCmdSysProcAttr(ecmd)

	if err := ecmd.Start(); err != nil {
		return nil, err
	}

	// ensure kill to child processes on context cancel
	// the ctx must be cancelable, otherwise it might kill the process on start
	go func() {
		select {
		case <-ctx.Done():
			_ = osutil.KillExecCmd(ecmd)
		}
	}()

	return ecmd, nil
}

//------------

func (cmd *Cmd) mkdirAllWriteAstFile(filename string, astFile *ast.File) error {
	buf := &bytes.Buffer{}
	if err := cmd.annset.Print(buf, astFile); err != nil {
		return err
	}
	return mkdirAllWriteFile(filename, buf.Bytes())
}

//------------

func (cmd *Cmd) writeGoDebugConfigFilesToTmpDir() error {
	// godebugconfig pkg: config.go
	filename := GoDebugConfigFilepathName("config.go")
	src := cmd.annset.ConfigContent()
	filenameAtTmp := cmd.tmpDirBasedFilename(filename)
	if err := mkdirAllWriteFile(filenameAtTmp, []byte(src)); err != nil {
		return err
	}
	if !cmd.NoModules {
		// godebugconfig pkg: go.mod
		filename2 := GoDebugConfigFilepathName("go.mod")
		src2 := cmd.annset.ConfigGoModuleContent()
		filenameAtTmp2 := cmd.tmpDirBasedFilename(filename2)
		if err := mkdirAllWriteFile(filenameAtTmp2, []byte(src2)); err != nil {
			return err
		}
	}
	return nil
}

func (cmd *Cmd) writeTestMainFilesToTmpDir() error {
	u := cmd.annset.TestMainSources()
	for i, tms := range u {
		name := fmt.Sprintf("godebug_testmain%v_test.go", i)
		filename := filepath.Join(tms.Dir, name)
		filenameAtTmp := cmd.tmpDirBasedFilename(filename)
		return mkdirAllWriteFile(filenameAtTmp, []byte(tms.Src))
	}
	return nil
}

//------------

func (cmd *Cmd) parseArgs(args []string) (done bool, _ error) {
	if len(args) > 0 {
		switch args[0] {
		case "run":
			cmd.flags.mode.run = true
			return cmd.parseRunArgs(args[1:])
		case "test":
			cmd.flags.mode.test = true
			return cmd.parseTestArgs(args[1:])
		case "build":
			cmd.flags.mode.build = true
			return cmd.parseBuildArgs(args[1:])
		case "connect":
			cmd.flags.mode.connect = true
			return cmd.parseConnectArgs(args[1:])
		}
	}
	fmt.Fprint(cmd.Stderr, cmdUsage())
	return true, nil
}

func (cmd *Cmd) parseRunArgs(args []string) (done bool, _ error) {
	f := &flag.FlagSet{}
	f.BoolVar(&cmd.flags.work, "work", false, "print workdir and don't cleanup on exit")
	f.BoolVar(&cmd.flags.verbose, "verbose", false, "verbose godebug")
	dirs := f.String("dirs", "", "comma-separated list of directories")
	files := f.String("files", "", "comma-separated list of files to avoid annotating big directories")
	env := f.String("env", "", "set env variables, separated by comma (ex: \"GOOS=linux,...\"'")

	if err := f.Parse(args); err != nil {
		if err == flag.ErrHelp {
			f.SetOutput(cmd.Stderr)
			f.PrintDefaults()
			return true, nil
		}
		return true, err
	}

	cmd.flags.dirs = splitCommaList(*dirs)
	cmd.flags.files = splitCommaList(*files)
	cmd.flags.env = strings.Split(*env, ",")
	cmd.flags.otherArgs = f.Args()

	if len(cmd.flags.otherArgs) > 0 {
		cmd.flags.filename = cmd.flags.otherArgs[0]
		cmd.flags.otherArgs = cmd.flags.otherArgs[1:]
	}

	return false, nil
}

func (cmd *Cmd) parseTestArgs(args []string) (done bool, _ error) {
	f := &flag.FlagSet{}
	f.BoolVar(&cmd.flags.work, "work", false, "print workdir and don't cleanup on exit")
	f.BoolVar(&cmd.flags.verbose, "verbose", false, "verbose godebug")
	dirs := f.String("dirs", "", "comma-separated list of directories")
	files := f.String("files", "", "comma-separated list of files to avoid annotating big directories")
	run := f.String("run", "", "run test")
	verboseTests := f.Bool("v", false, "verbose tests")
	env := f.String("env", "", "set env variables, separated by comma (ex: \"GOOS=linux,...\"'")

	if err := f.Parse(args); err != nil {
		if err == flag.ErrHelp {
			f.SetOutput(cmd.Stderr)
			f.PrintDefaults()
			return true, nil
		}
		return true, err
	}

	cmd.flags.dirs = splitCommaList(*dirs)
	cmd.flags.files = splitCommaList(*files)
	cmd.flags.env = strings.Split(*env, ",")
	cmd.flags.otherArgs = f.Args()

	// set test run flag at other flags to pass to the test exec
	if *run != "" {
		a := []string{"-test.run", *run}
		cmd.flags.runArgs = append(a, cmd.flags.runArgs...)
	}

	// verbose
	if *verboseTests {
		a := []string{"-test.v"}
		cmd.flags.runArgs = append(a, cmd.flags.runArgs...)
	}

	return false, nil
}

func (cmd *Cmd) parseBuildArgs(args []string) (done bool, _ error) {
	f := &flag.FlagSet{}
	f.BoolVar(&cmd.flags.work, "work", false, "print workdir and don't cleanup on exit")
	f.BoolVar(&cmd.flags.verbose, "verbose", false, "verbose godebug")
	dirs := f.String("dirs", "", "comma-separated list of directories")
	files := f.String("files", "", "comma-separated list of files to avoid annotating big directories")
	addr := f.String("addr", "", "address to serve from, built into the binary")
	env := f.String("env", "", "set env variables, separated by comma (ex: \"GOOS=linux,...\"'")

	if err := f.Parse(args); err != nil {
		if err == flag.ErrHelp {
			f.SetOutput(cmd.Stderr)
			f.PrintDefaults()
			return true, nil
		}
		return true, err
	}

	cmd.flags.dirs = splitCommaList(*dirs)
	cmd.flags.files = splitCommaList(*files)
	cmd.flags.address = *addr
	cmd.flags.env = strings.Split(*env, ",")
	cmd.flags.otherArgs = f.Args()
	if len(cmd.flags.otherArgs) > 0 {
		cmd.flags.filename = cmd.flags.otherArgs[0]
		cmd.flags.otherArgs = cmd.flags.otherArgs[1:]
	}

	return false, nil
}

func (cmd *Cmd) parseConnectArgs(args []string) (done bool, _ error) {
	f := &flag.FlagSet{}
	addr := f.String("addr", "", "address to connect to, built into the binary")

	if err := f.Parse(args); err != nil {
		if err == flag.ErrHelp {
			f.SetOutput(cmd.Stderr)
			f.PrintDefaults()
			return true, nil
		}
		return true, err
	}

	cmd.flags.address = *addr

	return false, nil
}

//------------

func cmdUsage() string {
	return `Usage:
	GoDebug <command> [arguments]
The commands are:
	run		build and run program with godebug data
	test		test packages compiled with godebug data
	build 	build binary with godebug data (allows remote debug)
	connect	connect to a binary built with godebug data (allows remote debug)
Examples:
	GoDebug -help
	GoDebug run -help
	GoDebug run main.go -arg1 -arg2
	GoDebug run -dirs=dir1,dir2 -files=f1.go,f2.go main.go -arg1 -arg2
	GoDebug test -help
	GoDebug test
	GoDebug test -run mytest
	GoDebug build -addr=:8080 main.go
	GoDebug connect -addr=:8080
`
}

//------------

func mkdirAllWriteFile(filename string, src []byte) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0770); err != nil {
		return err
	}
	return ioutil.WriteFile(filename, []byte(src), 0660)
}

func mkdirAllCopyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0770); err != nil {
		return err
	}
	return copyFile(src, dst)
}

func copyFile(src, dst string) error {
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()
	to, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer to.Close()
	_, err = io.Copy(to, from)
	return err
}

//------------

func mkdirAllCopyFileSync(src, dst string) error {
	// must exist in src
	info1, err := os.Stat(src)
	if os.IsNotExist(err) {
		return fmt.Errorf("not found in src: %v", src)
	}

	// already exists in dest with same modification time
	info2, err := os.Stat(dst)
	if !os.IsNotExist(err) {
		// compare modification time in src
		if info2.ModTime().Equal(info1.ModTime()) {
			return nil
		}
	}

	if err := mkdirAllCopyFile(src, dst); err != nil {
		return err
	}

	// set modtime equal to src to avoid copy next time
	t := info1.ModTime().Local()
	return os.Chtimes(dst, t, t)
}

//------------

func replaceExt(filename, ext string) string {
	// remove extension
	tmp := filename
	ext2 := filepath.Ext(tmp)
	if len(ext2) > 0 {
		tmp = tmp[:len(tmp)-len(ext2)]
	}
	// add new extension
	return tmp + ext
}

func normalizeFilenameForExec(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}

	// TODO: review
	if !strings.HasPrefix(filename, "./") {
		return "./" + filename
	}

	return filename
}

//------------

func splitCommaList(val string) []string {
	a := strings.Split(val, ",")
	seen := make(map[string]bool)
	u := []string{}
	for _, s := range a {
		// don't add empty strings
		s := strings.TrimSpace(s)
		if s == "" {
			continue
		}
		// don't add repeats
		if seen[s] {
			continue
		}
		seen[s] = true

		u = append(u, s)
	}
	return u
}

//------------

func setupServerNetAddr(addr string) {
	if addr != "" {
		debug.ServerNetwork = "tcp"
		debug.ServerAddress = addr
		return
	}

	// generate address: allows multiple editors to run debug sessions at the same time.

	seed := time.Now().UnixNano() + int64(os.Getpid())
	ra := rand.New(rand.NewSource(seed))
	r := ra.Intn(10000)

	switch runtime.GOOS {
	case "linux":
		debug.ServerNetwork = "unix"
		p := "editor_godebug.sock" + fmt.Sprintf("%v", r)
		debug.ServerAddress = filepath.Join(os.TempDir(), p)
	default:
		debug.ServerNetwork = "tcp"
		p := fmt.Sprintf("%v", 30071+r)
		debug.ServerAddress = "127.0.0.1:" + p
	}
}
