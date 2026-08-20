# apis-go

Shared wire-contract DTOs and Russian date formatting for the 005-bot
services. This module is a dependency-free (stdlib-only) library extracted
from [005-bot/monitor-go](https://github.com/005-bot/monitor-go), consumed by
tg-bot-go and future services.

## Packages

- `domain` - wire-contract DTOs for the outages Redis channel, matching the
  monitor-go publisher JSON schema: `Outage`, `OrganizationInfo`, `Street`,
  `Reason`, `WaterDelivery`, `OutageDetails`, `ResourceType` and
  `DetectResourceType`. `time.Time` values serialize as RFC3339Nano.
- `format` - Russian date formatting: `FormatDateRU` (e.g. `18 августа 14:30`,
  lowercase genitive month names, rendered in the timestamp's own offset) and
  `FormatDates` (space-joined `FormatDateRU`).

## Usage

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"

	domain "github.com/005-bot/apis-go"
	"github.com/005-bot/apis-go/format"
)

func main() {
	rt := domain.DetectResourceType("Холодное водоснабжение")
	msg := domain.Outage{
		Area: "Район 1",
		OrganizationInfo: domain.OrganizationInfo{
			ResourceType: rt,
			Resource:     "ХВС",
			Organization: "ООО УК",
			Phones:       []string{"+7 (111) 111-11-11"},
		},
		Details: domain.OutageDetails{
			Streets: []domain.Street{{Name: "ул. Ленина", Buildings: []string{"1", "2"}}},
		},
		Period: []time.Time{time.Now()},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	fmt.Println(format.FormatDateRU(time.Now()))
	fmt.Println(format.FormatDates([]time.Time{time.Now()}))
}
```

## Development

- `make test` - run tests
- `make lint` - run golangci-lint
- `make build` - build the root package

## Attribution

The `domain` package is extracted from `005-bot/monitor-go`
(`internal/domain/types.go`, Apache-2.0). The `format` package mirrors the RU
date formatting of `005-bot/monitor-go/internal/parser/date`. This project is
licensed under Apache-2.0; see [LICENSE](LICENSE).
