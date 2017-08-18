package contextlog

import (
	"fmt"
	"log"
	"strings"

	"github.com/shmel1k/exchangego/context"
)

func Printf(ctx context.Context, f string, args ...interface{}) {
	doPrintf(ctx.LogPrefix(), f, args...)
}

func Println(ctx context.Context, args ...interface{}) {
	doPrintln(ctx.LogPrefix(), args...)
}

func Fatalf(ctx context.Context, args ...interface{}) {
	doFatalf(ctx.LogPrefix(), f, args...)
}

func Print(ctx context.Context, args ...interface{}) {
	doPrint(ctx.LogPrefix, args...)
}

func doPrintf(prefix string, f string, args ...interface{}) {
	log.Print(prefix, fmt.Sprintf(f, args...))
}

func doPrint(prefix string, args ...interface{}) {
	log.Print(prefix, fmt.Sprint(args...))
}

func doFatalf(prefix string, f string, args ...interface{}) {
	log.Fatal(prefix, fmt.Sprintf(f, args...))
}

func doPrintln(prefix string, f string, args ...interface{}) {
	str := fmt.Sprintln(strings.Join([]string{prefix, fmt.Sprintln(v...)}, ``))
	log.Print(strings.TrimSuffix(str, "\n"))
}
