package repository

import (
	"context"
	"fmt"
	"github.com/nvbn/termonizer/internal/model"
	"log"
	"strconv"
	"time"
)

var defaultPeriodToAmount = map[model.Period]int{
	model.Year:    4,
	model.Quarter: 4,
	model.Week:    4,
	model.Day:     5,
}

const periodToAmountPrefix = "period_to_amount_"

type settingsStore interface {
	ReadSettings(ctx context.Context) ([]model.Setting, error)
	UpdateSetting(ctx context.Context, setting model.Setting) error
}

type Settings struct {
	timeNow func() time.Time
	storage settingsStore

	periodToAmount map[model.Period]int
}

func NewSettings(ctx context.Context, timeNow func() time.Time, storage settingsStore) (*Settings, error) {
	s := &Settings{
		timeNow:        timeNow,
		storage:        storage,
		periodToAmount: defaultPeriodToAmount,
	}

	if err := s.init(ctx); err != nil {
		return nil, err
	}

	return s, nil
}

// improve if there'll be too many settings
func (s *Settings) init(ctx context.Context) error {
	lowLevel, err := s.storage.ReadSettings(ctx)
	if err != nil {
		return err
	}

	kvLowLevel := make(map[string]string)
	for _, entry := range lowLevel {
		kvLowLevel[entry.ID] = entry.Value
	}

	for _, period := range model.Periods {
		key := fmt.Sprintf("%s%d", periodToAmountPrefix, period)
		if value, ok := kvLowLevel[key]; ok {
			intValue, err := strconv.Atoi(value)
			if err != nil {
				log.Printf("invalid setting %s value %s", key, value)
				continue
			}

			s.periodToAmount[period] = intValue
		}
	}

	return nil
}

func (s *Settings) GetAmountForPeriod(period model.Period) int {
	return s.periodToAmount[period]
}

func (s *Settings) SetAmountForPeriod(ctx context.Context, period model.Period, amount int) error {
	s.periodToAmount[period] = amount

	if err := s.storage.UpdateSetting(ctx, model.Setting{
		ID:      fmt.Sprintf("%s%d", periodToAmountPrefix, period),
		Value:   fmt.Sprintf("%d", amount),
		Updated: s.timeNow(),
	}); err != nil {
		return fmt.Errorf("unable to update setting: %w", err)
	}

	return nil
}
