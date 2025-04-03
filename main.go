package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	simplesms "github.com/trriplejay/simple-sms-go"
)

// structure of the API response:
/**
{
  menuSchedules: [{
		menuBlocks: [{
		  blockName: "Breakfast",
			cafeteriaLineList: {
				data:[{
					foodItemList: {
						data: [{
							item_Name: "Milk, 1% white"
						},
						{
							item_Name: "Cereal, cinnamon toast crunch"
						}]
					}
				}]
			}
		},{
			blockName: "Lunch",
			cafeteriaLineList: {
				data:[{
					foodItemList: {
						data: [{
							item_Name: "Milk, 1% white"
						},
						{
							item_Name: "Crispy chicken tenders"
						}]
					}
				}]
			}
		}]
	}]
}
*/
type Menu struct {
	MenuSchedules []MenuSchedule `json:"menuSchedules"`
}

type MenuSchedule struct {
	MenuBlocks []MenuBlock `json:"menuBlocks"`
}

type MenuBlock struct {
	BlockName         string            `json:"blockName"`
	CafeteriaLineList CafeteriaLineList `json:"cafeteriaLineList"`
}

type CafeteriaLineList struct {
	Data []CafeteriaLine `json:"data"`
}

type CafeteriaLine struct {
	FoodItemList FoodItemList `json:"foodItemList"`
}

type FoodItemList struct {
	Data []FoodItem `json:"data"`
}

type FoodItem struct {
	ItemName string `json:"item_Name"`
}

func main() {
	args := os.Args

	if (len(args) > 4) || (len(args) < 4) {
		log.Fatal("Expected format: 'lunch <phone_number> <provider_string> <apikey>'")
	}

	phone := args[1]
	provider := args[2]
	apikey := args[3]

	err := simplesms.CheckProvider(provider)
	if err != nil {
		log.Fatal(err)
	}

	phoneInt, err := strconv.Atoi(phone)
	if err != nil {
		log.Fatal(err)
	}

	today := time.Now().Format("01-02-2006")

	menuString, err := getMenu(today)
	if err != nil {
		log.Fatalf("Failed to get menu: %s\n", err)
	}
	log.Printf("THE MESSAGE:\n%s", menuString)

	emailClient := simplesms.NewClient("stockbauer@gmail.com", apikey, "smtp.gmail.com", "587")
	err = emailClient.Send(phoneInt, provider, "today's menu", menuString)
	if err != nil {
		log.Fatalf("failed to send email: %s", err)
	}

}

func getMenu(dateString string) (string, error) {
	url := fmt.Sprintf("https://api.mealviewer.com/api/v4/school/Bryant/%s/%s/0", dateString, dateString)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("Error creating request:", err)
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making request:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return "nil", err
	}

	var menu Menu
	err = json.Unmarshal(body, &menu)
	if err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return "", err
	}

	return parseMenu(menu), nil
}

func parseMenu(menu Menu) string {

	breakfastString := "Breakfast:\n"
	lunchString := "Lunch:\n"
	for _, value := range menu.MenuSchedules[0].MenuBlocks {
		if value.BlockName == "Breakfast" {
			for _, bf := range value.CafeteriaLineList.Data[0].FoodItemList.Data {
				breakfastString += bf.ItemName + "\n"
			}

		} else if value.BlockName == "Lunch" {
			for _, bf := range value.CafeteriaLineList.Data[0].FoodItemList.Data {
				lunchString += bf.ItemName + "\n"
			}
		}
	}

	return fmt.Sprintf("\n\n%s\n%s", breakfastString, lunchString)

}
