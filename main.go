package main

import (
	"fmt"
	"sort"
	"time"
)

type Cuisine int

const (
	SouthIndian Cuisine = iota
	NorthIndian
	Chinese
)

type Restaurant struct {
	restaurantId  string
	cuisine       Cuisine
	costBracket   int
	rating        float32
	isRecommended bool
	onboardedTime time.Time
}

type CuisineTracking struct {
	cuisine    int
	noOfOrders int
}

type CostTracking struct {
	costBracket int
	noOfOrders  int
}

type User struct {
	cuisines     []CuisineTracking
	costBrackets []CostTracking
}

func getRestaurantRecommendations(user User, availableRestaurants []Restaurant) []string {
	// Helper function to check if a restaurant is recommended
	isRecommended := func(restaurant Restaurant) bool {
		return restaurant.isRecommended
	}

	// Helper function to check if a restaurant is newly created
	isNewlyCreated := func(restaurant Restaurant) bool {
		return restaurant.onboardedTime.After(time.Now().Add(-48 * time.Hour))
	}

	// Helper function to sort restaurants by rating in descending order
	byRating := func(restaurants []Restaurant) {
		sort.Slice(restaurants, func(i, j int) bool {
			return restaurants[i].rating > restaurants[j].rating
		})
	}

	// Helper function to get top N restaurants by rating
	getTopNRatedRestaurants := func(restaurants []Restaurant, n int) []Restaurant {
		byRating(restaurants)
		if n > len(restaurants) {
			n = len(restaurants)
		}
		return restaurants[:n]
	}

	// Filter restaurants by primary cuisine and primary cost bracket
	filterByPrimaryCuisineAndCostBracket := func(restaurants []Restaurant, cuisine int, costBracket int) []Restaurant {
		filtered := make([]Restaurant, 0)
		for _, restaurant := range restaurants {
			if int(restaurant.cuisine) == cuisine && restaurant.costBracket == costBracket {
				filtered = append(filtered, restaurant)
			}
		}
		return filtered
	}

	// Filter restaurants by cuisine, cost bracket, and minimum rating
	filterByCuisineCostAndRating := func(restaurants []Restaurant, cuisine int, costBracket int, minRating float32) []Restaurant {
		filtered := make([]Restaurant, 0)
		for _, restaurant := range restaurants {
			if int(restaurant.cuisine) == cuisine && restaurant.costBracket == costBracket && restaurant.rating >= minRating {
				filtered = append(filtered, restaurant)
			}
		}
		return filtered
	}

	// Sort restaurants based on the specified logic
	sortedRestaurants := make([]Restaurant, 0)
	primaryCuisine := int(user.cuisines[0].cuisine)
	primaryCostBracket := user.costBrackets[0].costBracket

	// Condition 1
	featuredPrimary := filterByPrimaryCuisineAndCostBracket(availableRestaurants, primaryCuisine, primaryCostBracket)
	if len(featuredPrimary) > 0 {
		sortedRestaurants = append(sortedRestaurants, featuredPrimary...)
	} else {
		featuredSecondary := filterByPrimaryCuisineAndCostBracket(availableRestaurants, primaryCuisine, int(user.cuisines[1].cuisine))
		sortedRestaurants = append(sortedRestaurants, featuredSecondary...)
		if len(featuredSecondary) == 0 {
			secondaryPrimary := filterByPrimaryCuisineAndCostBracket(availableRestaurants, int(user.cuisines[1].cuisine), primaryCostBracket)
			sortedRestaurants = append(sortedRestaurants, secondaryPrimary...)
		}
	}

	// Condition 2
	primaryRated := filterByCuisineCostAndRating(availableRestaurants, primaryCuisine, primaryCostBracket, 4.0)
	sortedRestaurants = append(sortedRestaurants, primaryRated...)

	// Condition 3
	secondaryRated := filterByCuisineCostAndRating(availableRestaurants, primaryCuisine, int(user.cuisines[1].cuisine), 4.5)
	sortedRestaurants = append(sortedRestaurants, secondaryRated...)

	// Condition 4
	secondaryCuisinePrimary := filterByCuisineCostAndRating(availableRestaurants, int(user.cuisines[2].cuisine), primaryCostBracket, 4.5)
	sortedRestaurants = append(sortedRestaurants, secondaryCuisinePrimary...)

	// Condition 5
	newlyCreated := getTopNRatedRestaurants(filterByCuisineCostAndRating(availableRestaurants, primaryCuisine, primaryCostBracket, 0.0), 4)
	sortedRestaurants = append(sortedRestaurants, newlyCreated...)

	// Condition 6
	primaryLowRated := filterByCuisineCostAndRating(availableRestaurants, primaryCuisine, primaryCostBracket, 4.0)
	sortedRestaurants = append(sortedRestaurants, primaryLowRated...)

	// Condition 7
	primarySecondaryLowRated := filterByCuisineCostAndRating(availableRestaurants, primaryCuisine, int(user.cuisines[1].cuisine), 4.5)
	sortedRestaurants = append(sortedRestaurants, primarySecondaryLowRated...)

	// Condition 8
	secondaryCuisinePrimaryLowRated := filterByCuisineCostAndRating(availableRestaurants, int(user.cuisines[2].cuisine), primaryCostBracket, 4.5)
	sortedRestaurants = append(sortedRestaurants, secondaryCuisinePrimaryLowRated...)

	// Condition 9
	remainingRestaurants := make([]Restaurant, 0)
	for _, restaurant := range availableRestaurants {
		if !isRecommended(restaurant) && !isNewlyCreated(restaurant) {
			remainingRestaurants = append(remainingRestaurants, restaurant)
		}
	}
	sortedRestaurants = append(sortedRestaurants, remainingRestaurants...)

	// Extract restaurant IDs and return the sorted list
	restaurantIDs := make([]string, 0)
	for _, restaurant := range sortedRestaurants {
		restaurantIDs = append(restaurantIDs, restaurant.restaurantId)
	}
	return restaurantIDs
}

