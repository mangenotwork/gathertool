package test

import (
	gt "github.com/mangenotwork/gathertool"
	"testing"
	"time"
)

func TestTimeHelperTest(t *testing.T) {
	t.Log(gt.Timestamp2Date(gt.Timestamp()))
	t.Log(gt.TimestampStr())
	t.Log(gt.Timestamp2Date(gt.BeginDayUnix()))
	t.Log(gt.Timestamp2Date(gt.EndDayUnix()))
	t.Log(gt.Timestamp2Date(gt.MinuteAgo(10)))
	t.Log(gt.Timestamp2Date(gt.HourAgo(1)))
	t.Log(gt.Timestamp2Date(gt.DayAgo(2)))
	t.Log(gt.DayDiff("2023-11-23 00:00:00", "2023-12-23 00:00:00"))
	t.Log(gt.DayDiffAtUnix(gt.DayAgo(10), gt.Timestamp()))
	t.Log(gt.NowToEnd())
	t.Log(gt.IsToday(gt.DayAgo(10)))
	t.Log(gt.IsTodayList(gt.DayAgo(10)))
	t.Log(gt.Timestamp2Week(gt.Timestamp()))
	t.Log(gt.Timestamp2WeekXinQi(gt.Timestamp()))
	t.Log(gt.LatestDate(10))
	t.Log(gt.GetCurrentMonthRange())
	t.Log(gt.GetCurrentWeekRange())
	t.Log(gt.GetTodayRange())
}

func TestTimeHelperTickerRunTest(t *testing.T) {
	gt.TickerRun(1*time.Second, true, func() {
		gt.Info("test...")
	})
}
