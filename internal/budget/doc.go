// Package budget implements an error-budget tracker for gRPC methods.
//
// An error budget quantifies how much unreliability a service is permitted
// before breaching its Service Level Objective (SLO). For a 99% SLO the
// budget is 1% of all requests — once that 1% is consumed the budget is
// exhausted and further degradation should trigger alerts or policy changes.
//
// Usage:
//
//	s := budget.New(0.99) // 99 % SLO
//	s.Record("/pkg.Service/Method", codes.OK)
//	s.Record("/pkg.Service/Method", codes.Internal)
//
//	if s.Exhausted("/pkg.Service/Method") {
//		log.Println("error budget exhausted:", s.Summary("/pkg.Service/Method"))
//	}
package budget
