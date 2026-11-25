package services

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestInboxPartitionBoundary(t *testing.T) {
	boundary := time.Date(2026, 9, 29, 12, 10, 0, 0, time.UTC)
	if got := alignInboxPartitionStart(boundary.Add(-time.Millisecond).UnixMilli()); got != boundary.Add(-inboxPartitionInterval).UnixMilli() {
		t.Fatalf("before boundary aligned to %d", got)
	}
	if got := alignInboxPartitionStart(boundary.UnixMilli()); got != boundary.UnixMilli() {
		t.Fatalf("boundary aligned to %d", got)
	}
	if got := formatInboxPartitionName(boundary.UnixMilli()); got != "p202609291210" {
		t.Fatalf("partition name = %s", got)
	}
}

func TestPlanInboxPartitionsCreatesThreeDataPartitionsAndFutureCoverage(t *testing.T) {
	now := time.Date(2026, 9, 29, 12, 3, 0, 0, time.UTC)
	plan, err := planInboxPartitions(now, []inboxPartition{{
		Name:        inboxPartitionFutureName,
		Description: "MAXVALUE",
		Future:      true,
	}})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"p202609291140", "p202609291150", "p202609291200", "p202609291210",
	}
	if got := desiredInboxPartitionNames(plan.Create); strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("created partitions = %v, want %v", got, want)
	}
	if len(plan.Drop) != 0 {
		t.Fatalf("initial plan unexpectedly drops %v", plan.Drop)
	}
}

func TestPlanInboxPartitionsDropsAllPartitionsOutsideThreeDataIntervals(t *testing.T) {
	now := time.Date(2026, 9, 29, 12, 3, 0, 0, time.UTC)
	partitions := []inboxPartition{
		finiteInboxPartition(time.Date(2026, 9, 29, 10, 30, 0, 0, time.UTC)),
		finiteInboxPartition(time.Date(2026, 9, 29, 10, 40, 0, 0, time.UTC)),
		finiteInboxPartition(time.Date(2026, 9, 29, 10, 50, 0, 0, time.UTC)),
		finiteInboxPartition(time.Date(2026, 9, 29, 11, 0, 0, 0, time.UTC)),
		finiteInboxPartition(time.Date(2026, 9, 29, 11, 10, 0, 0, time.UTC)),
		finiteInboxPartition(time.Date(2026, 9, 29, 11, 20, 0, 0, time.UTC)),
		finiteInboxPartition(time.Date(2026, 9, 29, 11, 30, 0, 0, time.UTC)),
		finiteInboxPartition(time.Date(2026, 9, 29, 11, 40, 0, 0, time.UTC)),
		finiteInboxPartition(time.Date(2026, 9, 29, 11, 50, 0, 0, time.UTC)),
		finiteInboxPartition(time.Date(2026, 9, 29, 12, 0, 0, 0, time.UTC)),
		{Name: inboxPartitionFutureName, Description: "MAXVALUE", Future: true},
	}

	plan, err := planInboxPartitions(now, partitions)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := desiredInboxPartitionNames(plan.Create), []string{"p202609291210"}; strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("created partitions = %v, want %v", got, want)
	}
	if got, want := plan.Drop, []string{
		"p202609291030", "p202609291040", "p202609291050", "p202609291100",
		"p202609291110", "p202609291120", "p202609291130",
	}; strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("dropped partitions = %v, want %v", got, want)
	}
}

