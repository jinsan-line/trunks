package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"

	trunks "github.com/straightdave/trunks/lib"
	"github.com/straightdave/trunks/lib/plot"
)

const plotUsage = `Usage: trunks plot [options] [<file>...]

Outputs an HTML time series plot of request latencies over time.
The X axis represents elapsed time in seconds from the beginning
of the earliest attack in all input files. The Y axis represents
request latency in milliseconds.

Click and drag to select a region to zoom into. Double click to zoom out.
Choose a different number on the bottom left corner input field
to change the moving average window size (in data points).

Arguments:
  <file>  A file with trunks attack results encoded with one of
          the supported encodings (gob | json | csv) [default: stdin]

Options:
  --title      Title and header of the resulting HTML page.
               [default: Trunks Plot]
  --threshold  Threshold of data points to downsample series to.
               Series with less than --threshold number of data
               points are not downsampled. [default: 4000]

Examples:
  echo "GET http://:80" | trunks attack -name=50qps -rate=50 -duration=5s > results.50qps.bin
  cat results.50qps.bin | trunks plot > plot.50qps.html
  echo "GET http://:80" | trunks attack -name=100qps -rate=100 -duration=5s > results.100qps.bin
  trunks plot results.50qps.bin results.100qps.bin > plot.html
`

func plotCmd() command {
	fs := flag.NewFlagSet("trunks plot", flag.ExitOnError)
	title := fs.String("title", "Trunks Plot", "Title and header of the resulting HTML page")
	threshold := fs.Int("threshold", 4000, "Threshold of data points above which series are downsampled.")
	output := fs.String("output", "stdout", "Output file")

	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, plotUsage)
	}

	return command{fs, func(args []string) error {
		fs.Parse(args)
		files := fs.Args()
		if len(files) == 0 {
			files = append(files, "stdin")
		}
		return plotRun(files, *threshold, *title, *output)
	}}
}

func plotRun(files []string, threshold int, title, output string) error {
	dec, mc, err := decoder(files)
	defer mc.Close()
	if err != nil {
		return err
	}

	out, err := file(output, true)
	if err != nil {
		return err
	}
	defer out.Close()

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)

	p := plot.New(
		plot.Title(title),
		plot.Downsample(threshold),
		plot.Label(plot.ErrorLabeler),
	)

decode:
	for {
		select {
		case <-sigch:
			break decode
		default:
			var r trunks.Result
			if err = dec.Decode(&r); err != nil {
				if err == io.EOF {
					break decode
				}
				return err
			}

			if err = p.Add(&r); err != nil {
				return err
			}
		}
	}

	p.Close()

	_, err = p.WriteTo(out)
	return err
}
