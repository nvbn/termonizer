package main

import (
	"context"
	"flag"
	"github.com/google/uuid"
	"github.com/nvbn/termonizer/internal/model"
	"github.com/nvbn/termonizer/internal/storage"
	"github.com/nvbn/termonizer/internal/utils"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var dbPath = flag.String("db", "test.db", "path to the database")

var loremIpsum = []string{
	"figure out the approach for resolver",
	"refine project structure",
	"finalize scope for the current sprint",
	"draft technical documentation for the new feature",
	"create project timeline and milestones",
	"set up initial repository and branches",
	"investigate tooling for unit testing",
	"research libraries for input validation",
	"define api endpoints and their contracts",
	"identify potential bottlenecks in architecture",
	"write specifications for core components",
	"set up ci/cd pipelines",
	"review documentation standards",
	"conduct feasibility analysis for feature x",
	"schedule sprint planning meeting",
	"implement the base resolver logic",
	"create mock data for testing",
	"write unit tests for core modules",
	"refactor code to reduce technical debt",
	"fix issues from code review feedback",
	"develop error-handling strategies",
	"ensure all modules meet coding guidelines",
	"integrate api with the backend",
	"test database queries for optimization",
	"finalize feature flags for incremental rollout",
	"debug api response inconsistencies",
	"build ui components for feature x",
	"write integration tests for key workflows",
	"conduct peer code reviews",
	"update and refine logging mechanisms",
	"add advanced configurations to resolver",
	"document edge cases in technical documentation",
	"test system scalability under load",
	"optimize query performance for complex filters",
	"conduct exploratory testing on feature y",
	"finalize designs for new submodules",
	"set up notifications for error thresholds",
	"implement middleware for data validation",
	"collaborate with the design team on ui fixes",
	"add support for multi-language localization",
	"ensure proper versioning of apis",
	"write scripts for database migrations",
	"mock services for testing distributed components",
	"verify authentication and authorization flows",
	"prepare data for analytics tracking",
	"conduct a security audit for the project",
	"optimize resolver for edge cases",
	"set up monitoring tools for live environments",
	"test rollback mechanisms for deployment",
	"fix critical bugs from user testing",
	"finalize release notes for the sprint",
	"push final changes to staging",
	"conduct a dry run of deployment",
	"present progress in a stakeholder meeting",
	"ensure compliance with industry standards",
	"train team members on new updates",
	"resolve performance issues in staging",
	"validate analytics tracking implementation",
	"verify automated backup schedules",
	"finalize team retrospectives",
	"stretch goals (week 5",
	"research emerging technologies for project enhancement",
	"prepare a knowledge-sharing session for the team",
	"document lessons learned from the sprint",
	"improve test coverage to 95%",
	"automate data validation tests",
	"prepare a report on system reliability metrics",
	"evaluate cloud hosting alternatives",
	"investigate potential areas for microservice decomposition",
	"collaborate on the roadmap for the next quarter",
	"write scripts to automate repetitive tasks",
	"improve user feedback mechanisms",
	"test app responsiveness across devices",
	"reduce bundle size for client-side assets",
	"plan a team-building activity",
	"review dependencies for security updates",
	"administrative task",
	"update jira/asana with the latest tasks",
	"sync with product managers for alignment",
	"allocate team bandwidth for the next sprint",
	"review team kpis for the month",
	"address blockers raised in daily standups",
	"organize project files for archiving",
	"plan resource allocation for feature z",
	"approve vendor contracts for tools",
	"follow up on pending feedback from stakeholders",
	"prepare slides for a management update",
	"renew licenses for essential software tools",
	"conduct 1:1s with team members for feedback",
	"complete mandatory compliance training",
	"revise budget estimates for the project",
	"review team performance for appraisal cycles",
	"learning & growt",
	"watch a webinar on ai/ml advancements",
	"complete a course on advanced devops",
	"experiment with a new framework or library",
	"read research papers on distributed systems",
	"contribute to an open-source project",
	"attend a virtual tech conference",
	"host a team knowledge-sharing session",
	"update linkedin with recent accomplishments",
	"reflect on personal development goals",
	"celebrate team wins with a virtual happy hour",
}

func pickRandomN(n int) []string {
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = loremIpsum[rand.Intn(len(loremIpsum))]
	}
	return out
}

func generateContent() string {
	return "* " + strings.Join(pickRandomN(3+rand.Intn(12)), "\n* ")
}

func generateYears() []model.Goal {
	out := make([]model.Goal, 0)

	current := time.Now().Year()
	for year := current - 3; year <= current; year++ {
		start, err := time.Parse("2006", strconv.Itoa(year))
		if err != nil {
			panic(err)
		}
		out = append(out, model.Goal{
			ID:      uuid.New().String(),
			Period:  model.Year,
			Content: generateContent(),
			Start:   start,
			Updated: start,
		})
	}

	return out
}

func generateQuarters() []model.Goal {
	out := make([]model.Goal, 0)

	current := time.Now().Year()
	for year := current - 3; year <= current; year++ {
		for quarter := 1; quarter <= 4; quarter++ {
			start := time.Date(year, time.Month(quarter*3-2), 1, 0, 0, 0, 0, time.Local)
			if start.After(time.Now()) {
				break
			}

			out = append(out, model.Goal{
				ID:      uuid.New().String(),
				Period:  model.Quarter,
				Content: generateContent(),
				Start:   start,
				Updated: start,
			})
		}
	}

	return out
}

func generateWeeks() []model.Goal {
	out := make([]model.Goal, 0)

	start := time.Date(time.Now().Year()-3, 1, 1, 0, 0, 0, 0, time.Local)
	start = utils.WeekStart(start)
	for {
		if start.After(time.Now()) {
			break
		}

		if start.Weekday() == time.Sunday || start.Weekday() == time.Saturday {
			continue
		}

		out = append(out, model.Goal{
			ID:      uuid.New().String(),
			Period:  model.Week,
			Content: generateContent(),
			Start:   start,
			Updated: start,
		})
		start = start.AddDate(0, 0, 7)
	}

	return out
}

func generateDays() []model.Goal {
	out := make([]model.Goal, 0)

	start := time.Now().AddDate(-3, 0, 0)
	for {
		if start.After(time.Now()) {
			break
		}

		out = append(out, model.Goal{
			ID:      uuid.New().String(),
			Period:  model.Day,
			Content: generateContent(),
			Start:   start,
			Updated: start,
		})
		start = start.AddDate(0, 0, 1)
	}

	return out
}

func main() {
	flag.Parse()

	expanded := os.ExpandEnv(*dbPath)
	os.Remove(expanded)

	ctx := context.Background()
	goalsStorage, err := storage.NewSQLite(ctx, expanded)
	if err != nil {
		panic(err)
	}
	defer goalsStorage.Close()

	for _, goal := range generateYears() {
		if err := goalsStorage.UpdateGoal(ctx, goal); err != nil {
			panic(err)
		}
	}

	for _, goal := range generateQuarters() {
		if err := goalsStorage.UpdateGoal(ctx, goal); err != nil {
			panic(err)
		}
	}

	for _, goal := range generateWeeks() {
		if err := goalsStorage.UpdateGoal(ctx, goal); err != nil {
			panic(err)
		}
	}

	for _, goal := range generateDays() {
		if err := goalsStorage.UpdateGoal(ctx, goal); err != nil {
			panic(err)
		}
	}
}
