// Copyright 2021 Security Scorecard Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pkg

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/ossf/scorecard/v2/checker"
	"github.com/ossf/scorecard/v2/docs/checks"
)

//nolint
func TestSARIFOutput(t *testing.T) {
	t.Parallel()

	type Check struct {
		Risk        string   `yaml:"-"`
		Short       string   `yaml:"short"`
		Description string   `yaml:"description"`
		Remediation []string `yaml:"remediation"`
		Tags        string   `yaml:"tags"`
	}

	commit := "68bc59901773ab4c051dfcea0cc4201a1567ab32"
	date, e := time.Parse(time.RFC822Z, "17 Aug 21 18:57 +0000")
	if e != nil {
		panic(fmt.Errorf("time.Parse: %w", e))
	}

	tests := []struct {
		name        string
		expected    string
		showDetails bool
		logLevel    zapcore.Level
		result      ScorecardResult
		checkDocs   checks.Doc
		minScore    int
	}{
		{
			name:        "check-1",
			showDetails: true,
			expected:    "./testdata/check1.sarif",
			logLevel:    zapcore.DebugLevel,
			minScore:    checker.MaxResultScore,
			checkDocs: checks.Doc{
				Checks: map[string]checks.Check{
					"Check-Name": {
						Risk:        "risk not used",
						Short:       "short description",
						Description: "long description\n other line",
						Remediation: []string{"remediation not used"},
						Tags:        "tag1, tag2 ",
					},
				},
			},
			result: ScorecardResult{
				Repo:      "repo not used",
				Date:      date,
				CommitSHA: commit,
				Checks: []checker.CheckResult{
					{
						Details2: []checker.CheckDetail{
							{
								Type: checker.DetailWarn,
								Msg: checker.LogMessage{
									Text:    "warn message",
									Path:    "src/file1.cpp",
									Type:    checker.FileTypeSource,
									Offset:  5,
									Snippet: "if (bad) {BUG();}",
								},
							},
						},
						Score:  5,
						Reason: "half score reason",
						Name:   "Check-Name",
					},
				},
				Metadata: []string{},
			},
		},
		{
			name:        "check-2",
			showDetails: true,
			expected:    "./testdata/check2.sarif",
			logLevel:    zapcore.DebugLevel,
			minScore:    checker.MaxResultScore,
			checkDocs: checks.Doc{
				Checks: map[string]checks.Check{
					"Check-Name": {
						Risk:        "risk not used",
						Short:       "short description",
						Description: "long description\n other line",
						Remediation: []string{"remediation not used"},
						Tags:        "tag1, tag2 ",
					},
				},
			},
			result: ScorecardResult{
				Repo:      "repo not used",
				Date:      date,
				CommitSHA: commit,
				Checks: []checker.CheckResult{
					{
						Details2: []checker.CheckDetail{
							{
								Type: checker.DetailWarn,
								Msg: checker.LogMessage{
									Text:   "warn message",
									Path:   "bin/binary.elf",
									Type:   checker.FileTypeBinary,
									Offset: 0,
								},
							},
						},
						Score:  checker.MinResultScore,
						Reason: "min score reason",
						Name:   "Check-Name",
					},
				},
				Metadata: []string{},
			},
		},
		{
			name:        "check-3",
			showDetails: true,
			expected:    "./testdata/check3.sarif",
			logLevel:    zapcore.InfoLevel,
			minScore:    checker.MaxResultScore,
			checkDocs: checks.Doc{
				Checks: map[string]checks.Check{
					"Check-Name": {
						Risk:        "risk not used",
						Short:       "short description",
						Description: "long description\n other line",
						Remediation: []string{"remediation not used"},
						Tags:        "tag1, tag2 ",
					},
					"Check-Name2": {
						Risk:        "risk not used",
						Short:       "short description 2",
						Description: "long description\n other line 2",
						Remediation: []string{"remediation not used"},
						Tags:        "tag1, tag2, tag3 ",
					},
					"Check-Name3": {
						Risk:        "risk not used",
						Short:       "short description 3",
						Description: "long description\n other line 3",
						Remediation: []string{"remediation not used"},
						Tags:        "tag1, tag2, tag3, tag4 ",
					},
				},
			},
			result: ScorecardResult{
				Repo:      "repo not used",
				Date:      date,
				CommitSHA: commit,
				Checks: []checker.CheckResult{
					{
						Details2: []checker.CheckDetail{
							{
								Type: checker.DetailWarn,
								Msg: checker.LogMessage{
									Text:   "warn message",
									Path:   "bin/binary.elf",
									Type:   checker.FileTypeBinary,
									Offset: 0,
								},
							},
						},
						Score:  checker.MinResultScore,
						Reason: "min result reason",
						Name:   "Check-Name",
					},
					{
						Details2: []checker.CheckDetail{
							{
								Type: checker.DetailWarn,
								Msg: checker.LogMessage{
									Text:    "warn message",
									Path:    "src/doc.txt",
									Type:    checker.FileTypeText,
									Offset:  3,
									Snippet: "some text",
								},
							},
						},
						Score:  checker.MinResultScore,
						Reason: "min result reason",
						Name:   "Check-Name2",
					},
					{
						Details2: []checker.CheckDetail{
							{
								Type: checker.DetailInfo,
								Msg: checker.LogMessage{
									Text:    "info message",
									Path:    "some/path.js",
									Type:    checker.FileTypeSource,
									Offset:  3,
									Snippet: "if (bad) {BUG();}",
								},
							},
							{
								Type: checker.DetailWarn,
								Msg: checker.LogMessage{
									Text:    "warn message",
									Path:    "some/path.py",
									Type:    checker.FileTypeSource,
									Offset:  3,
									Snippet: "if (bad) {BUG2();}",
								},
							},
							{
								Type: checker.DetailDebug,
								Msg: checker.LogMessage{
									Text:    "debug message",
									Path:    "some/path.go",
									Type:    checker.FileTypeSource,
									Offset:  3,
									Snippet: "if (bad) {BUG5();}",
								},
							},
						},
						Score:  checker.InconclusiveResultScore,
						Reason: "inconclusive reason",
						Name:   "Check-Name3",
					},
				},
				Metadata: []string{},
			},
		},
		{
			name:        "check-4",
			showDetails: true,
			expected:    "./testdata/check4.sarif",
			logLevel:    zapcore.DebugLevel,
			minScore:    checker.MaxResultScore,
			checkDocs: checks.Doc{
				Checks: map[string]checks.Check{
					"Check-Name": {
						Risk:        "risk not used",
						Short:       "short description",
						Description: "long description\n other line",
						Remediation: []string{"remediation not used"},
						Tags:        "tag1, tag2 ",
					},
					"Check-Name2": {
						Risk:        "risk not used",
						Short:       "short description 2",
						Description: "long description\n other line 2",
						Remediation: []string{"remediation not used"},
						Tags:        "tag1, tag2, tag3 ",
					},
					"Check-Name3": {
						Risk:        "risk not used",
						Short:       "short description 3",
						Description: "long description\n other line 3",
						Remediation: []string{"remediation not used"},
						Tags:        "tag1, tag2, tag3, tag4 ",
					},
				},
			},
			result: ScorecardResult{
				Repo:      "repo not used",
				Date:      date,
				CommitSHA: commit,
				Checks: []checker.CheckResult{
					{
						Details2: []checker.CheckDetail{
							{
								Type: checker.DetailWarn,
								Msg: checker.LogMessage{
									Text:   "warn message",
									Path:   "bin/binary.elf",
									Type:   checker.FileTypeBinary,
									Offset: 0,
								},
							},
						},
						Score:  checker.MinResultScore,
						Reason: "min result reason",
						Name:   "Check-Name",
					},
					{
						Details2: []checker.CheckDetail{
							{
								Type: checker.DetailWarn,
								Msg: checker.LogMessage{
									Text:    "warn message",
									Path:    "src/doc.txt",
									Type:    checker.FileTypeText,
									Offset:  3,
									Snippet: "some text",
								},
							},
						},
						Score:  checker.MinResultScore,
						Reason: "min result reason",
						Name:   "Check-Name2",
					},
					{
						Details2: []checker.CheckDetail{
							{
								Type: checker.DetailInfo,
								Msg: checker.LogMessage{
									Text:    "info message",
									Path:    "some/path.js",
									Type:    checker.FileTypeSource,
									Offset:  3,
									Snippet: "if (bad) {BUG();}",
								},
							},
							{
								Type: checker.DetailWarn,
								Msg: checker.LogMessage{
									Text:    "warn message",
									Path:    "some/path.py",
									Type:    checker.FileTypeSource,
									Offset:  3,
									Snippet: "if (bad) {BUG2();}",
								},
							},
							{
								Type: checker.DetailDebug,
								Msg: checker.LogMessage{
									Text:    "debug message",
									Path:    "some/path.go",
									Type:    checker.FileTypeSource,
									Offset:  3,
									Snippet: "if (bad) {BUG5();}",
								},
							},
						},
						Score:  checker.InconclusiveResultScore,
						Reason: "inconclusive reason",
						Name:   "Check-Name3",
					},
				},
				Metadata: []string{},
			},
		},
		{
			name:        "check-5",
			showDetails: true,
			expected:    "./testdata/check5.sarif",
			logLevel:    zapcore.WarnLevel,
			minScore:    5,
			checkDocs: checks.Doc{
				Checks: map[string]checks.Check{
					"Check-Name": {
						Risk:        "risk not used",
						Short:       "short description",
						Description: "long description\n other line",
						Remediation: []string{"remediation not used"},
						Tags:        "tag1, tag2 ",
					},
				},
			},
			result: ScorecardResult{
				Repo:      "repo not used",
				Date:      date,
				CommitSHA: commit,
				Checks: []checker.CheckResult{
					{
						Details2: []checker.CheckDetail{
							{
								Type: checker.DetailWarn,
								Msg: checker.LogMessage{
									Text:    "warn message",
									Path:    "src/file1.cpp",
									Type:    checker.FileTypeSource,
									Offset:  5,
									Snippet: "if (bad) {BUG();}",
								},
							},
						},
						Score:  6,
						Reason: "six score reason",
						Name:   "Check-Name",
					},
				},
				Metadata: []string{},
			},
		},
		{
			name:        "check-6",
			showDetails: true,
			expected:    "./testdata/check6.sarif",
			logLevel:    zapcore.WarnLevel,
			minScore:    checker.MaxResultScore,
			checkDocs: checks.Doc{
				Checks: map[string]checks.Check{
					"Check-Name": {
						Risk:        "risk not used",
						Short:       "short description",
						Description: "long description\n other line",
						Remediation: []string{"remediation not used"},
						Tags:        "tag1, tag2 ",
					},
				},
			},
			result: ScorecardResult{
				Repo:      "repo not used",
				Date:      date,
				CommitSHA: commit,
				Checks: []checker.CheckResult{
					{
						Details2: []checker.CheckDetail{
							{
								Type: checker.DetailWarn,
								Msg: checker.LogMessage{
									Text: "warn message",
									Path: "https://domain.com/something",
									Type: checker.FileTypeURL,
								},
							},
						},
						Score:  6,
						Reason: "six score reason",
						Name:   "Check-Name",
					},
				},
				Metadata: []string{},
			},
		},
	}
	for _, tt := range tests {
		tt := tt // Re-initializing variable so it is not changed while executing the closure below
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var content []byte
			var err error
			content, err = ioutil.ReadFile(tt.expected)
			if err != nil {
				t.Fatalf("cannot read file: %v", err)
			}

			var expected bytes.Buffer
			n, err := expected.Write(content)
			if err != nil {
				t.Fatalf("cannot write buffer: %v", err)
			}
			if n != len(content) {
				t.Fatalf("write %d bytes but expected %d", n, len(content))
			}

			var result bytes.Buffer
			err = tt.result.AsSARIF("1.2.3", tt.showDetails, tt.logLevel, &result, tt.checkDocs, tt.minScore)
			if err != nil {
				t.Fatalf("AsSARIF: %v", err)
			}

			r := bytes.Compare(expected.Bytes(), result.Bytes())
			if r != 0 {
				t.Fatalf("invalid result for %s: %d", tt.name, r)
			}
		})
	}
}
