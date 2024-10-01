package test

import (
	"testing"

	order_business "github.com/CamiloC-pvt/re-challenge/app/order/business"
	order_domain "github.com/CamiloC-pvt/re-challenge/app/order/domain"
	pack_domain "github.com/CamiloC-pvt/re-challenge/app/pack/domain"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestCalculatePackaging_ReSuggested(t *testing.T) {
	// Packages
	testPackList := []pack_domain.Pack{
		{Size: 250},
		{Size: 500},
		{Size: 1000},
		{Size: 2000},
		{Size: 5000},
	}

	// Test Cases
	testCasesMap := map[int32][]order_domain.OrderPack{
		1: {
			{Amount: 1, Size: 250},
		},
		250: {
			{Amount: 1, Size: 250},
		},
		251: {
			{Amount: 1, Size: 500},
		},
		501: {
			{Amount: 1, Size: 500},
			{Amount: 1, Size: 250},
		},
		12001: {
			{Amount: 2, Size: 5000},
			{Amount: 1, Size: 2000},
			{Amount: 1, Size: 250},
		},
	}

	// Result Check
	comparePacks := func(a, b []order_domain.OrderPack) bool {
		same := true

		for _, aP := range a {
			found := false

			for _, bP := range b {
				if aP.Amount == bP.Amount && aP.Size == bP.Size {
					found = true
					break
				}
			}

			if !found {
				same = false
				break
			}
		}

		return same
	}

	// Run
	business := order_business.NewOrderBusiness(nil, nil)

	for sizeToTest, testResult := range testCasesMap {
		result := business.CalculatePackaging(testPackList, sizeToTest)

		if !comparePacks(result, testResult) {
			t.Errorf("Obtained packs for '%d' are '%+v', expected was: %+v", sizeToTest, result, testResult)
		}
	}
}

func TestCalculatePackaging_Others(t *testing.T) {
	// Packages
	testPackList := []pack_domain.Pack{
		{Size: 250},
		{Size: 500},
		{Size: 1000},
		{Size: 2000},
		{Size: 5000},
		{Size: 13001},
	}

	// Test Cases
	testCasesMap := map[int32][]order_domain.OrderPack{
		999: {
			{Amount: 1, Size: 1000},
		},
		800006: {
			{Amount: 6, Size: 13001},
			{Amount: 144, Size: 5000},
			{Amount: 1, Size: 2000},
		},
	}

	// Result Check
	comparePacks := func(a, b []order_domain.OrderPack) bool {
		same := true

		for _, aP := range a {
			found := false

			for _, bP := range b {
				if aP.Amount == bP.Amount && aP.Size == bP.Size {
					found = true
					break
				}
			}

			if !found {
				same = false
				break
			}
		}

		return same
	}

	// Run
	business := order_business.NewOrderBusiness(nil, nil)

	for sizeToTest, testResult := range testCasesMap {
		result := business.CalculatePackaging(testPackList, sizeToTest)

		if !comparePacks(result, testResult) {
			t.Errorf("Obtained packs for '%d' are '%+v', expected was: %+v", sizeToTest, result, testResult)
		}
	}
}

func TestCalculatePackaging_ExtraSuggested(t *testing.T) {
	// Packages
	testPackList := []pack_domain.Pack{
		{Size: 5},
		{Size: 12},
	}

	// Test Cases
	testCasesMap := map[int32][]order_domain.OrderPack{
		15: {
			{Amount: 3, Size: 5},
		},
	}

	// Result Check
	comparePacks := func(a, b []order_domain.OrderPack) bool {
		same := true

		for _, aP := range a {
			found := false

			for _, bP := range b {
				if aP.Amount == bP.Amount && aP.Size == bP.Size {
					found = true
					break
				}
			}

			if !found {
				same = false
				break
			}
		}

		return same
	}

	// Run
	business := order_business.NewOrderBusiness(nil, nil)

	for sizeToTest, testResult := range testCasesMap {
		result := business.CalculatePackaging(testPackList, sizeToTest)

		if !comparePacks(result, testResult) {
			t.Errorf("Obtained packs for '%d' are '%+v', expected was: %+v", sizeToTest, result, testResult)
		}
	}
}
