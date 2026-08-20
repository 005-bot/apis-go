package domain_test

import (
	"encoding/json"
	"testing"
	"time"

	domain "github.com/005-bot/apis-go"
)

const goldenOutage = `{"area":"Район 1","organization_info":{"resource_type":"Холодное водоснабжение","resource":"ХВС","organization":"ООО УК","phones":["+7 (111) 111-11-11","+7 (222) 222-22-22"]},"details":{"streets":[{"name":"ул. Ленина","buildings":["1","2"]},{"name":"пр. Мира"}],"reason":{"type":"Ремонт","description":"Аварийный ремонт"},"water_deliveries":[{"street":"ул. Ленина","buildings":"1, 2","time_start":"14:00","time_end":"18:00"}],"comments":"Без воды"},"period":["2026-08-18T14:30:00Z","2026-08-18T20:00:00Z"]}`

func TestOutageMarshalGolden(t *testing.T) {
	rt := domain.ResourceTypeColdWater
	out := domain.Outage{
		Area: "Район 1",
		OrganizationInfo: domain.OrganizationInfo{
			ResourceType: &rt,
			Resource:     "ХВС",
			Organization: "ООО УК",
			Phones:       []string{"+7 (111) 111-11-11", "+7 (222) 222-22-22"},
		},
		Details: domain.OutageDetails{
			Streets: []domain.Street{
				{Name: "ул. Ленина", Buildings: []string{"1", "2"}},
				{Name: "пр. Мира"},
			},
			Reason: &domain.Reason{Type: "Ремонт", Description: "Аварийный ремонт"},
			WaterDeliveries: []domain.WaterDelivery{
				{Street: "ул. Ленина", Buildings: "1, 2", TimeStart: "14:00", TimeEnd: "18:00"},
			},
			Comments: "Без воды",
		},
		Period: []time.Time{
			time.Date(2026, time.August, 18, 14, 30, 0, 0, time.UTC),
			time.Date(2026, time.August, 18, 20, 0, 0, 0, time.UTC),
		},
	}

	got, err := json.Marshal(out)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	if string(got) != goldenOutage {
		t.Fatalf("golden bytes mismatch:\n got: %s\nwant: %s", got, goldenOutage)
	}
}

func TestOutageMarshalOmitEmpty(t *testing.T) {
	out := domain.Outage{
		Area:             "Район 2",
		OrganizationInfo: domain.OrganizationInfo{},
		Details:          domain.OutageDetails{},
	}

	want := `{"area":"Район 2","organization_info":{"resource_type":null,"resource":"","organization":"","phones":null},"details":{"streets":null},"period":null}`

	got, err := json.Marshal(out)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	if string(got) != want {
		t.Fatalf("omitempty golden bytes mismatch:\n got: %s\nwant: %s", got, want)
	}
}

func TestOutageUnmarshalRoundTrip(t *testing.T) {
	var out domain.Outage
	if err := json.Unmarshal([]byte(goldenOutage), &out); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}

	if out.Area != "Район 1" {
		t.Fatalf("Area = %q, want %q", out.Area, "Район 1")
	}
	if out.OrganizationInfo.ResourceType == nil || *out.OrganizationInfo.ResourceType != domain.ResourceTypeColdWater {
		t.Fatalf("ResourceType = %v, want %q", out.OrganizationInfo.ResourceType, domain.ResourceTypeColdWater)
	}
	if len(out.Details.Streets) != 2 || len(out.Details.WaterDeliveries) != 1 || out.Details.Comments != "Без воды" {
		t.Fatalf("Details = %+v", out.Details)
	}
	if out.Period == nil || len(out.Period) != 2 {
		t.Fatalf("Period = %v", out.Period)
	}
	wantFirst := time.Date(2026, time.August, 18, 14, 30, 0, 0, time.UTC)
	if !out.Period[0].Equal(wantFirst) {
		t.Fatalf("Period[0] = %v, want %v", out.Period[0], wantFirst)
	}

	got, err := json.Marshal(out)
	if err != nil {
		t.Fatalf("json.Marshal after round trip: %v", err)
	}
	if string(got) != goldenOutage {
		t.Fatalf("round-trip bytes mismatch:\n got: %s\nwant: %s", got, goldenOutage)
	}
}

