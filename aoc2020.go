package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"
)

var solutions = make(map[int][]func(*problemContext))

func addSolutions(problem int, fns ...func(*problemContext)) {
	solutions[problem] = append(solutions[problem], fns...)
}

func findSolution(problem, solNumber int) (func(*problemContext), error) {
	solns, ok := solutions[problem]
	if !ok {
		return nil, fmt.Errorf("no solutions for problem %d", problem)
	}
	if solNumber > len(solns) {
		return nil, fmt.Errorf("problem %d only has %d solution(s)", problem, len(solns))
	}
	return solns[solNumber-1], nil
}

func parseProblem(name string) (problem, solNumber int, err error) {
	parts := strings.SplitN(name, ".", 2)
	solNumber = 1
	if len(parts) == 2 {
		var err error
		solNumber, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, err
		}
	}
	problem, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	return problem, solNumber, nil
}

func main() {
	log.SetFlags(0)

	cpuProfile := flag.String("cpuprofile", "", "write CPU profile to `file`")
	printTimings := flag.Bool("t", false, "print timings")
	readStdin := flag.Bool("stdin", false, "read from stdin instead of default file")
	flag.Parse()

	if *printTimings && *cpuProfile != "" {
		log.Fatal("-t and -cpuprofile are incompatible")
	}
	if flag.NArg() != 1 {
		log.Fatalf("Usage: %s [flags] problem", os.Args[0])
	}
	problem, solNumber, err := parseProblem(flag.Arg(0))
	if err != nil {
		log.Fatalln("Bad problem number:", err)
	}
	fn, err := findSolution(problem, solNumber)
	if err != nil {
		log.Fatal(err)
	}
	ctx, err := newProblemContext(problem, *readStdin)
	if err != nil {
		log.Fatal(err)
	}
	defer ctx.close()

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatalln("Error writing CPU profile:", err)
			}
		}()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatalln("Error starting CPU profile:", err)
		}
		defer pprof.StopCPUProfile()

		ctx.profiling = true
		fn(ctx)
		return
	}

	ctx.timings.start = time.Now()
	fn(ctx)
	ctx.timings.done = time.Now()
	if *printTimings {
		ctx.printTimings()
	}
}

type problemContext struct {
	f            *os.File
	needClose    bool
	profiling    bool
	profileStart time.Time
	l            *log.Logger

	timings struct {
		start time.Time
		load  time.Time
		part1 time.Time
		part2 time.Time
		done  time.Time
	}
}

func (ctx *problemContext) reportLoad() { ctx.timings.load = time.Now() }

func (ctx *problemContext) reportPart1(v ...interface{}) {
	ctx.timings.part1 = time.Now()
	args := append([]interface{}{"Part 1:"}, v...)
	ctx.l.Println(args...)
}

func (ctx *problemContext) reportPart2(v ...interface{}) {
	ctx.timings.part2 = time.Now()
	args := append([]interface{}{"Part 2:"}, v...)
	ctx.l.Println(args...)
}

func (ctx *problemContext) printTimings() {
	ctx.l.Println("Total:", ctx.timings.done.Sub(ctx.timings.start))
	t := ctx.timings.start
	if !ctx.timings.load.IsZero() {
		ctx.l.Println("  Load:", ctx.timings.load.Sub(t))
		t = ctx.timings.load
	}
	if !ctx.timings.part1.IsZero() {
		ctx.l.Println("  Part 1:", ctx.timings.part1.Sub(t))
		t = ctx.timings.part1
	}
	if !ctx.timings.part2.IsZero() {
		ctx.l.Println("  Part 2:", ctx.timings.part2.Sub(t))
		t = ctx.timings.part2
	}
}

func newProblemContext(n int, readStdin bool) (*problemContext, error) {
	ctx := &problemContext{
		l: log.New(os.Stdout, "", 0),
	}
	if readStdin {
		ctx.f = os.Stdin
	} else {
		name := fmt.Sprintf("%02d.txt", n)
		f, err := os.Open(name)
		if err != nil {
			return nil, err
		}
		ctx.f = f
		ctx.needClose = true
	}
	return ctx, nil
}

func (ctx *problemContext) close() {
	if ctx.needClose {
		ctx.f.Close()
	}
}

func (ctx *problemContext) loopForProfile() bool {
	if ctx.profileStart.IsZero() {
		ctx.profileStart = time.Now()
		return true
	}
	if !ctx.profiling {
		return false
	}
	return time.Since(ctx.profileStart) < 5*time.Second
}

func (ctx *problemContext) scanner() scanner {
	return newScanner(ctx.f)
}

func (ctx *problemContext) int64s() []int64 {
	var ns []int64
	s := ctx.scanner()
	for s.scan() {
		ns = append(ns, s.int64())
	}
	return ns
}

type scanner struct {
	s *bufio.Scanner
}

func newScanner(r io.Reader) scanner {
	return scanner{bufio.NewScanner(r)}
}

func (s scanner) scan() bool {
	if !s.s.Scan() {
		if err := s.s.Err(); err != nil {
			log.Fatalln("Scan error:", err)
		}
		return false
	}
	return true
}

func (s scanner) text() string {
	return s.s.Text()
}

func (s scanner) int64() int64 {
	n, err := strconv.ParseInt(s.text(), 10, 64)
	if err != nil {
		log.Fatalf("Bad int64 %q: %s", s.text(), err)
	}
	return n
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
