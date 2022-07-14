SELECT 
    state.state_id
    state.state_name
    state.single_deduction
    state.married_deduction
    state.single_exemption
    state.married_exemption
    state.dependent_exemption
    state_brackets.single_rate
    state_brackets.single_bracket
    state_brackets.married_rate
    state_brackets.married_bracket
FROM state INNER JOIN state_brackets ON state.state_id = state_bracket.state_id