// TestOutageMarshalMatchesPublisher mirrors how monitor-go's publisher builds
// the Outage payload from a ParsedRecord (see
// monitor-go/internal/publisher/service.go) and asserts the emitted bytes.
func TestOutageMarshalMatchesPublisher(t *testing.T) {
	rt := domain.DetectResourceType("Горячее водоснабжение")
	msg := domain.Outage{
		Area: "Центр",
		OrganizationInfo: domain.OrganizationInfo{
			ResourceType: rt,
			Resource:     "ГВС",
			Organization: "МУП Водоканал",
			Phones:       []string{"8 800 000-00-00"},
		},
		Details: domain.OutageDetails{
			Streets: []domain.Street{
				{Name: "ул. Гагарина", Buildings: []string{"3", "5"}},
			},
		},
		Period: []time.Time{
			time.Date(2026, time.August, 18, 9, 0, 0, 0, time.UTC),
		},
	}

	want := `{"area":"Центр","organization_info":{"resource_type":"Горячее водоснабжение","resource":"ГВС","organization":"МУП Водоканал","phones":["8 800 000-00-00"]},"details":{"streets":[{"name":"ул. Гагарина","buildings":["3","5"]}]},"period":["2026-08-18T09:00:00Z"]}`

	got, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	if string(got) != want {
		t.Fatalf("publisher schema mismatch:\n got: %s\nwant: %s", got, want)
	}
}

func TestDetectResourceType(t *testing.T) {
	cases := []struct {
		name     string
		resource string
		want     domain.ResourceType
	}{
		{"exact cold water", "Холодное водоснабжение", domain.ResourceTypeColdWater},
		{"case-insensitive substring", "Плановые работы: ГОРЯЧЕЕ ВОДОСНАБЖЕНИЕ отключено", domain.ResourceTypeHotWater},
		{"prefix electricity", "Электроснабжение в доме 5", domain.ResourceTypeElectricity},
		{"suffix gas", "Газоснабжение приостановлено", domain.ResourceTypeGas},
		{"embedded heating", "плановые работы: теплоснабжение приостановлено", domain.ResourceTypeHeating},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := domain.DetectResourceType(tc.resource)
			if got == nil || *got != tc.want {
				t.Fatalf("DetectResourceType(%q) = %v, want %q", tc.resource, got, tc.want)
			}
		})
	}
}

func TestDetectResourceTypeUnknown(t *testing.T) {
	for _, resource := range []string{"", "Водоотведение", "Капитальный ремонт дома"} {
		if got := domain.DetectResourceType(resource); got != nil {
			t.Fatalf("DetectResourceType(%q) = %v, want nil", resource, got)
		}
	}
}

func TestStreetString(t *testing.T) {
	cases := []struct {
		name string
		in   domain.Street
		want string
	}{
		{"name only", domain.Street{Name: "ул. Ленина"}, "ул. Ленина"},
		{"one building", domain.Street{Name: "ул. Ленина", Buildings: []string{"1"}}, "ул. Ленина 1"},
		{"multiple buildings", domain.Street{Name: "пр. Мира", Buildings: []string{"2", "4", "6"}}, "пр. Мира 2, 4, 6"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.in.String(); got != tc.want {
				t.Fatalf("Street.String() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestOutageDetailsAddress(t *testing.T) {
	d := domain.OutageDetails{
		Streets: []domain.Street{
			{Name: "ул. Ленина", Buildings: []string{"1", "2"}},
			{Name: "пр. Мира"},
		},
	}
	want := "ул. Ленина 1, 2\nпр. Мира"
	if got := d.Address(); got != want {
		t.Fatalf("OutageDetails.Address() = %q, want %q", got, want)
	}
}
