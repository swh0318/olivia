package start

import (
	"fmt"
	"strings"
	"time"

	"github.com/olivia-ai/olivia/user"
)

func init() {
	RegisterModule(Module{
		Action: CheckReminders,
	})
}

// CheckReminders will check the dates of the user's reminder and if they are outdated, remove them
func CheckReminders(token string) {
	reminders := user.GetUserInformation(token).Reminders
	var messages []string

	// Iterate through the reminders to check if they are outdated
	for i, reminder := range reminders {
		date, _ := time.Parse("01/02/2006 03:04", reminder.Date)

		now := time.Now()
		// If the date is today
		if date.Year() == now.Year() && date.Day() == now.Day() && date.Month() == now.Month() {
			messages = append(messages, fmt.Sprintf("“%s”", reminder.Reason))
		}

		// If the date of the reminder has already been passed
		if time.Now().After(date) {
			user.ChangeUserInformation(token, func(information user.Information) user.Information {
				reminders := information.Reminders

				// Removes the element from the reminders slice
				if len(reminders) == 1 {
					reminders = []user.Reminder{}
				} else {
					reminders[i] = reminders[len(reminders)-1]
					reminders = reminders[:len(reminders)-1]
				}

				// Set the updated slice
				information.Reminders = reminders

				return information
			})
		}
	}

	// Send the startup message!
	if len(messages) != 0 {
		SendMessage(fmt.Sprintf(
			"Hello %s! For today you have these reminders: %s.",
			user.GetUserInformation(token).Name,
			strings.Join(messages, ", "),
		))
	} else {
		SendMessage(fmt.Sprintf("Hello %s!", user.GetUserInformation(token).Name))
	}
}