package service

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/scalent.io/scalent-hrms/entity"

	"github.com/scalent.io/scalent-hrms/pkg/validation"
)

func validateAttendanceSequence(logs []entity.AttendanceLog) error {

	if len(logs) == 0 {
		return nil
	}

	sort.SliceStable(logs, func(i, j int) bool {
		return logs[i].Timestamp.Before(
			logs[j].Timestamp,
		)
	})

	expectedPunch := 0 // First must be CHECK_IN.

	for _, log := range logs {

		punch, valid := normalizePunchType(log)

		if !valid {
			continue
		}

		if punch != expectedPunch {
			return fmt.Errorf("invalid attendance sequence for employee %s: expected punch %d, got %d at %v", log.EmpID, expectedPunch, punch, log.Timestamp)
		}

		if expectedPunch == 0 {
			expectedPunch = 1
		} else {
			expectedPunch = 0
		}
	}

	// Last must be CHECK_OUT.
	if expectedPunch != 0 {
		return fmt.Errorf("attendance sequence does not end with CHECK_OUT")
	}

	return nil
}

func createSyntheticAttendanceLog(source entity.AttendanceLog, timestamp time.Time, punch int, state string) entity.AttendanceLog {

	return entity.AttendanceLog{
		UID:             source.UID,
		EmpID:           source.EmpID,
		EmpName:         source.EmpName,
		Timestamp:       timestamp,
		Status:          1,
		Punch:           punch,
		AttendanceState: state,
		DeviceName:      source.DeviceName,
		CreatedAt:       source.CreatedAt,
		Synthesized:     true,
	}
}
func synthesizeAttendanceLogs(logs []entity.AttendanceLog, officeEnd time.Time) []entity.AttendanceLog {

	if len(logs) == 0 {
		return nil
	}

	// Work on a copy so the original slice is not modified.
	originalLogs := append([]entity.AttendanceLog(nil), logs...)

	// Sort chronologically.
	sort.SliceStable(originalLogs, func(i, j int) bool {
		return originalLogs[i].Timestamp.Before(
			originalLogs[j].Timestamp,
		)
	})

	syntheticLogs := make([]entity.AttendanceLog, 0)

	var activeCheckIn *entity.AttendanceLog
	var previousCheckout *entity.AttendanceLog

	const syntheticDuration = 15 * time.Minute

	for i := 0; i < len(originalLogs); i++ {

		current := originalLogs[i]

		if current.Timestamp.IsZero() {
			continue
		}

		punch, valid := normalizePunchType(current)

		if !valid {
			continue
		}

		switch punch {

		// =========================================================
		// CHECK_IN
		// =========================================================
		case 0: // CHECK_IN

			// Current CHECK_IN becomes the active CHECK_IN.
			activeCheckIn = &originalLogs[i]

			// Check if the next ORIGINAL attendance record exists.
			if i+1 < len(originalLogs) {

				next := originalLogs[i+1]

				if !next.Timestamp.IsZero() {

					nextPunch, nextValid := normalizePunchType(next)

					// -------------------------------------------------
					// CASE:
					//
					// IN
					// IN
					//
					// Synthetic OUT should be:
					//
					// NEXT IN - 15 minutes
					// -------------------------------------------------
					if nextValid && nextPunch == 0 {

						syntheticCheckoutTime := next.Timestamp.Add(-15 * time.Minute)

						// Make sure synthetic checkout is after
						// the current CHECK_IN.
						if syntheticCheckoutTime.After(current.Timestamp) {

							syntheticLogs = append(syntheticLogs, createSyntheticAttendanceLog(current, syntheticCheckoutTime, 1, "CHECK_OUT"))

							previousCheckout =
								&syntheticLogs[len(syntheticLogs)-1]

							activeCheckIn = nil
						}
					}
				}
			}

		// =========================================================
		// CHECK_OUT
		// =========================================================
		case 1:

			// Normal case:
			//
			//     IN
			//     OUT
			//
			// Do not modify either record.
			if activeCheckIn != nil {

				previousCheckout = &originalLogs[i]
				activeCheckIn = nil

				continue
			}

			// We have:
			//
			//     OUT
			//
			// without an active IN.
			//
			// Need synthetic IN.
			if previousCheckout != nil {

				syntheticCheckInTime := previousCheckout.Timestamp.Add(syntheticDuration)

				// Only create if it is before current OUT.
				if syntheticCheckInTime.Before(current.Timestamp) {

					syntheticLogs = append(syntheticLogs, createSyntheticAttendanceLog(*previousCheckout, syntheticCheckInTime, 0, "CHECK_IN"))
				}

			} else {
				// First record is OUT.
				//
				// There is no previous checkout.
				// We still need the first record to be IN.
				//
				// Use 15 minutes before the first checkout.
				syntheticCheckInTime := current.Timestamp.Add(-syntheticDuration)

				syntheticLogs = append(syntheticLogs, createSyntheticAttendanceLog(current, syntheticCheckInTime, 0, "CHECK_IN"))
			}

			// Current original CHECK_OUT remains untouched.
			previousCheckout = &originalLogs[i]
		}
	}

	// =============================================================
	// FINAL UNMATCHED CHECK_IN
	// =============================================================
	if activeCheckIn != nil {

		// Final IN has no checkout.
		//
		// Close it at office end = 19:00.
		syntheticLogs = append(syntheticLogs, createSyntheticAttendanceLog(*activeCheckIn, officeEnd, 1, "CHECK_OUT"))
	}

	return syntheticLogs
}
func calculateWorkingHoursFromLogs(targetDate time.Time, employee entity.Employee, logs []entity.AttendanceLog) entity.WorkingHoursSummary {

	if len(logs) == 0 {
		return entity.WorkingHoursSummary{
			EmpID:            employee.EmpID,
			EmpName:          employee.EmpName,
			Date:             targetDate.Format("2006-01-02"),
			Status:           "ABSENT",
			WorkingHours:     "00:00",
			OutOfOfficeHours: "00:00",
		}
	}

	sort.SliceStable(logs, func(i, j int) bool {
		return logs[i].Timestamp.Before(
			logs[j].Timestamp,
		)
	})

	var (
		workingDuration     time.Duration
		outOfOfficeDuration time.Duration

		firstCheckIn *time.Time
		lastCheckOut *time.Time

		currentCheckIn *time.Time
	)

	for _, log := range logs {

		if log.Timestamp.IsZero() {
			continue
		}

		punch, valid := normalizePunchType(log)

		if !valid {
			continue
		}

		timestamp := log.Timestamp

		switch punch {

		case 0: // CHECK_IN

			if firstCheckIn == nil {
				t := timestamp
				firstCheckIn = &t
			}

			currentCheckIn = &timestamp

		case 1: // CHECK_OUT

			if currentCheckIn != nil {

				duration := timestamp.Sub(
					*currentCheckIn,
				)

				if duration > 0 {
					workingDuration += duration
				}

				currentCheckIn = nil
			}

			if lastCheckOut == nil ||
				timestamp.After(*lastCheckOut) {

				t := timestamp
				lastCheckOut = &t
			}
		}
	}

	// Calculate OUT -> IN periods separately.
	for i := 0; i < len(logs)-1; i++ {

		currentPunch, currentValid := normalizePunchType(logs[i])

		nextPunch, nextValid := normalizePunchType(logs[i+1])

		if !currentValid || !nextValid {
			continue
		}

		if currentPunch == 1 && nextPunch == 0 {

			outDuration := logs[i+1].Timestamp.Sub(logs[i].Timestamp)

			if outDuration > 0 {
				outOfOfficeDuration += outDuration
			}
		}
	}

	status := "ABSENT"

	if firstCheckIn != nil {
		status = "PRESENT"
	}

	return entity.WorkingHoursSummary{
		EmpID:   employee.EmpID,
		EmpName: employee.EmpName,
		Date:    targetDate.Format("2006-01-02"),

		CheckInTime:  firstCheckIn,
		CheckOutTime: lastCheckOut,

		WorkingHours: validation.FormatWorkingHours(workingDuration),

		OutOfOfficeHours: validation.FormatWorkingHours(outOfOfficeDuration),

		Status: status,
	}
}
func normalizePunchType(log entity.AttendanceLog) (int, bool) {

	attendanceState := strings.ToUpper(strings.TrimSpace(log.AttendanceState))

	normalizedState := strings.NewReplacer("_", "", "-", "", " ", "").Replace(attendanceState)

	switch normalizedState {
	case "CHECKIN", "IN":
		return 0, true

	case "CHECKOUT", "OUT":
		return 1, true
	}

	if log.Punch == 0 || log.Punch == 1 {
		return log.Punch, true
	}

	return -1, false
}
