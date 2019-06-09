// Copyright © 2018 Bitgrip <berlin@bitgrip.de>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ctl

import "testing"

func Test_isSupportedAPIVersion(t *testing.T) {
	type args struct {
		apiVersion string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "smaller_minor",
			args: args{apiVersion: "1.0"},
			want: true,
		},
		{
			name: "equal_minor",
			args: args{apiVersion: "1.2"},
			want: true,
		},
		{
			name: "bigger_minor",
			args: args{apiVersion: "1.3"},
			want: false,
		},
	}
	oldVersion := Version
	Version = "v1.2.0-local"
	defer func() {
		Version = oldVersion
	}()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSupportedAPIVersion(tt.args.apiVersion); got != tt.want {
				t.Errorf("isSupportedAPIVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
