package upd8

import (
	"context"
	"sync"
	"time"
)

// Scanner coordinates detection and update queries across managers.
type Scanner struct {
	Managers []Manager
	Runner   CommandRunner
}

// Scan detects available managers and fetches their outdated package lists.
func (s Scanner) Scan(ctx context.Context) []Result {
	if s.Runner == nil {
		s.Runner = ExecRunner{}
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	type resultTuple struct {
		idx int
		res Result
	}

	output := make([]Result, 0, len(s.Managers))
	resCh := make(chan resultTuple, len(s.Managers))
	var wg sync.WaitGroup

	for idx, mgr := range s.Managers {
		mgr := mgr
		idx := idx
		if !mgr.Detect(ctx) {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			start := time.Now()
			res := mgr.CheckUpdates(ctx)
			if res.DurationMs == 0 {
				res.DurationMs = time.Since(start).Milliseconds()
			}
			resCh <- resultTuple{idx: idx, res: res}
		}()
	}

	wg.Wait()
	close(resCh)

	indexed := make(map[int]Result)
	for res := range resCh {
		indexed[res.idx] = res.res
	}

	for idx := range s.Managers {
		if res, ok := indexed[idx]; ok {
			output = append(output, res)
		}
	}

	return output
}

// Watch repeatedly scans with a provided interval and invokes the callback with each result batch.
func (s Scanner) Watch(ctx context.Context, interval time.Duration, cb func([]Result)) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	select {
	case <-ctx.Done():
		return
	default:
		cb(s.Scan(ctx))
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			results := s.Scan(ctx)
			cb(results)
		}
	}
}