func TestInboxPartitionSQLValidation(t *testing.T) {
	start := time.Date(2026, 9, 29, 12, 0, 0, 0, time.UTC).UnixMilli()
	statement, err := buildCreateInboxPartitionsSQL([]desiredInboxPartition{{
		Name:        formatInboxPartitionName(start),
		StartMillis: start,
		UpperMillis: start + inboxPartitionInterval.Milliseconds(),
	}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(statement, "ALTER TABLE `inboxmsgs_split` REORGANIZE PARTITION `p_future`") {
		t.Fatalf("unexpected create statement: %s", statement)
	}
	if _, err = buildDropInboxPartitionsSQL([]string{inboxPartitionFutureName}); err == nil {
		t.Fatal("p_future must never be droppable")
	}
	if _, err = buildDropInboxPartitionsSQL([]string{"bad`; DROP TABLE inboxmsgs"}); err == nil {
		t.Fatal("invalid partition name was accepted")
	}
	dropStatement, err := buildDropInboxPartitionsSQL([]string{"p202609291100", "p202609291110"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(dropStatement, "DROP PARTITION `p202609291100`, `p202609291110`") {
		t.Fatalf("unexpected drop statement: %s", dropStatement)
	}
}

func TestSendboxPartitionSQLValidation(t *testing.T) {
	start := time.Date(2026, 9, 30, 8, 20, 0, 0, time.UTC).UnixMilli()
	partitions := []desiredInboxPartition{{
		Name:        formatInboxPartitionName(start),
		StartMillis: start,
		UpperMillis: start + inboxPartitionInterval.Milliseconds(),
	}}
	statement, err := buildCreateMessagePartitionsSQL(sendboxPartitionTable, partitions)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(statement, "ALTER TABLE `sendboxmsgs_split` REORGANIZE PARTITION `p_future`") {
		t.Fatalf("unexpected create statement: %s", statement)
	}

	dropStatement, err := buildDropMessagePartitionsSQL(sendboxPartitionTable, []string{"p202609300800"})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(dropStatement, "ALTER TABLE `sendboxmsgs_split` DROP PARTITION `p202609300800`") {
		t.Fatalf("unexpected drop statement: %s", dropStatement)
	}
	if _, err = buildCreateMessagePartitionsSQL("sendboxmsgs_split`; DROP TABLE sendboxmsgs", partitions); err == nil {
		t.Fatal("invalid partition table was accepted")
	}
}

func TestCmdPartitionBoundaryAndRetention(t *testing.T) {
	boundary := time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC)
	if got := alignMessagePartitionStart(boundary.Add(-time.Millisecond).UnixMilli(), cmdPartitionInterval); got != boundary.Add(-cmdPartitionInterval).UnixMilli() {
		t.Fatalf("before daily boundary aligned to %d", got)
	}
	if got := alignMessagePartitionStart(boundary.UnixMilli(), cmdPartitionInterval); got != boundary.UnixMilli() {
		t.Fatalf("daily boundary aligned to %d", got)
	}

	policy, err := messagePartitionPolicyForTable(cmdInboxPartitionTable)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := planMessagePartitions(boundary.Add(12*time.Hour), []inboxPartition{{
		Name:        inboxPartitionFutureName,
		Description: "MAXVALUE",
		Future:      true,
	}}, policy)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"p202609240000", "p202609250000", "p202609260000", "p202609270000",
		"p202609280000", "p202609290000", "p202609300000", "p202610010000",
	}
	if got := desiredInboxPartitionNames(plan.Create); strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("created cmd partitions = %v, want %v", got, want)
	}
	if len(plan.Drop) != 0 {
		t.Fatalf("initial cmd plan unexpectedly drops %v", plan.Drop)
	}
}

func TestCmdPartitionSQLValidation(t *testing.T) {
	start := time.Date(2026, 9, 30, 0, 0, 0, 0, time.UTC).UnixMilli()
	partition := desiredInboxPartition{
		Name:        formatInboxPartitionName(start),
		StartMillis: start,
		UpperMillis: start + cmdPartitionInterval.Milliseconds(),
	}
	for _, tableName := range []string{cmdInboxPartitionTable, cmdSendboxPartitionTable} {
		statement, err := buildCreateMessagePartitionsSQL(tableName, []desiredInboxPartition{partition})
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(statement, "ALTER TABLE `"+tableName+"` REORGANIZE PARTITION `p_future`") {
			t.Fatalf("unexpected create statement: %s", statement)
		}
	}

	invalid := partition
	invalid.UpperMillis = start + inboxPartitionInterval.Milliseconds()
	if _, err := buildCreateMessagePartitionsSQL(cmdInboxPartitionTable, []desiredInboxPartition{invalid}); err == nil {
		t.Fatal("ten-minute boundary was accepted for daily cmd partition")
	}
}

func TestCmdPartitionDropsDataOlderThanSevenDays(t *testing.T) {
	now := time.Date(2026, 9, 30, 12, 0, 0, 0, time.UTC)
	partitions := make([]inboxPartition, 0, 13)
	for day := time.Date(2026, 9, 20, 0, 0, 0, 0, time.UTC); !day.After(time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)); day = day.Add(24 * time.Hour) {
		partitions = append(partitions, finiteMessagePartition(day, cmdPartitionInterval))
	}
	partitions = append(partitions, inboxPartition{Name: inboxPartitionFutureName, Description: "MAXVALUE", Future: true})

	policy, err := messagePartitionPolicyForTable(cmdSendboxPartitionTable)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := planMessagePartitions(now, partitions, policy)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"p202609200000", "p202609210000", "p202609220000", "p202609230000"}
	if strings.Join(plan.Drop, ",") != strings.Join(want, ",") {
		t.Fatalf("dropped cmd partitions = %v, want %v", plan.Drop, want)
	}
	if len(plan.Create) != 0 {
		t.Fatalf("cmd plan unexpectedly creates %v", desiredInboxPartitionNames(plan.Create))
	}
}

func finiteInboxPartition(start time.Time) inboxPartition {
	return finiteMessagePartition(start, inboxPartitionInterval)
}

func finiteMessagePartition(start time.Time, interval time.Duration) inboxPartition {
	startMillis := start.UnixMilli()
	return inboxPartition{
		Name:        formatInboxPartitionName(startMillis),
		Description: fmt.Sprintf("%d", startMillis+interval.Milliseconds()),
		StartMillis: startMillis,
		UpperMillis: startMillis + interval.Milliseconds(),
	}
}

func desiredInboxPartitionNames(partitions []desiredInboxPartition) []string {
	names := make([]string, 0, len(partitions))
	for _, partition := range partitions {
		names = append(names, partition.Name)
	}
	return names
}
