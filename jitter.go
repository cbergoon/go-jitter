package jitter

import (
	"math"
	"time"

	ping "github.com/sparrc/go-ping"
)

// Statistics represents the jitter test results with corrected and uncorrected deviations
type Statistics struct {
	Start time.Time
	End   time.Time

	Host string
	RTTS []time.Duration

	UncorrectedSD    time.Duration
	CorrectedSD      time.Duration
	SquaredDeviation time.Duration

	RttRange time.Duration

	PingStatistics *ping.Statistics
}

// Jitterer represents the configuration and actors to test jitter
type Jitterer struct {
	Host string
	// blockSampleSize represents the number of measurements that will result in 1 jitter calculation
	blockSampleSize int
	// OnFinish callback function to receive results of test
	OnFinish func(stats *Statistics)
	// pinger used to execute consecutive ping requests
	pinger           *ping.Pinger
	pingerPrivileged bool
	pingerTimeout    time.Duration
}

// NewJitterer returns a new Jitterer for the host specified
func NewJitterer(targetHost string) (*Jitterer, error) {
	pngr, err := ping.NewPinger(targetHost)
	if err != nil {
		return nil, err
	}

	return &Jitterer{
		Host:             targetHost,
		blockSampleSize:  3,
		pingerPrivileged: false,
		pingerTimeout:    time.Second,
		pinger:           pngr,
	}, nil
}

// Run executes jitter test
func (j *Jitterer) Run() {
	startTime := time.Now()

	j.pinger.SetPrivileged(j.pingerPrivileged)
	j.pinger.Count = j.blockSampleSize
	j.pinger.Timeout = j.pingerTimeout

	j.pinger.OnRecv = nil
	j.pinger.OnFinish = func(stats *ping.Statistics) {
		endTime := time.Now()
		js := j.generateStatistics(stats, startTime, endTime)
		j.OnFinish(js)
	}

	j.pinger.Run()
}

// SetBlockSampleSize controls the number of test in the sample
func (j *Jitterer) SetBlockSampleSize(size int) {
	j.blockSampleSize = size
}

// SetPingerPrivileged indicates if application should use UDP or priveleged ICMP packets
func (j *Jitterer) SetPingerPrivileged(value bool) {
	j.pingerPrivileged = value
}

// SetPingerTimeout time for tests to complete
func (j *Jitterer) SetPingerTimeout(timeout time.Duration) {
	j.pingerTimeout = timeout
}

// generateStatistics calculates jitter
func (j *Jitterer) generateStatistics(pingStats *ping.Statistics, startedAt time.Time, endedAt time.Time) *Statistics {
	usd := time.Duration(calculateUncorrectedDeviation(pingStats.Rtts))
	csd := time.Duration(calculateCorrectedDeviation(pingStats.Rtts))
	sd := time.Duration(calculateSquaredDeviation(pingStats.Rtts))
	rng := calculateRange(pingStats.Rtts)

	return &Statistics{
		Host:             j.Host,
		Start:            startedAt,
		End:              endedAt,
		PingStatistics:   pingStats,
		RTTS:             pingStats.Rtts,
		UncorrectedSD:    time.Duration(usd),
		CorrectedSD:      time.Duration(csd),
		SquaredDeviation: time.Duration(sd),
		RttRange:         rng,
	}
}

// calculateRange finds the range of a slice of durations
func calculateRange(values []time.Duration) time.Duration {
	if len(values) <= 1 {
		return time.Duration(0)
	}
	min := values[0]
	max := time.Duration(0)
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return max - min
}

// calculateSquaredDeviation calculates the squared deviation
func calculateSquaredDeviation(values []time.Duration) float64 {
	avg := calculateAverageDuration(values)
	sd := 0.0
	for _, v := range values {
		sd += math.Pow((float64(v) - float64(avg)), 2.0)
	}
	return sd
}

// calculateUncorrectedDeviation calculates standard deviation without correction
func calculateUncorrectedDeviation(values []time.Duration) float64 {
	sd := calculateSquaredDeviation(values)
	return math.Sqrt(sd / float64(len(values)))
}

// calculateCorrectedDeviation calculates standard deviation using Bessel's correction which uses n-1 in the SD formula to correct bias of small sample size
func calculateCorrectedDeviation(values []time.Duration) float64 {
	sd := calculateSquaredDeviation(values)
	return math.Sqrt(sd / (float64(len(values)) - 1))
}

// calculateAverageDuration calculates the average of a slice of durations
func calculateAverageDuration(values []time.Duration) float64 {
	l := len(values)
	if l <= 0 {
		return float64(0.0)
	}
	s := time.Duration(0)
	for _, d := range values {
		s += d
	}
	return float64(s) / float64(l)
}
