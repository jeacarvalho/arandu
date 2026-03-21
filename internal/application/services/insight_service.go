package services

import (
	"context"

	"arandu/internal/domain/insight"
)

type InsightService struct {
	repo insight.Repository
}

func NewInsightService(repo insight.Repository) *InsightService {
	return &InsightService{repo: repo}
}

func (s *InsightService) CreateInsight(ctx context.Context, content, source string) (*insight.Insight, error) {
	ins := &insight.Insight{
		Content: content,
		Source:  source,
	}
	if err := s.repo.Save(ctx, ins); err != nil {
		return nil, err
	}
	return ins, nil
}

func (s *InsightService) GetInsight(ctx context.Context, id string) (*insight.Insight, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *InsightService) ListInsights(ctx context.Context) ([]*insight.Insight, error) {
	return s.repo.FindAll(ctx)
}

func (s *InsightService) DeleteInsight(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
