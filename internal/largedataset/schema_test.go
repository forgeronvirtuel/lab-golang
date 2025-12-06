package largedataset

import "testing"

// Example: Creating a custom CSV schema for user data
func ExampleCSVSchema_custom() {
	minAge := 18.0
	maxAge := 120.0

	userSchema := &CSVSchema{
		MinColumns:    5,
		StrictColumns: true,
		Columns: []ColumnDef{
			{
				Index:    0,
				Name:     "UserID",
				Type:     TypeInt,
				Required: true,
				Min:      &minAge, // reusing variable, means >= 18
			},
			{
				Index:     1,
				Name:      "Email",
				Type:      TypeEmail,
				Required:  true,
				MaxLength: 100,
			},
			{
				Index:    2,
				Name:     "Age",
				Type:     TypeInt,
				Required: true,
				Min:      &minAge,
				Max:      &maxAge,
			},
			{
				Index:       3,
				Name:        "Status",
				Type:        TypeString,
				Required:    true,
				AllowedVals: []string{"active", "inactive", "pending"},
			},
			{
				Index:      4,
				Name:       "JoinDate",
				Type:       TypeDate,
				Required:   true,
				DateFormat: "2006-01-02",
			},
		},
	}

	// Validate a record
	validRecord := []string{"1001", "user@example.com", "25", "active", "2023-01-15"}
	if err := userSchema.ValidateRecord(validRecord); err != nil {
		panic(err)
	}

	// This will fail validation
	invalidRecord := []string{"1002", "invalid-email", "15", "unknown", "2023-01-15"}
	_ = userSchema.ValidateRecord(invalidRecord) // Returns error
}

// TestSchemaValidation demonstrates various validation scenarios
func TestSchemaValidation(t *testing.T) {
	schema := NewStockDataSchema()

	tests := []struct {
		name    string
		record  []string
		wantErr bool
	}{
		{
			name: "valid record",
			record: []string{
				"1", "2020-01-01 09:30:01", "AAPL", "NYSE", "Technology",
				"Buy", "Market", "100", "150.50", "15050.00",
				"150.00", "150.50", "151.00", "149.50", "1000000",
				"2500000000000", "25.5", "1.5", "1.2", "180.00", "120.00", "0.33",
			},
			wantErr: false,
		},
		{
			name: "invalid exchange",
			record: []string{
				"1", "2020-01-01 09:30:01", "AAPL", "INVALID", "Technology",
				"Buy", "Market", "100", "150.50", "15050.00",
				"150.00", "150.50", "151.00", "149.50", "1000000",
				"2500000000000", "25.5", "1.5", "1.2", "180.00", "120.00", "0.33",
			},
			wantErr: true,
		},
		{
			name: "negative price",
			record: []string{
				"1", "2020-01-01 09:30:01", "AAPL", "NYSE", "Technology",
				"Buy", "Market", "100", "-150.50", "15050.00",
				"150.00", "150.50", "151.00", "149.50", "1000000",
				"2500000000000", "25.5", "1.5", "1.2", "180.00", "120.00", "0.33",
			},
			wantErr: true,
		},
		{
			name: "invalid trade type",
			record: []string{
				"1", "2020-01-01 09:30:01", "AAPL", "NYSE", "Technology",
				"Trade", "Market", "100", "150.50", "15050.00",
				"150.00", "150.50", "151.00", "149.50", "1000000",
				"2500000000000", "25.5", "1.5", "1.2", "180.00", "120.00", "0.33",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := schema.ValidateRecord(tt.record)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
