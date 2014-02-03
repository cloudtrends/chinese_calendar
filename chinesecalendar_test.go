package chinesecalendar

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type yearInfo2YearDayTest struct {
	in  int
	out int
}

var yearInfo2YearDayTests = []yearInfo2YearDayTest{
	yearInfo2YearDayTest{0, 348},                // no leap month, and every month has 29 days
	yearInfo2YearDayTest{1, 377},                // 1 leap month, and every month has 29 days.
	yearInfo2YearDayTest{(1<<12 - 1) * 16, 360}, // no leap month, and every month has 30 days.
	yearInfo2YearDayTest{(1<<13-1)*16 + 1, 390}, // 1 leap month, and every month has 30 days.
	yearInfo2YearDayTest{(1<<12-1)*16 + 1, 389}, // 1 leap month, and every normal month has 30 days, and leap month has 29 days.
}

func TestYearInfo_TotalDays(t *testing.T) {
	for _, dt := range yearInfo2YearDayTests {
		v := yearInfo{1900, dt.in}.TotalDays()
		if v != dt.out {
			t.Errorf("YearInfo{info:%d}.TotalDays() = %d, expect %d.", dt.in, v, dt.out)
		}
	}
}

type fromOffsetTest struct {
	in     int
	expect ChineseCalendar
}

var fromOffsetTests = []fromOffsetTest{
	fromOffsetTest{0, ChineseCalendar{1900, 1, 1, false}},
}

func TestFromOffset(t *testing.T) {
	for _, item := range fromOffsetTests {
		out, err := fromOffset(item.in)
		assert.NoError(t, err)
		assert.Equal(t, out, item.expect)
	}
}

func TestFromSolarDate(t *testing.T) {
	t1, err := FromSolarDate(1976, 10, 1)
	assert.NoError(t, err)
	assert.Equal(t, t1, ChineseCalendar{1976, 8, 8, true})
}

func TestChineseCalendar_Validate(t *testing.T) {
	assert.Equal(t, ChineseCalendar{1899, 12, 29, false}.Validate(), ErrYearOutOfRange)
	assert.Equal(t, ChineseCalendar{1899, 12, 30, false}.Validate(), ErrYearOutOfRange)
	assert.Equal(t, ChineseCalendar{1900, 1, 1, false}.Validate(), nil)
	assert.Equal(t, ChineseCalendar{2049, 1, 1, false}.Validate(), nil)
	assert.Equal(t, ChineseCalendar{2050, 1, 1, false}.Validate(), ErrYearOutOfRange)
	assert.Equal(t, ChineseCalendar{2051, 1, 1, false}.Validate(), ErrYearOutOfRange)

	assert.Equal(t, ChineseCalendar{1900, 1, 1, false}.Validate(), nil)
	assert.Equal(t, ChineseCalendar{1900, 1, 1, true}.Validate(), ErrNotLeapMonth)
	assert.Equal(t, ChineseCalendar{1976, 8, 1, false}.Validate(), nil)
	assert.Equal(t, ChineseCalendar{1976, 8, 1, true}.Validate(), nil)

	assert.Equal(t, ChineseCalendar{1900, 1, -1, false}.Validate(), ErrDayOutOfRange)
	assert.Equal(t, ChineseCalendar{1900, 1, 0, false}.Validate(), ErrDayOutOfRange)
	assert.Equal(t, ChineseCalendar{1900, 1, 1, false}.Validate(), nil)
	assert.Equal(t, ChineseCalendar{1900, 1, 29, false}.Validate(), nil)
	assert.Equal(t, ChineseCalendar{1900, 1, 30, false}.Validate(), ErrDayOutOfRange)
	assert.Equal(t, ChineseCalendar{1900, 2, 29, false}.Validate(), nil)
	assert.Equal(t, ChineseCalendar{1900, 2, 30, false}.Validate(), nil)
	assert.Equal(t, ChineseCalendar{1900, 2, 31, false}.Validate(), ErrDayOutOfRange)

	assert.Equal(t, ChineseCalendar{1900, 0, 1, false}.Validate(), ErrMonthOutOfRange)
	assert.Equal(t, ChineseCalendar{1900, 1, 1, false}.Validate(), nil)
	assert.Equal(t, ChineseCalendar{1900, 12, 1, false}.Validate(), nil)
	assert.Equal(t, ChineseCalendar{1900, 13, 1, false}.Validate(), ErrMonthOutOfRange)
}
func TestChineseCalendar_ToSolarDate(t *testing.T) {
	assert.Equal(t,
		ChineseCalendar{1976, 8, 8, true}.MustToSolarDate().Format("2006-01-02"),
		"1976-10-01",
	)
}

func TestChineseCalendar_Before(t *testing.T) {
	assert.True(t,
		ChineseCalendar{1982, 11, 20, false}.Before(ChineseCalendar{1982, 11, 21, false}))
}

func TestChineseCalendar_After(t *testing.T) {
	assert.False(t,
		ChineseCalendar{1982, 11, 20, false}.After(ChineseCalendar{1982, 11, 21, false}))
}

func TestChineseCalendar_NextDay(t *testing.T) {
	d, err := ChineseCalendar{1982, 11, 20, false}.NextDay()
	assert.NoError(t, err)
	assert.Equal(t, d, ChineseCalendar{1982, 11, 21, false})
}

func TestChineseCalendar_PrevDay(t *testing.T) {
	d, err := ChineseCalendar{1982, 11, 1, false}.PrevDay()
	assert.NoError(t, err)
	assert.Equal(t, d, ChineseCalendar{1982, 10, 30, false})
	d, err = ChineseCalendar{1976, 9, 1, false}.PrevDay()
	assert.NoError(t, err)
	assert.Equal(t, d, ChineseCalendar{1976, 8, 29, true})

}

func TestFromTime(t *testing.T) {
	t1 := time.Date(1983, time.January, 3, 0, 0, 0, 0, time.Local)

	t2, err := FromTime(t1)
	assert.NoError(t, err)
	assert.Equal(t, t2, ChineseCalendar{1982, 11, 20, false})
}

func TestToday(t *testing.T) {
	c := Today()
	assert.NotNil(t, c)
}
