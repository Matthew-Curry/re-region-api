package services

// function used by both the state and federal services to get
// taxable income based on income, deductions, exemptions, and dependents
func getTaxableIncome(income, deduction, exemptions, dependents int) int{
	income = income - deduction - dependents*exemptions
	if income < 0 {
		return 0
	}

	return income 
}