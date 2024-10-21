//
// This file is part of PathsHelper library.
//
// Copyright 2023 Arduino AG (http://www.arduino.cc/)
//
// PathsHelper library is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
//
// As a special exception, you may use this file as part of a free software
// library without restriction.  Specifically, if other files instantiate
// templates or use macros or inline functions from this file, or you compile
// this file and link it with other files to produce an executable, this
// file does not by itself cause the resulting executable to be covered by
// the GNU General Public License.  This exception does not however
// invalidate any other reasons why the executable file might be covered by
// the GNU General Public License.
//

package paths

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
)

// Process is representation of an external process run
type Process struct {
	cmd *exec.Cmd
}

// NewProcess creates a command with the provided command line arguments
// and environment variables (that will be added to the parent os.Environ).
// The argument args[0] is the path to the executable, the remainder are the
// arguments to the command.
func NewProcess(extraEnv []string, args ...string) (*Process, error) {
	if len(args) == 0 {
		return nil, errors.New("no executable specified")
	}
	p := &Process{
		cmd: exec.Command(args[0], args[1:]...),
	}
	p.cmd.Env = append(os.Environ(), extraEnv...)
	tellCommandNotToSpawnShell(p.cmd)          // windows specific
	tellCommandToStartOnNewProcessGroup(p.cmd) // linux specific

	// This is required because some tools detects if the program is running
	// from terminal by looking at the stdin/out bindings.
	// https://github.com/arduino/arduino-cli/issues/844
	p.cmd.Stdin = nullReaderInstance
	return p, nil
}

// TellCommandNotToSpawnShell avoids that the specified Cmd display a small
// command prompt while runnning on Windows. It has no effects on other OS.
func (p *Process) TellCommandNotToSpawnShell() {
	tellCommandNotToSpawnShell(p.cmd)
}

// NewProcessFromPath creates a command from the provided executable path,
// additional environment vars (in addition to the system default ones)
// and command line arguments.
func NewProcessFromPath(extraEnv []string, executable *Path, args ...string) (*Process, error) {
	processArgs := []string{executable.String()}
	processArgs = append(processArgs, args...)
	return NewProcess(extraEnv, processArgs...)
}

// RedirectStdoutTo will redirect the process' stdout to the specified
// writer. Any previous redirection will be overwritten.
func (p *Process) RedirectStdoutTo(out io.Writer) {
	p.cmd.Stdout = out
}

// RedirectStderrTo will redirect the process' stdout to the specified
// writer. Any previous redirection will be overwritten.
func (p *Process) RedirectStderrTo(out io.Writer) {
	p.cmd.Stderr = out
}

// StdinPipe returns a pipe that will be connected to the command's standard
// input when the command starts. The pipe will be closed automatically after
// Wait sees the command exit. A caller need only call Close to force the pipe
// to close sooner. For example, if the command being run will not exit until
// standard input is closed, the caller must close the pipe.
func (p *Process) StdinPipe() (io.WriteCloser, error) {
	if p.cmd.Stdin == nullReaderInstance {
		p.cmd.Stdin = nil
	}
	return p.cmd.StdinPipe()
}

// StdoutPipe returns a pipe that will be connected to the command's standard
// output when the command starts.
//
// Wait will close the pipe after seeing the command exit, so most callers
// don't need to close the pipe themselves. It is thus incorrect to call Wait
// before all reads from the pipe have completed.
// For the same reason, it is incorrect to call Run when using StdoutPipe.
func (p *Process) StdoutPipe() (io.ReadCloser, error) {
	return p.cmd.StdoutPipe()
}

// StderrPipe returns a pipe that will be connected to the command's standard
// error when the command starts.
//
// Wait will close the pipe after seeing the command exit, so most callers
// don't need to close the pipe themselves. It is thus incorrect to call Wait
// before all reads from the pipe have completed.
// For the same reason, it is incorrect to use Run when using StderrPipe.
func (p *Process) StderrPipe() (io.ReadCloser, error) {
	return p.cmd.StderrPipe()
}

// Start will start the underliyng process.
func (p *Process) Start() error {
	return p.cmd.Start()
}

// Wait waits for the command to exit and waits for any copying to stdin or copying
// from stdout or stderr to complete.
func (p *Process) Wait() error {
	// TODO: make some helpers to retrieve exit codes out of *ExitError.
	return p.cmd.Wait()
}

// Signal sends a signal to the Process. Sending Interrupt on Windows is not implemented.
func (p *Process) Signal(sig os.Signal) error {
	return p.cmd.Process.Signal(sig)
}

// Kill causes the Process to exit immediately. Kill does not wait until the Process has
// actually exited. This only kills the Process itself, not any other processes it may
// have started.
func (p *Process) Kill() error {
	return kill(p.cmd)
}

// SetDir sets the working directory of the command. If Dir is the empty string, Run
// runs the command in the calling process's current directory.
func (p *Process) SetDir(dir string) {
	p.cmd.Dir = dir
}

// GetDir gets the working directory of the command.
func (p *Process) GetDir() string {
	return p.cmd.Dir
}

// SetDirFromPath sets the working directory of the command. If path is nil, Run
// runs the command in the calling process's current directory.
func (p *Process) SetDirFromPath(path *Path) {
	if path == nil {
		p.cmd.Dir = ""
	} else {
		p.cmd.Dir = path.String()
	}
}

// Run starts the specified command and waits for it to complete.
func (p *Process) Run() error {
	return p.cmd.Run()
}

// SetEnvironment set the environment for the running process. Each entry is of the form "key=value".
// System default environments will be wiped out.
func (p *Process) SetEnvironment(values []string) {
	p.cmd.Env = append([]string{}, values...)
}

// RunWithinContext starts the specified command and waits for it to complete. If the given context
// is canceled before the normal process termination, the process is killed.
func (p *Process) RunWithinContext(ctx context.Context) error {
	if err := p.Start(); err != nil {
		return err
	}
	completed := make(chan struct{})
	defer close(completed)
	go func() {
		select {
		case <-ctx.Done():
			p.Kill()
		case <-completed:
		}
	}()
	return p.Wait()
}

// RunAndCaptureOutput starts the specified command and waits for it to complete. If the given context
// is canceled before the normal process termination, the process is killed. The standard output and
// standard error of the process are captured and returned at process termination.
func (p *Process) RunAndCaptureOutput(ctx context.Context) ([]byte, []byte, error) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	p.RedirectStdoutTo(stdout)
	p.RedirectStderrTo(stderr)
	err := p.RunWithinContext(ctx)
	return stdout.Bytes(), stderr.Bytes(), err
}

// GetArgs returns the command arguments
func (p *Process) GetArgs() []string {
	return p.cmd.Args
}

// nullReaderInstance is an io.Reader that will always return EOF
var nullReaderInstance = &nullReader{}

type nullReader struct{}

func (r *nullReader) Read(buff []byte) (int, error) {
	return 0, io.EOF
}
