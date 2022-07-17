SELECT 
    states.state_id,
    states.state_name,
    states.single_deduction,
    states.married_deduction,
    states.single_exemption,
    states.married_exemption,
    states.dependent_exemption,
    state_brackets.single_rate,
    state_brackets.single_bracket,
    state_brackets.married_rate,
    state_brackets.married_bracket
FROM states INNER JOIN state_brackets ON states.state_id = state_brackets.state_id;