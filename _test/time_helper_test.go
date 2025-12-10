package test

import (
	"github.com/mangenotwork/gathertool"
	"testing"
	"time"
)

func TestTimeHelperTest(t *testing.T) {
	t.Log(gathertool.Timestamp2Date(gathertool.Timestamp()))
	t.Log(gathertool.TimestampStr())
	t.Log(gathertool.Timestamp2Date(gathertool.BeginDayUnix()))
	t.Log(gathertool.Timestamp2Date(gathertool.EndDayUnix()))
	t.Log(gathertool.Timestamp2Date(gathertool.MinuteAgo(10)))
	t.Log(gathertool.Timestamp2Date(gathertool.HourAgo(1)))
	t.Log(gathertool.Timestamp2Date(gathertool.DayAgo(2)))
	t.Log(gathertool.DayDiff("2023-11-23 00:00:00", "2023-12-23 00:00:00"))
	t.Log(gathertool.DayDiffAtUnix(gathertool.DayAgo(10), gathertool.Timestamp()))
	t.Log(gathertool.NowToEnd())
	t.Log(gathertool.IsToday(gathertool.DayAgo(10)))
	t.Log(gathertool.IsTodayList(gathertool.DayAgo(10)))
	t.Log(gathertool.Timestamp2Week(gathertool.Timestamp()))
	t.Log(gathertool.Timestamp2WeekXinQi(gathertool.Timestamp()))
	t.Log(gathertool.LatestDate(10))
	t.Log(gathertool.GetCurrentMonthRange())
	t.Log(gathertool.GetCurrentWeekRange())
	t.Log(gathertool.GetTodayRange())
}

func TestTimeHelperTickerRunTest(t *testing.T) {
	gathertool.TickerRun(1*time.Second, true, func() {
		gathertool.Info("test...")
	})
}