func main() {
	// Example usage
	user := User{
		cuisines: []CuisineTracking{
			{cuisine: int(SouthIndian), noOfOrders: 3},
			{cuisine: int(NorthIndian), noOfOrders: 2},
			{cuisine: int(Chinese), noOfOrders: 1},
		},
		costBrackets: []CostTracking{
			{costBracket: 3, noOfOrders: 5},
			{costBracket: 2, noOfOrders: 4},
			{costBracket: 1, noOfOrders: 3},
		},
	}

	restaurants := []Restaurant{
		{restaurantId: "1", cuisine: SouthIndian, costBracket: 3, rating: 4.2, isRecommended: true, onboardedTime: time.Now().Add(-72 * time.Hour)},
		{restaurantId: "2", cuisine: NorthIndian, costBracket: 1, rating: 4.5, isRecommended: true, onboardedTime: time.Now().Add(-24 * time.Hour)},
		{restaurantId: "3", cuisine: Chinese, costBracket: 3, rating: 3.8, isRecommended: false, onboardedTime: time.Now().Add(-60 * time.Hour)},
		{restaurantId: "4", cuisine: Chinese, costBracket: 2, rating: 4.8, isRecommended: true, onboardedTime: time.Now().Add(-30 * time.Hour)},
		{restaurantId: "5", cuisine: SouthIndian, costBracket: 3, rating: 3.9, isRecommended: false, onboardedTime: time.Now().Add(-10 * time.Hour)},
		{restaurantId: "6", cuisine: NorthIndian, costBracket: 1, rating: 4.3, isRecommended: true, onboardedTime: time.Now().Add(-15 * time.Hour)},
		{restaurantId: "7", cuisine: SouthIndian, costBracket: 3, rating: 4.2, isRecommended: false, onboardedTime: time.Now().Add(-5 * time.Hour)},
	}

	recommendations := getRestaurantRecommendations(user, restaurants)
	fmt.Println(recommendations)
}
