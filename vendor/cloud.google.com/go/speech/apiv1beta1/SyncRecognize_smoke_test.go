// Copyright 2017, Google Inc. All rights reserved.
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

// AUTO-GENERATED CODE. DO NOT EDIT.

package speech

import (
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1beta1"
)

import (
	"testing"

	"cloud.google.com/go/internal/testutil"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var _ = iterator.Done

func TestSpeechSmoke(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping smoke test in short mode")
	}
	ctx := context.Background()
	ts := testutil.TokenSource(ctx, DefaultAuthScopes()...)
	if ts == nil {
		t.Skip("Integration tests skipped. See CONTRIBUTING.md for details")
	}

	projectId := testutil.ProjID()
	uidSpace := testutil.NewUIDSpace("TestSpeechSmoke")
	_, _ = projectId, uidSpace

	c, err := NewClient(ctx, option.WithTokenSource(ts))
	if err != nil {
		t.Fatal(err)
	}

	var languageCode string = "en-US"
	var sampleRate int32 = 44100
	var encoding speechpb.RecognitionConfig_AudioEncoding = speechpb.RecognitionConfig_FLAC
	var config = &speechpb.RecognitionConfig{
		LanguageCode: languageCode,
		SampleRate:   sampleRate,
		Encoding:     encoding,
	}
	var uri string = "gs://gapic-toolkit/hello.flac"
	var audio = &speechpb.RecognitionAudio{
		AudioSource: &speechpb.RecognitionAudio_Uri{
			Uri: uri,
		},
	}
	var request = &speechpb.SyncRecognizeRequest{
		Config: config,
		Audio:  audio,
	}

	resp, err := c.SyncRecognize(ctx, request)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}
