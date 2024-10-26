// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package modelprocessor_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/elastic/apm-data/model/modelpb"
	"github.com/elastic/apm-server/internal/model/modelprocessor"
)

func TestSetLibraryFrames(t *testing.T) {
	processor := modelprocessor.SetLibraryFrame{
		Pattern: regexp.MustCompile("foo"),
	}

	tests := []struct {
		input, output modelpb.Batch
	}{{
		input:  modelpb.Batch{{Error: &modelpb.Error{}}, {Transaction: &modelpb.Transaction{}}},
		output: modelpb.Batch{{Error: &modelpb.Error{}}, {Transaction: &modelpb.Transaction{}}},
	}, {
		input: modelpb.Batch{{
			Span: &modelpb.Span{
				Stacktrace: []*modelpb.StacktraceFrame{
					{LibraryFrame: true, Filename: "foo.go"},
					{LibraryFrame: false, AbsPath: "foobar.go"},
					{LibraryFrame: true, Filename: "bar.go"},
					{LibraryFrame: true},
				},
			},
		}},
		output: modelpb.Batch{{
			Span: &modelpb.Span{
				Stacktrace: []*modelpb.StacktraceFrame{
					{LibraryFrame: true, Filename: "foo.go", Original: &modelpb.Original{LibraryFrame: true}},
					{LibraryFrame: true, AbsPath: "foobar.go", Original: &modelpb.Original{LibraryFrame: false}},
					{LibraryFrame: false, Filename: "bar.go", Original: &modelpb.Original{LibraryFrame: true}},
					{LibraryFrame: false, Original: &modelpb.Original{LibraryFrame: true}},
				},
			},
		}},
	}, {
		input: modelpb.Batch{{
			Error: &modelpb.Error{
				Log: &modelpb.ErrorLog{
					Stacktrace: []*modelpb.StacktraceFrame{
						{LibraryFrame: true, Filename: "foo.go"},
						{LibraryFrame: false, AbsPath: "foobar.go"},
						{LibraryFrame: true, Filename: "bar.go"},
						{LibraryFrame: true},
					},
				},
			},
		}, {
			Error: &modelpb.Error{
				Exception: &modelpb.Exception{
					Stacktrace: []*modelpb.StacktraceFrame{
						{LibraryFrame: true, Filename: "foo.go"},
					},
					Cause: []*modelpb.Exception{{
						Stacktrace: []*modelpb.StacktraceFrame{
							{LibraryFrame: true, Filename: "foo.go"},
						},
					}},
				},
			},
		}},
		output: modelpb.Batch{{
			Error: &modelpb.Error{
				Log: &modelpb.ErrorLog{
					Stacktrace: []*modelpb.StacktraceFrame{
						{LibraryFrame: true, Filename: "foo.go", Original: &modelpb.Original{LibraryFrame: true}},
						{LibraryFrame: true, AbsPath: "foobar.go", Original: &modelpb.Original{LibraryFrame: false}},
						{LibraryFrame: false, Filename: "bar.go", Original: &modelpb.Original{LibraryFrame: true}},
						{LibraryFrame: false, Original: &modelpb.Original{LibraryFrame: true}},
					},
				},
			},
		}, {
			Error: &modelpb.Error{
				Exception: &modelpb.Exception{
					Stacktrace: []*modelpb.StacktraceFrame{
						{LibraryFrame: true, Filename: "foo.go", Original: &modelpb.Original{LibraryFrame: true}},
					},
					Cause: []*modelpb.Exception{{
						Stacktrace: []*modelpb.StacktraceFrame{
							{LibraryFrame: true, Filename: "foo.go", Original: &modelpb.Original{LibraryFrame: true}},
						},
					}},
				},
			},
		}},
	}}

	for _, test := range tests {
		err := processor.ProcessBatch(context.Background(), &test.input)
		assert.NoError(t, err)
		assert.Equal(t, test.output, test.input)
	}

}
