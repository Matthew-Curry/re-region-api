package test

import  (
	"github.com/Matthew-Curry/re-region-api/src/model"
)

var tl = append(make([]model.TaxLocale,0), model.TaxLocale{
	Locale_id:   3376,
	Locale_name: "New York City",
	Total_tax:   0,
	Federal_tax: 0,
	State_tax:   0,
	Locale_tax:  0,
})

var exCounty = &model.County{
	County_id:     36061,
	County_name:   "New York County",
	State_id:      36,
	State_name:    "New York",
	Pop:           1628706,
	Male_pop:      771278,
	Female_pop:    857428,
	Median_income: 93651,
	Average_rent:  1753,
	Commute:       81,
	Tax_locale:    tl,
}
var cmp = model.CountyMetricPair{
	County_id:    36061,
	County_name:  "New York County",
	State_id:     36,
	State_name:   "New York",
	Metric_value: 81,
}

var rl = append(make([]model.CountyMetricPair, 0), cmp)

var exCountyList  = &model.CountyList{
	Metric_name : "metric",
	Ranked_list: rl,
}

var tli = append(make([]model.TaxLocaleInfo,0), model.TaxLocaleInfo{
	Locale_id  : 3376,
	Local_name: "New York City",
	Resident_desc:           "3.078% - 3.876%",
	Resident_rate:           0,
	Resident_month_fee:      0,
	Resident_year_fee:       0,
	Resident_pay_period_fee: 0,
	Resident_state_rate:     0,
	Nonresident_desc:       "0.00%",
	Nonresident_rate:      0,
	Nonresident_month_fee:     0,
	Nonresident_year_fee:    0,
	Nonresident_pay_period_fee: 0,
	Nonresident_state_rate:    0,
})

var exCountyTaxList = &model.CountyTaxList{
	County_name :"New York County",
	County_id   :36061,
	State_name  :"New York",
	State_id    :36,
	Tax_locales :tli,
}

var exState = &model.State{
	State_id   :36,
	State_name :"New York",
	Pop           :18466230,
	Male_pop      :8953064,
	Female_pop    :9513166,
	Median_income :77578,
	Average_rent  :1381,
	Commute       :17,
	Total_tax   :0,
	State_tax   :0,
	Federal_tax :0,
}

var mpList = model.StateMetricPair{
	State_id: 36,
	State_name: "New York",
	Metric_value: 17,

}

var exStateList = model.GetMetricStateList("commute")

var bracket1 = model.StateBracket {
	Single_rate     :0.02,
	Single_bracket  :0,
	Married_rate    :0.02,
	Married_bracket :0,
}

var bracket2 = model.StateBracket {
	Single_rate     :0.12,
	Single_bracket  :500,
	Married_rate    :0.12,
	Married_bracket :1000,
}

var exStateTaxInfoId = &model.StateTaxInfo{
	State_id   :36,
	State_name :"New York",
	Single_deduction    :2500,
	Married_deduction   :7500,
	Single_exemption    :1500,
	Married_exemption   :3000,
	Dependent_exemption :1000,
}

var exStateTaxInfoName = &model.StateTaxInfo{
	State_id   :36,
	State_name :"New York",
	Single_deduction    :2500,
	Married_deduction   :7500,
	Single_exemption    :1500,
	Married_exemption   :3000,
	Dependent_exemption :1000,
}

var fb1 = model.FederalBracket{
	Rate            :0.1,
	Single_bracket  :0,
	Married_bracket :0,
	Head_bracket    :0,
}

var fb2 = model.FederalBracket{
	Rate            :0.12,
	Single_bracket  :10275,
	Married_bracket :20550,
	Head_bracket    :14650,
}

var fb3 = model.FederalBracket{
	Rate            :0.22,
	Single_bracket  :41775,
	Married_bracket :83550,
	Head_bracket    :55900,
}


var federalTaxInfo = &model.FederalTaxInfo{
	Single_deduction  :12950,
	Married_deduction :25900,
	Head_deduction    :19400,
}