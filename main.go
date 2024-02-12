package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/xuri/excelize/v2"
)

type Meal struct {
	Day   string
	Date  string
	Meal  string
	Items []string
}

func NewMeal(day, date, meal string, items []string) *Meal {
	return &Meal{
		Day:   day,
		Date:  date,
		Meal:  meal,
		Items: items,
	}
}

func (m *Meal) PrintDetails() {
	fmt.Printf("Day: %s\n", m.Day)
	fmt.Printf("Date: %s\n", m.Date)
	fmt.Printf("Meal: %s\n", m.Meal)
	fmt.Println("Items:")
	for _, item := range m.Items {
		fmt.Printf("  - %s\n", item)
	}

}

func ConvertToMeals(menu []map[string]interface{}) []*Meal {
	var meals []*Meal

	for _, day := range menu {
		dayName := day["Day"].(string)
		dayDate := day["Date"].(string)
		delete(day, "Day")
		delete(day, "Date")
		for mealType, items := range day {
			if mealType != "Day" {
				meal := strings.ToUpper(mealType) // Capitalize the meal type
				var mealItems []string
				for _, item := range items.([]string) { // Convert each item to string
					mealItems = append(mealItems, item)
				}
				mealStruct := NewMeal(dayName, dayDate, meal, mealItems)
				meals = append(meals, mealStruct)
			}
		}
	}

	return meals
}

func interfaceSliceToStringSlice(slice []interface{}) []string {
	strSlice := make([]string, len(slice))
	for i, v := range slice {
		strSlice[i] = v.(string)
	}
	return strSlice
}

