package enum

type Plan int

var (
	//PlanFree is used for tenants on the free plan
	PlanFree Plan = 1
	//PlanPro is used for tenants on the pro plan
	PlanPro Plan = 2
)

var planIDs = map[Plan]string{
	PlanFree: "free",
	PlanPro:  "pro",
}

// String returns the string version of the plan
func (p Plan) String() string {
	return planIDs[p]
}
