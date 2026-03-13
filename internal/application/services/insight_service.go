package services

import (
	"arandu/internal/domain/insight"
)

type InsightService struct {
	repo insight.Repository
}

func NewInsightService(repo insight.Repository) *InsightService {
	return &InsightService{repo: repo}
}

func (s *InsightService) CreateInsight(content, source string) (*insight.Insight, error) {
	ins := &insight.Insight{
		Content: content,
		Source:  source,
	}
	if err := s.repo.Save(ins); err != nil {
		return nil, err
	}
	return ins, nil
}

func (s *InsightService) GetInsight(id string) (*insight.Insight, error) {
	return s.repo.FindByID(id)
}

func (s *InsightService) ListInsights() ([]*insight.Insight, error) {
	return s.repo.FindAll()
}

func (s *InsightService) DeleteInsight(id string) error {
	return s.repo.Delete(id)
}
