package http

import (
	"fmt"
	"regexp"
	model "ticket-tracker/internal/http/dtos/tcdd"
	"time"
)

func ValidateRequest(request *model.SearchTrainRequest) error {
	for _, detail := range request.Request {
		if err := validateIdentityNumberLength(detail); err != nil {
			return err
		}
		if err := validateEmail(detail); err != nil {
			return err
		}
		if err := validatePhone(detail); err != nil {
			return err
		}
		if err := validateGender(detail); err != nil {
			return err
		}
		if err := validateNameAndSurname(detail); err != nil {
			return err
		}
		if err := validateBirthDate(detail); err != nil {
			return err
		}
		if err := validateDepartureDateFormat(detail); err != nil {
			return err
		}
	}
	return nil
}

func validateIdentityNumberLength(request model.SearchTrainRequestDetail) error {
	if len(request.IdentityNumber) != 11 {
		return fmt.Errorf("identity number length must be 11")
	}
	return nil
}

func validateEmail(request model.SearchTrainRequestDetail) error {
	if request.Email == "" {
		return fmt.Errorf("email is required")
	}
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	validationResult := regexp.MustCompile(emailRegex).MatchString(request.Email)

	if !validationResult {
		return fmt.Errorf("email is not valid")
	}
	return nil
}

func validatePhone(request model.SearchTrainRequestDetail) error {
	if request.Phone == "" {
		return fmt.Errorf("phone is required")
	}
	if len(request.Phone) != 10 {
		return fmt.Errorf("phone length must be 10")
	}
	if request.Phone[0] != '5' {
		return fmt.Errorf("phone must start with 5")
	}

	// check contains character
	for _, c := range request.Phone {
		if c < '0' || c > '9' {
			return fmt.Errorf("phone must contain only numbers")
		}
	}

	return nil
}

func validateGender(request model.SearchTrainRequestDetail) error {
	if request.Gender != "M" && request.Gender != "F" {
		return fmt.Errorf("gender must be M or F")
	}
	return nil
}

func validateNameAndSurname(request model.SearchTrainRequestDetail) error {
	if request.Name == "" {
		return fmt.Errorf("name is required")
	}
	if request.LastName == "" {
		return fmt.Errorf("last name is required")
	}

	return nil

}

func validateBirthDate(request model.SearchTrainRequestDetail) error {
	if request.BirthDate == "" {
		return fmt.Errorf("birth date is required")
	}
	// check contains character
	for _, c := range request.BirthDate {
		if c < '0' || c > '9' {
			return fmt.Errorf("birth date must contain only numbers")
		}
	}

	// check length
	if len(request.BirthDate) != 8 {
		return fmt.Errorf("birth date length must be 8")
	}
	return nil

}

func validateDepartureDateFormat(request model.SearchTrainRequestDetail) error {
	if request.DepartureDate == "" {
		return fmt.Errorf("departure date is required")
	}
	departureDate, err := time.Parse("2006-01-02 15:04:05", request.DepartureDate)
	if err != nil {
		return fmt.Errorf("departure date is not valid")
	}
	if time.Now().After(departureDate) {
		return fmt.Errorf("departure date must be greater than current date")
	}
	return nil
}
