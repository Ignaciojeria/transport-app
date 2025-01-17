package domain

import (
	"testing"
)

func TestCalculateServiceTime_SmallSingleItem(t *testing.T) {
	order := Order{
		Items: []Item{
			{
				ReferenceID: "small-item",
				Dimensions: Dimensions{
					Height: 0.1, // metros
					Width:  0.1, // metros
					Depth:  0.1, // metros
				},
				Weight: Weight{
					Value: 0.5, // kilogramos
				},
				Quantity: Quantity{
					QuantityNumber: 1,
				},
			},
		},
	}

	expectedServiceTime := 5.0 + (2.0 * 0.001) + (1.0 * 0.5 / 10)
	actualServiceTime, err := order.CalculateServiceTime()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if actualServiceTime != expectedServiceTime {
		t.Errorf("expected %.5f, got %.5f", expectedServiceTime, actualServiceTime)
	}
}

func TestCalculateServiceTime_LargeSingleItem(t *testing.T) {
	order := Order{
		Items: []Item{
			{
				ReferenceID: "large-item",
				Dimensions: Dimensions{
					Height: 2.0, // metros
					Width:  1.5, // metros
					Depth:  1.2, // metros
				},
				Weight: Weight{
					Value: 150.0, // kilogramos
				},
				Quantity: Quantity{
					QuantityNumber: 1,
				},
			},
		},
	}

	expectedServiceTime := 5.0 + (2.0 * 3.6) + (1.0 * 150.0 / 10)
	actualServiceTime, err := order.CalculateServiceTime()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if actualServiceTime != expectedServiceTime {
		t.Errorf("expected %.2f, got %.2f", expectedServiceTime, actualServiceTime)
	}
}

func TestCalculateServiceTime_SmallMultipleItems(t *testing.T) {
	order := Order{
		Items: []Item{
			{
				ReferenceID: "small-item-1",
				Dimensions: Dimensions{
					Height: 0.2, // metros
					Width:  0.2, // metros
					Depth:  0.2, // metros
				},
				Weight: Weight{
					Value: 1.0, // kilogramos
				},
				Quantity: Quantity{
					QuantityNumber: 3,
				},
			},
		},
	}

	expectedServiceTime := 5.0 + (2.0 * (0.008 * 3)) + (1.0 * (1.0 * 3) / 10)
	actualServiceTime, err := order.CalculateServiceTime()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if actualServiceTime != expectedServiceTime {
		t.Errorf("expected %.5f, got %.5f", expectedServiceTime, actualServiceTime)
	}
}

func TestCalculateServiceTime_LargeMultipleItems(t *testing.T) {
	order := Order{
		Items: []Item{
			{
				ReferenceID: "large-item-1",
				Dimensions: Dimensions{
					Height: 1.0, // metros
					Width:  1.0, // metros
					Depth:  1.0, // metros
				},
				Weight: Weight{
					Value: 50.0, // kilogramos
				},
				Quantity: Quantity{
					QuantityNumber: 2,
				},
			},
			{
				ReferenceID: "large-item-2",
				Dimensions: Dimensions{
					Height: 2.0, // metros
					Width:  2.0, // metros
					Depth:  2.0, // metros
				},
				Weight: Weight{
					Value: 200.0, // kilogramos
				},
				Quantity: Quantity{
					QuantityNumber: 1,
				},
			},
		},
	}

	expectedServiceTime := 5.0 +
		(2.0 * ((1.0 + 8.0) * 2)) +
		(1.0 * ((50.0*2 + 200.0) / 10))

	actualServiceTime, err := order.CalculateServiceTime()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if actualServiceTime != expectedServiceTime {
		t.Errorf("expected %.2f, got %.2f", expectedServiceTime, actualServiceTime)
	}
}

func TestCalculateServiceTime_InvalidItems(t *testing.T) {
	order := Order{
		Items: []Item{
			{
				ReferenceID: "invalid-item",
				Dimensions: Dimensions{
					Height: 0.0, // Altura inválida
					Width:  0.5,
					Depth:  0.5,
				},
				Weight: Weight{
					Value: 10.0, // kilogramos
				},
				Quantity: Quantity{
					QuantityNumber: 1,
				},
			},
		},
	}

	_, err := order.CalculateServiceTime()
	if err == nil {
		t.Fatalf("expected error for invalid dimensions, got nil")
	}
}
