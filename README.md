<h1 align="center">Go Jitter</h1>
<p align="center">
<a href="https://goreportcard.com/report/github.com/cbergoon/go-jitter"><img src="https://goreportcard.com/badge/github.com/cbergoon/go-jitter?1=1" alt="Report"></a>
<a href="https://godoc.org/github.com/cbergoon/go-jitter"><img src="https://img.shields.io/badge/godoc-reference-brightgreen.svg" alt="Docs"></a>
<a href="#"><img src="https://img.shields.io/badge/version-0.1.0-brightgreen.svg" alt="Version"></a>
</p>

Library to test and calculate network "jitter"

#### Features

* ICMP jitter test
* Uncorrected Standard Deviation jitter calculation
* Corrected Standard Deviation jitter calculation (Bessel's Correction)
* RTT Range

#### Installation

Get the source with ```go get```:

```bash
$ go get github.com/cbergoon/go-jitter
```

#### Example Usage

```go
package main

import (
	"fmt"
	"time"

	jitter "github.com/cbergoon/go-jitter"
)

func main() {
	j, err := jitter.NewJitterer("google.com")
	if err != nil {
		fmt.Println(err)
	}

	j.SetBlockSampleSize(10)
	j.SetPingerPrivileged(true)
	j.SetPingerTimeout(time.Second * 10)

	j.Run()

	s := j.Statistics()

	fmt.Println("Squared Deviation: ", s.SquaredDeviation)
	fmt.Println("Uncorrected Deviation: ", s.UncorrectedSD)
	fmt.Println("Corrected Deviation: ", s.CorrectedSD)
	fmt.Println("RTT Range: ", s.RttRange)
	fmt.Println("RTTs: ", s.RTTS)
}
```

#### License

This project is licensed under the MIT License.