package internal

import (
	"math"
	"sort"
	"time"
)

func ComputeFeatures(userID string, sessions []*LoginSession, transdate time.Time) *ModelFeatures {
	req := &ModelFeatures{
		CstDimID: userID,
	}

	prev := make([]*LoginSession, 0)
	for _, s := range sessions {
		if s.When.Before(transdate) || s.When.Equal(transdate) {
			prev = append(prev, s)
		}
	}

	if len(prev) == 0 {
		return req
	}

	sort.Slice(prev, func(i, j int) bool { return prev[i].When.Before(prev[j].When) })
	last := prev[len(prev)-1]
	req.LastPhoneModelCategorical = last.PhoneModel
	req.LastOSCategorical = last.OS

	cutoff7 := transdate.Add(-7 * 24 * time.Hour)
	cutoff30 := transdate.Add(-30 * 24 * time.Hour)
	var logins7, logins30 int
	phoneModels30 := make(map[string]struct{})
	os30 := make(map[string]struct{})

	times30 := make([]time.Time, 0)
	for _, s := range prev {
		if s.When.After(cutoff30) || s.When.Equal(cutoff30) {
			logins30++
			phoneModels30[s.PhoneModel] = struct{}{}
			os30[s.OS] = struct{}{}
			times30 = append(times30, s.When)
		}
		if s.When.After(cutoff7) || s.When.Equal(cutoff7) {
			logins7++
		}
	}

	req.MonthlyPhoneModelChanges = len(phoneModels30)
	req.MonthlyOSChanges = len(os30)
	req.LoginsLast7Days = logins7
	req.LoginsLast30Days = logins30

	if logins7 > 0 {
		req.LoginFrequency7d = float64(logins7) / 7.0
	}
	if logins30 > 0 {
		req.LoginFrequency30d = float64(logins30) / 30.0
	}

	if req.LoginFrequency30d != 0 {
		req.FreqChange7dVsMean = (req.LoginFrequency7d - req.LoginFrequency30d) / req.LoginFrequency30d
	}

	if logins30 != 0 {
		req.Logins7dOver30dRatio = float64(logins7) / float64(logins30)
	}

	if len(times30) >= 2 {
		sort.Slice(times30, func(i, j int) bool { return times30[i].Before(times30[j]) })
		intervals := make([]float64, 0, len(times30)-1)
		for i := 1; i < len(times30); i++ {
			intervals = append(intervals, times30[i].Sub(times30[i-1]).Seconds())
		}

		sum := 0.0
		for _, v := range intervals {
			sum += v
		}
		mean := sum / float64(len(intervals))

		varSum := 0.0
		for _, v := range intervals {
			varSum += (v - mean) * (v - mean)
		}
		variance := varSum / float64(len(intervals))
		std := math.Sqrt(variance)
		req.AvgLoginInterval30d = mean
		req.StdLoginInterval30d = std
		req.VarLoginInterval30d = variance

		var ewm float64
		alpha := 0.3
		count := 0
		for i := 1; i < len(times30); i++ {
			if times30[i].After(cutoff7) || times30[i].Equal(cutoff7) {
				iv := times30[i].Sub(times30[i-1]).Seconds()
				if count == 0 {
					ewm = iv
				} else {
					ewm = alpha*iv + (1-alpha)*ewm
				}
				count++
			}
		}
		if count > 0 {
			req.EwmLoginInterval7d = ewm
		}

		if mean+std != 0 {
			req.BurstinessLoginInterval = (std - mean) / (std + mean)
		}

		if mean != 0 {
			req.FanoFactorLoginInterval = variance / mean
		}

		times7 := make([]time.Time, 0)
		for _, t := range times30 {
			if t.After(transdate.Add(-7*24*time.Hour)) || t.Equal(transdate.Add(-7*24*time.Hour)) {
				times7 = append(times7, t)
			}
		}
		if len(times7) >= 2 {
			ints7 := make([]float64, 0)
			for i := 1; i < len(times7); i++ {
				ints7 = append(ints7, times7[i].Sub(times7[i-1]).Seconds())
			}
			sum7 := 0.0
			for _, v := range ints7 {
				sum7 += v
			}
			mean7 := sum7 / float64(len(ints7))
			if std != 0 {
				req.ZscoreAvgLoginInterval7d = (mean7 - mean) / std
			}
		}
	}

	return req
}
