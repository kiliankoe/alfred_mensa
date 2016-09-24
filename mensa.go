package main

import (
	"fmt"
	"os"

	"github.com/BenchR267/goalfred"
	"github.com/kiliankoe/swdd/speiseplan"
)

var searchMensa = os.Getenv("STANDARD_MENSA")
var isStudent = os.Getenv("STUDENTENPREIS")

func main() {
	defer goalfred.Print()

	arg, err := goalfred.NormalizedArguments()
	handleError(err)
	if arg[0] != "" {
		searchMensa = arg[0]
	}

	meals, canteenName, err := speiseplan.GetCurrentForCanteen(searchMensa)
	handleError(err)

	invalid := false
	if len(meals) > 0 {
		goalfred.Add(goalfred.Item{
			Title: fmt.Sprintf("Heute @ %s:", canteenName),
			Valid: &invalid,
		})
	} else {
		goalfred.Add(goalfred.Item{
			Title: fmt.Sprintf("Heute gibt's leider nichts @ %s 😢", canteenName),
			Valid: &invalid,
		})
		return
	}

	for _, meal := range meals {
		item := goalfred.Item{
			Title:    formatTitle(meal),
			Subtitle: meal.PageURL,
			Arg:      meal.PageURL,
		}
		item.Mod.Cmd = &goalfred.ModContent{
			Arg:      meal.ImageURL,
			Subtitle: "Bild der Mahlzeit öffnen (möglicherweise nur ein Platzhalter)",
		}
		goalfred.Add(item)
	}
}

func formatTitle(meal *speiseplan.Meal) (title string) {
	title += meal.Name
	price := ""
	if isStudent == "1" {
		price = formatPrice(meal.StudentPrice)
	} else {
		price = formatPrice(meal.EmployeePrice)
	}
	if price != "" {
		title += " - "
		title += price
	}
	return
}

func formatPrice(price float64) string {
	if price != 0 {
		return fmt.Sprintf("%.2f€", price)
	}
	return ""
}

func handleError(err error) {
	if err == nil {
		return
	}

	goalfred.Add(goalfred.Item{
		Title:    "Unerwarteter Fehler 😲",
		Subtitle: err.Error(),
	})
	goalfred.Print()

	os.Exit(1)
}
