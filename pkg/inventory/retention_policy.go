package inventory

import "github.com/almanac1631/scrubarr/pkg/domain"

type EvaluationReport struct {
	Result EvaluationReportPart

	Seasons map[int]EvaluationReportPart
	Files   map[int64]EvaluationReportPart
}

type EvaluationReportPart struct {
	Decision domain.Decision
	Tracker  *domain.Tracker
}

type RetentionPolicy interface {
	Evaluate(media LinkedMedia) (EvaluationReport, error)
}
