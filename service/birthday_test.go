package service

import (
  "encoding/json"
  "fmt"
  "testing"
  "time"

   "github.com/kinbiko/jsonassert"
)

const yyyymmdd = "2006-01-02"

func TestBirthdayDaysRemaining(t *testing.T) {
  defer resetClock()
  var clock  = func(y, m, d int) func() time.Time {
    return func() time.Time {
      return time.Date(y, time.Month(m), d, 13, 37, 42, 0, time.UTC)
    }
  }

  // 2020 leap, (2021, 2022, 2023) non-leap, 2024 leap

  cases := []struct {
    name      string
    user      string
    birthDate string
    expected  uint
    today     func() time.Time
  }{
    {
      name: "HappyBirthday",
      user: "alpha",
      birthDate: "1978-06-18",
      expected: 0,
      today: clock(2021, 6, 18),
    },
    {
      name: "DateBeforeToday",
      user: "alpha",
      birthDate: "1978-06-18",
      expected: 19,
      today: clock(2022, 5, 30),
    },
    {
      name: "DateAfterToday",
      user: "bravo",
      birthDate: "1997-11-28",
      expected: 182,
      today: clock(2022, 5, 30),
    },
    {
      name: "DateAfterTodayNextYear",
      user: "charly",
      birthDate: "2001-05-13",
      expected: 348,
      today: clock(2022, 5, 30),
    },
    {
      name: "DateAfterTodayLeapYear",
      user: "delta",
      birthDate: "2001-05-13",
      expected: 348,
      today: clock(2020, 5, 30),
    },
    {
      name: "LeapYearFeb29",
      user: "echo",
      birthDate: "1984-02-29",
      expected: 275,  // March 1st
      today: clock(2022, 5, 30),
    },
    {
      name: "LeapYearFeb29NextYearIsLeap",
      user: "echo",
      birthDate: "1984-02-29",
      expected: 275,  // Feb 29nd
      today: clock(2023, 5, 30),
    },
    {
      name: "LeapYearEdgeCase1",
      user: "foxtrot",
      birthDate: "1981-12-30",
      expected: 364,
      today: clock(2020, 12, 31),
    },
    {
      name: "LeapYearEdgeCase2",
      user: "foxtrot",
      birthDate: "1981-12-30",
      expected: 365,
      today: clock(2019, 12, 31),
    },
    {
      name: "LeapYearEdgeCase3",
      user: "echo",
      birthDate: "1984-02-29",
      expected: 2,
      today: clock(2020, 02, 27),
    },
    {
      name: "LeapYearEdgeCase4",
      user: "echo",
      birthDate: "1984-02-29",
      expected: 2,
      today: clock(2021, 02, 27),
    },
  }

  for _, tc := range cases {
    t.Run(tc.name, func(t *testing.T) {
      nowFun = tc.today
      parsed, _ := time.Parse(yyyymmdd, tc.birthDate)
      bday := Birthday{name: tc.user, Dob: parsed}
      if dr:= bday.daysRemaining(); dr != tc.expected {
        t.Errorf("%s: expected %d, got %d", tc.user, tc.expected, dr)
      }
    })
  }
}

func TestBirthdayGreeting(t *testing.T) {
  defer resetClock()
  nowFun = func() time.Time {
      return time.Date(2023, 10, 12, 13, 37, 42, 0, time.UTC)
  }

  ja := jsonassert.New(t)

  cases := []struct {
    name      string
    user      string
    birthDate string
    expected  string
  }{
    {
      name: "HappyBirthday",
      user: "alpha",
      birthDate: "1955-10-12",
      expected: `{"message": "Hello, alpha! Happy birthday!"}`,
    },
    {
      name: "DateBeforeToday",
      user: "bravo",
      birthDate: "1978-06-18",
      expected: `{"message": "Hello, bravo! Your birthday is in 250 day(s)"}`,
    },
    {
      name: "DateAfterToday",
      user: "charly",
      birthDate: "1997-11-28",
      expected: `{"message": "Hello, charly! Your birthday is in 47 day(s)"}`,
    },
  }

  for _, tc := range cases {
    t.Run(tc.name, func(t *testing.T) {
      parsed, _ := time.Parse(yyyymmdd, tc.birthDate)
      bday := Birthday{name: tc.user, Dob: parsed}

      jsonData, err := json.Marshal(bday)
      if err != nil {
        t.Fatal(err)
      }

      ja.Assertf(string(jsonData), tc.expected)

    })
  }
}

func BenchmarkBirthdayDaysRemaining(b *testing.B) {
  defer resetClock()
  nowFun = func() time.Time {
      return time.Date(2023, 10, 12, 13, 37, 42, 0, time.UTC)
  }

  cases := []struct {
    birthDate time.Time
  }{
    {
      birthDate: time.Date(1978,6,18,0,0,0,0,time.UTC),
    },
    {
      birthDate: time.Date(1997,11,28,0,0,0,0,time.UTC),
    },
    {
      birthDate: time.Date(2001,5,13,0,0,0,0,time.UTC),
    },
    {
      birthDate: time.Date(1984,2,29,0,0,0,0,time.UTC),
    },
    {
      birthDate: time.Date(1981,1,05,0,0,0,0,time.UTC),
    },
    {
      birthDate: time.Date(1984,10,12,0,0,0,0,time.UTC),
    },
  }

  for _, bc := range cases {
    b.Run(fmt.Sprintf("date_%s_", bc.birthDate.Format(yyyymmdd)), func(b *testing.B) {
      bday := Birthday{Dob: bc.birthDate}
      for i := 0; i < b.N; i++ {
        bday.daysRemaining()
      }
    })
  }
}