func WriteToJson(menu []map[string]interface{}) {

	jsonData1, err := json.MarshalIndent(menu, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Write JSON data to a file
	err = ioutil.WriteFile("data.json", jsonData1, 0644)
	if err != nil {
		fmt.Println("Error writing JSON to file:", err)
		return
	}

	fmt.Println("JSON data written to data.json")

}

func PrintMenu(menu []map[string]interface{}) {
	// Convert menu slice to JSON
	jsonData, err := json.MarshalIndent(menu, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print formatted JSON
	fmt.Println("Menu:")
	fmt.Println(string(jsonData))
	WriteToJson(menu)
}

func ConvertToJson(food map[string][]string) []map[string]interface{} {

	// Convert it into the following type
	/*
		[
			{
				"Day" = Sunday,
				"Breakfast" = [Choice of Egg, Milk, Sandwich]
				"Lunch" = [Choice of Egg, Milk, Sandwich]
				"Dinner" = [Choice of Egg, Milk, Sandwich]
			},

			{
				"Day" = Monday
				"Breakfast" = [Choice of Egg, Milk, Sandwich]
				"Lunch" = [Choice of Egg, Milk, Sandwich]
				"Dinner" = [Choice of Egg, Milk, Sandwich]
			},

			{
				"Day" = Tuesday
				"Breakfast" = [Choice of Egg, Milk, Sandwich]
				"Lunch" = [Choice of Egg, Milk, Sandwich]
				"Dinner" = [Choice of Egg, Milk, Sandwich]
			},
		]
	*/

	// Define a slice to hold the converted data
	var menu []map[string]interface{}

	// Iterate over the original data structure
	for key, meals := range food {
		// Extract the day of the week from the key
		day := key[:len(key)-1]

		// Initialize a map for the current day
		dayMap := make(map[string]interface{})
		dayMap["Day"] = day

		// Determine meal type based on the last character of the key
		mealType := key[len(key)-1:]
		if mealType == "b" {
			dayMap["Date"] = meals[0]
			dayMap["Breakfast"] = meals[1:]
		} else if mealType == "l" {
			dayMap["Lunch"] = meals
		} else if mealType == "d" {
			dayMap["Dinner"] = meals
		}

		// Check if the current day already exists in the menu slice
		existingIndex := -1
		for i, menuDay := range menu {
			if menuDay["Day"] == day {
				existingIndex = i
				break
			}
		}

		// If the day doesn't exist in the menu slice, add it
		if existingIndex == -1 {
			menu = append(menu, dayMap)
		} else {
			// Merge the current day's meals into the existing entry
			for meal, foods := range dayMap {
				menu[existingIndex][meal] = foods
			}
		}
	}

	// Print the converted menu
	PrintMenu(menu)
	// Write it to the Json file
	WriteToJson(menu)

	return menu

}

func extractxlsx() map[string][]string {

	// Declaring a map for the data
	var food map[string][]string
	food = make(map[string][]string)

	// Extracting data from the excel file
	f, err := excelize.OpenFile("wtf.xlsx")
	if err != nil {
		fmt.Println(err)
		return food
	}

	// Get all the columns from the Sheet1
	cols, err := f.GetCols("Sheet1")
	if err != nil {
		fmt.Println(err)
		return food
	}

	/* Data is structured as follows
	food = {"Sundayb": { "milk", "egg", "bread"},
			"Sundayb": { "milk", "egg", "bread"},
			"Sundayb": { "milk", "egg", "bread"},
			"Sundayb": { "milk", "egg", "bread"},
				}
	*/

	var daySet bool = false
	var mealno int = 0
	var dayname string = ""

	// Iterate over the columns and append food items for each day and meal

	for _, col := range cols {
		daySet = false
		mealno = 0
		for _, rowCell := range col {

			_, ok1 := food[rowCell+"l"]
			_, ok2 := food[rowCell+"d"]

			if daySet == false {
				// Setting the day
				food[rowCell+"b"] = []string{}
				daySet = true
				dayname = rowCell

			} else if rowCell == "" {
				// if empty cell is returned
				continue

			} else if !ok1 && rowCell == dayname {
				//Checking if lunch menu has started
				mealno += 1

			} else if !ok2 && rowCell == dayname {
				// Checking if dinner menu has started
				mealno += 1

			} else if mealno == 0 {
				// appending food for breakfast
				food[dayname+"b"] = append(food[dayname+"b"], rowCell)

			} else if mealno == 1 {
				// appending food for lunch
				food[dayname+"l"] = append(food[dayname+"l"], rowCell)

			} else if mealno == 2 {
				// appending food for dinner
				food[dayname+"d"] = append(food[dayname+"d"], rowCell)
			}

		}
		fmt.Println()
	}
	return food

}

func finditemno(food map[string][]string) {

	var day, meal string

	fmt.Print("Enter day: ")
	fmt.Scanln(&day)

	fmt.Print("Enter meal: ")
	fmt.Scanln(&meal)

	// Capitalize the strings
	day = strings.ToUpper(day)
	meal = strings.ToUpper(meal)
	fmt.Printf("The day is %v and the meal required is %v\n", day, meal)

	// Checking for which meal
	if meal == "BREAKFAST" {
		fmt.Printf("The number of items is %v\n", len(food[day+"b"])-2)
	} else if meal == "LUNCH" {
		fmt.Printf("The number of items is %v\n", len(food[day+"l"])-1)
	} else {
		fmt.Printf("The number of items is %v\n", len(food[day+"d"])-1)
	}

}

func showitems(food map[string][]string) {

	var day, meal string

	fmt.Print("Enter day: ")
	fmt.Scanln(&day)

	fmt.Print("Enter meal: ")
	fmt.Scanln(&meal)

	// Capitalize the strings
	day = strings.ToUpper(day)
	meal = strings.ToUpper(meal)
	fmt.Printf("\nThe day is %v and the meal required is %v\n", day, meal)

	// Checking for which meal
	if meal == "BREAKFAST" {
		fmt.Printf("The items are %v\n", food[day+"b"][2:])
	} else if meal == "LUNCH" {
		fmt.Printf("The items are %v\n", food[day+"l"][1:])
	} else {
		fmt.Printf("The items are %v\n", food[day+"d"][1:])
	}

}

func checkitem(food map[string][]string) {

	var day, meal, item string

	fmt.Print("Enter day: ")
	fmt.Scanln(&day)

	fmt.Print("Enter meal: ")
	fmt.Scanln(&meal)

	// To Take item input from user . item can include more than one word

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the item: ")
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	item = strings.TrimSpace(line)
	fmt.Println("You entered:", item)

	// Capitalize the strings
	day = strings.ToUpper(day)
	meal = strings.ToUpper(meal)
	item = strings.ToUpper(item)

	// Checking for item in required meal

	var found bool
	found = false
	if meal == "BREAKFAST" {
		for _, itemcheck := range food[day+"b"] {
			if itemcheck == item {
				found = true
				break
			}
		}
		if found {
			fmt.Printf("Your item %v was found in %v %v", item, day, meal)
		} else {
			fmt.Printf("Your item is not found!")
		}
	} else if meal == "LUNCH" {
		for _, itemcheck := range food[day+"l"] {
			if itemcheck == item {
				found = true
				break
			}
		}
		if found {
			fmt.Printf("Your item %v was found in %v %v", item, day, meal)
		} else {
			fmt.Printf("Your item is not found!")
		}
	} else {
		for _, itemcheck := range food[day+"d"] {
			if itemcheck == item {
				found = true
				break
			}
		}
		if found {
			fmt.Printf("Your item %v was found in %v %v", item, day, meal)
		} else {
			fmt.Printf("Your item is not found!")
		}
	}
}

func main() {

	var food map[string][]string
	food = make(map[string][]string)
	food = extractxlsx()

	for {
		fmt.Printf("\n Menu \n 1.Show Items \n 2.Find Number of Items \n 3.Check Items \n 4. Write to Json file and create a struct  \n")
		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "1":
			showitems(food)
		case "2":
			finditemno(food)
		case "3":
			checkitem(food)
		case "4":
			var menu []map[string]interface{} = ConvertToJson(food)
			meals := ConvertToMeals(menu)

			// Sort meals by day and meal type
			sort.Slice(meals, func(i, j int) bool {
				if meals[i].Day != meals[j].Day {
					return meals[i].Day < meals[j].Day
				}
				return meals[i].Meal < meals[j].Meal
			})

			// Print meals in sorted order
			for _, meal := range meals {
				meal.PrintDetails()
				fmt.Println("###################################")
			}

		default:
			fmt.Println("\n\n## Invalid choice. Please select a valid option.##")
		}

		fmt.Printf("\n\n")
	}

}
